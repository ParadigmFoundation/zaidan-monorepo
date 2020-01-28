package grpc

import (
	"database/sql/driver"
	"encoding/base64"

	"github.com/gogo/protobuf/proto"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/utils"
)

func dbValue(pb proto.Message) (string, error) {
	bz, err := proto.Marshal(pb)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bz), nil
}

func dbScan(bz []byte, pb proto.Message) error {
	dec, err := base64.StdEncoding.DecodeString(string(bz))
	if err != nil {
		return err
	}
	return proto.Unmarshal(dec, pb)
}

func (m *TradeInfo) Value() (driver.Value, error) { return dbValue(m) }
func (m *TradeInfo) Scan(v interface{}) error     { return dbScan(utils.AnyToBytes(v), m) }

func (m *QuoteInfo) Value() (driver.Value, error) { return dbValue(m) }
func (m *QuoteInfo) Scan(v interface{}) error     { return dbScan(utils.AnyToBytes(v), m) }

func (m *SignedOrder) Value() (driver.Value, error) { return dbValue(m) }
func (m *SignedOrder) Scan(v interface{}) error     { return dbScan(utils.AnyToBytes(v), m) }

func (m *ZeroExTransactionInfo) Value() (driver.Value, error) { return dbValue(m) }
func (m *ZeroExTransactionInfo) Scan(v interface{}) error     { return dbScan(utils.AnyToBytes(v), m) }
