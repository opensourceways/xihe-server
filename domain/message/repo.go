package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type RepoMessageProducer interface {
	SendResourceDownloaded(e domain.RepoDownload) error
	IncreaseDownload(*domain.ResourceObject) error
	AddOperateLogForDownloadFile(domain.Account, RepoFile) error
}
