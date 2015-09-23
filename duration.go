package vast

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Duration is a VAST duration expressed a hh:mm:ss
type Duration time.Duration

// MarshalText implements the encoding.TextMarshaler interface.
func (dur Duration) MarshalText() ([]byte, error) {
	s := int(dur / Duration(time.Second))
	ms := int(dur) - s*int(time.Second)
	h := int(s / 3600)
	s = s % 3600
	m := int(s / 60)
	s = s % 60
	if ms == 0 {
		return []byte(fmt.Sprintf("%02d:%02d:%02d", h, m, s)), nil
	}
	return []byte(fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (dur *Duration) UnmarshalText(data []byte) (err error) {
	parts := strings.SplitN(string(data), ":", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid duration: %s", data)
	}
	if i := strings.IndexByte(parts[2], '.'); i > 0 {
		ms, err := strconv.ParseInt(parts[2][i+1:], 10, 32)
		if err != nil || ms < 0 || ms > 999 {
			return fmt.Errorf("invalid duration: %s", data)
		}
		parts[2] = parts[2][:i]
		*dur += Duration(ms)
	}
	f := Duration(time.Second)
	for i := 2; i >= 0; i-- {
		n, err := strconv.ParseInt(parts[i], 10, 32)
		if err != nil || n < 0 || n > 59 {
			return fmt.Errorf("invalid duration: %s", data)
		}
		*dur += Duration(n) * f
		f *= 60
	}
	return nil
}
