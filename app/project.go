package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure"
	"github.com/opensourceways/xihe-server/infrastructure/mq"
	"github.com/opensourceways/xihe-server/infrastructure/redis"
)

type ProjectCmd struct {
	Owner     string
	Name      domain.ProjName
	Desc      domain.ProjDesc
	Type      domain.RepoType
	CoverId   domain.CoverId
	Protocol  domain.ProtocolName
	Training  domain.TrainingSDK
	Inference domain.InferenceSDK
}

type ProjectApp struct {
	ProjectInfra infrastructure.Project
}

func NewPorjectAPP() *ProjectApp {
	app := new(ProjectApp)

	return app
}

func (app *ProjectApp) LikeCountIncrease(wg *sync.WaitGroup, project_id, user_id string) (data map[string]interface{}, err error) {

	//--------------------使用mq 消息队列--------------------------
	data = make(map[string]interface{})
	data["project_id"] = project_id
	data["user_id"] = user_id
	err = mq.PushEvent(mq.ProjectLikeCountIncreaseEvent, data)
	if err != nil {
		return
	}
	if wg != nil {

		wg.Done()
	}
	return
}

func ReceiveFunction(ctx context.Context, event cloudevents.Event) {
	fmt.Println("--------receiveFunction--------")
	var imageStatusEvent map[string]interface{}
	err := json.Unmarshal(event.Data(), &imageStatusEvent)
	if err != nil {
		log.Printf(" handleDownloadStatusEvent error : %v     \n", err.Error())
		return
	}
	conn := redis.RedisOpen()
	defer conn.Close()
	redis.IncreKey(fmt.Sprintf("project:%s:likeCount", imageStatusEvent["project_id"]), conn)
	likeCount, err := redisgo.Ints(redis.HMGet(fmt.Sprintf("user:%s", imageStatusEvent["user_id"]), imageStatusEvent["project_id"], conn))
	if err != nil {
		log.Printf(" redis.HMGet error : %v     \n", err.Error())
		return
	}
	if likeCount[0] > 0 {
		redis.HMSet(fmt.Sprintf("user:%s", imageStatusEvent["user_id"]), imageStatusEvent["project_id"], 0, conn)
	} else {
		redis.HMSet(fmt.Sprintf("user:%s", imageStatusEvent["user_id"]), imageStatusEvent["project_id"], 1, conn)

	}
	// redis.HMHINCRBY(fmt.Sprintf("user:%s", imageStatusEvent["user_id"]), imageStatusEvent["project_id"], 1, conn)

}
