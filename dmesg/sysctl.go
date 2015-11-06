// +build linux

package dmesg

import "golang.org/x/sys/unix"

// dunno why these aren't exposed by the sys module
const syslogActionSizeBuffer = 10
const syslogActionReadAll = 3

// Current retrieves the current contents of the kernel message ring buffer.
func Current() ([]byte, error) {
	size, err := unix.Klogctl(syslogActionSizeBuffer, nil)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, size)
	bytesread, err := unix.Klogctl(syslogActionReadAll, buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:bytesread], nil
}
