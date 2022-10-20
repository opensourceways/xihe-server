package controller

import "github.com/opensourceways/xihe-server/app"

type RepoFileCreateRequest struct {
	Content       string `json:"content"`
	Base64Encoded bool   `json:"base64_encoded"`
}

func (req *RepoFileCreateRequest) toContent() (cmd app.RepoFileContent) {
	cmd.Content = &req.Content
	cmd.IsEncoded = req.Base64Encoded

	return
}

type RepoFileUpdateRequest = RepoFileCreateRequest
