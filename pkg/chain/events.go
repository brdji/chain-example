package chain

// func GetMintRoboEvents(instance *contract.Contract, from uint64) []*contract.ContractMintMultipleRobo {
// 	var fromAgain uint64 = from
// 	_ = fromAgain
// 	events, err := instance.FilterMintMultipleRobo(&bind.FilterOpts{
// 		Start: from,
// 		End:   nil,
// 	},
// 		// any address
// 		[]common.Address{})

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var eventList []*contract.ContractMintMultipleRobo = make([]*contract.ContractMintMultipleRobo, 0)
// 	for {
// 		eventList = append(eventList, events.Event)
// 		if !events.Next() {
// 			break
// 		}
// 	}

// 	return eventList

// }
