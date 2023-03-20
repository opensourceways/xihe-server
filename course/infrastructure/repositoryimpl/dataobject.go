package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/repository"
)

// Course
func (doc *DCourse) toCourse(c *domain.Course) (err error) {
	// courseSummary
	if err = doc.toCourseSummary(&c.CourseSummary); err != nil {
		return
	}

	// course
	c.Id = doc.Id

	if c.Teacher, err = domain.NewURL(doc.Teacher); err != nil {
		return
	}

	if c.Doc, err = domain.NewURL(doc.Doc); err != nil {
		return
	}

	if c.Forum, err = domain.NewURL(doc.Forum); err != nil {
		return
	}

	if c.PassScore, err = domain.NewCoursePassScore(doc.PassScore); err != nil {
		return
	}

	if c.Cert, err = domain.NewURL(doc.Cert); err != nil {
		return
	}

	// section
	c.Sections = make([]domain.Section, len(doc.Sections))
	for i := range doc.Sections {
		if err = doc.Sections[i].toSection(&c.Sections[i]); err != nil {
			return
		}
	}

	return
}

func (doc *DCourse) toCourseSummary(c *domain.CourseSummary) (err error) {
	c.Id = doc.Id

	if c.Name, err = domain.NewCourseName(doc.Name); err != nil {
		return
	}

	if c.Desc, err = domain.NewCourseDesc(doc.Desc); err != nil {
		return
	}

	if c.Host, err = domain.NewCourseHost(doc.Host); err != nil {
		return
	}

	if c.Hours, err = domain.NewCourseHours(doc.Hours); err != nil {
		return
	}

	if c.Type, err = domain.NewCourseType(doc.Type); err != nil {
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

	return
}

func (doc *dSection) toSection(s *domain.Section) (err error) {
	s.Id = doc.Id

	if s.Name, err = domain.NewSectionName(doc.Name); err != nil {
		return
	}

	// lesson
	s.Lessons = make([]domain.Lesson, len(doc.Lessons))
	for i := range doc.Lessons {
		if err = doc.Lessons[i].toLesson(&s.Lessons[i]); err != nil {
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

	if l.Video, err = domain.NewLessonURL(doc.Video); err != nil {
		return
	}

	// point
	l.Points = make([]domain.Point, len(doc.Points))
	for i := range doc.Points {
		if err = doc.Points[i].toPoint(&l.Points[i]); err != nil {
			return
		}
	}

	return
}

func (doc *dPoint) toPoint(p *domain.Point) (err error) {
	p.Id = doc.Id

	if p.Name, err = domain.NewPointName(doc.Name); err != nil {
		return
	}

	if p.Video, err = domain.NewURL(doc.Video); err != nil {
		return
	}

	return
}

// Assignments

func (doc *dAssignments) toAssignment(c *domain.Assignment) (err error) {
	c.Id = doc.Id

	if c.Name, err = domain.NewAsgName(doc.Name); err != nil {
		return
	}

	if c.Desc, err = domain.NewURL(doc.Desc); err != nil {
		return
	}

	if c.DeadLine, err = domain.NewAsgDeadLine(doc.DeadLine); err != nil {
		return
	}

	return
}

// player
func (doc *DCoursePlayer) toPlayerNoStudent(p *repository.PlayerVersion) (err error) {
	p.Player.Id = doc.Id

	p.Player.CourseId = doc.CourseId
	p.Player.RelatedProject = doc.Repo

	if p.CreatedAt, err = domain.NewCourseTime(doc.CreatedAt); err != nil {
		return
	}

	return
}

// work
func (doc *DCourseWork) toCourseWork(w *domain.Work) (err error) {
	w.Score = doc.Score
	w.AsgId = doc.AsgId
	w.CourseId = doc.CourseId
	w.PlayerId = doc.Account
	w.Status = doc.Status
	w.Version = doc.Version

	return
}
