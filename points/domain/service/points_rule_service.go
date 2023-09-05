package service

import (
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
)

var instance *pointsRuleService

func PointsRuleService() *pointsRuleService {
	return instance
}

func InitPointsRuleService(repo repository.PointsRule) (err error) {
	s := pointsRuleService{}

	if s.maxPointsOfDay, err = repo.PointsOfDay(); err != nil {
		return err
	}

	if s.rules, err = repo.FindPointsItemRules(); err != nil {
		return err
	}

	s.repo = repo

	instance = &s

	return nil
}

// pointsRuleService
type pointsRuleService struct {
	repo           repository.PointsRule
	rules          []domain.PointsItemRule
	maxPointsOfDay int
}

func (s *pointsRuleService) Calculator(ruleType string) domain.PointsCalculator {
	if s == nil {
		return nil
	}

	rs := s.rules

	for i := range rs {
		if rs[i].Type == ruleType {
			return calculator{&rs[i]}
		}
	}

	return nil
}

// calculator
type calculator struct {
	r *domain.PointsItemRule
}

// pointsOfDay is the total points that user has got that day
// pointsOfItem is the points that user has got on the item that day
func (c calculator) Calc(pointsOfDay, pointsOfItem int) int {
	if pointsOfDay >= instance.maxPointsOfDay {
		return 0
	}

	v := c.r.Calc(pointsOfItem)
	if v == 0 {
		return 0
	}

	if n := instance.maxPointsOfDay - pointsOfDay; v >= n {
		return n
	}

	return v
}
