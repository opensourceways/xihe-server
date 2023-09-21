package domain

// training
type TrainingCreatedEvent struct {
	Account        Account
	TrainingIndex  TrainingIndex
	TrainingInputs []Input
}

type UserSignedInEvent struct {
	Account Account
}

type RepoDownloadEvent struct {
	Account Account
	Type    ResourceType
	Name    string
}
