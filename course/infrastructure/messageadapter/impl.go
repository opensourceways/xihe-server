package messageadapter

import (
	"fmt"

	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func MessageAdapter(cfg *Config, p common.Publisher) *messageAdapter {
	return &messageAdapter{cfg: *cfg, publisher: p}
}

type messageAdapter struct {
	cfg       Config
	publisher common.Publisher
}

func (impl *messageAdapter) SendCourseAppliedEvent(v *domain.CourseAppliedEvent) error {
	cfg := &impl.cfg.CourseApplied

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      fmt.Sprintf("applied course of %s", v.CourseName),
		CreatedAt: utils.Now(),
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Config
type Config struct {
	CourseApplied common.TopicConfig `json:"course_applied"`
}
