package app

import (
	"github.com/opensourceways/xihe-server/domain/message"
	platform "github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

type UserService interface {
	// user
	Create(*UserCreateCmd) (UserDTO, error)
	CreatePlatformAccount(*CreatePlatformAccountCmd) (PlatformInfoDTO, error)
	UpdatePlateformInfo(*UpdatePlateformInfoCmd) error
	UpdateBasicInfo(domain.Account, UpdateUserBasicInfoCmd) error

	GetByAccount(domain.Account) (UserDTO, error)
	GetByFollower(owner, follower domain.Account) (UserDTO, bool, error)

	AddFollowing(*domain.FollowerInfo) error
	RemoveFollowing(*domain.FollowerInfo) error
	ListFollowing(*FollowsListCmd) (FollowsDTO, error)

	AddFollower(*domain.FollowerInfo) error
	RemoveFollower(*domain.FollowerInfo) error
	ListFollower(*FollowsListCmd) (FollowsDTO, error)
}

// ps: platform user service
func NewUserService(
	repo repository.User,
	ps platform.User,
	sender message.Sender,
) UserService {
	return userService{
		ps:     ps,
		repo:   repo,
		sender: sender,
	}
}

type userService struct {
	ps     platform.User
	repo   repository.User
	sender message.Sender
}

func (s userService) Create(cmd *UserCreateCmd) (dto UserDTO, err error) {
	// TODO keep transaction

	v := cmd.toUser()

	// update user
	u, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	s.toUserDTO(&u, &dto)

	_ = s.sender.AddOperateLogForNewUser(u.Account)

	return
}

func (s userService) GetByAccount(account domain.Account) (dto UserDTO, err error) {
	v, err := s.repo.GetByAccount(account)
	if err != nil {
		return
	}

	s.toUserDTO(&v, &dto)

	return
}

func (s userService) GetByFollower(owner, follower domain.Account) (
	dto UserDTO, isFollower bool, err error,
) {
	v, isFollower, err := s.repo.GetByFollower(owner, follower)
	if err != nil {
		return
	}

	s.toUserDTO(&v, &dto)

	return
}

func (s userService) CreatePlatformAccount(cmd *CreatePlatformAccountCmd) (dto PlatformInfoDTO, err error) {
	// create platform account
	pu, err := s.ps.New(platform.UserOption{
		Email:    cmd.email,
		Name:     cmd.account,
		Password: cmd.password,
	})
	if err != nil {
		return
	}

	dto.PlatformUser = pu

	// apply token
	token, err := s.ps.NewToken(pu)
	if err != nil {
		return
	}

	dto.PlatformToken = token

	return
}

func (s userService) UpdatePlateformInfo(cmd *UpdatePlateformInfoCmd) (err error) {
	// get userinfo
	u, err := s.repo.GetByAccount(cmd.User)
	if err != nil {
		return
	}

	// update some data
	u.PlatformUser = cmd.PlatformUser
	u.PlatformToken = cmd.PlatformToken
	u.Email = cmd.Email

	// update userinfo
	if _, err = s.repo.Save(&u); err != nil {
		return
	}

	return
}

func (s userService) toUserDTO(u *domain.User, dto *UserDTO) {
	*dto = UserDTO{
		Id:      u.Id,
		Email:   u.Email.Email(),
		Account: u.Account.Account(),
	}

	if u.Bio != nil {
		dto.Bio = u.Bio.Bio()
	}

	if u.AvatarId != nil {
		dto.AvatarId = u.AvatarId.AvatarId()
	}

	dto.FollowerCount = u.FollowerCount
	dto.FollowingCount = u.FollowingCount

	dto.Platform.Token = u.PlatformToken
	dto.Platform.UserId = u.PlatformUser.Id
	dto.Platform.NamespaceId = u.PlatformUser.NamespaceId
}
