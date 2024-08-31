package domain

type Inference struct {
	InferenceInfo

	// following fields is not under the controlling of version
	InferenceDetail
}

type InferenceInfo struct {
	InferenceIndex

	ProjectName   ResourceName
	ResourceLevel string
	Requester     string
}

type InferenceDetail struct {
	// Expiry stores the time when the inference instance will exit.
	Expiry int64

	// Error stores the message when the reference instance starts failed
	Error string

	// AccessURL stores the url to access the inference service.
	AccessURL string
}

type InferenceIndex struct {
	Project    ResourceIndex
	Id         string
	LastCommit string
}
