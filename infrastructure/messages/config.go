package messages

type Topics struct {
	Like            string `json:"like"             required:"true"`
	Fork            string `json:"fork"             required:"true"`
	Download        string `json:"download"         required:"true"`
	Training        string `json:"training"         required:"true"`
	Finetune        string `json:"finetune"         required:"true"`
	Following       string `json:"following"        required:"true"`
	Inference       string `json:"inference"        required:"true"`
	Submission      string `json:"submission"       required:"true"`
	OperateLog      string `json:"operate_log"      required:"true"`
	RelatedResource string `json:"related_resource" required:"true"`
	Cloud           string `json:"cloud"            required:"true"`
	ReleaseCloud    string `json:"release_cloud"    required:"true"`
	Async           string `json:"async"            required:"true"`
	AICCFinetune    string `json:"aiccfinetune"     required:"true"`
}
