package message

import (
	bmdomain "github.com/opensourceways/xihe-server/bigmodel/domain"
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
	AddOperateLogForAccessBigModel(domain.Account, bmdomain.BigmodelType) error

	AddLike(*domain.ResourceObject) error
	RemoveLike(*domain.ResourceObject) error

	CreateFinetune(*domain.FinetuneIndex) error

	CreateInference(*domain.InferenceInfo) error
	ExtendInferenceSurvivalTime(*InferenceExtendInfo) error

	CreateEvaluate(*EvaluateInfo) error

	CalcScore(*SubmissionInfo) error
}

type EventHandler interface {
	RelatedResourceHandler
	LikeHandler
	ForkHandler
	DownloadHandler
	FinetuneHandler
	InferenceHandler
	EvaluateHandler
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
