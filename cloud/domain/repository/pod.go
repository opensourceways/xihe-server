package repository

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type PodInfoList struct {
	PodInfos []domain.PodInfo
}

type Pod interface {
	GetRunningPod(cid string) (PodInfoList, error)
	GetPodInfo(pid string) (domain.PodInfo, error)
	GetUserPod(types.Account) (PodInfoList, error)
	GetUserCloudIdPod(user types.Account, cloudId string) (PodInfoList, error)
	AddStartingPod(*domain.PodInfo) (pid string, err error)
	UpdatePod(*domain.PodInfo) error
}
