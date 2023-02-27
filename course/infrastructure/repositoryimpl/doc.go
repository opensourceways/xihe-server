package repositoryimpl

const (
	fieldId      = "id"
	fieldAccount = "account"
)

// Course
type DCourse struct {
	Id        string  `bson:"id"              json:"id"`
	Name      string  `bson:"name"            json:"name"`
	Teacher   string  `bson:"teacher"         json:"teacher"`
	Desc      string  `bson:"desc"            json:"desc"`
	Host      string  `bson:"host"            json:"host"`
	Type      string  `bson:"type"            json:"type"`
	PassScore int     `bson:"pass_score"      json:"pass_score"`
	Status    string  `bson:"status"          json:"status"`
	Duration  string  `bson:"duration"        json:"duration"`
	Hours     float32 `bson:"hours"           json:"hours"`
	Doc       string  `bson:"doc"             json:"doc"`
	Forum     string  `bson:"forum"           json:"forum"`
	Poster    string  `bson:"poster"          json:"poster"`
	Cert      string  `bson:"cert"            json:"cert"`

	Assignments []dAssignments `bson:"assignments" json:"-"`
	Sections    []dSection     `bson:"lessons"    json:"-"`
}

type dSection struct {
	Id   string `bson:"cert"            json:"cert"`
	Name string `bson:"name"            json:"name"`

	Lessons []dLesson `bson:"lessons"            json:"-"`
}

type dLesson struct {
	Id    string `bson:"id"               json:"id"`
	Name  string `bson:"name"             json:"name"`
	Desc  string `bson:"desc"             json:"desc"`
	Video string `bson:"video"            json:"video"`
}

type dAssignments struct {
	Id       string `bson:"id"            json:"id"`
	Name     string `bson:"name"            json:"name"`
	Desc     string `bson:"desc"            json:"desc"`
	DeadLine string `bson:"deadline"            json:"deadline"`
}

// Course Player
type DCoursePlayer struct {
	Id        string `json:"id" bson:"id"`
	CourseId  string `json:"course_id" bson:"course_id"`
	Name      string `json:"name" bson:"name"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}
