package infrastructure

type fileScanDAO interface {
	GetRecord(filter, result any) error
}
