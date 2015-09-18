package vast

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationMarshaler(t *testing.T) {
	b, err := Duration(0).MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "00:00:00", string(b))
	}
	b, err = Duration(2 * time.Second).MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "00:00:02", string(b))
	}
	b, err = Duration(2 * time.Minute).MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "00:02:00", string(b))
	}
	b, err = Duration(2 * time.Hour).MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "02:00:00", string(b))
	}
}

func TestDurationUnmarshalXML(t *testing.T) {
	var d Duration
	if assert.NoError(t, d.UnmarshalText([]byte("00:00:00"))) {
		assert.Equal(t, Duration(0), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte("00:00:02"))) {
		assert.Equal(t, Duration(2*time.Second), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte("00:02:00"))) {
		assert.Equal(t, Duration(2*time.Minute), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte("02:00:00"))) {
		assert.Equal(t, Duration(2*time.Hour), d)
	}
	assert.EqualError(t, d.UnmarshalText([]byte("00:00:60")), "invalid duration: 00:00:60")
	assert.EqualError(t, d.UnmarshalText([]byte("00:60:00")), "invalid duration: 00:60:00")
	assert.EqualError(t, d.UnmarshalText([]byte("00h01m")), "invalid duration: 00h01m")
}
