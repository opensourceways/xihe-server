package domain

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/domain"
)

type Project struct {
	Id string

	Owner    domain.Account
	Type     domain.ProjType
	Protocol domain.ProtocolName
	Training domain.TrainingPlatform

	ProjectModifiableProperty

	RepoId string

	RelatedModels   domain.RelatedResources
	RelatedDatasets domain.RelatedResources

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fields is not under the controlling of version
	LikeCount     int
	ForkCount     int
	DownloadCount int

	Hardware  domain.Hardware
	BaseImage domain.BaseImage
}

func (p *Project) MaxRelatedResourceNum() int {
	return domain.DomainConfig.MaxRelatedResourceNum
}

func (p *Project) IsPrivate() bool {
	return p.RepoType.RepoType() == domain.RepoTypePrivate
}

func (p *Project) IsOnline() bool {
	return p.RepoType.RepoType() == domain.RepoTypeOnline
}

func (p *Project) ResourceIndex() domain.ResourceIndex {
	return domain.ResourceIndex{
		Owner: p.Owner,
		Id:    p.Id,
	}
}

func (p *Project) ResourceObject() (domain.ResourceObject, domain.RepoType) {
	return domain.ResourceObject{
		Type:          domain.ResourceTypeProject,
		ResourceIndex: p.ResourceIndex(),
	}, p.RepoType
}

func (p *Project) RelatedResources() []domain.ResourceObjects {
	r := make([]domain.ResourceObjects, 0, 2)

	if len(p.RelatedModels) > 0 {
		r = append(r, domain.ResourceObjects{
			Type:    domain.ResourceTypeModel,
			Objects: p.RelatedModels,
		})
	}

	if len(p.RelatedDatasets) > 0 {
		r = append(r, domain.ResourceObjects{
			Type:    domain.ResourceTypeDataset,
			Objects: p.RelatedDatasets,
		})
	}

	return r
}

// GetQuotaCount returns the quota count of the Space.
func (s *Project) GetQuotaCount() int {
	if s.Hardware.IsNpu() {
		return 1
	} else if s.Hardware.IsCpu() {
		return 0
	}

	return 0
}

// GetComputeType returns the compute type of the Space.
func (s *Project) GetComputeType() domain.ComputilityType {
	if s.Hardware.IsNpu() {
		return domain.CreateComputilityType("npu")
	} else if s.Hardware.IsCpu() {
		return domain.CreateComputilityType("cpu")
	}

	return nil
}

// SetSpaceCommitId for update space commitId.
func (s *Project) SetSpaceCommitId(commitId string) {
	s.CommitId = commitId
}

// SetNoApplicationFile for set NoApplicationFile and Exception.
func (s *Project) SetNoApplicationFile(noApplicationFile bool) {
	s.NoApplicationFile = noApplicationFile
	if noApplicationFile {
		s.Exception = domain.CreateException(domain.NoApplicationFile)
		return
	}
	if !noApplicationFile && s.Exception == domain.ExceptionNoApplicationFile {
		s.Exception = domain.CreateException("")
	}
}

func (m *Project) PreCheck() error {
	if m.Exception.Exception() != "" {
		e := xerrors.Errorf("spaceId:%s failed to create space failed, space has exception reason :%s",
			m.RepoId, m.Exception.Exception())

		err := allerror.New(allerror.ErrorCodeSpaceAppCreateFailed, e.Error(), e)
		logrus.Errorf("spaceId：%s create space failed, err:%s", m.RepoId, err)

		return err
	}
	return nil
}

type ProjectModifiableProperty struct {
	Name               domain.ResourceName
	Desc               domain.ResourceDesc
	Title              domain.ResourceTitle
	CoverId            domain.CoverId
	RepoType           domain.RepoType
	Tags               []string
	TagKinds           []string
	Level              domain.ResourceLevel
	CommitId           string
	NoApplicationFile  bool
	Exception          domain.Exception
	CompPowerAllocated bool
}

type ProjectSummary struct {
	Id            string
	Owner         domain.Account
	Name          domain.ResourceName
	Desc          domain.ResourceDesc
	Title         domain.ResourceTitle
	Level         domain.ResourceLevel
	CoverId       domain.CoverId
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	ForkCount     int
	DownloadCount int
}
