package repositoryimpl

// Tables is a struct that represents tables of space.
type Tables struct {
	Project string `json:"project" required:"true"`
	Tags    string `json:"tags"    required:"true"`
	Dataset string `json:"dataset" required:"true"`
	Model   string `json:"model"   required:"true"`
}
