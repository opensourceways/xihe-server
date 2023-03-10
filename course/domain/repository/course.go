package repository

import (
	"github.com/opensourceways/xihe-server/course/domain"
)

type Course interface {
	FindCourse(cid string) (domain.Course, error)
	FindCourses(*CourseListOption) ([]domain.CourseSummary, error)
}

type CourseSummary struct {
	domain.CourseSummary
	CompetitorCount int
}

type CourseListOption struct {
	Status domain.CourseStatus
	Type   domain.CourseType
}
