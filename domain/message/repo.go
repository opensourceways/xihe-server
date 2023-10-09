package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type RepoMessageProducer interface {
	SendRepoDownloaded(*domain.RepoDownloadedEvent) error
	IncreaseDownload(*domain.ResourceObject) error
	AddOperateLogForDownloadFile(domain.Account, RepoFile) error
}
