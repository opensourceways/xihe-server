package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type trainingCreateResp struct {
	Id string `json:"id"`
}

type trainingLogResp struct {
	LogURL string `json:"log_url"`
}

type trainingDetail struct {
	app.TrainingDTO
	Log string `json:"log"`
}

type TrainingCreateRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`

	CodeDir  string `json:"code_dir"`
	BootFile string `json:"boot_file"`

	Hyperparameters []KeyValue    `json:"hyperparameter"`
	Env             []KeyValue    `json:"env"`
	Models          []TrainingRef `json:"models"`
	Datasets        []TrainingRef `json:"datasets"`
	EnableAim       bool          `json:"enable_aim"`
	EnableOutput    bool          `json:"enable_output"`

	Compute Compute `json:"compute"`
}

func (req *TrainingCreateRequest) toCmd(cmd *app.TrainingCreateCmd) (err error) {
	if cmd.Name, err = domain.NewTrainingName(req.Name); err != nil {
		return
	}

	if cmd.Desc, err = domain.NewTrainingDesc(req.Desc); err != nil {
		return
	}

	if cmd.CodeDir, err = domain.NewDirectory(req.CodeDir); err != nil {
		return
	}

	if cmd.BootFile, err = domain.NewFilePath(req.BootFile); err != nil {
		return
	}

	if cmd.Compute, err = req.Compute.toCompute(); err != nil {
		return
	}

	if cmd.Env, err = req.toKeyValue(req.Env); err != nil {
		return
	}

	if cmd.Hyperparameters, err = req.toKeyValue(req.Hyperparameters); err != nil {
		return
	}

	cmd.EnableAim = req.EnableAim
	cmd.EnableOutput = req.EnableOutput

	return
}

func (req *TrainingCreateRequest) toKeyValue(kv []KeyValue) (r []domain.KeyValue, err error) {
	n := len(kv)
	if n == 0 {
		return nil, nil
	}

	r = make([]domain.KeyValue, n)
	for i := range kv {
		if r[i], err = kv[i].toKeyValue(); err != nil {
			return
		}
	}

	return
}

type Compute struct {
	Type    string `json:"type"`
	Flavor  string `json:"flavor"`
	Version string `json:"version"`
}

func (c *Compute) toCompute() (r domain.Compute, err error) {
	if c.Type == "" || c.Version == "" || c.Flavor == "" {
		err = errors.New("invalid compute info")

		return
	}

	if r.Type, err = domain.NewComputeType(c.Type); err != nil {
		return
	}

	if r.Flavor, r.Version, err = domain.NewComputeFlavorVersion(c.Flavor, c.Type, c.Version); err != nil {
		return
	}

	return
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (kv *KeyValue) toKeyValue() (r domain.KeyValue, err error) {
	if kv.Key == "" {
		err = errors.New("invalid key value")

		return
	}

	if r.Key, err = domain.NewCustomizedKey(kv.Key); err != nil {
		return
	}

	r.Value, err = domain.NewCustomizedValue(kv.Value)

	return
}

type TrainingRef struct {
	Key   string `json:"key"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
	File  string `json:"File"`
}

func (t *TrainingRef) toModelInput() (r domain.Input, name domain.ResourceName, err error) {
	if name, err = t.toInput(&r); err != nil {
		return
	}

	r.Type = domain.ResourceTypeModel

	return
}

func (t *TrainingRef) toDatasetInput() (r domain.Input, name domain.ResourceName, err error) {
	if name, err = t.toInput(&r); err != nil {
		return
	}

	r.Type = domain.ResourceTypeDataset

	return
}

