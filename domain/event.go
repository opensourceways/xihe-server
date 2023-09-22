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

type RepoDownloadedEvent struct {
	Account Account
	Type    ResourceType
	Name    string
}

type ResourceLikedEvent struct {
	Account Account
	Obj     ResourceObject
}
