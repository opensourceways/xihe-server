package repositoryimpl

type Config struct {
	Table Table `json:"table" required:"true"`
}

type Table struct {
	WukongRequest string `json:"wukong_request" required:"true"`
}
