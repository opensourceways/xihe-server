package messagequeue

import (
	"encoding/json"
	"strconv"

	kfk "github.com/opensourceways/kafka-lib/agent"
	asyncapp "github.com/opensourceways/xihe-extra-services/async-server/app"
	asyncdomain "github.com/opensourceways/xihe-extra-services/async-server/domain"
	asyncrepo "github.com/opensourceways/xihe-extra-services/async-server/domain/repository"
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
)

const (
	retryNum = 3

	handleNameWuKongInferenceStart  = "wukong_inference_start"
	handleNameWuKongInferenceError  = "wukong_inference_error"
	handleNameWuKongAsyncTaskStart  = "wukong_async_task_start"
	handleNameWuKongAsyncTaskFinish = "wukong_async_task_finish"
)

func Subscribe(s asyncapp.AsyncMessageService, topics *TopicConfig) (err error) {
	c := &consumer{s: s}

	// wukong inference start
	if err = kfk.SubscribeWithStrategyOfRetry(
		handleNameWuKongInferenceStart,
		c.handleEventBigModelWuKongInferenceStart,
		[]string{topics.InferenceStart}, retryNum,
	); err != nil {
		return
	}

	// wukong inference error
	if err = kfk.SubscribeWithStrategyOfRetry(
		handleNameWuKongInferenceError,
		c.handleEventBigModelWuKongInferenceError,
		[]string{topics.InferenceError}, retryNum,
	); err != nil {
		return
	}

	// wukong async task start
	if err = kfk.SubscribeWithStrategyOfRetry(
		handleNameWuKongAsyncTaskStart,
		c.handleEventBigModelWuKongAsyncTaskStart,
		[]string{topics.InferenceAsyncStart}, retryNum,
	); err != nil {
		return
	}

	// wukong async task finish
	err = kfk.SubscribeWithStrategyOfRetry(
		handleNameWuKongAsyncTaskFinish,
		c.handleEventBigModelWuKongAsyncTaskFinish,
		[]string{topics.InferenceAsyncFinish}, retryNum,
	)

	return
}

type consumer struct {
	s asyncapp.AsyncMessageService
}

func (c *consumer) handleEventBigModelWuKongInferenceStart(body []byte, h map[string]string) (err error) {
	b := comsg.MsgNormal{}
	if err = json.Unmarshal(body, &b); err != nil {
		return
	}

	user, err := domain.NewAccount(b.User)
	if err != nil {
		return err
	}

	desc, err := bigmodeldomain.NewWuKongPictureDesc(b.Details["desc"])
	if err != nil {
		return err
	}

	tt, err := asyncdomain.NewTaskType(b.Details["task_type"])
	if err != nil {
		return err
	}

	v := asyncdomain.WuKongRequest{
		User:     user,
		TaskType: tt,
		Style:    b.Details["style"],
		Desc:     desc,
	}

	return c.s.CreateWuKongTask(&v)
}

func (c *consumer) handleEventBigModelWuKongInferenceError(body []byte, h map[string]string) (err error) {
	b := comsg.MsgNormal{}
	if err = json.Unmarshal(body, &b); err != nil {
		return
	}

	status, err := asyncdomain.NewTaskStatus(b.Details["status"])
	if err != nil {
		return err
	}

	taskId, err := strconv.ParseUint(b.Details["task_id"], 10, 64)
	if err != nil {
		return err
	}
	v := asyncrepo.WuKongResp{
		WuKongTask: asyncrepo.WuKongTask{
			Id:     taskId,
			Status: status,
		},
	}

	if b.Details != nil {
		if v.Links, err = asyncdomain.NewLinks(b.Details["error"]); err != nil { // TODO do't use links to save error
			return err
		}
	}

	return c.s.UpdateWuKongTask(&v)
}

func (c *consumer) handleEventBigModelWuKongAsyncTaskStart(body []byte, h map[string]string) (err error) {
	b := comsg.MsgNormal{}
	if err = json.Unmarshal(body, &b); err != nil {
		return
	}

	status, err := asyncdomain.NewTaskStatus(b.Details["status"])
	if err != nil {
		return err
	}

	taskId, err := strconv.ParseUint(b.Details["task_id"], 10, 64)
	if err != nil {
		return err
	}

	v := asyncrepo.WuKongResp{
		WuKongTask: asyncrepo.WuKongTask{
			Id:     taskId,
			Status: status,
		},
	}

	return c.s.UpdateWuKongTask(&v)
}

func (c *consumer) handleEventBigModelWuKongAsyncTaskFinish(body []byte, h map[string]string) (err error) {
	b := comsg.MsgNormal{}
	if err = json.Unmarshal(body, &b); err != nil {
		return
	}

	status, err := asyncdomain.NewTaskStatus(b.Details["status"])
	if err != nil {
		return err
	}

	taskId, err := strconv.ParseUint(b.Details["task_id"], 10, 64)
	if err != nil {
		return err
	}
	v := asyncrepo.WuKongResp{
		WuKongTask: asyncrepo.WuKongTask{
			Id:     taskId,
			Status: status,
		},
	}

	if b.Details != nil {
		if v.Links, err = asyncdomain.NewLinks(b.Details["links"]); err != nil {
			return err
		}
	}

	return c.s.UpdateWuKongTask(&v)
}

type TopicConfig struct {
	BigModelFinished     string `json:"bigmodel_finished"        required:"true"`
	InferenceStart       string `json:"inference_start"`
	InferenceError       string `json:"inference_error"`
	InferenceAsyncStart  string `json:"inference_async_start"`
	InferenceAsyncFinish string `json:"inference_async_finish"`
}
