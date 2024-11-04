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
		toFollowDTO(&items[i], &dtos[i])
	}

	dto.Total = len(dtos)
	dto.Data = dtos

	return
}
