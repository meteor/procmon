// +build linux

package dmesg

import "bufio"
import "os"
import "bytes"
import "io"
import "golang.org/x/sys/unix"
import "strings"
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
	for {
		rune, _, err := buf.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if rune == ' ' {
			// continuation line
			line, err := buf.ReadString('\n')
			if err != nil {
				return nil, err
			}
			if lastMessage == nil {
				return nil, fmt.Errorf("First line in ring buffer was a continuation line!")
			}
			lastMessage.Message += line[:len(line)-1]
		} else {
			if err := buf.UnreadRune(); err != nil {
				return nil, err
			}
			line, err := buf.ReadString(']')
			if err != nil {
				return nil, err
			}
			var level, secs, nanosecs int64
			matched, err := fmt.Sscanf(line, "<%d>[%d.%d]", &level, &secs, &nanosecs)
			if err != nil {
				return nil, err
			}
			if matched != 3 {
				return nil, fmt.Errorf("Couldn't parse line %q for some reason", line)
			}
			message, err := buf.ReadString('\n')
			if err != nil && err != io.EOF {
				return nil, err
			}
			adjustedTime := time.Unix(secs, nanosecs).Add(time.Second * time.Duration(s.bootTime.Unix()))
			result = append(result, &Message{
				level,
				adjustedTime,
				strings.TrimSuffix(message[1:], "\n"),
			})
		}
	}
	return result, nil
}
