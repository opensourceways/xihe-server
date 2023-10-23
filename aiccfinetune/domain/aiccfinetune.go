package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type AICCFinetune struct {
	Id    string
	User  types.Account
	Model ModelName
	Task  FinetuneTask

	AICCFinetuneConfig

	CreatedAt int64

	// following fileds is not under the controlling of version
	Job       JobInfo
	JobDetail JobDetail
}

type AICCFinetuneConfig struct {
	Name FinetuneName
	Desc FinetuneDesc

	Hyperparameters []KeyValue
	Env             []KeyValue
}

type KeyValue struct {
	Key   CustomizedKey
	Value CustomizedValue
}

type JobInfo struct {
	Endpoint  string
	JobId     string
	LogDir    string
	OutputDir string
}

type JobDetail struct {
	Status     string
	Error      string
	LogPath    string
	OutputPath string
	Duration   int
}

type AICCFinetuneSummary struct {
	Id        string
	Name      FinetuneName
	Desc      FinetuneDesc
	Error     string
	Status    string
	Duration  int
	CreatedAt int64
	Task      string
}

type AICCFinetuneIndex struct {
	User       types.Account
	Model      ModelName
	FinetuneId string
}
