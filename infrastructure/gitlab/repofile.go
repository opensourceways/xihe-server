package gitlab

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	liburl "net/url"
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

func (impl *repoFile) Create(
	u *platform.UserInfo, info *platform.RepoFileInfo,
	content *platform.RepoFileContent,
) error {
	return impl.modify(u, info, http.MethodPost, "create", content)
}

func (impl *repoFile) Update(
	u *platform.UserInfo, info *platform.RepoFileInfo,
	content *platform.RepoFileContent,
) error {
	return impl.modify(u, info, http.MethodPut, "update", content)
}

func (impl *repoFile) Delete(u *platform.UserInfo, info *platform.RepoFileInfo) error {
	opt := impl.toCommitInfo(u, "delete file: "+info.Path.FilePath())

	req, err := impl.newRequest(u.Token, impl.baseURL(info), http.MethodDelete, &opt)
	if err != nil {
		return err
	}

	_, err = impl.cli.ForwardTo(req, nil)

	return err
}

func (impl *repoFile) Download(
	token string, info *platform.RepoFileInfo,
) (data []byte, notFound bool, err error) {
	req, err := http.NewRequest(
		http.MethodGet, impl.baseURL(info)+"/raw?ref="+defaultBranch, nil,
	)
	if err != nil {
		return
	}

	if token != "" {
		h := &req.Header
		h.Add("PRIVATE-TOKEN", token)
	}

	code := 0
	data, code, err = impl.cli.Download(req)
	if err != nil && code == 404 {
		notFound = true
	}

	return
}

func (impl *repoFile) StreamDownload(token string, info *platform.RepoFileInfo, handle func(io.Reader, int64)) error {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)

	if req, err = http.NewRequest(
		http.MethodGet, impl.baseURL(info)+"/raw?ref="+defaultBranch, nil,
	); err != nil {
		return err
	}

	if token != "" {
		h := &req.Header
		h.Add("PRIVATE-TOKEN", token)
	}

	if resp, err = http.DefaultClient.Do(req); err != nil {
		return err
	}

	defer resp.Body.Close()

	handle(resp.Body, resp.ContentLength)

	return nil
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
	return fmt.Sprintf(
		"%s/projects/%s/repository/files/%s",
		endpoint,
		liburl.PathEscape(info.RepoId),
		liburl.PathEscape(info.Path.FilePath()),
	)
}

func (impl *repoFile) newRequest(
	token, url, method string, param interface{},
) (*http.Request, error) {
	var body io.Reader

	if param != nil {
		v, err := utils.JsonMarshal(param)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(v)
	}

	req, err := http.NewRequest(method, url, body)
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
	method, action string, content *platform.RepoFileContent,
) error {
	opt := FileCreateOption{
		CommitInfo: impl.toCommitInfo(u, action+" file: "+info.Path.FilePath()),
		Content:    *content.Content,
	}
	if content.IsEncoded {
		opt.Encoding = "base64"
	}

	req, err := impl.newRequest(u.Token, impl.baseURL(info), method, &opt)
	if err != nil {
		return err
	}

	_, err = impl.cli.ForwardTo(req, nil)

	return err
}

func (impl *repoFile) List(u *platform.UserInfo, info *platform.RepoDir) (
	r []platform.RepoPathItem, err error,
) {
	body := `
{
	"query":"query {
		project(fullPath: \"%s\") {
			repository {
				tree(ref: \"%s\", path: \"%s\") {
					blobs {
						nodes {
							name path lfsOid
						}
					},
					trees {
						nodes {
							name path
						}
					}
				}
			}
		}
	}"
}
`
	data := fmt.Sprintf(
		body,
		u.User.Account()+"/"+info.RepoName.ResourceName(),
		defaultBranch, info.Path.Directory(),
	)

	data = strings.ReplaceAll(data, "\n", "")
	data = strings.ReplaceAll(data, "\t", "")

	req, err := http.NewRequest(http.MethodPost, graphqlEndpoint, strings.NewReader(data))
	if err != nil {
		return
	}

	h := &req.Header
	if u.Token != "" {
		h.Add("Authorization", "Bearer "+u.Token)
	}
	h.Add("Content-Type", "application/json")

	v := graphqlResult{}
	if _, err = impl.cli.ForwardTo(req, &v); err != nil {
		return
	}

	r = v.toRepoPathItems()

	return
}

