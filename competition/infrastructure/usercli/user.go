package usercli

import (
	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/user"
	types "github.com/opensourceways/xihe-server/domain"
	userApp "github.com/opensourceways/xihe-server/user/app"
	userDomain "github.com/opensourceways/xihe-server/user/domain"
)

func NewUserCli(c userApp.RegService) user.User {
	return &userImpl{c}
}

type userImpl struct {
	srv userApp.RegService
}

func (impl *userImpl) AddUserRegInfo(s *domain.Competitor) (err error) {
	cmd := new(userApp.UserRegisterInfoCmd)
	if err = toUserRegisterInfoCmd(s, cmd); err != nil {
		return
	}

	return impl.srv.UpsertUserRegInfo(cmd)
}

func (impl *userImpl) GetUserRegInfo(user types.Account) (s domain.Competitor, err error) {
	dto, err := impl.srv.GetUserRegInfo(user)
	if err != nil {
		return
	}

	if toCompetitort(&dto, &s) != nil {
		return
	}

	return
}

func toUserRegisterInfoCmd(s *domain.Competitor, cmd *userApp.UserRegisterInfoCmd) (err error) {
	cmd.Account = s.Account

	if cmd.Name, err = userDomain.NewName(s.Name.CompetitorName()); err != nil {
		return
	}

	if cmd.City, err = userDomain.NewCity(s.City.City()); err != nil {
		return
	}

	if cmd.Email, err = userDomain.NewEmail(s.Email.Email()); err != nil {
		return
	}

	if cmd.Phone, err = userDomain.NewPhone(s.Phone.Phone()); err != nil {
		return
	}

	if cmd.Identity, err = userDomain.NewIdentity(s.Identity.CompetitionIdentity()); err != nil {
		return
	}

	if cmd.Province, err = userDomain.NewProvince(s.Province.Province()); err != nil {
		return
	}

	cmd.Detail = s.Detail

	return
}

func toCompetitort(dto *userApp.UserRegisterInfoDTO, s *domain.Competitor) (err error) {
	if s.Account, err = types.NewAccount(dto.Account.Account()); err != nil {
		return
	}

	if s.Name, err = domain.NewCompetitorName(dto.Name.Name()); err != nil {
		return
	}

	if s.City, err = domain.NewCity(dto.City.City()); err != nil {
		return
	}

	if s.Email, err = types.NewEmail(dto.Email.Email()); err != nil {
		return
	}

	if s.Phone, err = domain.NewPhone(dto.Phone.Phone()); err != nil {
		return
	}

	if s.Identity, err = domain.NewcompetitionIdentity(dto.Identity.Identity()); err != nil {
		return
	}

	if s.Province, err = domain.NewProvince(dto.Province.Province()); err != nil {
		return
	}

	s.Detail = dto.Detail

	return
}
