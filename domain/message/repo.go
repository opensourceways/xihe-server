package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type RepoMessageProducer interface {
	SendResourceDownloaded(*domain.RepoDownloadEvent) error
	IncreaseDownload(*domain.ResourceObject) error
	AddOperateLogForDownloadFile(domain.Account, RepoFile) error
}
