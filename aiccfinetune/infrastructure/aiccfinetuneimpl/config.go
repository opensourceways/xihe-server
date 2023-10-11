package aiccfinetuneimpl

type Config struct {
	Endpoint           string    `json:"endpoint"              required:"true"`
	JobDoneStatus      []string  `json:"job_done_status"       required:"true"`
	CanTerminateStatus []string  `json:"can_terminate_status"  required:"true"`
	OBSConfig          OBSConfig `json:"obs_config"  required:"true"`
}

type OBSConfig struct {
	Prefix    string `json:"prefix"`
	Bucket    string `json:"bucket"         required:"true"`
	Endpoint  string `json:"endpoint"       required:"true"`
	AccessKey string `json:"access_key"     required:"true"`
	SecretKey string `json:"secret_key"     required:"true"`
}
