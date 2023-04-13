package messages

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

const (
	actionAdd    = "add"
	actionRemove = "remove"
	actionCreate = "create"
	actionExtend = "extend"
)

type msgOperateLog struct {
	When int64             `json:"when"`
	User string            `json:"user"`
	Type string            `json:"type"`
	Info map[string]string `json:"info,omitempty"`
}

type msgFollower struct {
	Action   string `json:"action"`
	User     string `json:"user"`
	Follower string `json:"follower"`
}

type msgLike struct {
	Action string `json:"action"`

	Resource resourceObject `json:"resource"`
}

type msgTraining struct {
	User       string `json:"user"`
	ProjectId  string `json:"pid"`
	TrainingId string `json:"rid"`
}

type msgFinetune struct {
	User string `json:"user"`
	Id   string `json:"id"`
}

type msgInference struct {
	Action       string `json:"action"`
	ProjectId    string `json:"pid"`
	LastCommit   string `json:"commit"`
	InferenceId  string `json:"id"`
	ProjectOwner string `json:"owner"`

	msgCreateInference
	msgExtendInference
}

type msgCreateInference struct {
	ProjectName string   `json:"name"`
	ProjectTags []string `json:"tags"`
}

type msgExtendInference struct {
	Expiry int64 `json:"expiry"`
}

type msgEvaluate struct {
	Type         string `json:"type"`
	OBSPath      string `json:"path"`
	ProjectId    string `json:"pid"`
	TrainingId   string `json:"tid"`
	EvaluateId   string `json:"id"`
	ProjectOwner string `json:"owner"`
}

type msgRelatedResources struct {
	Action    string            `json:"action"`
	Promoter  resourceObject    `json:"promoter"`
	Resources []resourceObjects `json:"resources"`
}

func (msg *msgRelatedResources) handle(f func(*message.RelatedResource) error) error {
	promoter := domain.ResourceObject{}
	if err := msg.Promoter.toResourceObject(&promoter); err != nil {
		return err
	}

	relatedResource := message.RelatedResource{
		Promoter: &promoter,
	}

	f1 := func(resource *domain.ResourceObject) error {
		relatedResource.Resource = resource

		return f(&relatedResource)
	}

	for i := range msg.Resources {
		if err := msg.Resources[i].handle(f1); err != nil {
			return err
		}
	}

	return nil
}

type resourceObjects struct {
	Type    string          `json:"type"`
	Objects []resourceIndex `json:"objects"`
}

func (r *resourceObjects) handle(f func(*domain.ResourceObject) error) error {
	t, err := domain.NewResourceType(r.Type)
	if err != nil {
		return err
	}

	obj := domain.ResourceObject{
		Type: t,
	}

	for i := range r.Objects {
		if err = r.Objects[i].toResourceIndex(&obj.ResourceIndex); err != nil {
			return err
		}

		if err := f(&obj); err != nil {
			return err
		}
	}

	return nil
}

func toMsgResourceObjects(v *domain.ResourceObjects, r *resourceObjects) {
	r.Type = v.Type.ResourceType()

	r.Objects = make([]resourceIndex, len(v.Objects))
	for i := range v.Objects {
		toMsgResourceIndex(&v.Objects[i], &r.Objects[i])
	}
}

type resourceObject struct {
	Type string `json:"type"`

	resourceIndex
}

func (r *resourceObject) toResourceObject(obj *domain.ResourceObject) (err error) {
	if err = r.resourceIndex.toResourceIndex(&obj.ResourceIndex); err != nil {
		return
	}

	if obj.Type, err = domain.NewResourceType(r.Type); err != nil {
		return
	}

	return
}

func toMsgResourceObject(v *domain.ResourceObject, r *resourceObject) {
	r.Type = v.Type.ResourceType()

	toMsgResourceIndex(&v.ResourceIndex, &r.resourceIndex)
}

type resourceIndex struct {
	Owner string `json:"owner"`
	Id    string `json:"id"`
}

func (r *resourceIndex) toResourceIndex(obj *domain.ResourceIndex) (err error) {
	obj.Id = r.Id
	obj.Owner, err = domain.NewAccount(r.Owner)

	return
}

func toMsgResourceIndex(v *domain.ResourceIndex, index *resourceIndex) {
	*index = resourceIndex{
		Owner: v.Owner.Account(),
		Id:    v.Id,
	}
}

type msgSubmission struct {
	CId   string `json:"competition_id"`
	Phase string `json:"phase"`
	SId   string `json:"submission_id"`
	Path  string `json:"path"`
}

type msgPodCreate struct {
	User      string `json:"user"`
	PodId     string `json:"pod_id"`
	CloudId   string `json:"cloud_id"`
	CloudName string `json:"cloud_name"`
}
