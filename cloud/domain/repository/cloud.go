package repository

import "github.com/opensourceways/xihe-server/cloud/domain"

type Cloud interface {
	ListCloudConf() ([]domain.CloudConf, error)
}
