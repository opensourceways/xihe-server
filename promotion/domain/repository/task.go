package repository

import "github.com/opensourceways/xihe-server/promotion/domain"

type Task interface {
	Find(id string) (domain.Task, error)
	FindAll() ([]domain.Task, error)
}
