package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

// PlayRecord
type Record struct {
	Cid         string
	User        types.Account
	SectionId   SectionId
	LessonId    LessonId
	PointId     string
	PlayCount   int
	FinishCount int
}
