package repositoryadapter

import "github.com/opensourceways/xihe-server/points/domain"

func PointsRuleAdapter() *pointsRuleAdapter {
	return &pointsRuleAdapter{}
}

type pointsRuleAdapter struct {
}

func (impl *pointsRuleAdapter) FindPointsItemRules() ([]domain.PointsItemRule, error) {
	return nil, nil
}

func (impl *pointsRuleAdapter) PointsOfDay() (int, error) {
	return 0, nil
}
