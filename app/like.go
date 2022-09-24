package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type LikeCreateCmd struct {
	ResourceOwner domain.Account
	ResourceType  domain.ResourceType
	ResourceId    string
}

func (cmd *LikeCreateCmd) Validate() error {
	b := cmd.ResourceOwner != nil &&
		cmd.ResourceType != nil &&
		cmd.ResourceId != ""

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

type LikeRemoveCmd = LikeCreateCmd

type LikeDTO struct {
	Time     string      `json:"time"`
	Resource ResourceDTO `json:"resource"`
}

type LikeService interface {
	Create(domain.Account, LikeCreateCmd) error
	Delete(domain.Account, LikeRemoveCmd) error
	List(domain.Account) ([]LikeDTO, error)
}

func NewLikeService(
	repo repository.Like,
	user repository.User,
	model repository.Model,
	project repository.Project,
	dataset repository.Dataset,
	activity repository.Activity,
	sender message.Sender,
) LikeService {
	return likeService{
		repo:     repo,
		activity: activity,
		sender:   sender,

		rs: resourceService{
			user:    user,
			model:   model,
			project: project,
			dataset: dataset,
		},
	}
}

type likeService struct {
	repo     repository.Like
	activity repository.Activity
	sender   message.Sender

	rs resourceService
}

func (s likeService) Create(owner domain.Account, cmd LikeCreateCmd) error {
	now := utils.Now()

	obj := domain.ResourceObject{Type: cmd.ResourceType}
	obj.Owner = cmd.ResourceOwner
	obj.Id = cmd.ResourceId

	v := domain.UserLike{
		Owner: owner,
		Like: domain.Like{
			CreatedAt:      now,
			ResourceObject: obj,
		},
	}

	if err := s.repo.Save(&v); err != nil {
		return err
	}

	ua := domain.UserActivity{
		Owner: owner,
		Activity: domain.Activity{
			Type:           domain.ActivityTypeLike,
			Time:           now,
			ResourceObject: v.ResourceObject,
		},
	}
	if err := s.activity.Save(&ua); err != nil {
		return err
	}

	// send event
	_ = s.sender.AddLike(v.Like)

	return nil
}

func (s likeService) Delete(owner domain.Account, cmd LikeRemoveCmd) error {
	obj := domain.ResourceObject{Type: cmd.ResourceType}
	obj.Owner = cmd.ResourceOwner
	obj.Id = cmd.ResourceId

	v := domain.UserLike{
		Owner: owner,
		Like:  domain.Like{ResourceObject: obj},
	}

	if err := s.repo.Remove(&v); err != nil {
		return err
	}

	// send event
	_ = s.sender.RemoveLike(v.Like)

	return nil
}

func (s likeService) List(owner domain.Account) (
	dtos []LikeDTO, err error,
) {
	likes, err := s.repo.Find(owner, repository.LikeFindOption{})
	if err != nil || len(likes) == 0 {
		return
	}

	total := len(likes)
	objs := make([]*domain.ResourceObject, total)
	orders := make([]orderByTime, total)
	for i := range likes {
		item := &likes[i]

		objs[i] = &item.ResourceObject
		orders[i] = orderByTime{t: item.CreatedAt, p: i}
	}

	resources, err := s.rs.list(objs)
	if err != nil {
		return
	}

	rm := make(map[string]*ResourceDTO)
	for i := range resources {
		item := &resources[i]

		rm[item.identity()] = item
	}

	dtos = make([]LikeDTO, len(likes))
	err = sortAndSet(orders, func(i, j int) error {
		item := &likes[i]

		r, ok := rm[item.String()]
		if !ok {
			return errors.New("no matched resource")
		}

		dtos[j] = LikeDTO{
			Time:     utils.ToDate(item.CreatedAt),
			Resource: *r,
		}

		return nil
	})

	return
}
