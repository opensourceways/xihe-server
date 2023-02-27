package repository

import "github.com/opensourceways/xihe-server/course/domain"

type Course interface {
	FindCourse(cid string) (domain.Course, error)
}