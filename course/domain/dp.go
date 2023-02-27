package domain

import (
	"errors"
	"net/url"
	"time"
)

const (
	courseTypeAI = "AI"

	studentIdentityStudent   = "student"
	studentIdentityTeacher   = "teacher"
	studentIdentityDeveloper = "developer"

	courseStatusOver       = "over"
	courseStatusPreparing  = "preparing"
	courseStatusInProgress = "in-progress"
)

// StudentName
type StudentName interface {
	StudentName() string
}

func NewStudentName(v string) (StudentName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return studentName(v), nil
}

type studentName string

func (r studentName) StudentName() string {
	return string(r)
}

// City
type City interface {
	City() string
}

func NewCity(v string) (City, error) {
	return city(v), nil
}

type city string

func (r city) City() string {
	return string(r)
}

// Phone
type Phone interface {
	Phone() string
}

func NewPhone(v string) (Phone, error) {
	return phone(v), nil
}

type phone string

func (r phone) Phone() string {
	return string(r)
}

// StudentIdentity
type StudentIdentity interface {
	StudentIdentity() string
}

func NewStudentIdentity(v string) (StudentIdentity, error) {
	b := v == studentIdentityStudent ||
		v == studentIdentityTeacher ||
		v == studentIdentityDeveloper ||
		v == ""

	if !b {
		return nil, errors.New("invalid student identity")
	}

	return studentIdentity(v), nil
}

type studentIdentity string

func (r studentIdentity) StudentIdentity() string {
	return string(r)
}

// Province
type Province interface {
	Province() string
}

func NewProvince(v string) (Province, error) {
	return province(v), nil
}

type province string

func (r province) Province() string {
	return string(r)
}

// URL
type URL interface {
	URL() string
}

func NewURL(v string) (URL, error) {
	if v == "" {
		return nil, errors.New("empty url")
	}

	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid url")
	}

	return dpURL(v), nil
}

type dpURL string

func (r dpURL) URL() string {
	return string(r)
}

// CourseType
type CourseType interface {
	CourseType() string
}

func NewCourseType(v string) (CourseType, error) {
	if v == "" || v == courseTypeAI {
		return courseType(v), nil
	}

	return nil, errors.New("invalid course type")
}

type courseType string

func (r courseType) CourseType() string {
	return string(r)
}

// CourseName
type CourseName interface {
	CourseName() string
}

func NewCourseName(v string) (CourseName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return courseName(v), nil
}

type courseName string

func (r courseName) CourseName() string {
	return string(r)
}

// CourseDesc
type CourseDesc interface {
	CourseDesc() string
}

func NewCourseDesc(v string) (CourseDesc, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return courseDesc(v), nil
}

type courseDesc string

func (r courseDesc) CourseDesc() string {
	return string(r)
}

// CourseHost
type CourseHost interface {
	CourseHost() string
}

func NewCourseHost(v string) (CourseHost, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return courseHost(v), nil
}

type courseHost string

func (r courseHost) CourseHost() string {
	return string(r)
}

// CoursePassScore
type CoursePassScore interface {
	CoursePassScore() int
}

func NewCoursePassScore(v int) (CoursePassScore, error) {
	if v == 0 {
		return nil, errors.New("zero value")
	}

	return coursePassScore(v), nil
}

type coursePassScore int

func (r coursePassScore) CoursePassScore() int {
	return int(r)
}

// CourseStatus
type CourseStatus interface {
	CourseStatus() string
	IsEnabled() bool
	IsOver() bool
	IsPreliminary() bool
}

func NewCourseStatus(v string) (CourseStatus, error) {
	b := v == courseStatusOver ||
		v == courseStatusPreparing ||
		v == courseStatusInProgress

	if b {
		return courseStatus(v), nil

	}

	return nil, errors.New("invalid course status")
}

type courseStatus string

func (r courseStatus) CourseStatus() string {
	return string(r)
}

func (r courseStatus) IsEnabled() bool {
	return string(r) == courseStatusInProgress
}

func (r courseStatus) IsOver() bool {
	return string(r) == courseStatusOver
}

func (r courseStatus) IsPreliminary() bool {
	return string(r) == courseStatusPreparing
}

// CourseDuration
type CourseDuration interface {
	CourseDuration() string
}

func NewCourseDuration(v string) (CourseDuration, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return courseDuration(v), nil
}

type courseDuration string

func (r courseDuration) CourseDuration() string {
	return string(r)
}

// Assignment Name
type AsgName interface {
	AsgName() string
}

func NewAsgName(v string) (AsgName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return asgName(v), nil
}

type asgName string

func (r asgName) AsgName() string {
	return string(r)
}

// Assignment Desc
type AsgDesc interface {
	AsgDesc() string
}

func NewAsgDesc(v string) (AsgDesc, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return asgDesc(v), nil
}

type asgDesc string

func (r asgDesc) AsgDesc() string {
	return string(r)
}

// Assignment DeadLine
type AsgDeadLine interface {
	AsgDeadLine() string
}

func NewAsgDeadLine(v string) (AsgDeadLine, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return asgDeadLine(v), nil
}

type asgDeadLine string

func (r asgDeadLine) AsgDeadLine() string {
	return string(r)
}

// SectionName
type SectionName interface {
	SectionName() string
}

func NewSectionName(v string) (SectionName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return sectionName(v), nil
}

type sectionName string

func (r sectionName) SectionName() string {
	return string(r)
}

// LessonName
type LessonName interface {
	LessonName() string
}

func NewLessonName(v string) (LessonName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return lessonName(v), nil
}

type lessonName string

func (r lessonName) LessonName() string {
	return string(r)
}

// LessonDesc
type LessonDesc interface {
	LessonDesc() string
}

func NewLessonDesc(v string) (LessonDesc, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return lessonDesc(v), nil
}

type lessonDesc string

func (r lessonDesc) LessonDesc() string {
	return string(r)
}

// Time
type CourseTime interface {
	CourseTimeStr() string
	CourseTimeInt() int64
}

func NewCourseTime(v string) (CourseTime, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return coursetime(v), nil
}

type coursetime string

func (r coursetime) CourseTimeStr() string {
	return string(r)
}

func (r coursetime) CourseTimeInt() int64 {
	t, _ := time.Parse(string(r), "2022-02-15")

	return t.Unix()
}
