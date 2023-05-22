package app

import "github.com/opensourceways/xihe-server/user/domain"

func (s userService) AddFollower(f *domain.FollowerInfo) error {
	return s.repo.AddFollower(f)
}

func (s userService) RemoveFollower(f *domain.FollowerInfo) error {
	return s.repo.RemoveFollower(f)
}

func (s userService) ListFollower(cmd *FollowsListCmd) (
	dto FollowsDTO, err error,
) {
	v, err := s.repo.FindFollower(cmd.User, &cmd.FollowFindOption)
	items := v.Users
	if err != nil || len(items) == 0 {
		return
	}

	dtos := make([]FollowDTO, len(items))
	for i := range items {
		s.toFollowDTO(&items[i], &dtos[i])
	}

	dto.Total = v.Total
	dto.Data = dtos

	return
}

func (s userService) toFollowDTO(f *domain.FollowerUserInfo, dto *FollowDTO) {
	*dto = FollowDTO{
		Account:    f.Account.Account(),
		IsFollower: f.IsFollower,
	}

	if f.Bio != nil {
		dto.Bio = f.Bio.Bio()
	}

	if f.AvatarId != nil {
		dto.AvatarId = f.AvatarId.AvatarId()
	}
}
