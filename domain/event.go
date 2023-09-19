package domain

// training
type TrainingCreatedEvent struct {
	Account        Account
	TrainingIndex  TrainingIndex
	TrainingInputs []Input
}
