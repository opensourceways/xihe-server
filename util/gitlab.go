package util

import (
	"log"

	"github.com/xanzy/go-gitlab"
)

//create a new gitlab client with accessToken
func InitGitlabClient() *gitlab.Client {
	client, err := gitlab.NewClient(GetConfig().GitLabConfig.AccessToken, gitlab.WithBaseURL(GetConfig().GitLabConfig.Baseurl))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil
	}
	_, _, err = client.Projects.ListProjects(nil, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return client
}
