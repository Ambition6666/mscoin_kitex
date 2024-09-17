package task

import (
	"github.com/go-co-op/gocron"
	"jobcenter/config"
	"jobcenter/market"
	"jobcenter/utils"
	"time"
)

type Task struct {
	s *gocron.Scheduler
}

func NewTask() *Task {
	return &Task{
		s: gocron.NewScheduler(time.UTC),
	}
}

func (t *Task) Run() {
	kline := market.NewKline(utils.GetMongoClient(), utils.GetRocketMQProducer(), config.GetConf().Okx, utils.GetRedis()) //开启kafka写
	rate := market.NewRate(utils.GetRedis())
	t.s.Every(60).Seconds().Do(func() {
		kline.Do("1m")
		kline.Do("1H")
		kline.Do("30m")
		kline.Do("15m")
		kline.Do("5m")
		kline.Do("1D")
		kline.Do("1W")
		kline.Do("1M")
		rate.Do()
	})
	t.s.Every(10).Minute().Do(func() {
		market.NewBitCoin().Do()
	})
}

func (t *Task) StartBlocking() {
	t.s.StartBlocking()
}

func (t *Task) Stop() {
	t.s.Stop()
}
