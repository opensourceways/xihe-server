package mongodb

import (
	"github.com/opensourceways/xihe-server/domain"
)

type cProject struct {
	User      string `bson:"user"      json:"user"`
	Name      string `bson:"name"      json:"name"`
	Desc      string `bson:"desc"      json:"desc"`
	Type      string `bson:"type"      json:"type"`
	CoverId   string `bson:"cover_id"  json:"cover_id"`
	Protocol  string `bson:"protocol"  json:"protocol"`
	Training  string `bson:"training"  json:"training"`
	Inference string `bson:"inference" json:"inference"`
}

type projectRepository struct {
	collectionName string
}

func (p projectRepository) Save(domain.Project) error {
	return nil
}
