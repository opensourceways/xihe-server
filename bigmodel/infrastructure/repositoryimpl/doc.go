package repositoryimpl

const (
	fieldId        = "id"
	fieldOwner     = "owner"
	fieldUser      = "user"
	fieldItems     = "items"
	fieldSamples   = "samples"
	fieldNum       = "num"
	fieldVersion   = "version"
	fieldLikes     = "likes"
	fieldPublics   = "publics"
	fieldModelName = "model_name"
	fieldEnabled   = "enabled"
	fieldCallCount = "call_count"
	fieldToken     = "token"
	fieldUpdateAt  = "update_at"
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

type dApiApply struct {
	User      string `bson:"user"        json:"user"`
	ModelName string `bson:"model_name"  json:"model_name"`
	ApplyAt   string `bson:"apply_at"    json:"apply_at"`
	UpdateAt  string `bson:"update_at"   json:"update_at"`
	Token     string `bson:"token"       json:"token"`
	Enabled   bool   `bson:"enabled"     json:"enabled"`
	CallCount int    `bson:"call_count"  json:"call_count"`
	Version   int    `bson:"version"     json:"-"`
}

type dApiInfo struct {
	Id       string `bson:"id"        json:"id"`
	Name     string `bson:"name"      json:"name"`
	Endpoint string `bson:"endpoint"  json:"endpoint"`
	Doc      string `bson:"doc"       json:"doc"`
}
