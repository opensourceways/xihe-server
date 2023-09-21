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

type RepoDownload struct {
	Account Account
	Type    ResourceType
	Name    string
}
