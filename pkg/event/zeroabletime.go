package event

import "time"

//ZeroableTime converts to "null" string for zero times when converting to JSON
type ZeroableTime struct {
	time.Time
}

//MarshalJSON converts struct to JSON type
func (t ZeroableTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}
