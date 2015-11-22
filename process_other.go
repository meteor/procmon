// +build !linux

package procmon

import "errors"

var errNotSupported = errors.New("Not supported on this OS")

func (m *Monitor) preflight() error {
	return errNotSupported
}

func (m *Monitor) fetchProcessUsage() (point, error) {
	return point{}, errNotSupported
}

func (m *Monitor) fetchProcessMemory() (uint64, error) {
	return 0, errNotSupported
}

func (m *Monitor) fetchTotalUsage() (point, error) {
	return point{}, errNotSupported
}
