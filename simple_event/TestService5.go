package simple_event

import (
	"fmt"
	"github.com/duanhf2012/origin/v2/event"
	"github.com/duanhf2012/origin/v2/node"
	"github.com/duanhf2012/origin/v2/service"
	"github.com/duanhf2012/origin/v2/util/timer"
	"time"
)

func init() {
	node.Setup(&TestService5{})
}

const (
	//自定义事件类型，必需从event.Sys_Event_User_Define开始
	//event.Sys_Event_User_Define以内给系统预留
	EVENT1 event.EventType = event.Sys_Event_User_Define + 1
)

type TestService5 struct {
	service.Service
}

type TestModule struct {
	service.Module
}

func (slf *TestModule) OnInit() error {
	//在TestModule中注册监听EVENT1事件
	slf.GetEventProcessor().RegEventReceiverFunc(EVENT1, slf.GetEventHandler(), slf.OnModuleEvent)

	return nil
}

// OnModuleEvent 模块监听事件回调
func (slf *TestModule) OnModuleEvent(ev event.IEvent) {
	event := ev.(*event.Event)
	fmt.Printf("OnModuleEvent type :%d data:%+v\n", event.GetEventType(), event.Data)
}

// OnInit 服务初始化函数，在安装服务时，服务将自动调用OnInit函数
func (slf *TestService5) OnInit() error {
	//在服务中注册监听EVENT1类型事件
	slf.RegEventReceiverFunc(EVENT1, slf.GetEventHandler(), slf.OnServiceEvent)
	slf.AddModule(&TestModule{})

	slf.AfterFunc(time.Second*10, slf.TriggerEvent)
	return nil
}

// OnServiceEvent 服务监听事件回调
func (slf *TestService5) OnServiceEvent(ev event.IEvent) {
	event := ev.(*event.Event)
	fmt.Printf("OnServiceEvent type :%d data:%+v\n", event.Type, event.Data)
}

func (slf *TestService5) TriggerEvent(t *timer.Timer) {
	//广播事件，传入event.Event对象，类型为EVENT1,Data可以自定义任何数据
	//这样，所有监听者都可以收到该事件
	slf.GetEventHandler().NotifyEvent(&event.Event{
		Type: EVENT1,
		Data: "event data.",
	})
}
