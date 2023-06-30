package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ActivityMapper interface {
	Insert(string, ActivityDO) error
	List(string, ActivityListDO) ([]ActivityDO, error)
}

func NewActivityRepository(mapper ActivityMapper) repository.Activity {
	return activity{mapper}
}

type activity struct {
	mapper ActivityMapper
}

func (impl activity) Save(ul *domain.UserActivity) error {
	err := impl.mapper.Insert(ul.Owner.Account(), impl.toActivityDO(&ul.Activity))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl activity) Find(owner domain.Account, opt repository.ActivityFindOption) (
	[]domain.Activity, error,
) {
	v, err := impl.mapper.List(owner.Account(), ActivityListDO{})
	if err != nil {
		if isErrorDataNotExists(err) {
			return nil, nil
		}

		return nil, convertError(err)
	}

	if len(v) == 0 {
		return nil, nil
	}

	r := make([]domain.Activity, len(v))
	for i := range v {
		if err := v[i].toActivity(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl activity) toActivityDO(v *domain.Activity) ActivityDO {
	return ActivityDO{
		Type:             v.Type.ActivityType(),
		Time:             v.Time,
		RepoType:         v.RepoType.RepoType(),
		ResourceObjectDO: toResourceObjectDO(&v.ResourceObject),
	}
}

type ActivityListDO struct {
}

type ActivityDO struct {
	Type string
	Time int64

	RepoType string

	ResourceObjectDO
}

func (do *ActivityDO) toActivity(r *domain.Activity) (err error) {
	if r.Type, err = domain.NewActivityType(do.Type); err != nil {
		return
	}

	if do.RepoType != "" {
		if r.RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
			return
		}
	}

	r.Time = do.Time

	return do.ResourceObjectDO.toResourceObject(&r.ResourceObject)
}
