package service

import (
	"sync"

	common "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
	"github.com/opensourceways/xihe-server/points/domain/taskdoc"
)

type TaskService interface {
	Doc(lang common.Language) ([]byte, error)
}

func InitTaskService(repo repository.Task, doc taskdoc.TaskDoc) (*taskService, error) {
	tm := &taskService{
		repo:    repo,
		taskDoc: doc,
		docs:    map[string][]byte{},
		tasks:   map[string]int{},
	}

	items := common.SupportedLanguages()

	for i := range items {
		if _, err := tm.Doc(items[i]); err != nil {
			return nil, err
		}
	}

	return tm, nil
}

// taskService
type taskService struct {
	repo    repository.Task
	taskDoc taskdoc.TaskDoc

	mutex sync.RWMutex
	tasks map[string]int    // map task id to version
	docs  map[string][]byte // map lanuage to doc path
}

func (tm *taskService) Doc(lang common.Language) ([]byte, error) {
	tasks, err := tm.repo.FindAllTasks()
	if err != nil {
		return nil, err
	}

	if v := tm.doc(tasks, lang); len(v) > 0 {
		return v, nil
	}

	v, err := tm.taskDoc.Doc(tasks, lang)
	if err != nil {
		return nil, err
	}

	tm.updateDoc(tasks, lang, v)

	return v, nil
}

func (tm *taskService) doc(tasks []domain.Task, lang common.Language) []byte {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	if tm.isSame(tasks) {
		return tm.docs[lang.Language()]
	}

	return nil
}

func (tm *taskService) updateDoc(tasks []domain.Task, lang common.Language, v []byte) {
	tm.mutex.RLock()
	b := tm.isSame(tasks)
	tm.mutex.RUnlock()

	if b {
		return
	}

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if tm.isSame(tasks) {
		return
	}

	for i := range tasks {
		item := &tasks[i]

		tm.tasks[item.Id] = item.Version
	}

	tm.docs[lang.Language()] = v
}

func (tm *taskService) isSame(tasks []domain.Task) bool {
	if len(tasks) != len(tm.tasks) {
		return false
	}

	for i := range tasks {
		if item := &tasks[i]; item.Version != tm.tasks[item.Id] {
			return false
		}
	}

	return true
}
