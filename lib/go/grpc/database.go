package grpc

import (
	"database/sql/driver"

	"github.com/gogo/protobuf/proto"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/utils"
)

func (m *TradeInfo) Value() (driver.Value, error) { return proto.Marshal(m) }
func (m *TradeInfo) Scan(v interface{}) error     { return proto.Unmarshal(utils.AnyToBytes(v), m) }

func (m *QuoteInfo) Value() (driver.Value, error) { return proto.Marshal(m) }
func (m *QuoteInfo) Scan(v interface{}) error     { return proto.Unmarshal(utils.AnyToBytes(v), m) }

func (m *SignedOrder) Value() (driver.Value, error) { return proto.Marshal(m) }
func (m *SignedOrder) Scan(v interface{}) error     { return proto.Unmarshal(utils.AnyToBytes(v), m) }
