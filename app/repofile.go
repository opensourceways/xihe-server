package app

import (
	"encoding/base64"
	"errors"
	"io"
	"sort"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
)

type RepoDir = platform.RepoDir
type UserInfo = platform.UserInfo
type RepoDirInfo = platform.RepoDirInfo
type RepoFileInfo = platform.RepoFileInfo
type RepoPathItem = platform.RepoPathItem
type RepoFileContent = platform.RepoFileContent

type RepoFileService interface {
	List(u *UserInfo, d *RepoDir) ([]RepoPathItem, error)
	Create(*UserInfo, *RepoFileCreateCmd) error
	Update(*UserInfo, *RepoFileUpdateCmd) error
	Delete(*UserInfo, *RepoFileDeleteCmd) error
	Preview(*UserInfo, *RepoFilePreviewCmd) ([]byte, error)
	DeleteDir(*UserInfo, *RepoDirDeleteCmd) (string, error)
	Download(domain.Account, *RepoFileDownloadCmd) (RepoFileDownloadDTO, error)
	DownloadRepo(u *UserInfo, repoId string, handle func(io.Reader, int64)) error
}

func NewRepoFileService(rf platform.RepoFile, sender message.Sender) RepoFileService {
	return &repoFileService{
		rf:     rf,
		sender: sender,
	}
}

type repoFileService struct {
	rf     platform.RepoFile
	sender message.Sender
}

type RepoFileListCmd = RepoDir
type RepoDirDeleteCmd = RepoDirInfo
type RepoFileDeleteCmd = RepoFileInfo
type RepoFilePreviewCmd = RepoFileInfo

type RepoFileDownloadCmd struct {
	Token    string
	Path     domain.FilePath
	Type     domain.ResourceType
	Resource domain.ResourceSummary
}

type RepoFileCreateCmd struct {
	RepoFileInfo

	RepoFileContent
}

type RepoFileUpdateCmd = RepoFileCreateCmd

func (s *repoFileService) Create(u *platform.UserInfo, cmd *RepoFileCreateCmd) error {
	return s.rf.Create(u, &cmd.RepoFileInfo, &cmd.RepoFileContent)
}

func (s *repoFileService) Update(u *platform.UserInfo, cmd *RepoFileUpdateCmd) error {
	data, _, err := s.rf.Download(u.Token, &cmd.RepoFileInfo)
	if err != nil {
		return err
	}

	if b, _ := s.rf.IsLFSFile(data); b {
		return ErrorUpdateLFSFile{
			errors.New("can't update lfs directly"),
		}
	}

	return s.rf.Update(u, &cmd.RepoFileInfo, &cmd.RepoFileContent)
}

func (s *repoFileService) Delete(u *platform.UserInfo, cmd *RepoFileDeleteCmd) error {
	return s.rf.Delete(u, cmd)
}

func (s *repoFileService) DeleteDir(u *platform.UserInfo, cmd *RepoDirDeleteCmd) (
	code string, err error,
) {
	if err = s.rf.DeleteDir(u, cmd); err == nil {
		return
	}

	if platform.IsErrorTooManyFilesToDelete(err) {
		code = ErrorRepoFileTooManyFilesToDelete
	}

	return
}

func (s *repoFileService) Download(who domain.Account, cmd *RepoFileDownloadCmd) (
	RepoFileDownloadDTO, error,
) {
	dto, err := s.download(cmd)
	if err == nil {
		r := &cmd.Resource

		_ = s.sender.AddOperateLogForDownloadFile(
			who, message.RepoFile{
				User: r.Owner,
				Name: r.Name,
				Path: cmd.Path,
			},
		)

		_ = s.sender.IncreaseDownload(&domain.ResourceObject{
			Type:          cmd.Type,
			ResourceIndex: r.ResourceIndex(),
		})
	}

	return dto, err
}

func (s *repoFileService) download(cmd *RepoFileDownloadCmd) (
	dto RepoFileDownloadDTO, err error,
) {
	data, notFound, err := s.rf.Download(cmd.Token, &RepoFileInfo{
		Path:   cmd.Path,
		RepoId: cmd.Resource.RepoId,
	})
	if err != nil {
		if notFound {
			err = ErrorUnavailableRepoFile{err}
		}

		return
	}

	if isLFS, sha := s.rf.IsLFSFile(data); !isLFS {
		dto.Content = base64.StdEncoding.EncodeToString(data)
	} else {
		dto.DownloadURL, err = s.rf.GenLFSDownloadURL(sha)
	}

	return
}

func (s *repoFileService) Preview(u *platform.UserInfo, cmd *RepoFilePreviewCmd) (
	content []byte, err error,
) {
	content, notFound, err := s.rf.Download(u.Token, cmd)
	if err != nil {
		if notFound {
			err = ErrorUnavailableRepoFile{err}
		}

		return
	}

	if isLFS, _ := s.rf.IsLFSFile(content); isLFS {
		err = ErrorPreviewLFSFile{
			errors.New("can't preview the lfs file, download it"),
		}
	}

	return
}

func (s *repoFileService) List(u *UserInfo, d *RepoFileListCmd) ([]RepoPathItem, error) {
	r, err := s.rf.List(u, d)
	if err != nil || len(r) == 0 {
		return nil, err
	}

	sort.Slice(r, func(i, j int) bool {
		a := &r[i]
		b := &r[j]

		if a.IsDir != b.IsDir {
			return a.IsDir
		}

		return a.Name < b.Name
	})

	return r, nil
}

func (s *repoFileService) DownloadRepo(u *UserInfo, repoId string, handle func(io.Reader, int64)) error {
	return s.rf.DownloadRepo(u, repoId, handle)
}

type RepoFileDownloadDTO struct {
	Content     string `json:"content"`
	DownloadURL string `json:"download_url"`
}
