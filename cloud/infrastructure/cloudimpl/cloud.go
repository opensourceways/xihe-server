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
	}

	return impl.cli.CreateCloudPod(opt)
}
