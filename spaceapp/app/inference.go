package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/common/domain/allerror"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
	"github.com/opensourceways/xihe-server/spaceapp/domain/inference"
	spacemesage "github.com/opensourceways/xihe-server/spaceapp/domain/message"
	spaceapprepo "github.com/opensourceways/xihe-server/spaceapp/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type InferenceIndex = domain.InferenceIndex
type InferenceDetail = domain.InferenceDetail

type InferenceCreateCmd struct {
	ProjectId     string
	ProjectName   types.ResourceName
	ProjectOwner  types.Account
	ResourceLevel string

	InferenceDir types.Directory
	BootFile     types.FilePath
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

func (cmd *InferenceCreateCmd) toInference(v *domain.Inference, lastCommit, requester string) {
	v.Project.Id = cmd.ProjectId
	v.LastCommit = lastCommit
	v.ProjectName = cmd.ProjectName
	v.ResourceLevel = cmd.ResourceLevel
	v.Project.Owner = cmd.ProjectOwner
	v.Requester = requester
}

type InferenceService interface {
	Create(string, *app.UserInfo, *InferenceCreateCmd) (InferenceDTO, string, error)
	Get(info *InferenceIndex) (InferenceDTO, error)
	CreateSpaceApp(CmdToCreateApp) error
	NotifyIsServing(ctx context.Context, cmd *CmdToNotifyServiceIsStarted) error
	NotifyIsBuilding(ctx context.Context, cmd *CmdToNotifyBuildIsStarted) error
	NotifyStarting(ctx context.Context, cmd *CmdToNotifyStarting) error
	NotifyIsBuildFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error
	NotifyIsStartFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error
}

func NewInferenceService(
	p platform.RepoFile,
	repo spaceapprepo.Inference,
	sender message.Sender,
	minSurvivalTime int,
	spacesender spacemesage.SpaceAppMessageProducer,
	spaceappRepo spaceapprepo.SpaceAppRepository,
	spaceRepo spacerepo.Project,
) InferenceService {
	return inferenceService{
		p:               p,
		repo:            repo,
		sender:          sender,
		spacesender:     spacesender,
		minSurvivalTime: int64(minSurvivalTime),
		spaceappRepo:    spaceappRepo,
		spaceRepo:       spaceRepo,
	}
}

type inferenceService struct {
	p               platform.RepoFile
	repo            spaceapprepo.Inference
	sender          message.Sender
	spacesender     spacemesage.SpaceAppMessageProducer
	minSurvivalTime int64
	spaceappRepo    spaceapprepo.SpaceAppRepository
	spaceRepo       spacerepo.Project
}

func (s inferenceService) Create(user string, owner *app.UserInfo, cmd *InferenceCreateCmd) (
	dto InferenceDTO, sha string, err error,
) {
	sha, b, err := s.p.GetDirFileInfo(owner, &platform.RepoDirFile{
		RepoName: cmd.ProjectName,
		Dir:      cmd.InferenceDir,
		File:     cmd.BootFile,
	})
	if err != nil {
		return
	}

	if !b {
		err = UnavailableRepoFileError{
			errors.New("no boot file"),
		}

		return
	}

	instance := new(domain.Inference)
	cmd.toInference(instance, sha, user)
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

	var target *spaceapprepo.InferenceSummary

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

func (s inferenceService) CreateSpaceApp(cmd CmdToCreateApp) error {
	if err := s.spacesender.SendSpaceAppCreateMsg(&domain.SpaceAppCreateEvent{
		Id:       cmd.SpaceId.Identity(),
		CommitId: cmd.CommitId,
	}); err != nil {
		return err
	}

	fmt.Println("success ====================================== send ======================= space ================================ create")

	return nil
}

type InferenceInternalService interface {
	UpdateDetail(*InferenceIndex, *InferenceDetail) error
}

func NewInferenceInternalService(repo spaceapprepo.Inference) InferenceInternalService {
	return inferenceInternalService{
		repo: repo,
	}
}

type inferenceInternalService struct {
	repo spaceapprepo.Inference
}

func (s inferenceInternalService) UpdateDetail(index *InferenceIndex, detail *InferenceDetail) error {
	return s.repo.UpdateDetail(index, detail)
}

type InferenceMessageService interface {
	CreateInferenceInstance(*domain.InferenceInfo) error
	ExtendSurvivalTime(*message.InferenceExtendInfo) error
}

func NewInferenceMessageService(
	repo spaceapprepo.Inference,
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
	repo    spaceapprepo.Inference
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
		UserToken:     v.PlatformToken.Token,
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

// NotifyIsServing notifies that a service of a SpaceApp has serving.
func (s inferenceService) NotifyIsServing(ctx context.Context, cmd *CmdToNotifyServiceIsStarted) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}

	if err := v.StartServing(cmd.AppURL, cmd.LogURL); err != nil {
		logrus.Errorf("spaceId:%s set space app serving failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	if err := s.spaceappRepo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify serving successful", cmd.SpaceId.Identity())

	return nil
}

// NotifyIsBuilding notifies that the build process of a SpaceApp has started.
func (s inferenceService) NotifyIsBuilding(ctx context.Context, cmd *CmdToNotifyBuildIsStarted) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}

	if err := v.StartBuilding(cmd.LogURL); err != nil {
		logrus.Errorf("spaceId:%s set space app building failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}
	if err := s.spaceappRepo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify building successful", cmd.SpaceId.Identity())
	return nil
}

// NotifyStarting notifies that the build process of a SpaceApp has finished.
func (s inferenceService) NotifyStarting(ctx context.Context, cmd *CmdToNotifyStarting) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}

	if err := v.SetStarting(); err != nil {
		logrus.Errorf("spaceId:%s set space app starting failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	if err := s.spaceappRepo.SaveWithBuildLog(&v, &domain.SpaceAppBuildLog{
		Logs: cmd.Logs,
	}); err != nil {
		logrus.Errorf("spaceId:%s save with build log db failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	logrus.Infof("spaceId:%s notify starting successful, save build logs:%d",
		cmd.SpaceId.Identity(), len(cmd.Logs))
	return nil
}

// NotifyIsBuildFailed notifies change SpaceApp status.
func (s inferenceService) NotifyIsBuildFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}

	if err := v.SetBuildFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Errorf("spaceId:%s set space app %s failed, err:%s",
			cmd.SpaceId.Identity(), cmd.Status.AppStatus(), err)
		return err
	}

	if err := s.spaceappRepo.SaveWithBuildLog(&v, &domain.SpaceAppBuildLog{
		Logs: cmd.Logs,
	}); err != nil {
		logrus.Errorf("spaceId:%s save with build log db failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	logrus.Infof("spaceId:%s notify build failed successful, save build logs:%d",
		cmd.SpaceId.Identity(), len(cmd.Logs))
	return nil
}

