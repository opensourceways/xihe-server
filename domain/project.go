package domain

type Project struct {
	NewOne bool
	UserId string

	Name      ProjName
	Desc      ProjDesc
	Type      RepoType
	CoverId   string
	Protocol  ProtocolName
	Training  TrainingSDK
	Inference InferenceSDK
}

func (p *Project) Validate() error {
	return nil
}
