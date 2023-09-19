package app

import (
	"github.com/opensourceways/xihe-server/user/domain"
)

func (s userService) AddFollowing(f *domain.FollowerInfo) error {
	err := s.repo.AddFollowing(f)
	if err != nil {
		return err
	}

	// TODO: activity

	// send event
	_ = s.sender.SendFollowingAddedEvent(f)

	return nil
}

func (s userService) RemoveFollowing(f *domain.FollowerInfo) error {
	err := s.repo.RemoveFollowing(f)
	if err != nil {
		return err
	}

	// send event
	_ = s.sender.SendFollowingRemovedEvent(f)

	return nil
}

func (s userService) ListFollowing(cmd *FollowsListCmd) (
	dto FollowsDTO, err error,
) {
	v, err := s.repo.FindFollowing(cmd.User, &cmd.FollowFindOption)
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
