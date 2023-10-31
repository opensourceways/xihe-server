package app

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	"github.com/opensourceways/xihe-server/cloud/domain/service"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
)

type CloudService interface {
	// cloud
	ListCloud(*GetCloudConfCmd) ([]CloudDTO, error)
}

var _ CloudService = (*cloudService)(nil)

func NewCloudService(
	cloudRepo repository.Cloud,
	podRepo repository.Pod,
) *cloudService {
	return &cloudService{
		cloudRepo:    cloudRepo,
		podRepo:      podRepo,
		cloudService: service.NewCloudService(podRepo),
	}
}

type cloudService struct {
	cloudRepo    repository.Cloud
	podRepo      repository.Pod
	cloudService service.CloudService
}

func (s *cloudService) ListCloud(cmd *GetCloudConfCmd) (dto []CloudDTO, err error) {
	// list cloud conf
	confs, err := s.cloudRepo.ListCloudConf()
	if err != nil {
		return
	}

	// to cloud
	c := make([]domain.Cloud, len(confs))
	for i := range confs {
		c[i].CloudConf = confs[i]
		if err = s.cloudService.ToCloud(&c[i]); err != nil {
			return
		}
	}

	// to dto without holding
	if cmd.IsVisitor {
		dto = make([]CloudDTO, len(c))
		for i := range c {
			dto[i].toCloudDTO(&c[i], c[i].HasIdle(), false)
		}

		return
	}

	// to dto
	dto = make([]CloudDTO, len(c))
	for i := range c {
		var b bool
		if b, err = s.cloudService.HasHolding(types.Account(cmd.User), &c[i].CloudConf); err != nil {
			if !commonrepo.IsErrorResourceNotExists(err) {
				return
			}

			err = nil
		}

		dto[i].toCloudDTO(&c[i], c[i].HasIdle(), b)
	}

	return
}
