package messages

import (
	"fmt"
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
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

// Following
func (s *sender) AddFollowing(msg *userdomain.FollowerInfo) error {
	return s.sendFollowing(msg, actionAdd)
}

func (s *sender) RemoveFollowing(msg *userdomain.FollowerInfo) error {
	return s.sendFollowing(msg, actionRemove)
}

func (s *sender) sendFollowing(msg *userdomain.FollowerInfo, action string) error {
	v := msgFollower{
		Action:   action,
		User:     msg.User.Account(),
		Follower: msg.Follower.Account(),
	}

	return s.send(s.topics.Following, &v)
}

// Like
func (s *sender) AddLike(msg *domain.ResourceObject) error {
	return s.sendLike(msg, actionAdd)
}

func (s *sender) RemoveLike(msg *domain.ResourceObject) error {
	return s.sendLike(msg, actionRemove)
}

func (s *sender) sendLike(msg *domain.ResourceObject, action string) error {
	v := msgLike{Action: action}

	toMsgResourceObject(msg, &v.Resource)

	return s.send(s.topics.Like, &v)
}

// Fork
func (s *sender) IncreaseFork(msg *domain.ResourceIndex) error {
	v := new(resourceIndex)
	toMsgResourceIndex(msg, v)

	return s.send(s.topics.Fork, v)
}

// Download
func (s *sender) IncreaseDownload(obj *domain.ResourceObject) error {
	v := new(resourceObject)
	toMsgResourceObject(obj, v)

	return s.send(s.topics.Download, v)
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
func (s *sender) CreateInference(info *domain.InferenceInfo) error {
	v := s.toInferenceMsg(&info.InferenceIndex)
	v.Action = actionCreate
	v.ProjectName = info.ProjectName.ResourceName()
	v.ResourceLevel = info.ResourceLevel

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

func (s *sender) toInferenceMsg(index *domain.InferenceIndex) msgInference {
	return msgInference{
		ProjectId:    index.Project.Id,
		LastCommit:   index.LastCommit,
		InferenceId:  index.Id,
		ProjectOwner: index.Project.Owner.Account(),
	}
}

// Evaluate
func (s *sender) CreateEvaluate(info *message.EvaluateInfo) error {
	v := msgEvaluate{
		Type:         info.Type,
		OBSPath:      info.OBSPath,
		ProjectId:    info.Project.Id,
		TrainingId:   info.TrainingId,
		EvaluateId:   info.Id,
		ProjectOwner: info.Project.Owner.Account(),
	}

	return s.send(s.topics.Evaluate, &v)
}

// RelatedResource
func (s *sender) AddRelatedResource(msg *message.RelatedResource) error {
	return s.sendRelatedResource(msg, actionAdd)
}

func (s *sender) RemoveRelatedResource(msg *message.RelatedResource) error {
	return s.sendRelatedResource(msg, actionRemove)
}

func (s *sender) RemoveRelatedResources(msg *message.RelatedResources) error {
	v := msgRelatedResources{Action: actionRemove}

	toMsgResourceObject(&msg.Promoter, &v.Promoter)

	v.Resources = make([]resourceObjects, len(msg.Resources))
	for i := range msg.Resources {
		toMsgResourceObjects(&msg.Resources[i], &v.Resources[i])
	}

	return s.send(s.topics.RelatedResource, &v)
}

func (s *sender) sendRelatedResource(msg *message.RelatedResource, action string) error {
	v := msgRelatedResources{Action: action}

	toMsgResourceObject(msg.Promoter, &v.Promoter)

	v.Resources = []resourceObjects{
		{
			Type: msg.Resource.Type.ResourceType(),
			Objects: []resourceIndex{
				{
					Owner: msg.Resource.Owner.Account(),
					Id:    msg.Resource.Id,
				},
			},
		},
	}

	return s.send(s.topics.RelatedResource, &v)
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

// Sign In
func (s *sender) SignIn(u domain.Account) error {
	return s.send(s.topics.SignIn.Topic, &common.MsgNormal{
		Type:      s.topics.SignIn.Name,
		User:      u.Account(),
		CreatedAt: utils.Now(),
	})
}

// Daily Like
func (s *sender) DailyLike(u domain.Account, n domain.ResourceType) error {
	desc := fmt.Sprintf("Like Rousource %s", n)

	return s.send(s.topics.DailyLike.Topic, &common.MsgNormal{
		Type:      s.topics.DailyLike.Name,
		User:      u.Account(),
		Desc:      desc,
		CreatedAt: utils.Now(),
	})
}

// operate log
func (s *sender) AddOperateLogForNewUser(u domain.Account) error {
	return s.sendOperateLog(u, "user", nil)
}

func (s *sender) AddOperateLogForAccessBigModel(u domain.Account, t bigmodeldomain.BigmodelType) error {
	return s.sendOperateLog(u, "bigmodel", map[string]string{
		"bigmodel": string(t),
	})
}

func (s *sender) AddOperateLogForCreateResource(
	obj domain.ResourceObject, name domain.ResourceName,
) error {
	return s.sendOperateLog(obj.Owner, "resource", map[string]string{
		"id":   obj.Id,
		"name": name.ResourceName(),
		"type": obj.Type.ResourceType(),
	})
}

func (s *sender) AddOperateLogForDownloadFile(u domain.Account, repo message.RepoFile) error {
	return s.sendOperateLog(u, "download", map[string]string{
		"user": repo.User.Account(),
		"repo": repo.Name.ResourceName(),
		"path": repo.Path.FilePath(),
	})
}

func (s *sender) sendOperateLog(u domain.Account, t string, info map[string]string) error {
	a := ""
	if u != nil {
		a = u.Account()
	}

	return s.send(s.topics.OperateLog, &MsgOperateLog{
		When: utils.Now(),
		User: a,
		Type: t,
		Info: info,
	})
}

// send
func (s *sender) send(topic string, v interface{}) error {
	return s.publisher.Publish(topic, v, nil)
}
