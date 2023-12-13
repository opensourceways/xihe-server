package app

import (
	common "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type PointsCmd struct {
	Promotionid string
	User        types.Account
	Lang        common.Language
}

type PointsDTO struct {
	Items []Item `json:"items"`
	Total int    `json:"total"`
}

type PointsRankDTO struct {
	User  string `json:"user"`
	Point int    `json:"point"`
}

func (dto *PointsRankDTO) toDTO(up domain.UserPoints) {
	*dto = PointsRankDTO{
		User:  up.User.Account(),
		Point: up.Total,
	}
}

type Item struct {
	TaskName string `json:"task_name"`
	Descs    string `json:"descs"`
	Points   int    `json:"points"`
	Time     string `json:"time"`
}

func toPointsDTO(p domain.UserPoints, lang common.Language) PointsDTO {
	items := make([]Item, len(p.Items))

	for i := range p.Items {
		items[i] = Item{
			TaskName: p.Items[i].TaskName.Sentence(lang),
			Descs:    p.Items[i].Descs.Sentence(lang),
			Points:   p.Items[i].Points,
			Time:     utils.ToDate(p.Items[i].Date),
		}
	}

	return PointsDTO{
		Items: items,
		Total: p.Total,
	}
}

type PromotionCmd struct {
	promotionId string
	User        types.Account
}

type PromotionDTO struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	IsRegister bool   `json:"is_register"`
	Duration   string `json:"duration"`
}

func (dto *PromotionDTO) toDTO(p *domain.Promotion, user *types.Account) {
	*dto = PromotionDTO{
		Id:         p.Id,
		Name:       p.Name.PromotionName(),
		Desc:       p.Desc.PromotionDesc(),
		IsRegister: p.HasRegister(*user),
		Duration:   p.Duration.PromotionDuration(),
	}
}

type UserRegistrationCmd struct {
	PromotionId string
	domain.UserRegistration
}
