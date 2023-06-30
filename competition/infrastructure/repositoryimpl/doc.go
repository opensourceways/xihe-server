package repositoryimpl

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	fieldId          = "id"
	fieldCid         = "cid"
	fieldPid         = "pid"
	fieldRepo        = "repo"
	fieldVersion     = "version"
	fieldFinal       = "final"
	fieldPreliminary = "preliminary"
	fieldEnabled     = "enabled"
	fieldTeamName    = "team_name"
	fieldUpdatedAt   = "updated_at"
	fieldAccount     = "account"
	fieldCompetitors = "competitors"
	fieldLeader      = "leader"
	fieldStatus      = "status"
	fieldTags        = "tags"
)

type dCompetition struct {
	Id         string   `bson:"id"              json:"id"`
	Name       string   `bson:"name"            json:"name"`
	Desc       string   `bson:"desc"            json:"desc"`
	Host       string   `bson:"host"            json:"host"`
	Type       string   `bson:"type"            json:"type"`
	Tags       []string `bson:"tags"            json:"tags"`
	Phase      string   `bson:"phase"           json:"phase"`
	Status     string   `bson:"status"          json:"status"`
	Duration   string   `bson:"duration"        json:"duration"`
	Doc        string   `bson:"doc"             json:"doc"`
	Forum      string   `bson:"forum"           json:"forum"`
	Poster     string   `bson:"poster"          json:"poster"`
	Winners    string   `bson:"winners"         json:"winners"`
	DatasetDoc string   `bson:"dataset_doc"     json:"dataset_doc"`
	DatasetURL string   `bson:"dataset_url"     json:"dataset_url"`
	Bonus      int      `bson:"bonus"           json:"bonus"`
	SmallerOk  bool     `bson:"order"           json:"order"`
}

type dWork struct {
	CompetitionId string        `bson:"cid"            json:"cid"`
	PlayerId      string        `bson:"pid"            json:"pid"`
	PlayerName    string        `bson:"pname"          json:"pname"`
	Repo          string        `bson:"repo"           json:"repo"`
	Final         []dSubmission `bson:"final"          json:"final"`
	Preliminary   []dSubmission `bson:"preliminary"    json:"preliminary"`
	Version       int           `bson:"version"        json:"-"`
}

type dSubmission struct {
	Id       string  `bson:"id"          json:"id"`
	Status   string  `bson:"status"      json:"status"`
	OBSPath  string  `bson:"path"        json:"path"`
	SubmitAt int64   `bson:"submit_at"   json:"submit_at"`
	Score    float64 `bson:"score"       json:"score"`
}

// dPlayer
// Leader stands for the user who is an individual competitor or a team leader
// Enabled: it will be set to false when the competitor become a team.
// Competitors: the first one is the leader
type dPlayer struct {
	Id            primitive.ObjectID `bson:"_id"            json:"-"`
	CompetitionId string             `bson:"cid"            json:"cid"`
	Leader        string             `bson:"leader"         json:"leader"`
	TeamName      string             `bson:"team_name"      json:"team_name"`
	Competitors   []dCompetitor      `bson:"competitors"    json:"competitors"`
	IsFinalist    bool               `bson:"is_finalist"    json:"is_finalist"`
	Enabled       bool               `bson:"enabled"        json:"enabled"`
	Version       int                `bson:"version"        json:"-"`
}

type dCompetitor struct {
	Name     string            `bson:"name"      json:"name,omitempty"`
	City     string            `bson:"city"      json:"city,omitempty"`
	Email    string            `bson:"email"     json:"email,omitempty"`
	Phone    string            `bson:"phone"     json:"phone,omitempty"`
	Account  string            `bson:"account"   json:"account,omitempty"`
	Identity string            `bson:"identity"  json:"identity,omitempty"`
	Province string            `bson:"province"  json:"province,omitempty"`
	Detail   map[string]string `bson:"detail"    json:"detail,omitempty"`
}
