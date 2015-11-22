package dmesg

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPlainLine(t *testing.T) {
	message, err := parseMessage("sample test message\n")
	if assert.NoError(t, err) {
		assert.Equal(t, int64(6), message.Level)
		assert.Equal(t, time.Unix(0, 0), message.Timestamp)
		assert.Equal(t, "sample test message", message.Message)
	}
}

func TestLineWithoutNL(t *testing.T) {
	message, err := parseMessage("sample test message")
	if assert.NoError(t, err) {
		assert.Equal(t, int64(6), message.Level)
		assert.Equal(t, time.Unix(0, 0), message.Timestamp)
		assert.Equal(t, "sample test message", message.Message)
	}
}

func TestPriority(t *testing.T) {
	message, err := parseMessage("<4>sample test message\n")
	if assert.NoError(t, err) {
		assert.Equal(t, int64(4), message.Level)
		assert.Equal(t, time.Unix(0, 0), message.Timestamp)
		assert.Equal(t, "sample test message", message.Message)
	}
}

func TestNonPriorityStupidity(t *testing.T) {
	message, err := parseMessage("<sample test message\n")
	if assert.NoError(t, err) {
		assert.Equal(t, int64(6), message.Level)
		assert.Equal(t, time.Unix(0, 0), message.Timestamp)
		assert.Equal(t, "<sample test message", message.Message)
	}
}

func TestTimestamp(t *testing.T) {
	message, err := parseMessage("[42.42]sample test message\n")
	if assert.NoError(t, err) {
		assert.Equal(t, int64(6), message.Level)
		assert.Equal(t, time.Unix(42, 42), message.Timestamp)
		assert.Equal(t, "sample test message", message.Message)
	}
}

func TestPriorityAndTimestamp(t *testing.T) {
	message, err := parseMessage("<4>[42.42]sample test message\n")
	if assert.NoError(t, err) {
		assert.Equal(t, int64(4), message.Level)
		assert.Equal(t, time.Unix(42, 42), message.Timestamp)
		assert.Equal(t, "sample test message", message.Message)
	}
}

func TestHMSTimestamp(t *testing.T) {
	message, err := parseMessage("42:42:42.42 sample test message\n")
	if assert.NoError(t, err) {
		assert.Equal(t, int64(6), message.Level)
		assert.Equal(t, time.Unix(42*60*60+42*60+42, 420000000), message.Timestamp)
		assert.Equal(t, "sample test message", message.Message)
	}
	message, err = parseMessage("42:41:40.001564 sample test message\n")
	if assert.NoError(t, err) {
		assert.Equal(t, int64(6), message.Level)
		assert.Equal(t, time.Unix(42*60*60+41*60+40, 1564000), message.Timestamp)
		assert.Equal(t, "sample test message", message.Message)
	}
}
