package message

import (
	bmdomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/domain"
)

type EvaluateInfo struct {
	domain.EvaluateIndex
	Type    string
	OBSPath string
}

type InferenceExtendInfo struct {
	domain.InferenceInfo
	Expiry int64
}

type SubmissionInfo struct {
	Index   domain.CompetitionIndex
	Id      string
	OBSPath string
}

type RepoFile struct {
	User domain.Account
	Name domain.ResourceName
	Path domain.FilePath
}

type Sender interface {
	AddOperateLogForNewUser(domain.Account) error
	AddOperateLogForCreateTraining(domain.TrainingIndex) error
	AddOperateLogForAccessBigModel(domain.Account, bmdomain.BigmodelType) error
	AddOperateLogForCreateResource(domain.ResourceObject, domain.ResourceName) error
	AddOperateLogForDownloadFile(domain.Account, RepoFile) error

	AddFollowing(*userdomain.FollowerInfo) error
	RemoveFollowing(*userdomain.FollowerInfo) error

	AddLike(*domain.ResourceObject) error
	RemoveLike(*domain.ResourceObject) error

	IncreaseFork(*domain.ResourceIndex) error
	IncreaseDownload(*domain.ResourceObject) error

	AddRelatedResource(*RelatedResource) error
	RemoveRelatedResource(*RelatedResource) error
	RemoveRelatedResources(*RelatedResources) error

	CreateTraining(*domain.TrainingIndex) error
	CreateFinetune(*domain.FinetuneIndex) error

	CreateInference(*domain.InferenceInfo) error
	ExtendInferenceSurvivalTime(*InferenceExtendInfo) error

	CreateEvaluate(*EvaluateInfo) error

	CalcScore(*SubmissionInfo) error
}

type EventHandler interface {
	RelatedResourceHandler
	FollowingHandler
	LikeHandler
	ForkHandler
	DownloadHandler
	TrainingHandler
	FinetuneHandler
	InferenceHandler
	EvaluateHandler
}

type FollowingHandler interface {
	HandleEventAddFollowing(*userdomain.FollowerInfo) error
	HandleEventRemoveFollowing(*userdomain.FollowerInfo) error
}

type LikeHandler interface {
	HandleEventAddLike(*domain.ResourceObject) error
	HandleEventRemoveLike(*domain.ResourceObject) error
}

type ForkHandler interface {
	HandleEventFork(*domain.ResourceIndex) error
}

type DownloadHandler interface {
	HandleEventDownload(*domain.ResourceObject) error
}

type RelatedResourceHandler interface {
	HandleEventAddRelatedResource(*RelatedResource) error
	HandleEventRemoveRelatedResource(*RelatedResource) error
}

type RelatedResource struct {
	Promoter *domain.ResourceObject
	Resource *domain.ResourceObject
}

type RelatedResources struct {
	Promoter  domain.ResourceObject
	Resources []domain.ResourceObjects
}

type TrainingHandler interface {
	HandleEventCreateTraining(*domain.TrainingIndex) error
}

type FinetuneHandler interface {
	HandleEventCreateFinetune(*domain.FinetuneIndex) error
}

type InferenceHandler interface {
	HandleEventCreateInference(*domain.InferenceInfo) error
	HandleEventExtendInferenceSurvivalTime(*InferenceExtendInfo) error
}

type EvaluateHandler interface {
	HandleEventCreateEvaluate(*EvaluateInfo) error
}
