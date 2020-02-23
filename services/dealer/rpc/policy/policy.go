package policy

import "github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store"

type Mode int

const (
	BlackListMode Mode = iota
	WhiteListMode
)

type Reason string

const (
	WhiteListed   Reason = "WHITELISTED"
	BlackListed          = "BLACKLISTED"
	ReasonUnknown        = "UNKNOWN"
)

type Policy struct {
	db   store.Policy
	mode Mode
}

func New(db store.Policy, mode Mode) *Policy {
	return &Policy{db: db, mode: mode}
}

func (p *Policy) Store() store.Policy { return p.db }

func (p *Policy) AuthStatus(addr string) (Reason, bool, error) {
	found, err := p.db.HasPolicy(addr)
	if err != nil {
		return ReasonUnknown, false, err
	}

	var auth bool
	switch p.mode {
	case BlackListMode:
		auth = !found
	case WhiteListMode:
		auth = found
	}

	if auth {
		return WhiteListed, auth, nil
	}
	return BlackListed, auth, nil
}
