package app

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/infrastructure/mq"
	"github.com/opensourceways/xihe-server/infrastructure/redis"
)

func TestLikeCount(t *testing.T) {
	cfg, err := config.LoadConfig("../conf/app.conf.yaml")
	if err != nil {
		t.Fatal(err)
	}
	mq.InitMQ(cfg)
	redis.InitRedis(cfg)
	app := NewPorjectAPP(cfg)

	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go app.LikeCountIncrease(&wg, fmt.Sprintf("project%d", 1), fmt.Sprintf("user%d", i))
	}
	wg.Wait()
	time.Sleep(time.Minute)
}

func TestReceive(t *testing.T) {
	cfg, err := config.LoadConfig("../conf/app.conf.yaml")
	if err != nil {
		t.Fatal(err)
	}
	mq.InitMQ(cfg)
	redis.InitRedis(cfg)
	mq.StartEventLinsten(mq.ProjectLikeCountIncreaseEvent, "thisguroup", ProjectLikeCountHandle)

}
