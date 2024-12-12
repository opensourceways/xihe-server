package app

import (
	common "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
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
	Id   string
	User types.Account
}

type PromotionDTO struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Desc           string `json:"desc"`
	Poster         string `json:"poster"`
	Status         string `json:"status"`
	IsRegister     bool   `json:"is_register"`
	Total          int    `json:"total"`
	Duration       string `json:"duration"`
	RegistrantsNum int    `json:"registrants_num"`
	Host           string `json:"host"`
	Way            string `json:"way"`
	Type           string `json:"type"`
	Intro          string `json:"intro"`
	IsStatic       bool   `json:"is_static"`
}

func (dto *PromotionDTO) toDTO(p *domain.Promotion, user types.Account, total int) error {
	var err error

	*dto = PromotionDTO{
		Id:             p.Id,
		Name:           p.Name.PromotionName(),
		Desc:           p.Desc.PromotionDesc(),
		Poster:         p.Poster,
		IsRegister:     p.HasRegister(user),
		Total:          total,
		RegistrantsNum: p.CountRegUsers(),
		Host:           p.Host,
		Intro:          p.Intro,
		IsStatic:       p.IsStatic,
	}

	if dto.Status, err = p.Status(); err != nil {
		return err
	}

	if dto.Duration, err = p.Duration(); err != nil {
		return err
	}

	if p.Type != nil {
		dto.Type = p.Type.PromotionType()
	}

	if p.Way != nil {
		dto.Way = p.Way.PromotionWay()
	}

	return nil
}

type UserRegistrationCmd struct {
	PromotionId string
	domain.UserRegistration
	Origin domain.Origin
}

type ListPromotionsCmd struct {
	User     types.Account
	Type     domain.PromotionType
	Status   domain.PromotionStatus
	Way      domain.PromotionWay
	PageNo   int
	PageSize int
}

func (cmd *ListPromotionsCmd) toPromotionsQuery() *repository.PromotionsQuery {
	query := &repository.PromotionsQuery{
		Offset: int64((cmd.PageNo - 1) * cmd.PageSize),
		Limit:  int64(cmd.PageSize),
		Sort: [][2]string{
			{repository.SortFieldPriority, repository.SortDesc},
			{repository.SortFieldStartTime, repository.SortDesc},
		},
	}

	query.Type = cmd.Type
	query.Status = cmd.Status
	query.Way = cmd.Way

	return query
}

type PromotionsDTO struct {
	Items []PromotionDTO `json:"items"`
	Total int64          `json:"total"`
}
