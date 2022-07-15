package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"runtime"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure"
	"github.com/opensourceways/xihe-server/infrastructure/git"
	"github.com/opensourceways/xihe-server/infrastructure/mq"
	"github.com/opensourceways/xihe-server/infrastructure/redis"
	"github.com/xanzy/go-gitlab"
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
	ProjectDB    infrastructure.ProjectDB
	ProjectRedis infrastructure.ProjectRedis
	ProjectGit   git.GitProjectClient
}

func NewPorjectAPP(cfg *config.Config) *ProjectApp {
	app := new(ProjectApp)
	gitlabConfig, err := git.NewGitlabClient(cfg)
	if err != nil {
		log.Fatalln("NewPorjectAPP Error:", err)
	}
	app.ProjectGit = *git.NewGitProjectClient(gitlabConfig)

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

func ProjectLikeCountHandle(ctx context.Context, event cloudevents.Event) {
	fmt.Println("--------receiveFunction--------")
	var imageStatusEvent map[string]interface{}
	err := json.Unmarshal(event.Data(), &imageStatusEvent)
	if err != nil {
		log.Printf(" handleDownloadStatusEvent error : %v     \n", err.Error())
		return
	}
	conn := redis.RedisOpen()
	defer conn.Close()

	likeCount, err := redisgo.Ints(redis.HMGet(fmt.Sprintf("user:%s", imageStatusEvent["user_id"]), imageStatusEvent["project_id"], conn))
	if err != nil {
		log.Printf(" redis.HMGet error : %v     \n", err.Error())
		return
	}
	if likeCount[0] > 0 {
		redis.HMSet(fmt.Sprintf("user:%s", imageStatusEvent["user_id"]), imageStatusEvent["project_id"], 0, conn)
		redis.HMHINCRBY(infrastructure.ProjectLikeCountTemp, imageStatusEvent["project_id"], -1, conn)
	} else {
		redis.HMSet(fmt.Sprintf("user:%s", imageStatusEvent["user_id"]), imageStatusEvent["project_id"], 1, conn)
		redis.HMHINCRBY(infrastructure.ProjectLikeCountTemp, imageStatusEvent["project_id"], 1, conn)

	}

}

func (app *ProjectApp) GetInfo(project_id interface{}) (data domain.Project, err error) {
	//分2步，第一步从db中获取基本信息，
	data, err = app.ProjectDB.GetBaseInfo(project_id)
	if err != nil {
		return
	}
	//--- 第二步从redis 获取likeCount ,downloads 等附加信息
	var likeAccount int
	likeAccount, err = app.ProjectRedis.GetLikeCount(project_id)
	if err != nil {
		return
	}
	data.LikeAccount = domain.LikeAccount(likeAccount)

	return
}
func (app *ProjectApp) Save(data domain.Project) (result domain.Project, err error) {
	//直接保存所有信息到数据库
	result, err = app.ProjectDB.Save(data)
	return
}

func (app *ProjectApp) MulitpleUpload(project_id interface{}, files []*multipart.FileHeader) (result []*gitlab.ProjectFile, err error) {
	//直接保存所有信息到数据库
	// var fileObj multipart.File
	var fileReulst *gitlab.ProjectFile
	for _, file := range files {
		var err error
		fileObj, err := file.Open()
		if err != nil {
			return result, err
		}
		fileReulst, err = app.ProjectGit.UploadFile(project_id, fileObj, file.Filename)
		if err != nil {
			return result, err
		}
		fileObj.Close()
		runtime.GC()
		result = append(result, fileReulst)
	}
	return
}
func (app *ProjectApp) DeleteFile(project_id interface{}, files []*gitlab.ProjectFile) (err error) {

	for _, fileURL := range files {
		_, err = app.ProjectGit.DeleteFile(project_id, fileURL.URL)
		if err != nil {
			return
		}
	}
	return
}
