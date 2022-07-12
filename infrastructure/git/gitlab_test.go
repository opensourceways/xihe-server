package git

import (
	"testing"

	"github.com/xanzy/go-gitlab"
)

func TestGitlab(t *testing.T) {
	userGitlabClient, err := NewUserGitlabClient("root", "luonancom@qq.com", "http://192.168.1.193:70")
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal("1---err", err)
	}
	opts := gitlab.ListProjectsOptions{}

	projectlist, _, err := userGitlabClient.Client.Projects.ListProjects(&opts)
	if err != nil {
		t.Log("2---err", err)
	}
	for _, project := range projectlist {
		t.Log("------", project.Name)

	}

}
