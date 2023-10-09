package message

import "github.com/opensourceways/xihe-server/course/domain"

type MessageProducer interface {
	SendCourseAppliedEvent(*domain.CourseAppliedEvent) error
}
