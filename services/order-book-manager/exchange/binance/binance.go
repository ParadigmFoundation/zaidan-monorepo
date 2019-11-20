package binance

import (
	"context"
	"fmt"
	"log"

	"github.com/adshao/go-binance"
	"github.com/kr/pretty"
)

type Exchange struct {
}

func New() *Exchange {
	return &Exchange{}
}

func (x *Exchange) depthHandler(event *binance.WsDepthEvent) {
	pretty.Print(event)
}

func (x *Exchange) Subscribe(ctx context.Context) error {
	errHandler := func(err error) {
		log.Printf("ERROR: %+v", err)
	}

	doneCh, stopCh, err := binance.WsDepthServe("ETHBTC", x.depthHandler, errHandler)
	if err != nil {
		return fmt.Errorf("Subscribe(): %w", err)
	}

	select {
	case <-ctx.Done():
		stopCh <- struct{}{}
		return ctx.Err()
	case <-doneCh:
		return nil
	}
}
