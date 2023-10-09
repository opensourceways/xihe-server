package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

// CourseAppliedEvent
type CourseAppliedEvent struct {
	Account    types.Account
	CourseName CourseName
}
