package rpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	svc := &Service{}

	delay := 2 * time.Millisecond
	refTime := time.Now().UnixNano() / 1e6
	time.Sleep(delay)
	gotTime := svc.Time(&refTime)

	// diff should be >= delay, as we slept for at least delay
	assert.GreaterOrEqual(t, *gotTime.diff, delay.Milliseconds())
}
