package rpc

import (
	"encoding/json"
	"time"
)

// timeResponse represents the response from the dealer_time endpoint
type timeResponse struct {
	dealerTime int64

	// pointer so the default is null instead of 0
	diff *int64
}

// MarshalJSON implements json.Marshaller. TimeResponse is marshalled as an array
func (tr *timeResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{tr.dealerTime, tr.diff})
}

// Time implements the dealer_time method
func (svc *Service) Time(clientTime *int64) *timeResponse {
	ts := time.Now().UnixNano() / 1e6 // conversion to ms
	var diff *int64

	// if client has specified their local timestamp, we calculate the difference
	if clientTime != nil {
		tmp := ts - *clientTime
		diff = &tmp
	}

	return &timeResponse{dealerTime: ts, diff: diff}
}
