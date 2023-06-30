package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ActivityDTO struct {
	Type     string      `json:"type"`
	Time     string      `json:"time"`
	Resource ResourceDTO `json:"resource"`
}

type ActivityService interface {
	List(domain.Account) ([]ActivityDTO, error)
}

func NewActivityService(
	repo repository.Activity,
	user userrepo.User,
	model repository.Model,
	project repository.Project,
	dataset repository.Dataset,
) ActivityService {
	return activityService{
		repo: repo,
		rs: resourceService{
			user:    user,
			model:   model,
			project: project,
			dataset: dataset,
		},
	}
}

type activityService struct {
	repo repository.Activity
	rs   resourceService
}

func (s activityService) List(owner domain.Account) (
	dtos []ActivityDTO, err error,
) {
	activities, err := s.repo.Find(owner, repository.ActivityFindOption{})
	if err != nil || len(activities) == 0 {
		return
	}

	total := len(activities)
	objs := make([]*domain.ResourceObject, total)
	orders := make([]orderByTime, total)
	for i := range activities {
		item := &activities[i]

		objs[i] = &item.ResourceObject
		orders[i] = orderByTime{t: item.Time, p: i}
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

	dtos = make([]ActivityDTO, total)
	j := 0
	_ = sortAndSet(orders, func(i int) error {
		item := &activities[i]

		if r, ok := rm[item.String()]; ok {
			dtos[j] = ActivityDTO{
				Type:     item.Type.ActivityType(),
				Time:     utils.ToDate(item.Time),
				Resource: *r,
			}

			j++
		}

		return nil
	})

	if j < total {
		dtos = dtos[:j]
	}

	return
}

func genActivityForCreatingResource(obj domain.ResourceObject) domain.UserActivity {
	return domain.UserActivity{
		Owner: obj.Owner,
		Activity: domain.Activity{
			Type:           domain.ActivityTypeCreate,
			Time:           utils.Now(),
			ResourceObject: obj,
		},
	}
}

func genActivityForDeletingResource(obj *domain.ResourceObject) domain.UserActivity {
	return domain.UserActivity{
		Owner: obj.Owner,
		Activity: domain.Activity{
			Type:           domain.ActivityTypeDelete,
			Time:           utils.Now(),
			ResourceObject: *obj,
		},
	}
}
