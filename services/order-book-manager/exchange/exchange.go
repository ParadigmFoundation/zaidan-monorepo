package exchange

import (
	"context"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store"
)

type Exchange interface {
	// Subscribe subscribes to one or more symbols with a given set of callbacks
	Subscribe(context.Context, store.Store, ...string) error
}
