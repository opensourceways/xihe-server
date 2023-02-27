package domain

// Course
type CourseSummary struct {
	Name    CourseName
	Desc    CourseDesc
	Host    CourseHost
	Teacher URL

	PassScore CoursePassScore
	Status    CourseStatus
	Duration  CourseDuration
	Poster    URL
	Cert      URL
}

type Course struct {
	CourseSummary

	Id    string
	Doc   URL
	Forum URL

	Type CourseType

	Sections []Section
}

// Assignment
type Assignment struct {
	Id       string
	Name     AsgName
	Desc     AsgDesc
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
	Video URL
}

func (c *Course) IsOver() bool {
	return c.Status != nil && c.Status.IsOver()
}

func (c *Course) IsPreliminary() bool {
	return c.Status != nil && c.Status.IsPreliminary()
}
