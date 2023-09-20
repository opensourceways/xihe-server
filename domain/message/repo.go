package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type DownloadEvent struct {
	Account domain.Account
	Type    domain.ResourceType
	Name    string
}

type RepoMessageProducer interface {
	DownloadRepo(e DownloadEvent) error
	IncreaseDownload(*domain.ResourceObject) error
	AddOperateLogForDownloadFile(domain.Account, RepoFile) error
}
