package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	fieldId             = "id"
	fieldPId            = "pid"
	fieldBio            = "bio"
	fieldJob            = "job"
	fieldDesc           = "desc"
	fieldDate           = "date"
	fieldLevel          = "level"
	fieldDetail         = "detail"
	fieldScore          = "score"
	fieldCoverId        = "cover_id"
	fieldCommit         = "commit"
	fieldExpiry         = "expiry"
	fieldRepoId         = "repo_id"
	fieldTags           = "tags"
	fieldKinds          = "kinds"
	fieldStatus         = "status"
	fieldName           = "name"
	fieldItems          = "items"
	fieldLikes          = "likes"
	fieldPublics        = "publics"
	fieldOwner          = "owner"
	fieldEmail          = "email"
	fieldAccount        = "account"
	fieldVersion        = "version"
	fieldRepoType       = "repo_type"
	fieldAvatarId       = "avatar_id"
	fieldFollower       = "follower"
	fieldFollowing      = "following"
	fieldLikeCount      = "like_count"
	fieldForkCount      = "fork_count"
	fieldIsFollower     = "is_follower"
	fieldFollowerCount  = "follower_count"
	fieldFollowingCount = "following_count"
	fieldCreatedAt      = "created_at"
	fieldType           = "type"
	fieldModels         = "models"
	fieldDatasets       = "datasets"
	fieldProjects       = "projects"
	fieldRId            = "rid"
	fieldTId            = "tid"
	fieldROwner         = "rowner"
	fieldRType          = "rtype"
	fieldUpdatedAt      = "updated_at"
	fieldDownloadCount  = "download_count"
	fieldFirstLetter    = "fl"
	fieldPhase          = "phase"
	fieldTeams          = "teams"
	fieldRepos          = "repos"
	fieldOrder          = "order"
	fieldEnabled        = "enabled"
	fieldCompetitors    = "competitors"
	fieldSubmissions    = "submissions"
	fieldNum            = "num"
	fieldSamples        = "samples"
	fieldPictures       = "pictures"
	fieldChoices        = "choices"
	fieldCompletions    = "completions"
)

type dProject struct {
	Owner string        `bson:"owner" json:"owner"`
	Items []projectItem `bson:"items" json:"-"`
}

type projectItem struct {
	Id        string `bson:"id"         json:"id"`
	Type      string `bson:"type"       json:"type"`
	Protocol  string `bson:"protocol"   json:"protocol"`
	Training  string `bson:"training"   json:"training"`
	RepoId    string `bson:"repo_id"    json:"repo_id"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	UpdatedAt int64  `bson:"updated_at" json:"updated_at"`

	ProjectPropertyItem `bson:",inline"`

	// These two items are not allowd to be set,
	// So, don't marshal it to avoid setting it occasionally.
	RelatedModels   []ResourceIndex `bson:"models"   json:"-"`
	RelatedDatasets []ResourceIndex `bson:"datasets" json:"-"`

	// Version, LikeCount and ForkCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version       int `bson:"version"           json:"-"`
	LikeCount     int `bson:"like_count"        json:"-"`
	ForkCount     int `bson:"fork_count"        json:"-"`
	DownloadCount int `bson:"download_count"    json:"-"`
}

type ProjectPropertyItem struct {
	Level    int    `bson:"level"      json:"level"`
	Name     string `bson:"name"       json:"name"`
	FL       byte   `bson:"fl"         json:"fl"`
	Desc     string `bson:"desc"       json:"desc"`
	Title    string `bson:"title"       json:"title"`
	CoverId  string `bson:"cover_id"   json:"cover_id"`
	RepoType string `bson:"repo_type"  json:"repo_type"`
	// set omitempty to avoid set it to null occasionally.
	// don't depend on the magic to guarantee the correctness.
	Tags     []string `bson:"tags"     json:"tags"`
	TagKinds []string `bson:"kinds"    json:"kinds"`
}

func (doc *ProjectPropertyItem) setDefault() {
	// The serach by the tag want the tags exist.
	if doc.Tags == nil {
		doc.Tags = []string{}
	}

	if doc.TagKinds == nil {
		doc.TagKinds = []string{}
	}
}

type dModel struct {
	Owner string      `bson:"owner" json:"owner"`
	Items []modelItem `bson:"items" json:"-"`
}

type modelItem struct {
	Id        string `bson:"id"         json:"id"`
	RepoId    string `bson:"repo_id"    json:"repo_id"`
	Protocol  string `bson:"protocol"   json:"protocol"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	UpdatedAt int64  `bson:"updated_at" json:"updated_at"`

	ModelPropertyItem `bson:",inline"`

	// RelatedDatasets is not allowd to be set,
	// So, don't marshal it to avoid setting it occasionally.
	RelatedDatasets []ResourceIndex `bson:"datasets" json:"-"`
	RelatedProjects []ResourceIndex `bson:"projects" json:"-"`

	// Version, LikeCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version       int `bson:"version"           json:"-"`
	LikeCount     int `bson:"like_count"        json:"-"`
	DownloadCount int `bson:"download_count"    json:"-"`
}

