package bigmodels

import (
	"errors"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	moderation "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3/model"

	moderationv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v2"
	modelv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v2/model"
	regionv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v2/region"

	"github.com/opensourceways/xihe-server/domain/bigmodel"
)

func initTextCheck(cfg *Moderation) textCheckService {
	auth := basic.NewCredentialsBuilder().
		WithAk(cfg.AccessKey).
		WithSk(cfg.SecretKey).
		WithIamEndpointOverride(cfg.IAMEndpint).
		Build()

	cli := moderation.NewModerationClient(
		moderation.ModerationClientBuilder().
			WithRegion(region.NewRegion(cfg.Region, cfg.Endpoint)).
			WithCredential(auth).
			Build(),
	)

	authv2 := basic.NewCredentialsBuilder().
		WithAk(cfg.AccessKey).
		WithSk(cfg.SecretKey).
		Build()

	cliv2 := moderationv2.NewModerationClient(
		moderation.ModerationClientBuilder().
			WithRegion(regionv2.ValueOf(cfg.Region)).
			WithCredential(authv2).
			Build(),
	)

	return textCheckService{cli, cliv2}
}

type textCheckService struct {
	cli   *moderation.ModerationClient
	cliv2 *moderationv2.ModerationClient
}

func (s *textCheckService) check(content string) error {
	request := &model.RunTextModerationRequest{
		Body: &model.TextDetectionReq{
			Data: &model.TextDetectionDataReq{
				Text: content,
			},
			EventType: "comment",
		},
	}

	resp, err := s.cli.RunTextModeration(request)
	if err != nil {
		return err
	}

	if *resp.Result.Suggestion != "pass" {
		return bigmodel.NewErrorSensitiveInfo(errors.New("invalid text"))
	}

	return nil
}

func (s *textCheckService) checkImages(urls []string) error {
	request := &modelv2.RunImageBatchModerationRequest{}
	var listCategoriesbody = []modelv2.ImageBatchModerationReqCategories{
		modelv2.GetImageBatchModerationReqCategoriesEnum().ALL,
	}
	var listUrlsbody = urls
	rule := "default"
	request.Body = &modelv2.ImageBatchModerationReq{
		ModerationRule: &rule,
		Categories:     &listCategoriesbody,
		Urls:           listUrlsbody,
	}
	resp, err := s.cliv2.RunImageBatchModeration(request)
	if err != nil {
		return err
	}

	results := resp.Result
	for _, res := range *results {
		if *res.Suggestion != "pass" {
			return bigmodel.NewErrorSensitiveInfo(errors.New("the generated image is illegal, please try again"))
		}
	}

	return nil
}
