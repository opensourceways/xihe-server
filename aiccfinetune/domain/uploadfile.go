package domain

import (
	"fmt"
	"io"

	"github.com/opensourceways/xihe-server/aiccfinetune/domain/uploader"
)

type UploadService struct {
	uploader uploader.DataFileUploader
}

func NewUploadService(v uploader.DataFileUploader) UploadService {
	return UploadService{v}
}

func (s *UploadService) Upload(
	data io.Reader, fileName string, user string, model string, task string,
) (err error) {
	obspath := fmt.Sprintf(
		"%s/input/%s/%s_%s",
		model,
		task,
		user,
		fileName,
	)

	if err = s.uploader.UploadAICC(data, obspath); err != nil {
		return
	}
	return
}
