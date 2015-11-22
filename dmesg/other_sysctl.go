// +build !linux

package dmesg

import "errors"

var errNotSupported = errors.New("Not supported on this operating system")

// State is an opaque datatype representing whatever state is needed
// for the dmesg parser.
type State struct {
}

// Current retrieves the current contents of the kernel message ring
// buffer.  On this operating system it is not supported.
func (s *State) Current() ([]byte, error) {
	return nil, errNotSupported
}

// New creates a new State.
func New() (*State, error) {
	return nil, errNotSupported
}

// ParseMessages reads dmesg type messages out of buffer.  On this
// operating system, it is not supported.
func (s *State) ParseMessages([]byte) ([]*Message, error) {
	return nil, errNotSupported
}
