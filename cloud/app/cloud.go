package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/message"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	"github.com/opensourceways/xihe-server/cloud/domain/service"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
)

type CloudService interface {
	// cloud
	ListCloud(*GetCloudConfCmd) ([]CloudDTO, error)
	SubscribeCloud(*SubscribeCloudCmd) (code string, err error)

	// pod
	Get(*PodInfoCmd) (PodInfoDTO, error)
	ReleaseCloud(*ReleaseCloudCmd) error
	GetReleasedPod(*GetReleasedPodCmd) (PodInfoDTO, error)
}

var _ CloudService = (*cloudService)(nil)

func NewCloudService(
	cloudRepo repository.Cloud,
	podRepo repository.Pod,
	producer message.CloudMessageProducer,
	whitelistRepo userrepo.WhiteList,
) *cloudService {
	return &cloudService{
		cloudRepo:        cloudRepo,
		podRepo:          podRepo,
		producer:         producer,
		cloudService:     service.NewCloudService(podRepo, producer),
		whitelistService: userapp.NewWhiteListService(whitelistRepo),
	}
}

type cloudService struct {
	cloudRepo        repository.Cloud
	podRepo          repository.Pod
	producer         message.CloudMessageProducer
	cloudService     service.CloudService
	whitelistService userapp.WhiteListService
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
			dto[i].toCloudDTO(&c[i], c[i].HasSingleCardIdle() || c[i].HasMultiCardsIdle(0), false)
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

		dto[i].toCloudDTO(&c[i], c[i].HasSingleCardIdle() || c[i].HasMultiCardsIdle(0), b)
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
		useNPU, useMultiNPU, err := s.whitelistService.CheckCloudWhitelist(cmd.User)
		if err != nil {
			return "", err
		}

		if (cmd.CardsNum.CloudSpecCardsNum() > 1 && !useMultiNPU) ||
			(cmd.CardsNum.CloudSpecCardsNum() == 1 && !useNPU) {
			return errorWhitelistNotAllowed, errors.New("not in cloud whitelist")
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

	deduction := cmd.CardsNum.CloudSpecCardsNum()

	singleCardBusy := deduction == 1 && !c.HasSingleCardIdle()
	if singleCardBusy {
		code = errorResourceBusy
		err = errors.New("no idle resource remain")

		return
	}

	multiCardsBusy := deduction > 1 && !c.HasMultiCardsIdle(deduction)
	if multiCardsBusy {
		code = errorResourceBusy
		err = errors.New("no idle multiple cards remain")

		return
	}

	// subscribe
	err = s.cloudService.SubscribeCloud(&c.CloudConf, cmd.User, cmd.ImageAlias, cmd.CardsNum)

	return
}

func (s *cloudService) ReleaseCloud(cmd *ReleaseCloudCmd) error {
	podInfo, err := s.podRepo.GetPodInfo(cmd.PodId)
	if err != nil {
		return err
	}

	if !podInfo.IsOnwer(cmd.User) {
		return ErrCloudNotAllowed
	}

	if !podInfo.CanRelease() {
		return ErrCloudReleased
	}

	podInfo.StatusSetTerminating()
	if err := s.podRepo.UpdatePod(&podInfo); err != nil {
		return err
	}

	return s.producer.ReleaseCloud(&message.ReleaseCloudEvent{
		PodId:     podInfo.Id,
		CloudType: podInfo.GetCloudType(),
	})
}

func (s *cloudService) GetReleasedPod(cmd *GetReleasedPodCmd) (PodInfoDTO, error) {
	podInfoDto := PodInfoDTO{}
	podInfo, err := s.podRepo.GetPodInfo(cmd.PodId)
	if err != nil {
		return podInfoDto, err
	}

	if !podInfo.IsTerminated() {
		return podInfoDto, ErrPodNotFound
	}

	cloudConf, err := s.cloudRepo.GetCloudConf(podInfo.CloudId)
	if err != nil {
		return podInfoDto, err
	}

	podInfoDto.toPodInfoDTO(&podInfo, &cloudConf)

	return podInfoDto, nil
}
