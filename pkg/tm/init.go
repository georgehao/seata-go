package tm

import (
	"sync"

	"github.com/seata/seata-go/pkg/remoting/getty"
)

var onceInit sync.Once

// InitTM init seata tm client
func InitTM() {
	onceInit.Do(func() {
		initConfig()
		initRemoting()
	})
}

// todo
// initConfig init config processor
func initConfig() {
}

// initRemoting init rpc client
func initRemoting() {
	getty.InitRpcClient()
}
