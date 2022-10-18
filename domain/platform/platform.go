package platform

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserOption struct {
	Name     domain.Account
	Email    domain.Email
	Password domain.Password
}

type User interface {
	New(UserOption) (domain.PlatformUser, error)
	NewToken(domain.PlatformUser) (string, error)
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

type RepoFile interface {
	List(u *UserInfo, info *RepoFileInfo) error
	Create(u *UserInfo, f *RepoFileInfo, content *string) error
	Update(u *UserInfo, f *RepoFileInfo, content *string) error
	Delete(u *UserInfo, f *RepoFileInfo) error
	Download(u *UserInfo, f *RepoFileInfo) (data []byte, notFound bool, err error)
	IsLFSFile(data []byte) (is bool, sha string)
	GenLFSDownloadURL(sha string) (string, error)
}
