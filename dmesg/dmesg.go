package dmesg

import "time"

// Message encapsulates kernel ring buffer messages.  The timestamp is
// absolute, not relative to boot time.
type Message struct {
	Level     int64
	Timestamp time.Time
	Message   string
}

// Messages retrieves all kernel ring buffer messages and returns
// them, or a reason why it could not.
func (s *State) Messages() ([]*Message, error) {
	buffer, err := s.Current()
	if err != nil {
		return nil, err
	}
	messages, err := s.ParseMessages(buffer)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
