package repositoryimpl

const (
	fieldId        = "id"
	fieldName      = "name"
	fieldDesc      = "desc"
	fieldDetail    = "detail"
	fieldJob       = "job"
	fieldCreatedAt = "created_at"
	fieldUser      = "user"
	fieldModel     = "model"
	fieldVersion   = "version"
	fieldStatus    = "status"
	fieldItems     = "items"
	fieldTask      = "task"
)

type dAICCFinetune struct {
	User    string `bson:"user"    json:"user"`
	Model   string `bson:"model"   json:"model"`
	Version int    `bson:"version" json:"-"`

	Items []aiccFinetuneItem `bson:"items"   json:"-"`
}

type aiccFinetuneItem struct {
	Id              string      `bson:"id"            json:"id"`
	Name            string      `bson:"name"          json:"name"`
	Desc            string      `bson:"desc"          json:"desc"`
	Model           string      `bson:"model"         json:"model"`
	Task            string      `bson:"task"         json:"task"`
	Env             []dKeyValue `bson:"env"           json:"env"`
	Hyperparameters []dKeyValue `bson:"parameters"    json:"parameters"`
	CreatedAt       int64       `bson:"created_at"    json:"created_at"`
	Job             dJobInfo    `bson:"job"           json:"-"`
	JobDetail       dJobDetail  `bson:"detail"        json:"-"`
}

type dKeyValue struct {
	Key   string `bson:"key"             json:"key"`
	Value string `bson:"value"           json:"value"`
}

type dJobInfo struct {
	Endpoint  string `bson:"endpoint"    json:"endpoint"`
	JobId     string `bson:"job_id"      json:"job_id"`
	LogDir    string `bson:"log"         json:"log"`
	OutputDir string `bson:"output"      json:"output"`
}

type dJobDetail struct {
	Duration   int    `bson:"duration"   json:"duration,omitempty"`
	Error      string `bson:"error"      json:"error,omitempty"`
	Status     string `bson:"status"     json:"status,omitempty"`
	LogPath    string `bson:"log"        json:"log,omitempty"`
	OutputPath string `bson:"output"     json:"output,omitempty"`
}
