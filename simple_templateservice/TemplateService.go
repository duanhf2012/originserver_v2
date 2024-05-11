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
	node.SetupTemplate(func() service.IService {
		return &TemplateService{}
	})
}

type TemplateService struct {
	service.Service
}

func (ts *TemplateService) OnInit() error {
	ts.RegDiscoverListener(ts)

	ts.AfterFunc(10*time.Second, func(timer *timer.Timer) {
		mapNodeId := cluster.GetNodeByServiceName("TestTemplateService1")
		for nodeId := range mapNodeId {
			var input InputData
			input.A = 1
			input.B = 2

			var output OutputData
			err := ts.CallNode(nodeId, "TestTemplateService1.RPC_Add", &input, &output)
			if err != nil {
				log.Error(fmt.Sprintf("I'm %s CallNode error", ts.GetName()), log.Any("err", err))
			} else {
				log.Debug(fmt.Sprintf("I'm %s CallNode success", ts.GetName()), log.Any("output", output))
			}

		}
	})

	return nil
}

func (ts *TemplateService) RPC_Sub(input *InputData, output *OutputData) error {
	output.C = input.A - input.B
	output.ServiceName = ts.GetName()
	return nil
}

func (ts *TemplateService) OnDiscoveryService(nodeId string, serviceName []string) {
	log.Debug(fmt.Sprintf("i'm %s,discovery service", ts.GetName()), log.Any("nodeId", nodeId), log.Any("serviceName", serviceName))
}

func (ts *TemplateService) OnUnDiscoveryService(nodeId string, serviceName []string) {
	log.Debug(fmt.Sprintf("i'm %s,undiscovery service", ts.GetName()), log.Any("nodeId", nodeId), log.Any("serviceName", serviceName))
}
