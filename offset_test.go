package vast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOffsetMarshaler(t *testing.T) {
	b, err := Offset{}.MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "0%", string(b))
	}
	b, err = Offset{Percent: .1}.MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "10%", string(b))
	}
	d := Duration(0)
	b, err = Offset{Duration: &d}.MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "00:00:00", string(b))
	}
}

func TestOffsetUnmarshaler(t *testing.T) {
	var o Offset
	if assert.NoError(t, o.UnmarshalText([]byte("0%"))) {
		assert.Nil(t, o.Duration)
		assert.Equal(t, float32(0.0), o.Percent)
	}
	o = Offset{}
	if assert.NoError(t, o.UnmarshalText([]byte("10%"))) {
		assert.Nil(t, o.Duration)
		assert.Equal(t, float32(0.1), o.Percent)
	}
	o = Offset{}
	if assert.NoError(t, o.UnmarshalText([]byte("00:00:00"))) {
		if assert.NotNil(t, o.Duration) {
			assert.Equal(t, Duration(0), *o.Duration)
		}
		assert.Equal(t, float32(0), o.Percent)
	}
	o = Offset{}
	assert.EqualError(t, o.UnmarshalText([]byte("abc%")), "invalid offset: abc%")
}
