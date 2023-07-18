package platform

import (
	"encoding/base64"
	"io"

	"github.com/opensourceways/xihe-server/domain"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
)

type UserOption struct {
	Name     domain.Account
	Email    domain.Email
	Password domain.Password
}

type User interface {
	New(UserOption) (userdomain.PlatformUser, error)
	NewToken(userdomain.PlatformUser) (string, error)
	RefreshToken(string) (string, error)
}

type RepoOption struct {
	Name     domain.ResourceName
	RepoType domain.RepoType
}

func (r *RepoOption) IsNotEmpty() bool {
	return r.Name != nil || r.RepoType != nil
}

type Repository interface {
	New(repo *RepoOption) (string, error)
	Delete(string) error
	Fork(srcRepoId string, Name domain.ResourceName) (string, error)
	Update(repoId string, repo *RepoOption) error
}

type UserInfo struct {
	User  domain.Account
	Email domain.Email
	Token string
}

type RepoFileInfo struct {
	//Namespace string
	RepoId string
	Path   domain.FilePath
}

type RepoFileContent struct {
	Content   *string
	IsEncoded bool
}

type RepoDir struct {
	RepoName domain.ResourceName
	Path     domain.Directory
}

type RepoDirInfo struct {
	RepoDir
	RepoId string
}

type RepoDirFile struct {
	RepoName domain.ResourceName
	Dir      domain.Directory
	File     domain.FilePath
}

type RepoPathItem struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	IsDir     bool   `json:"is_dir"`
	IsLFSFile bool   `json:"is_lfs_file"`
}

type RepoFile interface {
	List(u *UserInfo, d *RepoDir) ([]RepoPathItem, error)
	Create(u *UserInfo, f *RepoFileInfo, content *RepoFileContent) error
	Update(u *UserInfo, f *RepoFileInfo, content *RepoFileContent) error
	Delete(u *UserInfo, f *RepoFileInfo) error
	DeleteDir(u *UserInfo, f *RepoDirInfo) error
	Download(token string, f *RepoFileInfo) (data []byte, notFound bool, err error)
	IsLFSFile(data []byte) (is bool, sha string)
	GenLFSDownloadURL(sha string) (string, error)
	GetDirFileInfo(u *UserInfo, d *RepoDirFile) (sha string, exist bool, err error)
	DownloadRepo(u *UserInfo, repoId string, handle func(io.Reader, int64)) error
}

func (r *RepoFileContent) IsOverSize() bool {
	var decodeSize int
	if r.IsEncoded {
		decodeSize = base64.StdEncoding.DecodedLen(len(*r.Content))
	} else {
		decodeSize = len(*r.Content)
	}

	return decodeSize > 200*1024	// TODO 200KB to config
}
