package testground

import (
	"github.com/testground/sdk-go/run"
)

var testcases = map[string]interface{}{
	"simple": run.InitializedTestCaseFn(pingpong),
}

func main() {
	run.InvokeMap(testcases)
}
