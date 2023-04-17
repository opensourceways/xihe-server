package inferenceimpl

import (
	"github.com/opensourceways/xihe-inference-evaluate/sdk"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/inference"
)

func NewInference(cfg *Config) *inferenceImpl {
	v := sdk.NewInferenceEvaluate(cfg.ContainerManagerEndpoint)

	m := map[string]bool{}
	for _, item := range cfg.ProjectTagsForOfficial {
		m[item] = true
	}

	return &inferenceImpl{
		cli:                     &v,
		survivalTimeForNormal:   cfg.SurvivalTimeForNormal,
		survivalTimeForOfficial: cfg.SurvivalTimeForOfficial,
		projectTagsForOfficial:  m,
	}
}

type inferenceImpl struct {
	cli *sdk.InferenceEvaluate

	survivalTimeForNormal   int
	survivalTimeForOfficial int
	projectTagsForOfficial  map[string]bool
}

func (impl *inferenceImpl) GetSurvivalTime(info *domain.InferenceInfo) int {
	if impl.projectTagsForOfficial[info.ResourceLevel] {
		return impl.survivalTimeForOfficial
	}

	return impl.survivalTimeForNormal
}

func (impl *inferenceImpl) Create(info *inference.InferenceInfo) (int, error) {
	survivalTime := impl.GetSurvivalTime(info.InferenceInfo)

	opt := sdk.InferenceCreateOption{
		UserToken:    info.UserToken,
		LastCommit:   info.LastCommit,
		ProjectName:  info.ProjectName.ResourceName(),
		SurvivalTime: survivalTime,
	}
	opt.User = info.Project.Owner.Account()
	opt.ProjectId = info.Project.Id
	opt.InferenceId = info.Id

	return survivalTime, impl.cli.CreateInference(&opt)
}

func (impl *inferenceImpl) ExtendSurvivalTime(index *domain.InferenceIndex, timeToExtend int) error {
	opt := sdk.InferenceUpdateOption{
		TimeToExtend: timeToExtend,
	}
	opt.User = index.Project.Owner.Account()
	opt.ProjectId = index.Project.Id
	opt.InferenceId = index.Id

	return impl.cli.ExtendExpiryOfInference(&opt)
}
