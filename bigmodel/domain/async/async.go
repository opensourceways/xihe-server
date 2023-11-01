package async

import (
	asyncrepo "github.com/opensourceways/xihe-extra-services/async-server/domain/repository"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type AsyncTask interface {
	GetWaitingTaskRank(types.Account, commondomain.Time, []string) (int, error)
	GetLastFinishedTask(types.Account, []string) (asyncrepo.WuKongResp, error)
}
