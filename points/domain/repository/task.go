package repository

import "github.com/opensourceways/xihe-server/points/domain"

type Task interface {
	FindAllTasks() ([]domain.Task, error)
	Find(string) (domain.Task, error)
}
