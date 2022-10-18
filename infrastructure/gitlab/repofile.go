package gitlab

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain/platform"
)

const shaLen = 64

func NewRepoFile() platform.RepoFile {
	return &repoFile{
		cli: utils.NewHttpClient(3),
	}
}

type repoFile struct {
	cli utils.HttpClient
}

func (impl *repoFile) Create(u *platform.UserInfo, info *platform.RepoFileInfo, content *string) error {
	return impl.modify(u, info, http.MethodPost, "create", content)
}

func (impl *repoFile) Update(u *platform.UserInfo, info *platform.RepoFileInfo, content *string) error {
	return impl.modify(u, info, http.MethodPut, "update", content)
}

func (impl *repoFile) Delete(u *platform.UserInfo, info *platform.RepoFileInfo) error {
	opt := impl.toCommitInfo(u, "delete file: "+info.Path.FilePath())

	req, err := impl.newRequest(u.Token, info, http.MethodDelete, &opt)
	if err != nil {
		return err
	}

	_, err = impl.cli.ForwardTo(req, nil)

	return err
}

func (impl *repoFile) Download(
	u *platform.UserInfo, info *platform.RepoFileInfo,
) (data []byte, notFound bool, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		impl.baseURL(info)+"/raw?ref="+defaultBranch, nil,
	)
	if err != nil {
		return
	}

	h := &req.Header
	h.Add("PRIVATE-TOKEN", u.Token)

	code := 0
	data, code, err = impl.cli.Download(req)
	if err != nil && code == 404 {
		notFound = true
	}

	return
}

func (impl *repoFile) IsLFSFile(data []byte) (is bool, sha string) {
	if len(data) > 200 {
		return
	}

	v := strings.Split(string(data), "\n")
	if len(v) < 3 {
		return
	}

	line := v[1]
	if !strings.HasPrefix(line, "oid sha256:") {
		return
	}

	sha = strings.TrimPrefix(line, "oid sha256:")
	is = len(sha) == shaLen

	return
}

func (impl *repoFile) GenLFSDownloadURL(sha string) (string, error) {
	if len(sha) != shaLen {
		return "", errors.New("invalid sha")
	}

	return obsHelper.GenObjectDownloadURL(
		filepath.Join(sha[:2], sha[2:4], sha[4:]),
	)
}

func (impl *repoFile) baseURL(info *platform.RepoFileInfo) string {
	return endpoint + fmt.Sprintf(
		"/projects/%s/repository/files/%s", info.RepoId,
		strings.Replace(info.Path.FilePath(), "/", "%2F", -1),
	)
}

func (impl *repoFile) newRequest(
	token string, info *platform.RepoFileInfo,
	method string, param interface{},
) (*http.Request, error) {
	var body io.Reader

	if param != nil {
		v, err := utils.JsonMarshal(param)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(v)
	}

	req, err := http.NewRequest(method, impl.baseURL(info), body)
	if err != nil {
		return nil, err
	}

	h := &req.Header
	h.Add("PRIVATE-TOKEN", token)
	h.Add("Content-Type", "application/json")

	return req, nil
}

func (impl *repoFile) toCommitInfo(u *platform.UserInfo, commitMsg string) CommitInfo {
	return CommitInfo{
		Branch:      defaultBranch,
		Message:     commitMsg,
		AuthorName:  u.User.Account(),
		AuthorEmail: u.Email.Email(),
	}
}

func (impl *repoFile) modify(
	u *platform.UserInfo, info *platform.RepoFileInfo,
	method, action string, content *string,
) error {
	opt := FileCreateOption{
		CommitInfo: impl.toCommitInfo(u, action+" file: "+info.Path.FilePath()),
		Content:    *content,
	}

	req, err := impl.newRequest(u.Token, info, method, &opt)
	if err != nil {
		return err
	}

	_, err = impl.cli.ForwardTo(req, nil)

	return err
}

func (impl *repoFile) List(u *platform.UserInfo, info *platform.RepoFileInfo) error {
	return nil
}
