package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/course/domain"
	projdomain "github.com/opensourceways/xihe-server/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

// player
type PlayerApplyCmd domain.Player

func (cmd *PlayerApplyCmd) Validate() error {
	b := cmd.Student.Account != nil &&
		cmd.Student.Name != nil &&
		cmd.Student.Email != nil &&
		cmd.Student.Identity != nil

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *PlayerApplyCmd) toPlayer() (p domain.Player) {
	return *(*domain.Player)(cmd)
}

// course
type CourseListCmd struct {
	Status domain.CourseStatus
	Type   domain.CourseType
	User   types.Account
}

type CourseGetCmd struct {
	User types.Account
	Cid  string
}

type CourseAddReleatedProjectCmd struct {
	Cid     string
	User    types.Account
	Project types.ResourceSummary
}

func (cmd *CourseAddReleatedProjectCmd) repo() string {
	return cmd.User.Account() + "/" + cmd.Project.Name.ResourceName()
}

type CourseDTO struct {
	CourseSummaryDTO

	IsApply bool `json:"is_apply"`

	Teacher  string       `json:"teacher"`
	Doc      string       `json:"doc"`
	Forum    string       `json:"forum"`
	Sections []SectionDTO `json:"sections"`
}

func (dto *CourseDTO) toCourseDTO(c *domain.Course, apply bool, count int) {
	toCourseSummaryDTO(&c.CourseSummary, count, &dto.CourseSummaryDTO)

	dto.IsApply = apply

	dto.Teacher = c.Teacher.URL()

	dto.Doc = c.Doc.URL()

	dto.Forum = c.Forum.URL()

	dto.Sections = make([]SectionDTO, len(c.Sections))
	for i := range dto.Sections {
		dto.Sections[i].toSectionDTO(&c.Sections[i])
	}
}

func (dto *CourseDTO) toCourseNoVideoDTO(c *domain.Course, apply bool, count int) {
	toCourseSummaryDTO(&c.CourseSummary, count, &dto.CourseSummaryDTO)

	dto.IsApply = apply

	dto.Teacher = c.Teacher.URL()

	dto.Doc = c.Doc.URL()

	dto.Forum = c.Forum.URL()

	dto.Sections = make([]SectionDTO, len(c.Sections))
	for i := range dto.Sections {
		dto.Sections[i].toSectionNoVideoDTO(&c.Sections[i])
	}
}

