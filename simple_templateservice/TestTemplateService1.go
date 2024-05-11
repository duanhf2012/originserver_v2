package simple_templateservice

import (
	"fmt"
	"github.com/duanhf2012/origin/v2/cluster"
	"github.com/duanhf2012/origin/v2/log"
	"github.com/duanhf2012/origin/v2/node"
	"github.com/duanhf2012/origin/v2/service"
	"github.com/duanhf2012/origin/v2/util/timer"
	"time"
)

func init() {
	node.Setup(&TestTemplateService1{})
}

type TestTemplateService1 struct {
	service.Service
}

func (ts *TestTemplateService1) OnInit() error {
	ts.RegDiscoverListener(ts)
	ts.AfterFunc(10*time.Second, func(timer *timer.Timer) {
		mapServiceNode := cluster.GetNodeByTemplateServiceName("TemplateService")
		for serviceName, nodeId := range mapServiceNode {
			var input InputData
			input.A = 1
			input.B = 2
			input.ServiceName = serviceName
			var output OutputData
			err := ts.CallNode(nodeId, serviceName+".RPC_Sub", &input, &output)
			if err != nil {
				log.Error(fmt.Sprintf("I'm %s CallNode error", ts.GetName()), log.Any("err", err))
			} else {
				log.Debug(fmt.Sprintf("I'm %s CallNode success", ts.GetName()), log.Any("output", output))
			}

		}
	})
	return nil
}

func (ts *TestTemplateService1) OnDiscoveryService(nodeId string, serviceName []string) {
	log.Debug(fmt.Sprintf("i'm %s,discovery service", ts.GetName()), log.Any("nodeId", nodeId), log.Any("serviceName", serviceName))
}

func (ts *TestTemplateService1) OnUnDiscoveryService(nodeId string, serviceName []string) {
	log.Debug(fmt.Sprintf("i'm %s,undiscovery service", ts.GetName()), log.Any("nodeId", nodeId), log.Any("serviceName", serviceName))
}

func (ts *TestTemplateService1) RPC_Add(input *InputData, output *OutputData) error {
	output.C = input.A + input.B
	output.ServiceName = input.ServiceName
	return nil
}

type InputData struct {
	A int
	B int
	C int

	ServiceName string
}

type OutputData struct {
	C           int
	ServiceName string
}
