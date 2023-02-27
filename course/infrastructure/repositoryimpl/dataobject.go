package repositoryimpl

import "github.com/opensourceways/xihe-server/course/domain"

// Course
func (doc *DCourse) toCourse(c *domain.Course) (err error) {
	// courseSummary
	if err = doc.toCourseSummary(&c.CourseSummary); err != nil {
		return
	}

	// course
	c.Id = doc.Id

	if c.Doc, err = domain.NewURL(doc.Doc); err != nil {
		return
	}

	if c.Forum, err = domain.NewURL(doc.Forum); err != nil {
		return
	}

	if c.Type, err = domain.NewCourseType(doc.Type); err != nil {
		return
	}

	// section
	sections := make([]domain.Section, len(doc.Sections))
	for i := range doc.Sections {
		if err = doc.Sections[i].toSection(&sections[i]); err != nil {
			return
		}
	}

	return
}

func (doc *DCourse) toCourseSummary(c *domain.CourseSummary) (err error) {
	if c.Name, err = domain.NewCourseName(doc.Name); err != nil {
		return
	}

	if c.Desc, err = domain.NewCourseDesc(doc.Desc); err != nil {
		return
	}

	if c.Host, err = domain.NewCourseHost(doc.Host); err != nil {
		return
	}

	if c.Teacher, err = domain.NewURL(doc.Teacher); err != nil {
		return
	}

	if c.PassScore, err = domain.NewCoursePassScore(doc.PassScore); err != nil {
		return
	}

	if c.Status, err = domain.NewCourseStatus(doc.Status); err != nil {
		return
	}

	if c.Duration, err = domain.NewCourseDuration(doc.Duration); err != nil {
		return
	}

	if c.Poster, err = domain.NewURL(doc.Poster); err != nil {
		return
	}

	if c.Cert, err = domain.NewURL(doc.Cert); err != nil {
		return
	}

	return
}

func (doc *dSection) toSection(s *domain.Section) (err error) {
	s.Id = doc.Id

	if s.Name, err = domain.NewSectionName(doc.Name); err != nil {
		return
	}

	// lesson
	lessons := make([]domain.Lesson, len(doc.Lessons))
	for i := range doc.Lessons {
		if err = doc.Lessons[i].toLesson(&lessons[i]); err != nil {
			return
		}
	}

	return
}

func (doc *dLesson) toLesson(l *domain.Lesson) (err error) {
	l.Id = doc.Id

	if l.Name, err = domain.NewLessonName(doc.Name); err != nil {
		return
	}

	if l.Desc, err = domain.NewLessonDesc(doc.Desc); err != nil {
		return
	}

	if l.Video, err = domain.NewURL(doc.Video); err != nil {
		return
	}

	return
}
