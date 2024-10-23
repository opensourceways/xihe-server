package repositoryimpl

type SpaceAppDAO interface {
	UpdateWithOmitingSpecificFields(filter, values any, columns ...string) error
}
