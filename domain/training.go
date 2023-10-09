package domain

type UserTraining struct {
	Id        string
	Owner     Account
	ProjectId string

	TrainingConfig

	CreatedAt int64

	// following fileds is not under the controlling of version
	Job       JobInfo
	JobDetail JobDetail
}

type TrainingConfig struct {
	ProjectName   ResourceName
	ProjectRepoId string

	Name TrainingName
	Desc TrainingDesc

	CodeDir  Directory
	BootFile FilePath

	Hyperparameters []KeyValue
	Env             []KeyValue
	Inputs          []Input
	EnableOutput    bool
	EnableAim       bool

	Compute Compute
}

type Compute struct {
	Type    ComputeType
	Flavor  ComputeFlavor
	Version ComputeVersion
}

type KeyValue struct {
	Key   CustomizedKey
	Value CustomizedValue
}

type Input struct {
	Key CustomizedKey
	ResourceRef
}

type ResourceRef struct {
	User   Account
	Type   ResourceType
	RepoId string
	File   InputeFilePath
}

type JobInfo struct {
	Endpoint  string
	JobId     string
	LogDir    string
	OutputDir string
	AimDir    string
}

type JobDetail struct {
	Status     string
	Error      string
	LogPath    string
	AimPath    string
	OutputPath string
	Duration   int
}

type TrainingSummary struct {
	Id        string
	Name      TrainingName
	Desc      TrainingDesc
	Error     string
	Status    string
	Duration  int
	CreatedAt int64
}

type TrainingIndex struct {
	Project    ResourceIndex
	TrainingId string
}
