package message

type Topics struct {
	CloudCreate topicConfig `json:"cloud_create"`
}

type topicConfig struct {
	Name  string `json:"name"   required:"true"`
	Topic string `json:"topic"  required:"true"`
}
