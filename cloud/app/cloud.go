package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/message"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	"github.com/opensourceways/xihe-server/cloud/domain/service"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
)

type CloudService interface {
	// cloud
	ListCloud(*GetCloudConfCmd) ([]CloudDTO, error)
	SubscribeCloud(*SubscribeCloudCmd) (code string, err error)

	// pod
	Get(*PodInfoCmd) (PodInfoDTO, error)
}

var _ CloudService = (*cloudService)(nil)

func NewCloudService(
	cloudRepo repository.Cloud,
	podRepo repository.Pod,
	producer message.CloudMessageProducer,
	whitelistRepo userrepo.WhiteList,
) *cloudService {
	return &cloudService{
		cloudRepo:     cloudRepo,
		podRepo:       podRepo,
		producer:      producer,
		cloudService:  service.NewCloudService(podRepo, producer),
		whitelistRepo: whitelistRepo,
	}
}

type cloudService struct {
	cloudRepo     repository.Cloud
	podRepo       repository.Pod
	producer      message.CloudMessageProducer
	cloudService  service.CloudService
	whitelistRepo userrepo.WhiteList
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

func (s *cloudService) SubscribeCloud(cmd *SubscribeCloudCmd) (code string, err error) {
	// get cloud conf
	cloudConf, err := s.cloudRepo.GetCloudConf(cmd.CloudId)
	if err != nil {
		return
	}

	// whitelist
	if cloudConf.IsNPU() {
		const whitelistTypeCloud = "cloud"
		whitelist, err := s.whitelistRepo.GetWhiteListInfo(cmd.User, whitelistTypeCloud)
		if err != nil {
			return "", err
		}

		if !whitelist.Enable() {
			return errorWhitelistNotAllowed, errors.New("not allowed for this module")
		}
	}

	// check
	_, ok, err := s.cloudService.CheckUserCanSubsribe(cmd.User, cmd.CloudId)
	if err != nil {
		return
	}

	if !ok {
		code = errorNotAllowed
		err = errors.New("starting or running pod exist")

		return
	}

	c := new(domain.Cloud)
	c.CloudConf = cloudConf

	// check
	if err = s.cloudService.ToCloud(c); err != nil {
		return
	}

	if !c.HasIdle() {
		code = errorResourceBusy
		err = errors.New("no idle resource remain")

		return
	}

	// subscribe
	err = s.cloudService.SubscribeCloud(&c.CloudConf, cmd.User)

	return
}
