package cloudimpl

import (
	"github.com/opensourceways/xihe-inference-evaluate/sdk"
	"github.com/opensourceways/xihe-server/cloud/domain/cloud"
)

func NewCloud(cfg *Config) cloud.CloudPod {
	v := sdk.NewInferenceEvaluate(cfg.ContainerManagerEndpoint)

	return &cloudpodImpl{
		cli: &v,
	}
}

type cloudpodImpl struct {
	cli *sdk.InferenceEvaluate
}

func (impl *cloudpodImpl) Create(info *cloud.CloudPodCreateInfo) error {
	opt := &sdk.CloudPodCreateOption{
		PodId:        info.PodId,
		User:         info.User,
		SurvivalTime: info.SurvivalTime,
		CloudType:    info.CloudType,
		CloudImage:   info.CloudImage,
	}

	return impl.cli.CreateCloudPod(opt)
}

func (impl *cloudpodImpl) Release(podId, cloudType string) error {
	return impl.cli.ReleaseCloudPod(podId, cloudType)
}
