package dmesg

import "bytes"
import "fmt"
import "io"
import "time"
import "strings"
import "os"
import "bufio"

type Message struct {
	Level     int64
	Timestamp time.Time
	Message   string
}

func ParseMessages(bootTime time.Time, buffer []byte) ([]*Message, error) {
	buf := bytes.NewBuffer(buffer)
	result := make([]*Message, 0)
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
				return nil, fmt.Errorf("Couldn't parse line for some reason")
			}
			message, err := buf.ReadString('\n')
			if err != nil && err != io.EOF {
				return nil, err
			}
			adjustedTime := bootTime.Add(time.Second*time.Duration(secs) + time.Nanosecond*time.Duration(nanosecs))
			result = append(result, &Message{
				level,
				adjustedTime,
				strings.TrimSuffix(message, "\n"),
			})
		}
	}
	return result, nil
}

func Uptime() (time.Time, error) {
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
