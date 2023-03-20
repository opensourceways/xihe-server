package repositoryimpl

const (
	fieldId       = "id"
	fieldCourseId = "course_id"
	fieldAccount  = "account"
	fieldStatus   = "status"
	fieldType     = "type"
	fieldVersion  = "version"
	fieldRepo     = "repo"
	fieldAsgId    = "asg_id"
)

// Course
type DCourse struct {
	Id        string  `bson:"id"              json:"id"`
	Name      string  `bson:"name"            json:"name"`
	Teacher   string  `bson:"teacher"         json:"teacher"`
	Desc      string  `bson:"desc"            json:"desc"`
	Host      string  `bson:"host"            json:"host"`
	Type      string  `bson:"type"            json:"type"`
	PassScore float32 `bson:"pass_score"      json:"pass_score"`
	Status    string  `bson:"status"          json:"status"`
	Duration  string  `bson:"duration"        json:"duration"`
	Hours     int     `bson:"hours"           json:"hours"`
	Doc       string  `bson:"doc"             json:"doc"`
	Forum     string  `bson:"forum"           json:"forum"`
	Poster    string  `bson:"poster"          json:"poster"`
	Cert      string  `bson:"cert"            json:"cert"`

	Assignments []dAssignments `bson:"assignments"  json:"-"`
	Sections    []dSection     `bson:"sections"     json:"-"`
}

type dSection struct {
	Id   string `bson:"id"              json:"id"`
	Name string `bson:"name"            json:"name"`

	Lessons []dLesson `bson:"lessons"   json:"-"`
}

type dLesson struct {
	Id    string `bson:"id"               json:"id"`
	Name  string `bson:"name"             json:"name"`
	Desc  string `bson:"desc"             json:"desc"`
	Video string `bson:"video"            json:"video"`

	Points []dPoint `bson:"points"        json:"-"`
}

type dPoint struct {
	Id    string `bson:"id"    json:"id"`
	Name  string `bson:"name"  json:"name"`
	Video string `bson:"video" json:"video"`
}

type dAssignments struct {
	Id       string `bson:"id"            json:"id"`
	Name     string `bson:"name"          json:"name"`
	Desc     string `bson:"desc"          json:"desc"`
	DeadLine string `bson:"deadline"      json:"deadline"`
}

// Course Player
type DCoursePlayer struct {
	Id        string `bson:"id"         json:"id"`
	CourseId  string `bson:"course_id"  json:"course_id"`
	Name      string `bson:"name"       json:"name"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	Repo      string `bson:"repo"       json:"repo"`
	Version   int    `bson:"version"    json:"-"`
}

type DCourseWork struct {
	Id       string  `bson:"id"         json:"id"`
	CourseId string  `bson:"course_id"  json:"course_id"`
	Account  string  `bson:"account"    json:"account"`
	AsgId    string  `bson:"asg_id"     json:"asg_id"`
	Score    float32 `bson:"score"      json:"score"`
	Status   string  `bson:"status"     json:"status"`
	Version  int     `bson:"version"    json:"-"`
}
