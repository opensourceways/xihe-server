package infrastructure

type largeFileScanDAO interface {
	GetRecord(filter, result any) error
}
