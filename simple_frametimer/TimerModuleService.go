package simple_frametimer

import (
	"context"
	"github.com/duanhf2012/origin/v2/log"
	"github.com/duanhf2012/origin/v2/node"
	"github.com/duanhf2012/origin/v2/service"
	"github.com/duanhf2012/origin/v2/sysmodule/frametimer"
	"time"
)

func init() {
	node.Setup(&TimerModuleService{})
}

type TimerModuleService struct {
	service.Service
	ft  *frametimer.FrameTimer
	tg1 *frametimer.FrameGroup

	afterTimerID frametimer.FrameTimerID
	tickerID     frametimer.FrameTimerID
}

func (ts *TimerModuleService) OnInit() error {
	ts.ft = &frametimer.FrameTimer{}

	// 如果有加速需求，建议将fps设置的更大，就会更精准，否则按正常设置，设置100时，误差会有10ms
	ts.ft.SetFps(100)
	ts.ft.SetAccuracyInterval(3 * time.Millisecond)

	ts.AddModule(ts.ft)

	return nil
}

var i = 0

func (ts *TimerModuleService) OnStart() {
	// 新建定时器组1
	ts.tg1 = ts.ft.NewGroup()
	// 新建定时器组2
	g2 := ts.ft.NewGroup()

	// AfterFunc定时器
	log.Debug("FrameAfterFunc start")
	now := time.Now()
	ts.tg1.FrameAfterFunc(&ts.afterTimerID, 2*time.Second, nil, func(ctx context.Context, timerID frametimer.FrameTimerID) {
		log.Debug("FrameAfterFunc trigger", log.Any("left:", time.Now().Sub(now).Milliseconds()))
	})

	// Tick定时器
	log.Debug("FrameNewTicker start")
	now = time.Now()

	ts.tg1.SetMultiple(3)
	ts.tg1.FrameNewTicker(&ts.tickerID, 100*time.Millisecond, nil, func(ctx context.Context, tickerID frametimer.FrameTimerID) {
		i++
		log.Debug("FrameNewTicker trigger", log.Any("index", i), log.Any("left:", time.Now().Sub(now).Milliseconds()))
		now = time.Now()

		if i == 100 {
			log.Debug("FrameNewTicker Pause")
			ts.tg1.Pause()

			var id frametimer.FrameTimerID
			g2.FrameAfterFunc(&id, 10*time.Second, nil, func(ctx context.Context, timerID frametimer.FrameTimerID) {
				ts.tg1.Resume()
				log.Debug("FrameNewTicker Resume")

				g2.Close()
			})

		}
	})
}
