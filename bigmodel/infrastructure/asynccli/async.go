package asynccli

import (
	asyncapp "github.com/opensourceways/xihe-extra-services/async-server/app"
	asyncrepo "github.com/opensourceways/xihe-extra-services/async-server/domain/repository"
	"github.com/opensourceways/xihe-server/bigmodel/domain/async"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

func NewAsyncCli(c asyncapp.TaskService) async.AsyncTask {
	return &asyncImpl{c}
}

type asyncImpl struct {
	srv asyncapp.TaskService
}

func (impl *asyncImpl) GetWaitingTaskRank(user types.Account, time commondomain.Time, taskType []string) (rank int, err error) {
	return impl.srv.GetWaitingTaskRank(user, time, taskType)
}

func (impl *asyncImpl) GetLastFinishedTask(user types.Account, taskType []string) (resp asyncrepo.WuKongResp, err error) {
	return impl.srv.GetLastFinishedTask(user, taskType)
}
