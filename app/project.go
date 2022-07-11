package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/mq"
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

func (cmd *ProjectCmd) LikeCountIncrease(project_id, user_id string) {
	data := make(map[string]interface{})
	data["project_id"] = project_id
	data["user_id"] = user_id
	go mq.PushEvent(mq.ProjectLikeCountIncreaseEvent, data)
	//-------------
	//to do other logic

}
