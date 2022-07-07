package util

import "testing"

func TestGitlab(t *testing.T) {
	InitConfig("../conf/app.json")
	InitGitlabClient()

}
