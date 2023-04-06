package repositoryimpl

const (
	fieldId      = "id"
	fieldOwner   = "owner"
	fieldItems   = "items"
	fieldSamples = "samples"
	fieldNum     = "num"
	fieldVersion = "version"
	fieldLikes   = "likes"
	fieldPublics = "publics"
)

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
