package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
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
	user userrepo.User,
	model repository.Model,
	project spacerepo.Project,
	dataset repository.Dataset,
	activity repository.Activity,
	sender message.LikeMessageProducer,
) LikeService {
	return likeService{
		repo:     repo,
		activity: activity,
		sender:   sender,

		rs: ResourceService{
			User:    user,
			Model:   model,
			Project: project,
			Dataset: dataset,
		},
	}
}

type likeService struct {
	repo     repository.Like
	activity repository.Activity
	sender   message.LikeMessageProducer

	rs ResourceService
}

// owner is the user who like the resource
func (s likeService) Create(owner domain.Account, cmd LikeCreateCmd) error {
	// check if resource is private
	var repotype domain.RepoType
	if isprivate, ok := s.rs.IsPrivate(
		cmd.ResourceOwner, cmd.ResourceType, cmd.ResourceId,
	); !ok || isprivate {
		return errors.New("cannot like private or not exist resource")
	} else {
		repotype, _ = domain.NewRepoType(domain.RepoTypePublic)
	}

	// check if resource has liked
	hasLiked, err := s.repo.HasLike(owner, &domain.ResourceObject{
		Type: cmd.ResourceType,
		ResourceIndex: domain.ResourceIndex{
			Owner: cmd.ResourceOwner,
			Id:    cmd.ResourceId,
		},
	})
	if err != nil {
		return err
	}
	if hasLiked {
		return errors.New("cannot like resource you had liked")
	}

	// add like to like repo
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
			RepoType:       repotype,
			ResourceObject: v.ResourceObject,
		},
	}
	if err := s.activity.Save(&ua); err != nil {
		return err
	}

	// increase like in resource
	_ = s.sender.AddLike(&domain.ResourceLikedEvent{
		Account: owner,
		Obj:     v.ResourceObject,
	})

	return nil
}

func (s likeService) Delete(owner domain.Account, cmd LikeRemoveCmd) error {
	obj := domain.ResourceObject{Type: cmd.ResourceType}
	obj.Owner = cmd.ResourceOwner
	obj.Id = cmd.ResourceId

	// check if resource is private
	if isprivate, ok := s.rs.IsPrivate(
		cmd.ResourceOwner, cmd.ResourceType, cmd.ResourceId,
	); !ok || isprivate {
		return errors.New("cannot like private or not exist resource")
	}

	// check if resource has liked
	hasLiked, err := s.repo.HasLike(owner, &domain.ResourceObject{
		Type: cmd.ResourceType,
		ResourceIndex: domain.ResourceIndex{
			Owner: cmd.ResourceOwner,
			Id:    cmd.ResourceId,
		},
	})
	if err != nil {
		return err
	}
	if !hasLiked {
		return errors.New("cannot remove like resource you are not like")
	}

	// remove like
	v := domain.UserLike{
		Owner: owner,
		Like:  domain.Like{ResourceObject: obj},
	}

	if err := s.repo.Remove(&v); err != nil {
		return err
	}

	// reduce like count in resource
	_ = s.sender.RemoveLike(&domain.ResourceLikedEvent{
		Account: owner,
		Obj:     obj,
	})

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
	j := 0
	_ = sortAndSet(orders, func(i int) error {
		item := &likes[i]

		if r, ok := rm[item.String()]; ok {
			dtos[j] = LikeDTO{
				Time:     utils.ToDate(item.CreatedAt),
				Resource: *r,
			}

			j++
		}

		return nil
	})

	if j < len(dtos) {
		dtos = dtos[:j]
	}

	return
}
