// +build linux

package dmesg

import "bufio"
import "os"
import "bytes"
import "io"
import "golang.org/x/sys/unix"
import "time"
import "fmt"

// dunno why these aren't exposed by the sys module
const syslogActionSizeBuffer = 10
const syslogActionReadAll = 3

// State is an opaque datatype representing whatever state is needed
// for the dmesg parser.
type State struct {
	bootTime time.Time
}

// Current retrieves the current contents of the kernel message ring buffer.
func (s *State) Current() ([]byte, error) {
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

func uptime() (time.Time, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return time.Unix(0, 0), err
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		var btime int64
		matched, err := fmt.Sscanf(s.Text(), "btime %d\n", &btime)
		if err == nil && matched == 1 {
			// aha, that's the line.
			return time.Unix(btime, 0), nil
		}
	}
	return time.Unix(0, 0), fmt.Errorf("Did not find btime declaration in /proc/stat")
}

// New creates a new State.
func New() (*State, error) {
	var err error
	s := State{}
	s.bootTime, err = uptime()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// ParseMessages reads dmesg type messages out of buffer.  To do so it
// must read the system boot time out of /proc/stat, because dmesg
// timestamps are relative to when the system booted.
func (s *State) ParseMessages(buffer []byte) ([]*Message, error) {
	buf := bytes.NewBuffer(buffer)
	var result []*Message
	var lastMessage *Message
	lastTimestamp := time.Unix(0, 0)
	for {
		rune, _, err := buf.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return result, err
		}
		if rune == ' ' {
			// continuation line
			line, err := buf.ReadString('\n')
			if err != nil {
				return result, err
			}
			if lastMessage == nil {
				return result, fmt.Errorf("First line in ring buffer was a continuation line!")
			}
			lastMessage.Message += line[:len(line)-1]
		} else {
			if err := buf.UnreadRune(); err != nil {
				return result, err
			}
			line, err := buf.ReadString('\n')
			if err != nil {
				return result, err
			}

			message, err := parseMessage(line)
			if err != nil {
				return result, err
			}

			if message.Timestamp == time.Unix(0, 0) {
				message.Timestamp = lastTimestamp
			} else {
				adjustedTime := message.Timestamp.Add(time.Second * time.Duration(s.bootTime.Unix()))
				message.Timestamp = adjustedTime
				lastTimestamp = message.Timestamp
			}

			result = append(result, message)
		}
	}
	return result, nil
}
