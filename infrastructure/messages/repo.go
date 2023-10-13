package messages

import (
	"fmt"

	"github.com/sirupsen/logrus"

	commsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/utils"
)

func NewDownloadMessageAdapter(
	topic string,
	cfg *DownloadProducerConfig,
	p commsg.Publisher,
	o commsg.OperateLogPublisher,
) *downloadMessageAdapter {
	return &downloadMessageAdapter{topic: topic, cfg: *cfg, publisher: p, operateLog: o}
}

type downloadMessageAdapter struct {
	topic      string
	cfg        DownloadProducerConfig
	publisher  commsg.Publisher
	operateLog commsg.OperateLogPublisher
}

type DownloadProducerConfig struct {
	ModelDownload   commsg.TopicConfig `json:"model_download" required:"true"`
	DatasetDownload commsg.TopicConfig `json:"dataset_download" required:"true"`
	ProjectDownload commsg.TopicConfig `json:"project_download" required:"true"`
}

func (s *downloadMessageAdapter) AddOperateLogForDownloadFile(u domain.Account, repo message.RepoFile) error {
	return s.operateLog.SendOperateLog(u, "download", map[string]string{
		"user": repo.User.Account(),
		"repo": repo.Name.ResourceName(),
		"path": repo.Path.FilePath(),
	})
}

// Download
func (s *downloadMessageAdapter) IncreaseDownload(obj *domain.ResourceObject) error {
	v := new(resourceObject)
	toMsgResourceObject(obj, v)

	return s.publisher.Publish(s.topic, v, nil)
}

func (s *downloadMessageAdapter) SendRepoDownloaded(e *domain.RepoDownloadedEvent) (err error) {
	switch e.Obj.Type {
	case domain.ResourceTypeDataset:
		err = s.downloadDataset(e)
	case domain.ResourceTypeModel:
		err = s.downloadModel(e)
	case domain.ResourceTypeProject:
		err = s.downloadProject(e)
	}

	return
}

// Download project/model/dataset
func (s *downloadMessageAdapter) downloadModel(e *domain.RepoDownloadedEvent) error {
	v := &commsg.MsgNormal{
		User:      e.Account.Account(),
		Type:      s.cfg.ModelDownload.Name,
		CreatedAt: utils.Now(),
		Desc:      fmt.Sprintf("Downloaded a model %s", e.Name),
	}

	logrus.Debugf("Send a downloaded model message: %+v", v)

	return s.publisher.Publish(s.cfg.ModelDownload.Topic, v, nil)
}

func (s *downloadMessageAdapter) downloadDataset(e *domain.RepoDownloadedEvent) error {
	v := &commsg.MsgNormal{
		User:      e.Account.Account(),
		Type:      s.cfg.DatasetDownload.Name,
		CreatedAt: utils.Now(),
		Desc:      fmt.Sprintf("Downloaded a dataset %s", e.Name),
	}

	logrus.Debugf("Send a downloaded dataset message: %+v", v)

	return s.publisher.Publish(s.cfg.DatasetDownload.Topic, v, nil)
}

func (s *downloadMessageAdapter) downloadProject(e *domain.RepoDownloadedEvent) error {
	v := &commsg.MsgNormal{
		User:      e.Account.Account(),
		Type:      s.cfg.ProjectDownload.Name,
		CreatedAt: utils.Now(),
		Desc:      fmt.Sprintf("Downloaded a project %s", e.Name),
	}

	logrus.Debugf("Send a downloaded project message: %+v", v)

	return s.publisher.Publish(s.cfg.ProjectDownload.Topic, v, nil)
}
