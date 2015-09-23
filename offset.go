package vast

import (
	"fmt"
	"strconv"
	"strings"
)

// Offset represents either a vast.Duration or a percentage of the video duration.
type Offset struct {
	// If not nil, the Offset is duration based
	Duration *Duration
	// If Duration is nil, the Offset is percent based
	Percent float32
}

// MarshalText implements the encoding.TextMarshaler interface.
func (o Offset) MarshalText() ([]byte, error) {
	if o.Duration != nil {
		return o.Duration.MarshalText()
	}
	return []byte(fmt.Sprintf("%d%%", int(o.Percent*100))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (o *Offset) UnmarshalText(data []byte) error {
	if strings.HasSuffix(string(data), "%") {
		p, err := strconv.ParseInt(string(data[:len(data)-1]), 10, 8)
		if err != nil {
			return fmt.Errorf("invalid offset: %s", data)
		}
		o.Percent = float32(p) / 100
		return nil
	}
	var d Duration
	o.Duration = &d
	return o.Duration.UnmarshalText(data)
}
