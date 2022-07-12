package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/opensourceways/xihe-server/config"
	log "github.com/sirupsen/logrus"
)

func TestMQ(t *testing.T) {
	cfg, err := config.LoadConfig("../../conf/app.conf.yaml")
	if err != nil {
		t.Fatalf("LoadConfig error :%v", err)
	}

	InitMQ(cfg)

	// go StartEventLinsten(ProjectLikeCountIncreaseEvent, "thisguroup", testReceiveFunction)

	for i := 0; i < 10; i++ {
		projectLikeCount := make(map[string]interface{})
		projectLikeCount["p_id"] = 2 + i
		projectLikeCount["u_id"] = 3 + i
		err = PushEvent(ProjectLikeCountIncreaseEvent, projectLikeCount)
		if err != nil {
			t.Fatalf("PushEvent error: %v", err.Error())
		}
		time.Sleep(time.Second)
	}
	time.Sleep(time.Minute)
}

func testReceiveFunction(ctx context.Context, event event.Event) {
	var imageStatusEvent map[string]interface{}
	err := json.Unmarshal(event.Data(), &imageStatusEvent)
	if err != nil {
		log.Printf(" handleDownloadStatusEvent error : %v     \n", err.Error())
		return
	}
	fmt.Println("---receive:", imageStatusEvent)

}
