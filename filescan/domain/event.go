package domain

type ModerationEvent struct {
	ID       int64  `json:"id"`
	Owner    string `json:"owner"`
	RepoId   int64  `json:"repo_id"`
	RepoName string `json:"repo_name"`
	Dir      string `json:"dir"`
	File     string `json:"file"`
}
