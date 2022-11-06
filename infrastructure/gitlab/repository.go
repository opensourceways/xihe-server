package gitlab

import (
	"strconv"

	sdk "github.com/xanzy/go-gitlab"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
)

type UserInfo struct {
	Token     string
	Namespace string
}

func NewRepositoryService(v UserInfo) platform.Repository {
	return &repository{user: v}
}

type repository struct {
	user UserInfo
}

// TODO admin create repo instead
func (r *repository) New(repo *platform.RepoOption) (string, error) {
	cli, err := sdk.NewClient(r.user.Token, sdk.WithBaseURL(endpoint))
	if err != nil {
		return "", err
	}

	ns, err := strconv.Atoi(r.user.Namespace)
	if err != nil {
		return "", err
	}

	var visibility sdk.VisibilityValue
	switch repo.RepoType.RepoType() {
	case domain.RepoTypePrivate:
		visibility = sdk.PrivateVisibility

	default:
		visibility = sdk.PublicVisibility
	}

	name := repo.Name.ResourceName()
	b := true

	v, _, err := cli.Projects.CreateProject(&sdk.CreateProjectOptions{
		Name:                 &name,
		Path:                 &name,
		NamespaceID:          &ns,
		Visibility:           &visibility,
		InitializeWithReadme: &b,
	})

	if err != nil {
		return "", err
	}

	return strconv.Itoa(v.ID), nil
}

// TODO admin fork repo instead
func (r *repository) Fork(srcRepoId string, repoName domain.ResourceName) (string, error) {
	cli, err := sdk.NewClient(r.user.Token, sdk.WithBaseURL(endpoint))
	if err != nil {
		return "", err
	}

	ns, err := strconv.Atoi(r.user.Namespace)
	if err != nil {
		return "", err
	}

	repoId, err := strconv.Atoi(srcRepoId)
	if err != nil {
		return "", err
	}

	name := repoName.ResourceName()
	b := true

	v, _, err := cli.Projects.ForkProject(repoId, &sdk.ForkProjectOptions{
		Name:                          &name,
		Path:                          &name,
		NamespaceID:                   &ns,
		MergeRequestDefaultTargetSelf: &b,
	})

	if err != nil {
		return "", err
	}

	return strconv.Itoa(v.ID), nil
}

// TODO admin update repo instead
func (r *repository) Update(repoId string, repo *platform.RepoOption) error {
	cli, err := sdk.NewClient(r.user.Token, sdk.WithBaseURL(endpoint))
	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(repoId)
	if err != nil {
		return err
	}

	opts := &sdk.EditProjectOptions{}

	if repo.Name != nil {
		n := repo.Name.ResourceName()
		opts.Name = &n
		opts.Path = &n
	}

	if repo.RepoType != nil {
		var v sdk.VisibilityValue

		switch repo.RepoType.RepoType() {
		case domain.RepoTypePrivate:
			v = sdk.PrivateVisibility

		default:
			v = sdk.PublicVisibility
		}

		opts.Visibility = &v
	}

	_, _, err = cli.Projects.EditProject(pid, opts)

	return err
}