// NotifyIsBuildFailed notifies change SpaceApp status.
func (s inferenceService) NotifyIsStartFailed(ctx context.Context, cmd *CmdToNotifyFailedStatus) error {
	v, err := s.getSpaceApp(cmd.SpaceAppIndex)
	if err != nil {
		return err
	}
	if err := v.SetStartFailed(cmd.Status, cmd.Reason); err != nil {
		logrus.Errorf("spaceId:%s set space app %s failed, err:%s",
			cmd.SpaceId.Identity(), cmd.Status.AppStatus(), err)
		return err
	}

	if err := s.spaceappRepo.Save(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify start failed successful", cmd.SpaceId.Identity())
	return nil
}

func (s inferenceService) getSpaceApp(cmd CmdToCreateApp) (domain.SpaceApp, error) {
	space, err := s.spaceRepo.GetById(cmd.SpaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "space not found", err)
		} else {
			err = xerrors.Errorf("failed to get space, err:%w", err)
		}
		logrus.Errorf("spaceId:%s get space failed, err:%s", cmd.SpaceId.Identity(), err)
		return domain.SpaceApp{}, err
	}

	// if space.CommitId != cmd.CommitId {
	// 	err = allerror.New(allerror.ErrorCodeSpaceCommitConflict, "commit conflict",
	// 		xerrors.Errorf("spaceId:%s commit conflict", space.Id.Identity()))
	// 	logrus.Errorf("spaceId:%s latest commitId:%s, old commitId:%s, err:%s",
	// 		cmd.SpaceId.Identity(), space.CommitId, cmd.CommitId, err)
	// 	return domain.SpaceApp{}, err
	// }

	spaceId, err := types.NewIdentity(space.RepoId)
	if err != nil {
		return domain.SpaceApp{}, err
	}

	v, err := s.spaceappRepo.FindBySpaceId(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceAppNotFound, "space app not found", err)
		}
		logrus.Errorf("spaceId:%s get space app failed, err:%s", space.RepoId, err)
		return domain.SpaceApp{}, err
	}
	return v, nil
}
