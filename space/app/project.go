package app

import (
	"errors"

	sdk "github.com/opensourceways/xihe-sdk/space"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-server/app"
	computilityapp "github.com/opensourceways/xihe-server/computility/app"
	computilitydomain "github.com/opensourceways/xihe-server/computility/domain"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ProjectCreateCmd struct {
	Owner     domain.Account
	Name      domain.ResourceName
	Desc      domain.ResourceDesc
	Title     domain.ResourceTitle
	Type      domain.ProjType
	CoverId   domain.CoverId
	RepoType  domain.RepoType
	Protocol  domain.ProtocolName
	Training  domain.TrainingPlatform
	Tags      []string
	TagKinds  []string
	All       []domain.DomainTags
	Hardware  domain.Hardware
	BaseImage domain.BaseImage
}

func (cmd *ProjectCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.Type != nil &&
		cmd.CoverId != nil &&
		cmd.RepoType != nil &&
		cmd.Protocol != nil &&
		cmd.Training != nil

	if !b {
		return errors.New("invalid cmd of creating project")
	}

	return nil
}

func (cmd *ProjectCreateCmd) genTagKinds(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	r := make([]string, 0, len(cmd.All))

	for i := range cmd.All {
		if v := cmd.All[i].GetKindsOfTags(tags); len(v) > 0 {
			r = append(r, v...)
		}
	}

	return r
}

func (cmd *ProjectCreateCmd) toProject(r *spacedomain.Project) {
	now := utils.Now()
	normTags := []string{cmd.Type.ProjType(),
		cmd.Protocol.ProtocolName(),
		cmd.Training.TrainingPlatform()}

	*r = spacedomain.Project{
		Owner:     cmd.Owner,
		Type:      cmd.Type,
		Protocol:  cmd.Protocol,
		Training:  cmd.Training,
		CreatedAt: now,
		UpdatedAt: now,
		Hardware:  cmd.Hardware,
		BaseImage: cmd.BaseImage,
		ProjectModifiableProperty: spacedomain.ProjectModifiableProperty{
			Name:              cmd.Name,
			Desc:              cmd.Desc,
			Title:             cmd.Title,
			CoverId:           cmd.CoverId,
			RepoType:          cmd.RepoType,
			Tags:              append(cmd.Tags, normTags...),
			TagKinds:          cmd.genTagKinds(cmd.Tags),
			CommitId:          "",
			NoApplicationFile: true,
		},
	}
}

type ProjectService interface {
	CanApplyResourceName(domain.Account, domain.ResourceName) bool
	Create(*ProjectCreateCmd, platform.Repository) (ProjectDTO, error)
	Delete(*spacedomain.Project, platform.Repository) error
	GetByName(domain.Account, domain.ResourceName, bool) (ProjectDetailDTO, error)
	GetByRepoId(domain.Identity) (sdk.SpaceMetaDTO, error)
	List(domain.Account, *app.ResourceListCmd) (ProjectsDTO, error)
	ListGlobal(*app.GlobalResourceListCmd) (GlobalProjectsDTO, error)
	Update(*spacedomain.Project, *ProjectUpdateCmd, platform.Repository) (ProjectDTO, error)
	Fork(*ProjectForkCmd, platform.Repository) (ProjectDTO, error)

	AddRelatedModel(*spacedomain.Project, *domain.ResourceIndex) error
	RemoveRelatedModel(*spacedomain.Project, *domain.ResourceIndex) error

	AddRelatedDataset(*spacedomain.Project, *domain.ResourceIndex) error
	RemoveRelatedDataset(*spacedomain.Project, *domain.ResourceIndex) error

	SetTags(*spacedomain.Project, *app.ResourceTagsUpdateCmd) error

	NotifyUpdateCodes(domain.Identity, *CmdToNotifyUpdateCode) error
}

func NewProjectService(
	user userrepo.User,
	repo spacerepo.Project,
	model repository.Model,
	dataset repository.Dataset,
	activity repository.Activity,
	pr platform.Repository,
	sender message.ResourceProducer,
	computilityApp computilityapp.ComputilityInternalAppService,
	spaceProducer spacedomain.SpaceEventProducer,
) ProjectService {
	return projectService{
		repo:     repo,
		activity: activity,
		sender:   sender,
		rs: app.ResourceService{
			User:    user,
			Model:   model,
			Project: repo,
			Dataset: dataset,
		},
		computilityApp: computilityApp,
		spaceProducer:  spaceProducer,
	}
}

type projectService struct {
	repo spacerepo.Project
	//pr       platform.Repository
	activity       repository.Activity
	sender         message.ResourceProducer
	rs             app.ResourceService
	computilityApp computilityapp.ComputilityInternalAppService
	spaceProducer  spacedomain.SpaceEventProducer
}

func (s projectService) CanApplyResourceName(owner domain.Account, name domain.ResourceName) bool {
	return s.rs.CanApplyResourceName(owner, name)
}

func (s projectService) Create(cmd *ProjectCreateCmd, pr platform.Repository) (dto ProjectDTO, err error) {
	v := new(spacedomain.Project)
	cmd.toProject(v)
	count := v.GetQuotaCount()
	hdType := v.GetComputeType()
	id := domain.CreateIdentity(domain.GetId())
	compCmd := computilityapp.CmdToUserQuotaUpdate{
		Index: computilitydomain.ComputilityAccountRecordIndex{
			UserName:    cmd.Owner,
			ComputeType: hdType,
			SpaceId:     id,
		},
		QuotaCount: count,
	}

	err = s.computilityApp.UserQuotaConsume(compCmd)
	if err != nil {
		logrus.Errorf("space create error | call api for quota consume failed | user:%s ,err: %s", cmd.Owner, err)

		return ProjectDTO{}, err
	}

	// step1: create repo on gitlab
	pid, err := pr.New(&platform.RepoOption{
		Name:     cmd.Name,
		RepoType: cmd.RepoType,
	})
	if err != nil {
		return
	}

	repoId, err := domain.NewIdentity(pid)
	if err != nil {
		return
	}

	// step2: save
	v.RepoId = pid

	if v.Hardware.IsNpu() {
		v.CompPowerAllocated = true
	}

	p, err := s.repo.Save(v)
	if err != nil {
		return
	}

	if err = s.computilityApp.SpaceCreateSupply(computilityapp.CmdToSupplyRecord{
		Index: computilitydomain.ComputilityAccountRecordIndex{
			UserName:    cmd.Owner,
			ComputeType: hdType,
			SpaceId:     id,
		},
		QuotaCount: count,
		NewSpaceId: repoId,
	}); err != nil {
		logrus.Errorf("add space id supplyment failed | user: %s, err: %s", cmd.Owner, err)

		err = s.Delete(v, pr)
		if err != nil {
			logrus.Errorf("delete space after add space id supplyment failed | user: %s, err: %s",
				cmd.Owner.Account(), err)

			return ProjectDTO{}, xerrors.Errorf("add space id supplyment failed: %w", err)
		}

		err = s.computilityApp.UserQuotaRelease(compCmd)
		if err != nil {
			logrus.Errorf("release quota after add space id supplyment failed | user: %s, err: %s",
				cmd.Owner.Account(), err)

			return ProjectDTO{}, xerrors.Errorf("add space id supplyment failed: %w", err)
		}

		return ProjectDTO{}, xerrors.Errorf("add space id supplyment failed: %w", err)
	}

	s.toProjectDTO(&p, &dto)

	// add activity
	r, repoType := p.ResourceObject()
	ua := app.GenActivityForCreatingResource(r, repoType)
	_ = s.activity.Save(&ua)

	_ = s.sender.AddOperateLogForCreateResource(r, p.Name)

	_ = s.sender.CreateProject(message.ProjectCreatedEvent{
		Account:     r.Owner,
		ProjectName: dto.Name,
	})

	return
}

func (s projectService) Delete(r *spacedomain.Project, pr platform.Repository) (err error) {
	// step1: delete repo on gitlab
	if err = pr.Delete(r.RepoId); err != nil {
		return
	}

	obj, repoType := r.ResourceObject()

	// step2:
	if resources := r.RelatedResources(); len(resources) > 0 {
		err = s.sender.RemoveRelatedResources(&message.RelatedResources{
			Promoter:  obj,
			Resources: resources,
		})
		if err != nil {
			return
		}
	}

	if r.Hardware != nil {
		if r.Hardware.IsNpu() {
			logrus.Infof("release quota after user:%s npu space:%s delete", r.Owner.Account(), r.Id)

			rid, err := domain.NewIdentity(r.RepoId)
			if err != nil {
				return err
			}

			c := computilityapp.CmdToUserQuotaUpdate{
				Index: computilitydomain.ComputilityAccountRecordIndex{
					UserName:    r.Owner,
					ComputeType: r.GetComputeType(),
					SpaceId:     rid,
				},
				QuotaCount: r.GetQuotaCount(),
			}

			err = s.computilityApp.UserQuotaRelease(c)
			if err != nil {
				logrus.Errorf("failed to release user:%s quota after space:%s delete: %s",
					r.Owner.Account(), r.RepoId, err)

				return nil
			}
		}

	}

	// step3: delete
	if err = s.repo.Delete(&obj.ResourceIndex); err != nil {
		return
	}

	if err = s.spaceProducer.SendDeletedEvent(&spacedomain.DeleteSpaceEvent{
		Time:      utils.Now(),
		Owner:     r.Owner.Account(),
		SpaceName: r.Name.ResourceName(),
		DeletedBy: r.Owner.Account(),
		SpaceId:   r.RepoId,
	}); err != nil {
		return
	}

	// add activity
	ua := app.GenActivityForDeletingResource(&obj, repoType)

	// ignore the error
	_ = s.activity.Save(&ua)

	return
}

func (s projectService) GetByName(
	owner domain.Account, name domain.ResourceName,
	allowPrivacy bool,
) (dto ProjectDetailDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	if !allowPrivacy && v.IsPrivate() {
		err = PrivateRepoError{
			errors.New("private repo"),
		}

		return
	}

	m, err := s.rs.ListModels(v.RelatedModels)
	if err != nil {
		return
	}
	dto.RelatedModels = m

	d, err := s.rs.ListDatasets(v.RelatedDatasets)
	if err != nil {
		return
	}
	dto.RelatedDatasets = d

	s.toProjectDTO(&v, &dto.ProjectDTO)

	return
}

func (s projectService) GetByRepoId(id domain.Identity) (sdk.SpaceMetaDTO, error) {
	v, err := s.repo.GetByRepoId(id)

	if err != nil {
		return sdk.SpaceMetaDTO{}, err
	}

	res := s.toSpaceMetaDTO(v)
	return res, err
}

func (s projectService) ListGlobal(cmd *app.GlobalResourceListCmd) (
	dto GlobalProjectsDTO, err error,
) {
	option := cmd.ToResourceListOption()

	var v spacerepo.UserProjectsInfo

	if cmd.SortType == nil {
		v, err = s.repo.ListGlobalAndSortByUpdateTime(&option)
	} else {
		switch cmd.SortType.SortType() {
		case domain.SortTypeUpdateTime:
			v, err = s.repo.ListGlobalAndSortByUpdateTime(&option)

		case domain.SortTypeFirstLetter:
			v, err = s.repo.ListGlobalAndSortByFirstLetter(&option)

		case domain.SortTypeDownloadCount:
			v, err = s.repo.ListGlobalAndSortByDownloadCount(&option)
		}
	}

	items := v.Projects

	if err != nil || len(items) == 0 {
		return
	}

	// find avatars
	users := make([]userdomain.Account, len(items))
	for i := range items {
		users[i] = items[i].Owner
	}

	avatars, err := s.rs.FindUserAvater(users)
	if err != nil {
		return
	}

	// gen result
	dtos := make([]GlobalProjectDTO, len(items))
	for i := range items {
		s.toProjectSummaryDTO(&items[i], &dtos[i].ProjectSummaryDTO)
		dtos[i].AvatarId = avatars[i]
	}

	dto.Total = v.Total
	dto.Projects = dtos

	return
}

func (s projectService) List(owner domain.Account, cmd *app.ResourceListCmd) (
	dto ProjectsDTO, err error,
) {
	option := cmd.ToResourceListOption()

	var v spacerepo.UserProjectsInfo

	if cmd.SortType == nil {
		v, err = s.repo.ListAndSortByUpdateTime(owner, &option)
	} else {
		switch cmd.SortType.SortType() {
		case domain.SortTypeUpdateTime:
			v, err = s.repo.ListAndSortByUpdateTime(owner, &option)

		case domain.SortTypeFirstLetter:
			v, err = s.repo.ListAndSortByFirstLetter(owner, &option)

		case domain.SortTypeDownloadCount:
			v, err = s.repo.ListAndSortByDownloadCount(owner, &option)
		}
	}

	items := v.Projects

	if err != nil || len(items) == 0 {
		return
	}

	dtos := make([]ProjectSummaryDTO, len(items))
	for i := range items {
		s.toProjectSummaryDTO(&items[i], &dtos[i])
	}

	dto.Total = v.Total
	dto.Projects = dtos

	return
}

func (s projectService) NotifyUpdateCodes(id domain.Identity, cmd *CmdToNotifyUpdateCode) (
	err error,
) {
	space, err := s.repo.GetByRepoId(id)
	if err != nil {
		return err
	}

	space.SetSpaceCommitId(cmd.CommitId)
	space.SetNoApplicationFile(cmd.NoApplicationFile)

	// step2
	info := spacerepo.ProjectPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(&space),
		Property:         space.ProjectModifiableProperty,
	}
	err = s.repo.UpdateProperty(&info)

	if err != nil {
		err = xerrors.Errorf("save space failed, err: %w", err)

		return err
	}
	logrus.Infof("spaceId:%s set notify commitId:%s, req:%v success",
		id, cmd.CommitId, cmd)

	return nil
}
