package domain

import (
	"errors"
)

const (
	activityTypeFork   = "fork"
	activityTypeLike   = "like"
	activityTypeCreate = "create"
	activityTypeDelete = "delete"
)

var (
	ActivityTypeFork   = activityType(activityTypeFork)
	ActivityTypeLike   = activityType(activityTypeLike)
	ActivityTypeCreate = activityType(activityTypeCreate)
	ActivityTypeDelete = activityType(activityTypeDelete)
)

// ActivityType
type ActivityType interface {
	ActivityType() string
}

func NewActivityType(v string) (ActivityType, error) {
	b := v != activityTypeLike &&
		v != activityTypeCreate &&
		v != activityTypeDelete &&
		v != activityTypeFork

	if b {
		return nil, errors.New("unknown activity type")
	}

	return activityType(v), nil
}

type activityType string

func (r activityType) ActivityType() string {
	return string(r)
}

// UserActivity
type UserActivity struct {
	Owner Account

	Activity
}

type Activity struct {
	Type ActivityType
	Time int64

	RepoType RepoType

	ResourceObject
}

func (r Activity) IsPublic() bool {
	if r.RepoType == nil {
		return false
	}

	return r.RepoType.RepoType() == RepoTypePublic
}
