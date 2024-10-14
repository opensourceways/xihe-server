package messages

import (
	"fmt"

	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	spaceappdomain "github.com/opensourceways/xihe-server/spaceapp/domain"
	"github.com/opensourceways/xihe-server/utils"
)

var _ message.Sender = (*sender)(nil)

func NewMessageSender(topic *Topics, p common.Publisher) *sender {
	return &sender{topics: *topic, publisher: p}
}

type sender struct {
	topics    Topics
	publisher common.Publisher
}

// Finetune
func (s *sender) CreateFinetune(info *domain.FinetuneIndex) error {
	v := msgFinetune{
		User: info.Owner.Account(),
		Id:   info.Id,
	}

	return s.send(s.topics.Finetune, &v)
}

// Inference
func (s *sender) CreateInference(info *spaceappdomain.InferenceInfo) error {
	v := s.toInferenceMsg(&info.InferenceIndex)
	v.Action = actionCreate
	v.ProjectName = info.ProjectName.ResourceName()
	v.ResourceLevel = info.ResourceLevel
	v.MsgNormal = common.MsgNormal{
		User:      info.Requester,
		Type:      "",
		Desc:      fmt.Sprintf("%s start a inference job for %s", info.Requester, info.ProjectName.ResourceName()),
		CreatedAt: utils.Now(),
	}

	return s.send(s.topics.Inference, &v)

}

func (s *sender) ExtendInferenceSurvivalTime(info *message.InferenceExtendInfo) error {
	v := s.toInferenceMsg(&info.InferenceIndex)
	v.Action = actionExtend
	v.Expiry = info.Expiry
	v.ProjectName = info.ProjectName.ResourceName()
	v.ResourceLevel = info.ResourceLevel

	return s.send(s.topics.Inference, &v)
}

func (s *sender) toInferenceMsg(index *spaceappdomain.InferenceIndex) msgInference {
	return msgInference{
		ProjectId:    index.Project.Id,
		LastCommit:   index.LastCommit,
		InferenceId:  index.Id,
		ProjectOwner: index.Project.Owner.Account(),
	}
}

// Competition
func (s *sender) CalcScore(info *message.SubmissionInfo) error {
	v := msgSubmission{
		CId:   info.Index.Id,
		Phase: info.Index.Phase.CompetitionPhase(),
		SId:   info.Id,
		Path:  info.OBSPath,
	}

	return s.send(s.topics.Submission, &v)
}

// operate log
func (s *sender) AddOperateLogForAccessBigModel(u domain.Account, t bigmodeldomain.BigmodelType) error {
	return s.sendOperateLog(u, "bigmodel", map[string]string{
		"bigmodel": string(t),
	})
}

func (s *sender) AddOperateLogForCloudCreated(u domain.Account, cloudId string) error {
	const operateType = "cloud"
	return s.sendOperateLog(u, operateType, map[string]string{
		"cloud_id": cloudId,
	})
}

func (s *sender) sendOperateLog(u domain.Account, t string, info map[string]string) error {
	return s.send(s.topics.OperateLog, operateLog(u, t, info))
}

// send
func (s *sender) send(topic string, v interface{}) error {
	return s.publisher.Publish(topic, v, nil)
}

func operateLog(u domain.Account, t string, info map[string]string) common.MsgOperateLog {
	a := ""
	if u != nil {
		a = u.Account()
	}

	return common.MsgOperateLog{
		When: utils.Now(),
		User: a,
		Type: t,
		Info: info,
	}
}
