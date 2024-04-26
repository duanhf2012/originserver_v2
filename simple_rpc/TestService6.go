package simple_rpc

import (
	"encoding/json"
	"fmt"
	"github.com/duanhf2012/origin/v2/log"
	"github.com/duanhf2012/origin/v2/node"
	"github.com/duanhf2012/origin/v2/rpc"
	"github.com/duanhf2012/origin/v2/service"
	"time"
)

func init() {
	node.Setup(&TestService6{})
}

type TestService6 struct {
	service.Service
}

func (slf *TestService6) OnInit() error {
	slf.RegDiscoverListener(slf)
	slf.RegRawRpc(1, slf.RawRpcCallBack)
	return nil
}

type InputData struct {
	A int
	B int
}

func (slf *TestService6) RPC_Sum(input *InputData, output *int) error {
	*output = input.A + input.B
	//等待1.5s
	time.Sleep(1500 * time.Millisecond)
	return nil
}

func (slf *TestService6) RPC_SyncTest(resp rpc.Responder, input *int, out *int) error {
	go func() {
		time.Sleep(3 * time.Second)
		var output int = *input
		resp(&output, rpc.NilError)
	}()

	return nil
}

func (slf *TestService6) RawRpcCallBack(data []byte) {
	retData := InputData{}
	err := json.Unmarshal(data, &retData)
	fmt.Println(err, retData)
}

func (slf *TestService6) OnDiscoveryService(nodeId string, serviceName []string) {
	log.Debug(">>>> OnDiscoveryService", log.String("nodeId", nodeId), log.Any("serviceName", serviceName))
}

func (slf *TestService6) OnUnDiscoveryService(nodeId string, serviceName []string) {
	log.Debug(">>>> OnUnDiscoveryService", log.String("nodeId", nodeId), log.Any("serviceName", serviceName))
}
