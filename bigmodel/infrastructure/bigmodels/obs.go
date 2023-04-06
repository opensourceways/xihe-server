package bigmodels

import (
	"io"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func initOBS(cfg *OBSAuthInfo) (s obsService, err error) {
	cli, err := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return
	}

	s.cli = cli

	return
}

type obsService struct {
	cli *obs.ObsClient
}

func (s *obsService) createObject(f io.Reader, bucket, path string) error {
	input := &obs.PutObjectInput{}
	input.Bucket = bucket
	input.Key = path
	input.Body = f

	_, err := s.cli.PutObject(input)

	return err
}

func (s *obsService) genFileDownloadURL(bucket, p string, downloadExpiry int) (string, error) {
	input := &obs.CreateSignedUrlInput{}
	input.Method = obs.HttpMethodGet
	input.Bucket = bucket
	input.Key = p
	input.Expires = downloadExpiry

	output, err := s.cli.CreateSignedUrl(input)
	if err != nil {
		return "", err
	}

	return output.SignedUrl, nil
}

func (s *obsService) copyObject(bucket, dst, src string) error {
	input := &obs.CopyObjectInput{}
	input.Bucket = bucket
	input.Key = dst
	input.CopySourceBucket = bucket
	input.CopySourceKey = src
	_, err := s.cli.CopyObject(input)

	return err
}

func (s *obsService) deleteObject(bucket string, path string) error {
	input := &obs.DeleteObjectInput{}
	input.Bucket = bucket
	input.Key = path

	v, err := s.cli.DeleteObject(input)
	if err != nil && v != nil && v.StatusCode == 404 {
		err = nil
	}

	return err
}
