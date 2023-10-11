package uploader

import "io"

type DataFileUploader interface {
	UploadAICC(data io.Reader, path string) error
}
