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
	Name    string
	RepoId  string
	Obj     ResourceObject
}

type ResourceLikedEvent struct {
	Account Account
	Obj     ResourceObject
}