type CourseSummaryDTO struct {
	PlayerCount int    `json:"count"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Hours       int    `json:"hours"`
	Host        string `json:"host"`
	Desc        string `json:"desc"`
	Status      string `json:"status"`
	Poster      string `json:"poster"`
	Duration    string `json:"duration"`
	Type        string `json:"type"`
}

func toCourseSummaryDTO(
	c *domain.CourseSummary, playerCount int, dto *CourseSummaryDTO,
) {
	*dto = CourseSummaryDTO{
		PlayerCount: playerCount,
		Id:          c.Id,
		Name:        c.Name.CourseName(),
		Hours:       c.Hours.CourseHours(),
		Host:        c.Host.CourseHost(),
		Desc:        c.Desc.CourseDesc(),
		Type:        c.Type.CourseType(),
		Status:      c.Status.CourseStatus(),
		Poster:      c.Poster.URL(),
		Duration:    c.Duration.CourseDuration(),
	}
}

// Section
type SectionDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Lessons []LessonDTO `json:"lessons"`
}

func (dto *SectionDTO) toSectionDTO(s *domain.Section) {
	dto.Id = s.Id

	dto.Name = s.Name.SectionName()

	dto.Lessons = make([]LessonDTO, len(s.Lessons))
	for i := range dto.Lessons {
		dto.Lessons[i].toLessonDTO(&s.Lessons[i])
	}
}

func (dto *SectionDTO) toSectionNoVideoDTO(s *domain.Section) {
	dto.Id = s.Id

	dto.Name = s.Name.SectionName()

	dto.Lessons = make([]LessonDTO, len(s.Lessons))
	for i := range dto.Lessons {
		dto.Lessons[i].toLessonNoVideoDTO(&s.Lessons[i])
	}
}

// Lesson
type LessonDTO struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Video string `json:"video"`

	Points []PointDTO `json:"points"`
}

func (dto *LessonDTO) toLessonDTO(l *domain.Lesson) {
	dto.Id = l.Id

	dto.Name = l.Name.LessonName()

	dto.Desc = l.Desc.LessonDesc()

	if !l.HasPoints() {
		dto.Video = l.Video.LessonURL()
	} else {
		dto.Points = make([]PointDTO, len(l.Points))
		for i := range dto.Points {
			dto.Points[i].toPointDTO(&l.Points[i])
		}
	}
}

func (dto *LessonDTO) toLessonNoVideoDTO(l *domain.Lesson) {
	dto.Id = l.Id

	dto.Name = l.Name.LessonName()

	dto.Desc = l.Desc.LessonDesc()

	if l.HasPoints() {
		dto.Points = make([]PointDTO, len(l.Points))
		for i := range dto.Points {
			dto.Points[i].toPointNoVideoDTO(&l.Points[i])
		}
	}
}

// Point
type PointDTO struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Video string `json:"video"`
}

func (dto *PointDTO) toPointDTO(p *domain.Point) {
	dto.Id = p.Id

	dto.Name = p.Name.PointName()

	dto.Video = p.Video.URL()
}

func (dto *PointDTO) toPointNoVideoDTO(p *domain.Point) {
	dto.Id = p.Id

	dto.Name = p.Name.PointName()
}

// Assignment
type AsgListCmd struct {
	User   types.Account
	Cid    string
	Status domain.WorkStatus
}

type AsgGetCmd struct {
	User  types.Account
	Cid   string
	AsgId string
}

type GetSubmissionCmd struct {
	User   types.Account
	Cid    string
	Status domain.WorkStatus
}

type RecordAddCmd domain.Record

func (cmd *RecordAddCmd) toRecord() (p domain.Record) {
	return *(*domain.Record)(cmd)
}

type AsgDTO struct {
	AsgId      string `json:"asg_id"`
	CourseName string `json:"course_name"`
	Deadline   string `json:"deadline"`
	AsgName    string `json:"asg_name"`
	Desc       string `json:"desc"`
}

func toAsgDTO(
	a *domain.Assignment, c *domain.Course, dto *AsgDTO,
) {
	*dto = AsgDTO{
		AsgId:      a.Id,
		AsgName:    a.Name.AsgName(),
		CourseName: c.Name.CourseName(),
		Deadline:   a.DeadLine.AsgDeadLine(),
		Desc:       a.Desc.URL(),
	}

}

type AsgWorkDTO struct {
	Id       string  `json:"id"`
	Deadline string  `json:"deadline"`
	Name     string  `json:"name"`
	Desc     string  `json:"desc"`
	Score    float32 `json:"score"`
	Status   string  `json:"status"`
}

func toAsgWorkDTO(
	c *domain.Assignment, score float32, status string, dto *AsgWorkDTO,
) {
	*dto = AsgWorkDTO{
		Id:       c.Id,
		Deadline: c.DeadLine.AsgDeadLine(),
		Name:     c.Name.AsgName(),
		Desc:     c.Desc.URL(),
		Score:    score,
		Status:   status,
	}
}

type ProjectSummuryDTO struct {
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	CoverId       string `json:"cover_id"`
	CreatedAt     string `json:"create_at"`
	LikeCount     int    `json:"like_count"`
	ForkCount     int    `json:"fork_count"`
	DownloadCount int    `json:"download_count"`
}
type RelateProjectDTO struct {
	RelatedProject []ProjectSummuryDTO `json:"related_project"`
}

func toProjectSummuryDTO(p *projdomain.Project, dto *ProjectSummuryDTO) {
	*dto = ProjectSummuryDTO{
		Owner:         p.Owner.Account(),
		Name:          p.Name.ResourceName(),
		CoverId:       p.CoverId.CoverId(),
		CreatedAt:     utils.ToDate(p.CreatedAt),
		LikeCount:     p.LikeCount,
		ForkCount:     p.ForkCount,
		DownloadCount: p.DownloadCount,
	}
}

type CertInfoDTO struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Cert   string `json:"cert"`
	IsPass bool   `json:"is_pass"`
}

func toCertInfoDTO(user types.Account, c *domain.Course, pass bool, dto *CertInfoDTO) {
	*dto = CertInfoDTO{
		Owner:  user.Account(),
		Name:   c.Name.CourseName(),
		Cert:   c.Cert.URL(),
		IsPass: pass,
	}
}
