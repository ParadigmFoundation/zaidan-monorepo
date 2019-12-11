package rpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/levenlabs/golib/timeutil"
)

func TestTime(t *testing.T) {
	svc := &Service{}

	delay := 50 * time.Microsecond
	refTime := timeutil.TimestampNow().Float64()
	time.Sleep(delay)
	gotTime := svc.Time(&refTime)

	// -diff should be greater than delay, as we slept for at least delay
	assert.GreaterOrEqual(t, -*gotTime.diff, delay.Seconds())
}
