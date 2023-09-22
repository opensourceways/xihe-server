package messageadapter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	bigmodelType = "bigmodel_type"
)

func NewMessageAdapter(cfg *Config, p common.Publisher) *messageAdapter {
	return &messageAdapter{cfg: *cfg, publisher: p}
}

type messageAdapter struct {
	cfg       Config
	publisher common.Publisher
}

func (impl *messageAdapter) SendWuKongInferenceStart(v *domain.WuKongInferenceStartEvent) error {
	cfg := &impl.cfg.InferenceStart

	msg := common.MsgNormal{
		User: v.Account.Account(),
		Type: cfg.Name,
		Desc: fmt.Sprintf("Tried a WuKong with %s pictures inference", v.EsStyle),
		Details: map[string]string{
			"status":    "waiting",
			"task_type": v.EsStyle,
			"style":     v.Style,
			"desc":      v.Desc.WuKongPictureDesc(),
		},
	}

	logrus.Debugf("Send WuKongInferenceStart: %v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

func (impl *messageAdapter) SendWuKongInferenceError(v *domain.WuKongInferenceErrorEvent) error {
	cfg := &impl.cfg.InferenceError

	msg := common.MsgNormal{
		User: v.Account.Account(),
		Details: map[string]string{
			"task_id": strconv.Itoa(int(v.TaskId)),
			"status":  "error",
			"error":   v.ErrMsg,
		},
	}

	logrus.Debugf("Send WuKongInferenceError: %v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

func (impl *messageAdapter) SendWuKongAsyncTaskStart(v *domain.WuKongAsyncTaskStartEvent) error {
	cfg := &impl.cfg.InferenceAsyncStart

	msg := common.MsgNormal{
		User: v.Account.Account(),
		Details: map[string]string{
			"status":  "running",
			"task_id": strconv.Itoa(int(v.TaskId)),
		},
	}

	logrus.Debugf("Send WuKongAsyncTaskStart: %v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

func (impl *messageAdapter) SendWuKongAsyncInferenceFinish(
	v *domain.WuKongAsyncInferenceFinishEvent,
) error {
	cfg := &impl.cfg.InferenceAsyncFinish

	var ls string
	for k := range v.Links { // TODO: Move it into domain.service
		ls += v.Links[k] + ","
	}

	msg := common.MsgNormal{
		User: v.Account.Account(),
		Details: map[string]string{
			"task_id": strconv.Itoa(int(v.TaskId)),
			"status":  "finished",
			"links":   strings.TrimRight(ls, ","),
		},
	}

	logrus.Debugf("Send WuKongAsyncInferenceFinish: %v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

func (impl *messageAdapter) SendBigModelStarted(v *domain.BigModelStartedEvent) error {
	cfg := &impl.cfg.BigModelStarted

	msg := common.MsgNormal{
		User: v.Account.Account(),
		Type: cfg.Name,
		Desc: fmt.Sprintf("Tried ad bigmodel %s", string(v.BigModelType)),
		Details: map[string]string{
			bigmodelType: string(v.BigModelType),
		},
		CreatedAt: utils.Now(),
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

func (impl *messageAdapter) SendWuKongPicturePublicized(v *domain.WuKongPicturePublicizedEvent) error {
	cfg := &impl.cfg.PicturePublicized

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		CreatedAt: utils.Now(),
	}

	logrus.Debugf("Send WuKongPicturePublicized: %v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Picture Liked
func (impl *messageAdapter) SendWuKongPictureLiked(v *domain.WuKongPictureLikedEvent) error {
	cfg := &impl.cfg.PictureLiked

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      "AI Picture Liked",
		CreatedAt: utils.Now(),
	}

	logrus.Debugf("Send WuKongPictureLiked: %v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Config
type Config struct {
	// wukong
	InferenceStart       common.TopicConfig `json:"inference_start"`
	InferenceError       common.TopicConfig `json:"inference_error"`
	InferenceAsyncStart  common.TopicConfig `json:"inference_async_start"`
	InferenceAsyncFinish common.TopicConfig `json:"inference_async_finish"`
	PicturePublicized    common.TopicConfig `json:"picture_publicized"`
	PictureLiked         common.TopicConfig `json:"picture_liked"`

	// common
	BigModelStarted common.TopicConfig `json:"bigmodel_started"`
}
