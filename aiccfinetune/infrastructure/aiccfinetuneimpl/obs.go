package aiccfinetuneimpl

import (
	"io"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func initOBS(cfg *OBSConfig) (s obsService, err error) {
	cli, err := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return
	}

	s.cli = cli
	s.bucket = cfg.Bucket
	s.prefix = cfg.Prefix

	return
}

type obsService struct {
	cli    *obs.ObsClient
	bucket string
	prefix string
}

func (s *obsService) genPath(path string) string {
	if s.prefix == "" {
		return path
	}

	return s.prefix + "/" + path
}

func (s *obsService) createObject(f io.Reader, path string) error {
	input := &obs.PutObjectInput{}
	input.Bucket = s.bucket
	input.Key = s.genPath(path)
	input.Body = f

	_, err := s.cli.PutObject(input)

	return err
}
