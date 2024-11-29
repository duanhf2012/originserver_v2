package simple_gin

import (
	"context"
	"github.com/duanhf2012/origin/v2/log"
	"github.com/duanhf2012/origin/v2/node"
	"github.com/duanhf2012/origin/v2/service"
	"github.com/duanhf2012/origin/v2/sysmodule/netmodule/ginmodule"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 新建自定义服务TestService1
type TestGinService struct {
	service.Service
	ginModule ginmodule.GinModule
}

func init() {
	// 注册服务
	node.Setup(&TestGinService{})
}

func (ts *TestGinService) OnInit() error {
	ts.ginModule.Init(":9090", 15*time.Second, nil)
	ts.AddModule(&ts.ginModule)

	//以下是在gin协程中
	ts.ginModule.GET("/test1", func(c *gin.Context) {
		c.String(200, "GET Hello World!")
	})

	//以下是在服务协程中
	ts.ginModule.SetupDataProcessor(&AuthProcessor{})

	//http://127.0.0.1:9090/test2/1234
	ts.ginModule.SafeGET("/test2/:userid", func(c *ginmodule.SafeContext) {
		log.Debug("get userid", log.String("userid", c.Param("userid")))

		//注意:这里都要使用AndDone接口
		c.StringAndDone(200, "Safe GET Hello World!")

		//或者使用以下方式
		//c.String(200, "Safe GET Hello World!")
		//c.Done()
	})

	ts.ginModule.SafePOST("/test3", func(c *ginmodule.SafeContext) {

		var resp struct {
			Code int
			Msg  string
		}

		c.ShouldBindBodyWithJSON(&resp)
		log.Debug("post", log.Any("body", resp))
		resp.Code = 0
		resp.Msg = "ok"

		//注意:这里都要使用AndDone接口
		c.JSONAndDone(http.StatusOK, &resp)
	})

	ts.ginModule.Start()

	return nil
}

func (ts *TestGinService) stopGin() {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	ts.ginModule.Stop(ctx)
}

func (ts *TestGinService) OnRetire() {
	ts.stopGin()
}

func (ts *TestGinService) OnRelease() {
	ts.stopGin()
}
