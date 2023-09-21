package message

import "github.com/opensourceways/xihe-server/domain"

type ProjectCreatedEvent struct {
	Account     domain.Account
	ProjectName string
}

type ModelCreatedEvent struct {
	Account   domain.Account
	ModelName string
}

type DatasetCreatedEvent struct {
	Account     domain.Account
	DatasetName string
}

type ResourceProducer interface {
	AddOperateLogForCreateResource(domain.ResourceObject, domain.ResourceName) error
	CreateProject(e ProjectCreatedEvent) error
	CreateModel(e ModelCreatedEvent) error
	CreateDataset(e DatasetCreatedEvent) error
	AddRelatedResource(*RelatedResource) error
	RemoveRelatedResource(*RelatedResource) error
	RemoveRelatedResources(*RelatedResources) error
	IncreaseFork(*domain.ResourceIndex) error
}
