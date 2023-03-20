package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

// Course
type CourseSummary struct {
	Id    string
	Name  CourseName
	Desc  CourseDesc
	Host  CourseHost
	Hours CourseHours

	Type     CourseType
	Status   CourseStatus
	Duration CourseDuration
	Poster   URL
}

type Course struct {
	CourseSummary

	Teacher   URL
	Doc       URL
	Forum     URL
	PassScore CoursePassScore
	Cert      URL
	Sections  []Section
}

// CourseRepo
type CourseProject struct {
	Owner       types.Account
	RepoRouting string
}

func NewCourseProject(owner types.Account, r string) CourseProject {
	return CourseProject{
		Owner:       owner,
		RepoRouting: r,
	}
}

// Assignment
type Assignment struct {
	Id       string
	Name     AsgName
	Desc     URL
	DeadLine AsgDeadLine
}

// Section
type Section struct {
	Id   string
	Name SectionName

	Lessons []Lesson
}

// Lesson
type Lesson struct {
	Id    string
	Name  LessonName
	Desc  LessonDesc
	Video LessonURL

	Points []Point
}

// Point
type Point struct {
	Id    string
	Name  PointName
	Video URL
}

func (c *Course) IsOver() bool {
	return c.Status != nil && c.Status.IsOver()
}

func (c *Course) IsPreliminary() bool {
	return c.Status != nil && c.Status.IsPreliminary()
}

func (c *Course) IsApplyed(p *Player) bool {
	return p.CourseId == c.Id
}

func (l *Lesson) HasPoints() bool {
	return len(l.Points) > 0
}
