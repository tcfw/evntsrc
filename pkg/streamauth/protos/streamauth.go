package evntsrc_streamauth

import "errors"

//Validate ensures a correct key
func (key *StreamKey) Validate(isCreate bool) error {
	if isCreate && key.GetId() != "" {
		return errors.New("New keys cannot have an id")
	}
	if isCreate && key.GetKey() != "" {
		return errors.New("New keys Cannot have a key")
	}
	if isCreate && key.GetSecret() != "" {
		return errors.New("New keys Cannot have a secret")
	}
	if key.GetStream() == 0 {
		return errors.New("Must have stream set")
	}

	return nil
}
