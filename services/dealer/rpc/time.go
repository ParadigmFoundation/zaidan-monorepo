package rpc

import (
	"encoding/json"

	"github.com/levenlabs/golib/timeutil"
)

// TimeResponse represents the response from the dealer_time endpoint
type TimeResponse struct {
	dealerTime float64

	// pointer so the default is null instead of 0
	diff *float64
}

// MarshalJSON implements json.Marshaller. TimeResponse is marshalled as an array
func (tr *TimeResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"time": tr.dealerTime,
		"diff": tr.diff,
	})
}

// Time implements the dealer_time method
func (svc *Service) Time(clientTime *float64) *TimeResponse {
	diff := new(float64)
	ts := timeutil.TimestampNow()

	// if client has specified their local timestamp, we calculate the difference
	if clientTime != nil {
		tmp := timeutil.TimestampFromFloat64(*clientTime).Time.Sub(ts.Time).Seconds()
		diff = &tmp
	}

	return &TimeResponse{dealerTime: ts.Float64(), diff: diff}
}
