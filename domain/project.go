package domain

import "fmt"

type Project struct {
	Id    string
	Owner string

	Name      ProjName
	Desc      ProjDesc
	Type      RepoType
	CoverId   CoverId
	Protocol  ProtocolName
	Training  TrainingSDK
	Inference InferenceSDK

	LikeAccount LikeAccount
	Downloads   ProjectDownloads
}

func (p Project) ValidateID() error {
	if len(p.Id) == 0 {
		return fmt.Errorf("Project id is inValidate")
	}
	return nil
}

type LikeAccount int

func (la LikeAccount) Validate() error {
	if la < 0 {
		return fmt.Errorf("LikeAccount must great than zero")
	}
	return nil
}

type ProjectDownloads map[string]int