func (t *TrainingRef) toInput(r *domain.Input) (name domain.ResourceName, err error) {
	if r.Key, err = domain.NewCustomizedKey(t.Key); err != nil {
		return
	}

	if r.User, err = domain.NewAccount(t.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewResourceName(t.Name); err != nil {
		return
	}
	name = r.Name

	if r.File, err = domain.NewInputeFilePath(t.File); err != nil {
		return
	}

	return
}

func (ctl *TrainingController) setProjectInfo(
	ctx *gin.Context, cmd *app.TrainingCreateCmd,
	user domain.Account, projectId string,
) (ok bool) {
	v, err := ctl.projectPg.GetSummary(user, projectId)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	cmd.User = user
	cmd.ProjectId = projectId
	cmd.ProjectName = v.Name
	cmd.ProjectRepoId = v.RepoId
	ok = true

	return
}

func (ctl *TrainingController) setModelsInput(
	ctx *gin.Context, cmd *app.TrainingCreateCmd, a domain.Account,
	inputs []TrainingRef,
) (ok bool) {
	p := make([]repository.ResourceSummaryListOption, 0, len(inputs))
	m := sets.NewString()
	tinputs := make([]domain.Input, len(inputs))

	index := func(u domain.Account, n string) string {
		return u.Account() + n
	}

	for i := range inputs {
		v := &inputs[i]

		ti, name, err := v.toModelInput()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))
			return
		}
		tinputs[i] = ti

		s := index(ti.User, v.Name)
		if m.Has(s) {
			continue
		}
		m.Insert(s)

		p = append(p, repository.ResourceSummaryListOption{
			Owner: ti.User,
			Name:  name,
		})

	}

	v, err := ctl.model.ListSummary(p)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ids := map[string]string{}
	for i := range v {
		item := &v[i]

		ids[index(item.Owner, item.Name.ResourceName())] = item.RepoId
	}

	for i := range inputs {
		v := &tinputs[i]

		if v.RepoId, ok = ids[index(v.User, inputs[i].Name)]; !ok {
			ctx.JSON(http.StatusBadRequest, respBadRequestParam(
				errors.New("can't find repo id"),
			))

			return
		}
	}

	for i := range v {
		item := &v[i]
		rname := domain.ResourceName(item.Name)

		model, err := ctl.model.GetByName(item.Owner, rname)
		if err != nil {
			return
		}

		if model.IsPrivate() && item.Owner.Account() != a.Account() {
			ok = false
			ctx.JSON(http.StatusBadRequest, respBadRequestParam(
				errors.New("invalid model"),
			))

			return
		}
	}

	cmd.Inputs = append(cmd.Inputs, tinputs...)
	ok = true

	return
}

func (ctl *TrainingController) setDatasetsInput(
	ctx *gin.Context, cmd *app.TrainingCreateCmd, a domain.Account,
	inputs []TrainingRef,
) (ok bool) {
	p := make([]repository.ResourceSummaryListOption, 0, len(inputs))
	m := sets.NewString()
	tinputs := make([]domain.Input, len(inputs))

	index := func(u domain.Account, n string) string {
		return u.Account() + n
	}

	for i := range inputs {
		v := &inputs[i]

		ti, name, err := v.toDatasetInput()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))
			return
		}
		tinputs[i] = ti

		s := index(ti.User, v.Name)
		if m.Has(s) {
			continue
		}
		m.Insert(s)

		p = append(p, repository.ResourceSummaryListOption{
			Owner: ti.User,
			Name:  name,
		})

	}

	v, err := ctl.dataset.ListSummary(p)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ids := map[string]string{}
	for i := range v {
		item := &v[i]

		ids[index(item.Owner, item.Name.ResourceName())] = item.RepoId
	}

	for i := range inputs {
		v := &tinputs[i]

		if v.RepoId, ok = ids[index(v.User, inputs[i].Name)]; !ok {
			ctx.JSON(http.StatusBadRequest, respBadRequestParam(
				errors.New("can't find repo id"),
			))

			return
		}
	}

	for i := range v {
		item := &v[i]
		rname := domain.ResourceName(item.Name)

		dataset, err := ctl.dataset.GetByName(item.Owner, rname)
		if err != nil {
			return
		}

		if dataset.IsPrivate() && item.Owner.Account() != a.Account() {
			ok = false
			ctx.JSON(http.StatusBadRequest, respBadRequestParam(
				errors.New("invalid dataset"),
			))

			return
		}
	}

	cmd.Inputs = append(cmd.Inputs, tinputs...)
	ok = true

	return
}
