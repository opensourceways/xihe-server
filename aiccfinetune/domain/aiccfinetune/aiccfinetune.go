package aiccfinetune

import (
	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
)

type AICCFinetuneServer interface {
	CreateJob(endpoint string, info *domain.AICCFinetuneIndex, t *domain.AICCFinetune) (domain.JobInfo, error)
	DeleteJob(endpoint, jobId string) error
	TerminateJob(jobId string) error
	GetLogPreviewURL(endpoint, jobId string) (string, error)
	IsJobDone(status string) bool
	GetFileDownloadURL(endpoint, file string) (string, error)
}
