package simple_pbrpc

import (
	"errors"
	"fmt"
	"github.com/duanhf2012/origin/v2/log"
	"github.com/duanhf2012/origin/v2/node"
	"github.com/duanhf2012/origin/v2/service"
	"github.com/duanhf2012/origin/v2/util/timer"
	"github.com/duanhf2012/origin/v2/util/uuid"
	"math/rand"
	"originserver/common/proto/rpc"
	"time"
)

func init() {
	node.Setup(&TestService9{})
}

type TestService9 struct {
	service.Service
}

func (slf *TestService9) OnInit() error {
	pCron, errCron := timer.NewCronExpr("*/10 * * * * *")
	if errCron != nil {
		return errCron
	}

	pCronCall, errCallCron := timer.NewCronExpr("*/30 * * * * *")
	if errCallCron != nil {
		return errCallCron
	}

	//开始定时器
	slf.CronFunc(pCron, slf.AsyncCallServer8TestOne)
	slf.CronFunc(pCron, slf.AsyncCallServer8TestTwo)
	slf.CronFunc(pCronCall, slf.CallServer8TestOne)
	slf.CronFunc(pCronCall, slf.CallServer8TestTwo)
	slf.AfterFunc(5*time.Second, slf.PrintMsg)
	return nil
}

func (slf *TestService9) PrintMsg(t *timer.Timer) {
	slf.AfterFunc(5*time.Second, slf.PrintMsg)
}

func (slf *TestService9) RPC_Service9TestOne(arg *rpc.TestOne, ret *rpc.TestOneRet) error {
	log.Info("RPC_Service9TestOne", log.Any("arg", arg))
	ret.Msg = arg.Msg
	return nil
}

func (slf *TestService9) RPC_Service9TestTwo(arg *rpc.TestTwo, ret *rpc.TestTwoRet) error {
	log.Info("RPC_Service9TestTwo", log.Any("arg", arg))
	ret.Msg = arg.Msg
	ret.Data = arg.Data
	return nil
}

func (slf *TestService9) RPC_Service9TestThree(arg *rpc.TestOne) error {
	go func() {
		time.Sleep(10 * time.Second)
		log.Info("RPC_Service9TestThree", log.Any("arg", arg))
	}()
	return nil
}

func (slf *TestService9) RPC_Service9TestFour(arg *rpc.TestOne, ret *rpc.TestOneRet) error {
	log.Info("RPC_Service9TestOne", log.Any("arg", arg))
	return errors.New("test error")
}

func (slf *TestService9) RPC_Service9TestFive(arg *rpc.TestOne, ret *rpc.TestOneRet) error {
	panic("test panic")
	return errors.New("test error")
}

func (slf *TestService9) RPC_Service9TestSix(arg *rpc.TestThree) error {
	log.Info("RPC_Service9TestSix", log.Any("arg", arg))

	return nil
}

func (slf *TestService9) AsyncCallServer8TestOne(cron *timer.Cron) {
	for i := 0; i < 10; i++ {
		go func() {
			arg := rpc.TestOne{Msg: uuid.Rand().HexEx()}
			errCall := slf.AsyncCall("TestService8.RPC_Service8TestOne", &arg, func(ret *rpc.TestOneRet, err error) {
				if err != nil || ret.Msg != arg.Msg {
					log.Error(fmt.Sprintf("TestService9 AsyncCallServer8TestOne err[%+v], arg[%+v], ret[%+v]", err, arg, ret))
				}
			})
			if errCall != nil {
				log.Error(fmt.Sprintf("TestService9 AsyncCallServer8TestOne err[%+v]", errCall))
			}
		}()
	}
}

func (slf *TestService9) AsyncCallServer8TestTwo(cron *timer.Cron) {
	for i := 0; i < 10; i++ {
		go func() {
			arg := rpc.TestTwo{Msg: uuid.Rand().HexEx(), Data: int32(rand.Int())}
			errCall := slf.AsyncCall("TestService8.RPC_Service8TestTwo", &arg, func(ret *rpc.TestTwoRet, err error) {
				if err != nil || ret.Msg != arg.Msg || ret.Data != arg.Data {
					log.Errorf("TestService9 AsyncCallServer8TestTwo err[%+v], arg[%+v], ret[%+v]", err, arg, ret)
				}
			})
			if errCall != nil {
				log.Error(fmt.Sprintf("TestService9 AsyncCallServer8TestTwo err[%+v]", errCall))
			}
		}()
	}
}

func (slf *TestService9) CallServer8TestOne(cron *timer.Cron) {
	for i := 0; i < 10; i++ {
		go func() {
			arg := rpc.TestOne{Msg: uuid.Rand().HexEx()}
			ret := rpc.TestOneRet{}
			errCall := slf.Call("TestService8.RPC_Service8TestOne", &arg, &ret)
			if errCall != nil || arg.Msg != ret.Msg {
				log.Error(fmt.Sprintf("TestService9 CallServer8TestOne err[%+v], arg[%+v], ret[%+v]", errCall, &arg, &ret))
			}
		}()
	}
}

func (slf *TestService9) CallServer8TestTwo(cron *timer.Cron) {
	for i := 0; i < 10; i++ {
		go func() {
			arg := rpc.TestTwo{Msg: uuid.Rand().HexEx(), Data: int32(rand.Int())}
			ret := rpc.TestTwoRet{}
			errCall := slf.Call("TestService8.RPC_Service8TestTwo", &arg, &ret)
			if errCall != nil || ret.Msg != arg.Msg || ret.Data != arg.Data {
				log.Errorf("TestService9 CallServer8TestTwo err[%+v], arg[%+v], ret[%+v]", errCall, &arg, &ret)
			}
		}()
	}
}
