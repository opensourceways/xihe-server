// Package repositoryadapter provides an adapter for working with computility repositories.
package repositoryimpl

// Tables is a struct that represents tables of computility.
type Tables struct {
	Project string `json:"project" required:"true"`
	Tags    string `json:"tags" required:"true"`
}
