package rpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	svc := &Service{}

	delay := 100 * time.Microsecond
	refTime := time.Now().UnixNano() / 1e6
	time.Sleep(delay)
	gotTime := svc.Time(&refTime)

	// -diff should be greater than delay, as we slept for at least delay
	assert.GreaterOrEqual(t, -*gotTime.diff, delay.Milliseconds())
}
