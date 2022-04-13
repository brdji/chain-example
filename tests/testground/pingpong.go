package testground

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/brdji/chain-listener/pkg/consumer"
	"github.com/brdji/chain-listener/pkg/producer"

	"github.com/testground/sdk-go/network"
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
	"github.com/testground/sdk-go/sync"
)

func pingpong(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	runenv.RecordMessage("before sync.MustBoundClient")
	client := initCtx.SyncClient
	netclient := initCtx.NetClient

	oldAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	config := &network.Config{
		// Control the "default" network. At the moment, this is the only network.
		Network: "default",

		// Enable this network. Setting this to false will disconnect this test
		// instance from this network. You probably don't want to do that.
		Enable: true,
		Default: network.LinkShape{
			Latency:   100 * time.Millisecond,
			Bandwidth: 1 << 20, // 1Mib
		},
		CallbackState: "network-configured",
		RoutingPolicy: network.DenyAll,
	}

	runenv.RecordMessage("before netclient.MustConfigureNetwork")
	netclient.MustConfigureNetwork(ctx, config)

	seq := client.MustSignalAndWait(ctx, "ip-allocation", runenv.TestInstanceCount)

	// Make sure that the IP addresses don't change unless we request it.
	if newAddrs, err := net.InterfaceAddrs(); err != nil {
		return err
	} else if !sameAddrs(oldAddrs, newAddrs) {
		return fmt.Errorf("interfaces changed")
	}

	runenv.RecordMessage("I am %d", seq)

	ipC := byte((seq >> 8) + 1)
	ipD := byte(seq)

	config.IPv4 = runenv.TestSubnet
	config.IPv4.IP = append(config.IPv4.IP[0:2:2], ipC, ipD)
	config.IPv4.Mask = []byte{255, 255, 255, 0}
	config.CallbackState = "ip-changed"

	var prod producer.Producer
	var consm consumer.Consumer
	if seq == 1 {
		prod = &producer.DummyProducer{
			IdGen: 1,
		}
	} else if seq == 2 {
		consm = &consumer.DummyConsumer{}
	} else {
		return fmt.Errorf("Expected at most two test instances")
	}

	_ = prod
	_ = consm

	runenv.RecordMessage("before reconfiguring network")
	netclient.MustConfigureNetwork(ctx, config)

	switch seq {
	case 1:
		runenv.RecordMessage("Doing nothing")
	case 2:
		runenv.RecordMessage("Listening for messages")
		go func() { consm.ListenForMessages() }()
	default:
		return fmt.Errorf("expected at most two test instances")
	}
	if err != nil {
		return err
	}

	testFunc := func(test string, rttMin, rttMax time.Duration) error {
		buf := make([]byte, 1)

		if seq == 1 {
			prod.ProduceMessage("Test message")
		} else if seq == 2 {

		}

		runenv.RecordMessage("done")

		// check the sequence number.
		if buf[0] != byte(seq) {
			return fmt.Errorf("read unexpected value")
		}

		// Don't reconfigure the network until we're done with the first test.
		state := sync.State("ping-pong-" + test)
		client.MustSignalAndWait(ctx, state, runenv.TestInstanceCount)

		return nil
	}
	err = testFunc("200", 200*time.Millisecond, 215*time.Millisecond)
	if err != nil {
		return err
	}

	config.Default.Latency = 10 * time.Millisecond
	config.CallbackState = "latency-reduced"
	netclient.MustConfigureNetwork(ctx, config)

	return nil
}

func sameAddrs(a, b []net.Addr) bool {
	if len(a) != len(b) {
		return false
	}
	aset := make(map[string]bool, len(a))
	for _, addr := range a {
		aset[addr.String()] = true
	}
	for _, addr := range b {
		if !aset[addr.String()] {
			return false
		}
	}
	return true
}
