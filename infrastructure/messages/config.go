package messages

var topics Topics

type Topics struct {
	Like               string      `json:"like"             required:"true"`
	Fork               string      `json:"fork"             required:"true"`
	Download           string      `json:"download"         required:"true"`
	Training           string      `json:"training"         required:"true"`
	Finetune           string      `json:"finetune"         required:"true"`
	Following          string      `json:"following"        required:"true"`
	Inference          string      `json:"inference"        required:"true"`
	Evaluate           string      `json:"evaluate"         required:"true"`
	Submission         string      `json:"submission"       required:"true"`
	OperateLog         string      `json:"operate_log"      required:"true"`
	RelatedResource    string      `json:"related_resource" required:"true"`
	Cloud              string      `json:"cloud"            required:"true"`
	Async              string      `json:"async"            required:"true"`
	BigModel           string      `json:"bigmodel"         required:"true"`
	SignIn             topicConfig `json:"signin"`
	DailyLike          topicConfig `json:"daily_like"`
	DailyDownload      topicConfig `json:"daily_download"`
	DailyCreate        topicConfig `json:"daily_create"`
	ForkProject        topicConfig `json:"fork_project "`
	ExperienceBigmodel topicConfig `json:"experience_bigmodel "`
	ApplyCourse        topicConfig `json:"apply_course "`
	Register           topicConfig `json:"register "`
	BindEmail          topicConfig `json:"bind_email "`
	SetAvatarId        topicConfig `json:"set_avatar "`
	SetBio             topicConfig `json:"set_bio"`
}

type topicConfig struct {
	Name  string `json:"name"   required:"true"`
	Topic string `json:"topic"  required:"true"`
}