type ModelPropertyItem struct {
	FL       byte     `bson:"fl"         json:"fl"`
	Level    int      `bson:"level"      json:"level"`
	Name     string   `bson:"name"       json:"name"`
	Desc     string   `bson:"desc"       json:"desc"`
	Title    string   `bson:"title"      json:"title"`
	RepoType string   `bson:"repo_type"  json:"repo_type"`
	Tags     []string `bson:"tags"       json:"tags"`
	TagKinds []string `bson:"kinds"      json:"kinds"`
}

func (doc *ModelPropertyItem) setDefault() {
	// The serach by the tag want the tags exist.
	if doc.Tags == nil {
		doc.Tags = []string{}
	}

	if doc.TagKinds == nil {
		doc.TagKinds = []string{}
	}
}

type dDataset struct {
	Owner string        `bson:"owner" json:"owner"`
	Items []datasetItem `bson:"items" json:"-"`
}

type datasetItem struct {
	Id        string `bson:"id"         json:"id"`
	RepoId    string `bson:"repo_id"    json:"repo_id"`
	Protocol  string `bson:"protocol"   json:"protocol"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	UpdatedAt int64  `bson:"updated_at" json:"updated_at"`

	DatasetPropertyItem `bson:",inline"`

	RelatedModels   []ResourceIndex `bson:"models"   json:"-"`
	RelatedProjects []ResourceIndex `bson:"projects" json:"-"`

	// Version, LikeCount will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version       int `bson:"version"               json:"-"`
	LikeCount     int `bson:"like_count"            json:"-"`
	DownloadCount int `bson:"download_count"        json:"-"`
}

type DatasetPropertyItem struct {
	FL       byte     `bson:"fl"         json:"fl"`
	Level    int      `bson:"level"      json:"level"`
	Name     string   `bson:"name"       json:"name"`
	Desc     string   `bson:"desc"       json:"desc"`
	RepoType string   `bson:"repo_type"  json:"repo_type"`
	Tags     []string `bson:"tags"       json:"tags"`
	TagKinds []string `bson:"kinds"      json:"kinds"`
}

func (doc *DatasetPropertyItem) setDefault() {
	// The serach by the tag want the tags exist.
	if doc.Tags == nil {
		doc.Tags = []string{}
	}

	if doc.TagKinds == nil {
		doc.TagKinds = []string{}
	}
}

type DUser struct {
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Name                    string `bson:"name"       json:"name"`
	Email                   string `bson:"email"      json:"email"`
	Bio                     string `bson:"bio"        json:"bio"`
	AvatarId                string `bson:"avatar_id"  json:"avatar_id"`
	PlatformToken           string `bson:"token"      json:"token"`
	PlatformUserId          string `bson:"uid"        json:"uid"`
	PlatformUserNamespaceId string `bson:"nid"        json:"nid"`

	Follower  []string `bson:"follower"   json:"-"`
	Following []string `bson:"following"  json:"-"`

	// Version will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version int `bson:"version"    json:"-"`
}

type dLogin struct {
	Account     string `bson:"account"   json:"account"`
	Info        string `bson:"info"      json:"info"`
	AccessToken string `bson:"access"    json:"access"`
	UserId      string `bson:"user_id"   json:"user_id"`
}

type dLike struct {
	Owner string     `bson:"owner" json:"owner"`
	Items []likeItem `bson:"items" json:"-"`
}

type likeItem struct {
	CreatedAt int64 `bson:"created_at" json:"created_at"`

	ResourceObject `bson:",inline"`
}

type dActivity struct {
	Owner string         `bson:"owner" json:"owner"`
	Items []activityItem `bson:"items" json:"-"`
}

type activityItem struct {
	Type string `bson:"type" json:"type"`
	Time int64  `bson:"time" json:"time"`

	ResourceObject `bson:",inline"`
}

type ResourceObject struct {
	Id    string `bson:"rid"     json:"rid"`
	Type  string `bson:"rtype"   json:"rtype"`
	Owner string `bson:"rowner"  json:"rowner"`
}

type ResourceIndex struct {
	Id    string `bson:"rid"     json:"rid"`
	Owner string `bson:"rowner"  json:"rowner"`
}

type dResourceTags struct {
	Items []dDomainTags `bson:"items"    json:"items"`
}

type dDomainTags struct {
	Name   string  `bson:"name"          json:"name"`
	Domain string  `bson:"domain"        json:"domain"`
	Tags   []dTags `bson:"tags"          json:"tags"`
}

type dTags struct {
	Kind string   `bson:"kind"           json:"kind"`
	Tags []string `bson:"tags"           json:"tags"`
}

type dTraining struct {
	Owner         string `bson:"owner"   json:"owner"`
	ProjectId     string `bson:"pid"     json:"pid"`
	ProjectName   string `bson:"name"    json:"name"`
	ProjectRepoId string `bson:"rid"     json:"rid"`
	Version       int    `bson:"version" json:"-"`

	Items []trainingItem `bson:"items"   json:"-"`
}

type trainingItem struct {
	Id              string      `bson:"id"            json:"id"`
	Name            string      `bson:"name"          json:"name"`
	Desc            string      `bson:"desc"          json:"desc"`
	CodeDir         string      `bson:"code_dir"      json:"code_dir"`
	BootFile        string      `bson:"boot_file"     json:"boot_file"`
	Compute         dCompute    `bson:"compute"       json:"compute"`
	Inputs          []dInput    `bson:"inputs"        json:"inputs"`
	EnableAim       bool        `bson:"aim"           json:"aim"`
	EnableOutput    bool        `bson:"output"        json:"output"`
	Env             []dKeyValue `bson:"env"           json:"env"`
	Hyperparameters []dKeyValue `bson:"parameters"    json:"parameters"`
	CreatedAt       int64       `bson:"created_at"    json:"created_at"`
	Job             dJobInfo    `bson:"job"           json:"-"`
	JobDetail       dJobDetail  `bson:"detail"        json:"-"`
}

type dCompute struct {
	Type    string `bson:"type"          json:"type"`
	Flavor  string `bson:"flavor"        json:"flavor"`
	Version string `bson:"version"       json:"version"`
}

type dKeyValue struct {
	Key   string `bson:"key"             json:"key"`
	Value string `bson:"value"           json:"value"`
}

type dInput struct {
	Key    string `bson:"key"            json:"key"`
	Type   string `bson:"type"           json:"type"`
	User   string `bson:"user"           json:"user"`
	File   string `bson:"file"           json:"file"`
	RepoId string `bson:"rid"            json:"rid"`
}

type dJobInfo struct {
	Endpoint  string `bson:"endpoint"    json:"endpoint"`
	JobId     string `bson:"job_id"      json:"job_id"`
	LogDir    string `bson:"log"         json:"log"`
	AimDir    string `bson:"aim"         json:"aim"`
	OutputDir string `bson:"output"      json:"output"`
}

type dJobDetail struct {
	Duration   int    `bson:"duration"   json:"duration,omitempty"`
	Error      string `bson:"error"      json:"error,omitempty"`
	Status     string `bson:"status"     json:"status,omitempty"`
	LogPath    string `bson:"log"        json:"log,omitempty"`
	AimPath    string `bson:"aim"        json:"aim,omitempty"`
	OutputPath string `bson:"output"     json:"output,omitempty"`
}

type dInference struct {
	Owner       string `bson:"owner"   json:"owner"`
	ProjectId   string `bson:"pid"     json:"pid"`
	ProjectName string `bson:"name"    json:"name"`
	LastCommit  string `bson:"commit"  json:"commit"`
	Version     int    `bson:"version" json:"-"`

	Items []inferenceItem `bson:"items"  json:"-"`
}

type inferenceItem struct {
	Id        string `bson:"id"          json:"id,omitempty"`
	Expiry    int64  `bson:"expiry"      json:"expiry,omitempty"`
	Error     string `bson:"error"       json:"error,omitempty"`
	AccessURL string `bson:"url"         json:"url,omitempty"`
}

type dEvaluate struct {
	Owner      string `bson:"owner"       json:"owner"`
	ProjectId  string `bson:"pid"         json:"pid"`
	TrainingId string `bson:"tid"         json:"tid"`
	Version    int    `bson:"version"     json:"-"`

	Items []evaluateItem `bson:"items"  json:"-"`
}

type evaluateItem struct {
	// must set omitempty, otherwise it will be override to empty
	Id                string   `bson:"id"          json:"id,omitempty"`
	Type              string   `bson:"type"        json:"type,omitempty"`
	MomentumScope     []string `bson:"momentum"    json:"momentum,omitempty"`
	BatchSizeScope    []string `bson:"bsize"       json:"bsize,omitempty"`
	LearningRateScope []string `bson:"rate"        json:"rate,omitempty"`
	Expiry            int64    `bson:"expiry"      json:"expiry,omitempty"`
	Error             string   `bson:"error"       json:"error,omitempty"`
	AccessURL         string   `bson:"url"         json:"url,omitempty"`
}

type DCompetition struct {
	Id         string `bson:"id"              json:"id"`
	Name       string `bson:"name"            json:"name"`
	Desc       string `bson:"desc"            json:"desc"`
	Host       string `bson:"host"            json:"host"`
	Type       string `bson:"type"            json:"type"`
	Phase      string `bson:"phase"           json:"phase"`
	Bonus      int    `bson:"bonus"           json:"bonus"`
	Status     string `bson:"status"          json:"status"`
	Duration   string `bson:"duration"        json:"duration"`
	Doc        string `bson:"doc"             json:"doc"`
	Forum      string `bson:"forum"           json:"forum"`
	Poster     string `bson:"poster"          json:"poster"`
	Winners    string `bson:"winners"         json:"winners"`
	DatasetDoc string `bson:"dataset_doc"     json:"dataset_doc"`
	DatasetURL string `bson:"dataset_url"     json:"dataset_url"`
	SmallerOk  bool   `bson:"order"           json:"order"`
	Enabled    bool   `bson:"enabled"         json:"enabled"`

	Teams       []dTeam            `bson:"teams"       json:"-"`
	Repos       []dCompetitionRepo `bson:"repos"       json:"-"`
	Competitors []dCompetitor      `bson:"competitors" json:"-"`
	Submissions []dSubmission      `bson:"submissions" json:"-"`
}

type DCompetitorInfo struct {
	Name     string            `bson:"name"      json:"name,omitempty"`
	City     string            `bson:"city"      json:"city,omitempty"`
	Email    string            `bson:"email"     json:"email,omitempty"`
	Phone    string            `bson:"phone"     json:"phone,omitempty"`
	Account  string            `bson:"account"   json:"account,omitempty"`
	Identity string            `bson:"identity"  json:"identity,omitempty"`
	Province string            `bson:"province"  json:"province,omitempty"`
	Detail   map[string]string `bson:"detail"    json:"detail,omitempty"`
}

type dCompetitor struct {
	DCompetitorInfo `bson:",inline"`

	TeamId   string `bson:"tid"       json:"tid,omitempty"`
	TeamRole string `bson:"role"      json:"role,omitempty"`
}

type dTeam struct {
	Id   string `bson:"id"      json:"id"`
	Name string `bson:"name"    json:"name"`
}

type dSubmission struct {
	Id         string  `bson:"id"          json:"id"`
	TeamId     string  `bson:"tid"         json:"tid"`     // if it is submitted by team, set it.
	Individual string  `bson:"account"     json:"account"` // if it is submitted by individual, set it.
	Status     string  `bson:"status"      json:"status"`
	OBSPath    string  `bson:"path"        json:"path"`
	SubmitAt   int64   `bson:"submit_at"   json:"submit_at"`
	Score      float64 `bson:"score"       json:"score"`
	Date       string  `bson:"date"        json:"date"`
}

type dCompetitionRepo struct {
	TeamId     string `bson:"tid"     json:"tid,omitempty"`     // if it is submitted by team, set it.
	Individual string `bson:"account" json:"account,omitempty"` // if it is submitted by individual, set it.
	Owner      string `bson:"owner"   json:"owner,omitempty"`
	Repo       string `bson:"repo"    json:"repo,omitempty"`
}

type dAIQuestion struct {
	Competitors []DCompetitorInfo     `bson:"competitors"   json:"-"`
	Submissions []dQuestionSubmission `bson:"submissions"   json:"-"`
}

type dQuestionSubmission struct {
	Id      string `bson:"id"          json:"id,omitempty"`
	Date    string `bson:"date"        json:"date,omitempty"`
	Status  string `bson:"status"      json:"status,omitempty"`
	Account string `bson:"account"     json:"account,omitempty"`
	Expiry  int64  `bson:"expiry"      json:"expiry,omitempty"`
	Score   int    `bson:"score"       json:"score,omitempty"`
	Times   int    `bson:"times"       json:"times,omitempty"`
	Version int    `bson:"version"     json:"-"`
}

type dQuestionPool struct {
	Choices     []dChoiceQuestion     `bson:"choices"       json:"choices"`
	Completions []dCompletionQuestion `bson:"completions"   json:"completions"`
}

type dChoiceQuestion struct {
	Num     int      `bson:"num"       json:"num"`
	Desc    string   `bson:"desc"      json:"desc"`
	Answer  string   `bson:"answer"    json:"answer"`
	Options []string `bson:"options"   json:"options"`
}

type dCompletionQuestion struct {
	Num    int    `bson:"num"          json:"num"`
	Desc   string `bson:"desc"         json:"desc"`
	Info   string `bson:"info"         json:"info"`
	Answer string `bson:"answer"       json:"answer"`
}

type dLuoJia struct {
	Owner string       `bson:"owner"   json:"owner"`
	Items []luojiaItem `bson:"items"   json:"-"`
}

type luojiaItem struct {
	Id        string `bson:"id"         json:"id"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
}

