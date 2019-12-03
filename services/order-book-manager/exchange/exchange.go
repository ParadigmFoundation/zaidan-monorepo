package exchange

import (
	"context"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
)

type Subscriber interface {
	OnSnapshot(string, *obm.Update) error
	OnUpdate(string, *obm.Update) error
}

type Exchange interface {
	// Subscribe subscribes to one or more symbols with a given set of callbacks
	Subscribe(context.Context, Subscriber, ...string) error
}
