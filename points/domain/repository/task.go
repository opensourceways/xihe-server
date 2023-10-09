package repository

import "github.com/opensourceways/xihe-server/points/domain"

type Task interface {
	Add(*domain.Task) error
	Find(string) (domain.Task, error)
	FindAllTasks() ([]domain.Task, error)
}
