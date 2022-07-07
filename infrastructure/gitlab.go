package infrastructure

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type GitlabClient struct {
}

func NewGitlabClient(mongodDB *mongo.Database) *GitlabClient {
	repo := new(GitlabClient)

	return repo
}
