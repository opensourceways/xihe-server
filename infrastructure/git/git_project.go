package git

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/xanzy/go-gitlab"
)

type GitProjectClient struct {
	*GitLabClient
}

func NewGitProjectClient(gitlabclient *GitLabClient) (client *GitProjectClient) {
	client = new(GitProjectClient)
	client.GitLabClient = gitlabclient
	return
}

//CreateProject create project
func (u *GitProjectClient) CreateProject(name, desc, visibility string, mergeRequestsEnabled, snippetsEnabled bool) error {
	var newGitlabProject domain.GitlabProject
	createOpts, err := newGitlabProject.MakeCreateOpt(name, desc, visibility, mergeRequestsEnabled, snippetsEnabled)
	if err != nil {
		return err
	}
	_, _, err = u.Client.Projects.CreateProject(createOpts)

	return err
}

//ListMyStarProject 列出我自己的项目
func (u *GitProjectClient) ListMyProject(userid string, page int) (result []*gitlab.Project, err error) {
	var opts gitlab.ListProjectsOptions
	opts.Page = page
	result, _, err = u.GitLabClient.Projects.ListUserProjects(userid, &opts)
	return
}

//ListMyStarProject 列出 我star了的项目
func (u *GitProjectClient) ListMyStarProject(userid string, page int) (result []*gitlab.Project, err error) {
	var opts gitlab.ListProjectsOptions
	opts.Page = page
	result, _, err = u.GitLabClient.Projects.ListUserStarredProjects(userid, &opts)
	return
}

//UpdateProject
func (u *GitProjectClient) UpdateProject(projectid interface{}, name, desc *string, visibility *gitlab.VisibilityValue, mergeRequestsEnabled, snippetsEnabled *bool) (result *gitlab.Project, err error) {
	var opts gitlab.EditProjectOptions
	opts.Name = name
	opts.Description = desc
	opts.Visibility = visibility
	opts.MergeRequestsEnabled = mergeRequestsEnabled
	opts.SnippetsEnabled = snippetsEnabled
	result, _, err = u.GitLabClient.Projects.EditProject(projectid, &opts)
	return
}

//DeleteProject
func (u *GitProjectClient) DeleteProject(projectid string) (err error) {
	_, err = u.GitLabClient.Projects.DeleteProject(projectid)
	return
}

func parseID(id interface{}) (string, error) {
	switch v := id.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid ID type %#v, the ID must be an int or a string", id)
	}
}

//UploadFile 文件上传
func (u *GitProjectClient) UploadFile(pid interface{}, file multipart.File, filename string) (result *gitlab.ProjectFile, err error) {

	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	URL := fmt.Sprintf("projects/%s/uploads", gitlab.PathEscape(project))

	req, err := u.UploadRequest(
		http.MethodPost,
		URL,
		file,
		filename,
		"file",
	)
	if err != nil {
		return nil, err
	}
	result = new(gitlab.ProjectFile)
	_, err = u.GitLabClient.Do(req, result)

	if err != nil {
		return nil, err
	}
	return result, nil

	// result, _, err = u.GitLabClient.Projects.UploadFile(pid, file, filename)
	// return
}

func (u *GitProjectClient) UploadRequest(method, path string, content io.Reader, filename string, uploadType gitlab.UploadType) (*retryablehttp.Request, error) {
	URL := u.GitLabClient.Client.BaseURL()
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	// Set the encoded path data
	URL.RawPath = URL.Path + path
	URL.Path = URL.Path + unescaped

	// Create a request specific headers map.
	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")
	if u.GitLabClient.Client.UserAgent != "" {
		reqHeaders.Set("User-Agent", u.GitLabClient.Client.UserAgent)
	}
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)

	fw, err := w.CreateFormFile(string(uploadType), filename)
	if err != nil {
		return nil, err
	}
	readSize := 0
	temp := make([]byte, 1024)
	for err == nil {
		readSize, err = io.ReadFull(content, temp)
		fw.Write(temp[:readSize])
	}

	if err = w.Close(); err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", w.FormDataContentType())

	req, err := retryablehttp.NewRequest(method, URL.String(), b)
	if err != nil {
		return nil, err
	}

	// Set the request specific headers.
	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

//DeleteFile 文件删除
func (u *GitProjectClient) DeleteFile(projectid interface{}, fileURL string) (result *gitlab.Response, err error) {
	var opts gitlab.DeleteFileOptions
	opts.CommitMessage = &fileURL
	result, err = u.GitLabClient.RepositoryFiles.DeleteFile(projectid, fileURL, &opts)
	return
}
