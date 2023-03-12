package repositoryimpl

type Config struct {
	Table Table `json:"table" required:"true"`
}

type Table struct {
	Pod string `json:"pod" reuired:"true"`
}
