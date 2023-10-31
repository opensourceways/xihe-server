package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/sirupsen/logrus"
)

func NewPodRepo(cfg *Config) repository.Pod {
	return &podRepoImpl{
		cli: pgsql.NewDBTable(cfg.Table.Pod),
	}
}

type podRepoImpl struct {
	cli pgsqlClient
}

func (impl *podRepoImpl) GetRunningPod(cid string) (
	pods repository.PodInfoList, err error,
) {
	filter := map[string]interface{}{
		fieldCloudId: cid,
		fieldStatus:  "running",
	}

	return impl.getFilterPods(filter)
}

func (impl *podRepoImpl) getFilterPods(filter interface{}) (
	pods repository.PodInfoList, err error,
) {

	var tpods []TPod
	if err = impl.cli.Filter(filter, &tpods); err != nil {
		return
	}

	podinfos := make([]domain.PodInfo, len(tpods))
	for i := range tpods {
		if err = tpods[i].toPodInfo(&podinfos[i]); err != nil {
			return
		}
	}

	pods.PodInfos = podinfos

	return
}

func (impl *podRepoImpl) getOrderOnePod(filter, order interface{}) (
	pod domain.PodInfo, err error,
) {
	var tpod TPod
	if err = impl.cli.GetOrderOneRecord(filter, order, &tpod); err != nil {
		if impl.cli.IsRowNotFound(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}
	} else {
		if err = tpod.toPodInfo(&pod); err != nil {
			return
		}
	}

	return
}

func (impl *podRepoImpl) GetPodInfo(pid string) (
	pod domain.PodInfo, err error,
) {
	filter := map[string]interface{}{
		fieldId: pid,
	}

	var tpod TPod
	if err = impl.cli.First(filter, &tpod); err != nil {
		return
	}

	if err = tpod.toPodInfo(&pod); err != nil {
		return
	}

	return
}

func (impl *podRepoImpl) GetUserPod(user types.Account) (
	pods repository.PodInfoList, err error,
) {
	filter := map[string]interface{}{
		fieldOwner: user.Account(),
	}

	return impl.getFilterPods(filter)
}

func (impl *podRepoImpl) GetUserCloudIdLastPod(
	user types.Account, cloudId string,
) (domain.PodInfo, error) {
	filter := map[string]interface{}{
		fieldOwner:   user.Account(),
		fieldCloudId: cloudId,
	}

	order := "created_at DESC"

	return impl.getOrderOnePod(filter, order)
}

func (impl *podRepoImpl) AddStartingPod(p *domain.PodInfo) (pid string, err error) {
	pod := new(TPod)
	pod.toTPod(p)

	err = impl.cli.Create(pod)
	pid = pod.Id

	return
}

func (impl *podRepoImpl) UpdatePod(p *domain.PodInfo) error {
	pod := new(TPod)
	pod.toTPod(p)

	logrus.Debugf(
		"update pod(%s/%s) to %v",
		p.Pod.Id, p.Pod.CloudId, p.AccessURL,
	)

	filter := map[string]interface{}{
		fieldId: pod.Id,
	}

	return impl.cli.Updates(filter, pod)
}
