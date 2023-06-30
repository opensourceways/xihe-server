package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/inference"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

type InferenceIndex = domain.InferenceIndex
type InferenceDetail = domain.InferenceDetail

type InferenceCreateCmd struct {
	ProjectId     string
	ProjectName   domain.ResourceName
	ProjectOwner  domain.Account
	ResourceLevel string

	InferenceDir domain.Directory
	BootFile     domain.FilePath
}

func (cmd *InferenceCreateCmd) Validate() error {
	b := cmd.ProjectId != "" &&
		cmd.ProjectName != nil &&
		cmd.ProjectOwner != nil &&
		cmd.InferenceDir != nil &&
		cmd.BootFile != nil

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *InferenceCreateCmd) toInference(v *domain.Inference, lastCommit string) {
	v.Project.Id = cmd.ProjectId
	v.LastCommit = lastCommit
	v.ProjectName = cmd.ProjectName
	v.ResourceLevel = cmd.ResourceLevel
	v.Project.Owner = cmd.ProjectOwner
}

type InferenceService interface {
	Create(*UserInfo, *InferenceCreateCmd) (InferenceDTO, string, error)
	Get(info *InferenceIndex) (InferenceDTO, error)
}

func NewInferenceService(
	p platform.RepoFile,
	repo repository.Inference,
	sender message.Sender,
	minSurvivalTime int,
) InferenceService {
	return inferenceService{
		p:               p,
		repo:            repo,
		sender:          sender,
		minSurvivalTime: int64(minSurvivalTime),
	}
}

type inferenceService struct {
	p               platform.RepoFile
	repo            repository.Inference
	sender          message.Sender
	minSurvivalTime int64
}

type InferenceDTO struct {
	expiry     int64
	Error      string `json:"error"`
	AccessURL  string `json:"access_url"`
	InstanceId string `json:"inference_id"`
}

func (dto *InferenceDTO) hasResult() bool {
	return dto.InstanceId != ""
}

func (dto *InferenceDTO) canReuseCurrent() bool {
	return dto.AccessURL != ""
}

func (s inferenceService) Create(u *UserInfo, cmd *InferenceCreateCmd) (
	dto InferenceDTO, sha string, err error,
) {
	sha, b, err := s.p.GetDirFileInfo(u, &platform.RepoDirFile{
		RepoName: cmd.ProjectName,
		Dir:      cmd.InferenceDir,
		File:     cmd.BootFile,
	})
	if err != nil {
		return
	}

	if !b {
		err = ErrorUnavailableRepoFile{
			errors.New("no boot file"),
		}
	}

	instance := new(domain.Inference)
	cmd.toInference(instance, sha)

	dto, version, err := s.check(instance)
	if err != nil {
		return
	}

	if dto.hasResult() {
		if dto.canReuseCurrent() {
			instance.Id = dto.InstanceId
			logrus.Debugf("will reuse the inference instance(%s)", dto.InstanceId)

			err1 := s.sender.ExtendInferenceSurvivalTime(&message.InferenceExtendInfo{
				InferenceInfo: instance.InferenceInfo,
				Expiry:        dto.expiry,
			})
			if err1 != nil {
				logrus.Errorf(
					"extend instance(%s) failed, err:%s",
					dto.InstanceId, err1.Error(),
				)
			}
		}

		return
	}

	if dto.InstanceId, err = s.repo.Save(instance, version); err == nil {
		instance.Id = dto.InstanceId
		err = s.sender.CreateInference(&instance.InferenceInfo)

		return
	}

	if repository.IsErrorDuplicateCreating(err) {
		dto, _, err = s.check(instance)
	}

	return
}

func (s inferenceService) Get(index *InferenceIndex) (dto InferenceDTO, err error) {
	v, err := s.repo.FindInstance(index)

	dto.Error = v.Error
	dto.AccessURL = v.AccessURL
	dto.InstanceId = v.Id

	return
}

func (s inferenceService) check(instance *domain.Inference) (
	dto InferenceDTO, version int, err error,
) {
	v, version, err := s.repo.FindInstances(&instance.Project, instance.LastCommit)
	if err != nil || len(v) == 0 {
		return
	}

	var target *repository.InferenceSummary

	for i := range v {
		item := &v[i]

		if item.Error != "" {
			dto.Error = item.Error
			dto.InstanceId = item.Id

			return
		}

		if target == nil || item.Expiry > target.Expiry {
			target = item
		}
	}

	if target == nil {
		return
	}

	e, n := target.Expiry, utils.Now()
	if n < e && n+s.minSurvivalTime <= e {
		dto.expiry = target.Expiry
		dto.AccessURL = target.AccessURL
		dto.InstanceId = target.Id
	}

	return
}

type InferenceInternalService interface {
	UpdateDetail(*InferenceIndex, *InferenceDetail) error
}

func NewInferenceInternalService(repo repository.Inference) InferenceInternalService {
	return inferenceInternalService{
		repo: repo,
	}
}

type inferenceInternalService struct {
	repo repository.Inference
}

func (s inferenceInternalService) UpdateDetail(index *InferenceIndex, detail *InferenceDetail) error {
	return s.repo.UpdateDetail(index, detail)
}

type InferenceMessageService interface {
	CreateInferenceInstance(*domain.InferenceInfo) error
	ExtendSurvivalTime(*message.InferenceExtendInfo) error
}

func NewInferenceMessageService(
	repo repository.Inference,
	user userrepo.User,
	manager inference.Inference,
) InferenceMessageService {
	return inferenceMessageService{
		repo:    repo,
		user:    user,
		manager: manager,
	}
}

type inferenceMessageService struct {
	repo    repository.Inference
	user    userrepo.User
	manager inference.Inference
}

func (s inferenceMessageService) CreateInferenceInstance(info *domain.InferenceInfo) error {
	v, err := s.user.GetByAccount(info.Project.Owner)
	if err != nil {
		return err
	}

	survivaltime, err := s.manager.Create(&inference.InferenceInfo{
		InferenceInfo: info,
		UserToken:     v.PlatformToken,
	})
	if err != nil {
		return err
	}

	return s.repo.UpdateDetail(
		&info.InferenceIndex,
		&domain.InferenceDetail{Expiry: utils.Now() + int64(survivaltime)},
	)
}

func (s inferenceMessageService) ExtendSurvivalTime(info *message.InferenceExtendInfo) error {
	expiry, n := info.Expiry, utils.Now()
	if expiry < n {
		logrus.Errorf(
			"extend survival time for inference instance(%s) failed, it is timeout.",
			info.Id,
		)

		return nil
	}

	n += int64(s.manager.GetSurvivalTime(&info.InferenceInfo))

	v := int(n - expiry)
	if v < 10 {
		logrus.Debugf(
			"no need to extend survival time for inference instance(%s) in a small range",
			info.Id,
		)

		return nil
	}

	if err := s.manager.ExtendSurvivalTime(&info.InferenceIndex, v); err != nil {
		return err
	}

	return s.repo.UpdateDetail(&info.InferenceIndex, &domain.InferenceDetail{Expiry: n})
}