func (impl *repoFile) DeleteDir(u *platform.UserInfo, info *platform.RepoDirInfo) (err error) {
	body := `
{
	"query":"query {
		project(fullPath: \"%s\") {
			repository {
				tree(ref: \"%s\", recursive: true, path: \"%s\") {
					blobs {
						nodes {
							path
						}
					},
				}
			}
		}
	}"
}
`
	data := fmt.Sprintf(
		body,
		u.User.Account()+"/"+info.RepoName.ResourceName(),
		defaultBranch, info.Path.Directory(),
	)

	data = strings.ReplaceAll(data, "\n", "")
	data = strings.ReplaceAll(data, "\t", "")

	req, err := http.NewRequest(http.MethodPost, graphqlEndpoint, strings.NewReader(data))
	if err != nil {
		return
	}

	h := &req.Header
	if u.Token != "" {
		h.Add("Authorization", "Bearer "+u.Token)
	}
	h.Add("Content-Type", "application/json")

	v := graphqlResult{}
	if _, err = impl.cli.ForwardTo(req, &v); err != nil {
		return
	}

	if v.allFilesCount() >= maxFileCount {
		err = platform.NewErrorTooManyFilesToDelete(
			errors.New("too many files to delete"),
		)

		return
	}

	files := v.allFiles()
	if len(files) == 0 {
		return
	}

	return impl.deleteMultiFiles(u, info, files)
}

func (impl *repoFile) deleteMultiFiles(
	u *platform.UserInfo, info *platform.RepoDirInfo, files []string,
) error {
	actions := make([]Action, len(files))
	for i := range files {
		actions[i] = Action{
			Action:   "delete",
			FilePath: files[i],
		}
	}

	opt := impl.toCommitInfo(u, "delete dir: "+info.Path.Directory())

	commit := Commits{
		CommitInfo: opt,
		Actions:    actions,
	}

	url := fmt.Sprintf(
		"%s/projects/%s/repository/commits",
		endpoint, liburl.PathEscape(info.RepoId),
	)

	req, err := impl.newRequest(u.Token, url, http.MethodPost, &commit)
	if err != nil {
		return err
	}

	_, err = impl.cli.ForwardTo(req, nil)

	return err
}

func (impl *repoFile) GetDirFileInfo(u *platform.UserInfo, info *platform.RepoDirFile) (
	sha string, exist bool, err error,
) {
	body := `
{
	"query":"query {
		project(fullPath: \"%s\") {
			repository {
				tree(ref: \"%s\", path: \"%s\") {
					lastCommit {
						sha
					},
				},
				blobs(ref: \"%s\", paths: [\"%s\"]) {
					nodes {
						path
					},
				}

			}
		}
	}"
}
`
	data := fmt.Sprintf(
		body,
		u.User.Account()+"/"+info.RepoName.ResourceName(),
		defaultBranch, info.Dir.Directory(),
		defaultBranch, info.File.FilePath(),
	)

	data = strings.ReplaceAll(data, "\n", "")
	data = strings.ReplaceAll(data, "\t", "")

	req, err := http.NewRequest(http.MethodPost, graphqlEndpoint, strings.NewReader(data))
	if err != nil {
		return
	}

	h := &req.Header
	if u.Token != "" {
		h.Add("Authorization", "Bearer "+u.Token)
	}
	h.Add("Content-Type", "application/json")

	v := graphqlResult{}
	if _, err = impl.cli.ForwardTo(req, &v); err != nil {
		return
	}

	sha = v.Data.Project.Repo.Tree.LastCommit.SHA
	exist = len(v.Data.Project.Repo.Blobs.Nodes) > 0

	return
}

func (impl *repoFile) DownloadRepo(
	u *platform.UserInfo, repoId string,
	handle func(io.Reader, int64),
) error {
	url := fmt.Sprintf(
		"%s/projects/%s/repository/archive.zip",
		endpoint, liburl.PathEscape(repoId),
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if u.Token != "" {
		h := &req.Header
		h.Add("PRIVATE-TOKEN", u.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	handle(resp.Body, resp.ContentLength)

	return nil
}
