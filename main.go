package main

import (
	"github.com/duanhf2012/origin/v2/node"
	_ "originserver/simple_asyncdo"
	_ "originserver/simple_event"
	_ "originserver/simple_gin"
	_ "originserver/simple_http"
	_ "originserver/simple_module"
	_ "originserver/simple_pbrpc"
	_ "originserver/simple_rpc"
	"time"
	//导入simple_service模块
	_ "originserver/simple_service"

	_ "originserver/simple_frametimer"
	_ "originserver/simple_templateservice"
)

func main() {

	//打开性能分析报告功能，并设置10秒汇报一次
	node.OpenProfilerReport(time.Second * 10)
	node.Start()
}