type dWuKong struct {
	Id      string    `bson:"id"      json:"id"`
	Samples []dSample `bson:"samples" json:"samples"`
}

type dSample struct {
	Num  int    `bson:"num"  json:"num"`
	Name string `bson:"name" json:"name"`
}

type dWuKongPicture struct {
	Owner   string        `bson:"owner"   json:"owner"`
	Version int           `bson:"version" json:"-"`
	Likes   []pictureItem `bson:"likes"   json:"-"` // like picture
	Publics []pictureItem `bson:"publics" json:"-"` // public picture
}

type pictureItem struct {
	Id        string   `bson:"id"         json:"id"`
	Owner     string   `bson:"owner"      json:"owner"`
	Desc      string   `bson:"desc"       json:"desc"`
	Style     string   `bson:"style"      json:"style"`
	OBSPath   string   `bson:"obspath"    json:"obspath"`
	Level     int      `bson:"level"      json:"level"`
	Diggs     []string `bson:"diggs"      json:"diggs"`
	DiggCount int      `bson:"digg_count" json:"digg_count"`
	Version   int      `bson:"version"    json:"-"`
	CreatedAt string   `bson:"created_at" json:"created_at"`
}

type dFinetune struct {
	Owner   string `bson:"owner"         json:"owner"`
	Expiry  int64  `bson:"expiry"        json:"expiry"`
	Version int    `bson:"version"       json:"-"`

	Items []finetuneItem `bson:"items"   json:"-"`
}

type finetuneItem struct {
	Id              string             `bson:"id"            json:"id"`
	Name            string             `bson:"name"          json:"name"`
	Task            string             `bson:"task"          json:"task"`
	Model           string             `bson:"model"         json:"model"`
	CreatedAt       int64              `bson:"created_at"    json:"created_at"`
	Hyperparameters map[string]string  `bson:"parameters"    json:"parameters,omitempty"`
	Job             dFinetuneJobInfo   `bson:"job"           json:"-"`
	JobDetail       dFinetuneJobDetail `bson:"detail"        json:"-"`
}

type dFinetuneJobInfo struct {
	Endpoint string `bson:"endpoint"    json:"endpoint"`
	JobId    string `bson:"job_id"      json:"job_id"`
}

type dFinetuneJobDetail struct {
	Duration int    `bson:"duration"   json:"duration,omitempty"`
	Error    string `bson:"error"      json:"error,omitempty"`
	Status   string `bson:"status"     json:"status,omitempty"`
}
