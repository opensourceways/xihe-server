package controller

type projectModel struct {
	Name      string `json:"name" required:"true"`
	Desc      string `json:"desc" required:"true"`
	Type      string `json:"type" required:"true"`
	CoverId   string `json:"cover_id" required:"true"`
	Protocol  string `json:"protocol" required:"true"`
	Training  string `json:"training" required:"true"`
	Inference string `json:"inference" required:"true"`
}
