package message

import (
	bmdomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	coursedomain "github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/domain"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
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

	CreateTraining(*MsgTraining) error
	CreateFinetune(*domain.FinetuneIndex) error

	CreateInference(*domain.InferenceInfo) error
	ExtendInferenceSurvivalTime(*InferenceExtendInfo) error

	CreateEvaluate(*EvaluateInfo) error

	CalcScore(*SubmissionInfo) error

	SignIn(domain.Account) error
	ApplyCourse(domain.Account, coursedomain.CourseSummary) error
	DailyDownload(domain.Account, domain.ResourceName) error
	DailyLike(domain.Account, domain.ResourceName) error
	LikePicture(domain.Account) error
	DailyCreate(domain.Account, domain.ResourceName) error
	ExperienceBigmodel(domain.Account, bmdomain.BigmodelType) error
	Register(domain.Account) error
	BindEmail(domain.Account, domain.Email) error
	SetAvatarId(domain.Account, userdomain.AvatarId) error
	SetBio(domain.Account, userdomain.Bio) error
	StartJupyter(domain.Account) error
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
