package repository

import (
	"github.com/opensourceways/xihe-server/course/domain"
)

type Course interface {
	FindCourse(cid string) (domain.Course, error)
	FindCourses(*CourseListOption) ([]domain.CourseSummary, error)
	FindAssignments(cid string) ([]domain.Assignment, error)
	FindAssignment(cid string, asgId string) (domain.Assignment, error)
}

type CourseSummary struct {
	domain.CourseSummary
	CompetitorCount int
}

type CourseListOption struct {
	CourseIds []string
	Status    domain.CourseStatus
	Type      domain.CourseType
}
