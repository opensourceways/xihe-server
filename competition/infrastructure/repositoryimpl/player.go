package repositoryimpl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

func NewPlayerRepo(m mongodbClient) repository.Player {
	return playerRepoImpl{m}
}

type playerRepoImpl struct {
	cli mongodbClient
}

func (impl playerRepoImpl) docFilter(cid string) bson.M {
	return bson.M{
		fieldCid:     cid,
		fieldEnabled: true,
	}
}

func (impl playerRepoImpl) playerFilter(p *domain.Player) (bson.M, error) {
	v, err := impl.cli.ObjectIdFilter(p.Id)
	if err != nil {
		return nil, err
	}

	filter := impl.docFilter(p.CompetitionId)
	for k := range v {
		filter[k] = v[k]
	}

	return filter, nil
}

func (impl playerRepoImpl) docFilterByUser(cid string, a types.Account) bson.M {
	filter := impl.docFilter(cid)

	impl.cli.AppendElemMatchToFilter(
		fieldCompetitors, true,
		bson.M{fieldAccount: a.Account()}, filter,
	)

	return filter
}

func (impl playerRepoImpl) disabledPlayerFilter(cid string, a types.Account) bson.M {
	return bson.M{
		fieldCid:     cid,
		fieldEnabled: false,
		fieldLeader:  a.Account(),
	}
}

// AddPlayer
func (impl playerRepoImpl) AddPlayer(p *domain.Player) error {
	return impl.insertPlayer(p)
}

func (repo playerRepoImpl) genPlayerDoc(p *domain.Player) (bson.M, error) {
	cs := make([]dCompetitor, p.CompetitorsCount())
	for i := range p.Members() {
		cs[i+1] = toCompetitorDoc(&p.Members()[i])
	}
	cs[0] = toCompetitorDoc(&p.Leader)

	obj := dPlayer{
		CompetitionId: p.CompetitionId,
		Competitors:   cs,
		Leader:        p.Leader.Account.Account(),
		Enabled:       true,
	}
	if p.IsATeam() {
		obj.TeamName = p.Team.Name.TeamName()
	}

	doc, err := genDoc(&obj)

	return doc, err
}

func (impl playerRepoImpl) insertPlayer(p *domain.Player) error {
	doc, err := impl.genPlayerDoc(p)
	if err != nil {
		return err
	}
	doc[fieldVersion] = 0

	f := func(ctx context.Context) error {
		filter := impl.docFilterByUser(p.CompetitionId, p.Leader.Account)

		if p.IsATeam() {
			v := impl.docFilter(p.CompetitionId)
			v[fieldTeamName] = p.Team.Name.TeamName()

			filter = bson.M{"$or": bson.A{filter, v}}
		}

		_, err := impl.cli.NewDocIfNotExist(ctx, filter, doc)

		return err
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}
	}

	return err
}

func (impl playerRepoImpl) DeletePlayer(p *domain.Player, version int) error {
	filter, err := impl.playerFilter(p)
	if err != nil {
		return err
	}

	return impl.update(filter, bson.M{fieldEnabled: false}, version)
}

func (impl playerRepoImpl) update(filter, doc bson.M, version int) error {
	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(ctx, filter, doc, mongoCmdSet, version)
	}

	err := withContext(f)
	if err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return err
}

// SaveTeamName
func (impl playerRepoImpl) SaveTeamName(p *domain.Player, version int) error {
	filter, err := impl.playerFilter(p)
	if err != nil {
		return err
	}

	return impl.update(filter, bson.M{fieldTeamName: p.Team.Name.TeamName()}, version)
}

// FindPlayer
func (impl playerRepoImpl) FindPlayer(cid string, a types.Account) (
	p domain.Player, version int, err error,
) {
	var v dPlayer

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(ctx, impl.docFilterByUser(cid, a), nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}
	} else {
		if err = v.toPlayer(&p); err == nil {
			p.SetCurrentUser(a)

			version = v.Version
		}
	}

	return
}

// FindCompetitionsUserApplied
func (impl playerRepoImpl) FindCompetitionsUserApplied(a types.Account) (
	r []string, err error,
) {
	var v []dPlayer

	f := func(ctx context.Context) error {
		filter := impl.docFilterByUser("", a)
		delete(filter, fieldCid)

		opts := options.FindOptions{}
		return impl.cli.GetDocs(ctx, filter, opts.SetProjection(bson.M{fieldCid: 1}), &v)
	}

	if err = withContext(f); err != nil || len(v) == 0 {
		return
	}

	r = make([]string, len(v))
	for i := range v {
		r[i] = v[i].CompetitionId
	}

	return
}

// CompetitorsCount
func (impl playerRepoImpl) CompetitorsCount(cid string) (int, error) {
	var v []struct {
		Total int `bson:"total"`
	}

	f := func(ctx context.Context) error {
		key := "$" + fieldCompetitors

		fields := bson.M{
			"num": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$isArray": key},
					"then": bson.M{"$size": key},
					"else": 0,
				},
			},
		}

		pipeline := bson.A{
			bson.M{"$match": impl.docFilter(cid)},
			bson.M{"$addFields": fields},
			bson.M{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$num"}}},
		}

		cursor, err := impl.cli.Collection().Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	}

	if err := withContext(f); err != nil || len(v) == 0 {
		return 0, err
	}

	return v[0].Total, nil
}

// AddMember
func (impl playerRepoImpl) AddMember(
	team repository.PlayerVersion,
	member repository.PlayerVersion,
) error {
	filter, err := impl.playerFilter(team.Player)
	if err != nil {
		return err
	}

	c := toCompetitorDoc(&member.Player.Leader)
	doc, err := genDoc(&c)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(
			ctx, filter,
			bson.M{fieldCompetitors: doc}, mongoCmdPush, team.Version,
		)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return err
}

// SavePlayer
func (impl playerRepoImpl) SavePlayer(p *domain.Player, version int) error {
	filter, err := impl.playerFilter(p)
	if err != nil {
		return err
	}

	doc, err := impl.genPlayerDoc(p)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(ctx, filter, doc, mongoCmdSet, version)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return err
}

func (impl playerRepoImpl) getVersion(cid string, a types.Account) (int, error) {
	var v dPlayer
	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, impl.disabledPlayerFilter(cid, a), bson.M{fieldVersion: 1}, &v,
		)
	}

	err := withContext(f)

	return v.Version, err
}

func (impl playerRepoImpl) ResumePlayer(cid string, a types.Account) error {
	version, err := impl.getVersion(cid, a)
	if err != nil {
		return err
	}

	return impl.update(
		impl.disabledPlayerFilter(cid, a), bson.M{fieldEnabled: true}, version,
	)
}
