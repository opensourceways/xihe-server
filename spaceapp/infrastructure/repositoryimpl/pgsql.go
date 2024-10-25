package repositoryimpl

type SpaceAppDAO interface {
	UpdateWithOmittingSpecificFields(filter, values any, columns ...string) error
	GetRecord(filter, result any) error
	GetByPrimaryKey(row any) error
	Update(filter, values any) error
}
