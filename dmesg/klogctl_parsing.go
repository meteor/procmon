package dmesg

import "time"
import "fmt"
import "regexp"
import "strconv"
import "strings"

var hmsRegex = regexp.MustCompile(`(?:<(\d+)>)?(\d+):(\d+):(\d+).(\d+) (.*)\n?`)
var bracketRegex = regexp.MustCompile(`(?:<(\d+)>)?(?:\[(\d+).(\d+)\])?(.*)\n?`)

func parseMessage(message string) (*Message, error) {
	parts := hmsRegex.FindStringSubmatch(message)
	if parts != nil {
		return parseHMS(parts)
	}
	parts = bracketRegex.FindStringSubmatch(message)
	if parts != nil {
		return parseBracketed(parts)
	}

	return nil, fmt.Errorf("Could not parse string %q", message)
}

func parseHMS(parts []string) (*Message, error) {
	var level int64
	var err error
	if parts[1] == "" {
		level = 6
	} else {
		level, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	hours, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, err
	}

	minutes, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return nil, err
	}

	secs, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return nil, err
	}

	// what we're getting from the regex is 0.xxxxxx.  In the event
	// that that is not the full 9 digit nanosecond range, it needs to
	// be adjusted so that `ParseInt` behaves well.
	var prensecs string
	if len(parts[5]) > 9 {
		prensecs = parts[5][0:9]
	} else {
		prensecs = parts[5] + strings.Repeat("0", 9-len(parts[5]))
	}

	nsecs, err := strconv.ParseInt(prensecs, 10, 64)
	if err != nil {
		return nil, err
	}

	// prensecs is of course a decimal.  True value is prensecs /
	// 10^(len((parts[2])); conversion to nanoseconds is that times 1e9.

	return &Message{level, time.Unix(hours*60*60+minutes*60+secs, nsecs), parts[6]}, nil
}

func parseBracketed(parts []string) (*Message, error) {
	var level int64
	var err error
	if parts[1] == "" {
		level = 6
	} else {
		level, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	var secs, nsecs int64

	if parts[2] == "" {
		secs = 0
	} else {
		secs, err = strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if parts[3] == "" {
		nsecs = 0
	} else {
		nsecs, err = strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &Message{level, time.Unix(secs, nsecs), parts[4]}, nil
}
