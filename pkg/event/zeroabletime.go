package event

import "time"

type ZeroableTime struct {
	time.Time
}

func (t ZeroableTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}
