package repositoryadapter

import "github.com/opensourceways/xihe-server/points/domain"

func TaskAdapter() *taskAdapter {
	return &taskAdapter{}
}

type taskAdapter struct {
}

func (impl *taskAdapter) FindAllTasks() ([]domain.Task, error) {
	return nil, nil
}

func (impl *taskAdapter) Find(string) (domain.Task, error) {
	return domain.Task{}, nil
}
