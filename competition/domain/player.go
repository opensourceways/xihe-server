package domain

import (
	"errors"

	types "github.com/opensourceways/xihe-server/domain"
)

var errorTeamMembersEnough = errors.New("members are enough")

func IsErrorTeamMembersEnough(err error) bool {
	return errors.Is(err, errorTeamMembersEnough)
}

type Competitor struct {
	Account  types.Account
	Name     CompetitorName
	City     City
	Email    types.Email
	Phone    Phone
	Identity CompetitionIdentity
	Province Province
	Detail   map[string]string
}

// Team
type Team struct {
	Name    TeamName
	Members []Competitor
}

func (t *Team) isMember(a types.Account) bool {
	return t.indexOfMember(a) >= 0
}

func (t *Team) indexOfMember(a types.Account) int {
	for i := range t.Members {
		item := &t.Members[i]
		if item.Account.Account() == a.Account() {
			return i
		}
	}

	return -1
}

func (t *Team) remove(a types.Account) error {
	i := t.indexOfMember(a)
	if i < 0 {
		return errors.New("not a member")
	}

	n := len(t.Members) - 1

	if i == 0 {
		if n == 0 {
			t.Members = nil
		} else {
			t.Members = t.Members[1:]
		}
	} else {
		if i != n {
			t.Members[i] = t.Members[n]
		}
		t.Members = t.Members[:n]
	}

	return nil
}

func (t *Team) join(c *Competitor) error {
	// TODO config
	if len(t.Members) >= 2 {
		return errorTeamMembersEnough
	}

	if t.isMember(c.Account) {
		return errors.New("already joined")
	}

	t.Members = append(t.Members, *c)

	return nil
}

// Player
type PlayerIndex struct {
	Id            string
	CompetitionId string
}

func NewPlayerIndex(cid, pid string) PlayerIndex {
	return PlayerIndex{
		Id:            pid,
		CompetitionId: cid,
	}
}

type Player struct {
	PlayerIndex
	IsFinalist bool

	Leader Competitor
	Team   Team

	// user is the current competitor who maybe is just a member of team or the leader.
	user types.Account
}

func (p *Player) SetCurrentUser(a types.Account) {
	p.user = a
}

func (p *Player) CompetitorsCount() int {
	return len(p.Team.Members) + 1
}

func (p *Player) IsIndividual() bool {
	return p.Team.Name == nil
}

func (p *Player) IsATeam() bool {
	return !p.IsIndividual()
}

func (p *Player) isUserTheLeader() bool {
	return p.user != nil && p.user.Account() == p.Leader.Account.Account()
}

func (p *Player) IsIndividualOrLeader() bool {
	return p.IsIndividual() || p.isUserTheLeader()
}

func (p *Player) Name() string {
	if p.IsIndividual() {
		return p.Leader.Account.Account()
	}

	return p.Team.Name.TeamName()
}

func (p *Player) Members() []Competitor {
	return p.Team.Members
}

func (p *Player) Has(u types.Account) bool {
	return p.Leader.Account.Account() == u.Account() || p.Team.isMember(u)
}

func (p *Player) RoleOfCurrentCompetitor() string {
	if p.IsATeam() && p.isUserTheLeader() {
		return competitionTeamRoleLeader
	}

	return ""
}

func (p *Player) CreateTeam(name TeamName) error {
	if p.IsATeam() {
		return errors.New("team is ready, no need to create team again")
	}

	p.Team.Name = name

	return nil
}

func (p *Player) ChangeTeamName(name TeamName) error {
	if !p.IsATeam() {
		return errors.New("I am not a team")
	}

	if !p.isUserTheLeader() {
		return errors.New("I am not leader")
	}

	p.Team.Name = name

	return nil
}

func (p *Player) JoinTo(team *Player) error {
	if !p.IsIndividual() {
		return errors.New("you are not an individual competitor")
	}

	if !team.IsATeam() {
		return errors.New("it is not a team")
	}

	return team.join(&p.Leader)
}

func (p *Player) join(c *Competitor) error {
	if p.Leader.Account.Account() == c.Account.Account() {
		return errors.New("invalid operation")
	}

	return p.Team.join(c)
}

func (p *Player) Quit() error {
	if !p.IsATeam() {
		return errors.New("invalid operation")
	}

	if p.isUserTheLeader() {
		return errors.New("leader can't quit directly")
	}

	return p.Team.remove(p.user)
}

func (p *Player) Delete(c types.Account) error {
	if !p.IsATeam() {
		return errors.New("invalid operation")
	}

	if !p.isUserTheLeader() {
		return errors.New("only leader can delete a member")
	}

	return p.Team.remove(c)
}

func (p *Player) TransferLeader(newOne types.Account) error {
	if !p.IsATeam() {
		return errors.New("invalid operation")
	}

	if !p.isUserTheLeader() {
		return errors.New("only leader can transfer leader")
	}

	t := &p.Team
	i := t.indexOfMember(newOne)
	if i < 0 {
		return errors.New("not a member")
	}

	p.Leader, t.Members[i] = t.Members[i], p.Leader

	return nil
}
