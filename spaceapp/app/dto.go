package app

type InferenceDTO struct {
	expiry     int64
	Error      string `json:"error"`
	AccessURL  string `json:"access_url"`
	InstanceId string `json:"inference_id"`
}

func (dto *InferenceDTO) hasResult() bool {
	return dto.InstanceId != ""
}

func (dto *InferenceDTO) canReuseCurrent() bool {
	return dto.AccessURL != ""
}
