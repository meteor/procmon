// +build linux

package procmon

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func (m *Monitor) preflight() error {
	_, err := m.fetchProcessUsage()
	if err != nil {
		return err
	}
	_, err = m.fetchTotalUsage()
	if err != nil {
		return err
	}
	return nil
}

func parseProcStat(in io.Reader) (point, error) {
	// per proc(5) we are after fields utime and stime, numbers 14
	// and 15.
	contents, err := ioutil.ReadAll(in)
	if err != nil {
		return point{}, err
	}
	fields := strings.Split(string(contents), " ")
	if len(fields) < 15 {
		return point{}, fmt.Errorf("Not enough fields")
	}
	cutime, err := strconv.ParseUint(fields[13], 10, 64)
	if err != nil {
		return point{}, err
	}
	cstime, err := strconv.ParseUint(fields[14], 10, 64)
	if err != nil {
		return point{}, err
	}
	return point{cutime, cstime, 0}, nil
}

func parseMemStat(in io.Reader) (uint64, error) {
	// per proc(5) we are after field rss, number 2.
	contents, err := ioutil.ReadAll(in)
	if err != nil {
		return 0, err
	}
	fields := strings.Split(string(contents), " ")
	if len(fields) < 2 {
		return 0, fmt.Errorf("Not enough fields")
	}
	rss, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return 0, err
	}
	return rss, nil
}

func parseGlobalStat(in io.Reader) (point, error) {
	// again per proc(5) we are after the user and system fields from
	// the cpu line, which are fields 1 and 3.  This is complicated
	// because /proc/stat gives lots of information we don't want or
	// need.
	s := bufio.NewScanner(in)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "cpu ") {
			fs := bufio.NewScanner(strings.NewReader(s.Text()))
			fs.Split(bufio.ScanWords)
			if !fs.Scan() {
				return point{}, fmt.Errorf("cpu line ended before user data seen: %q", s.Text())
			}
			if fs.Text() != "cpu" {
				return point{}, fmt.Errorf("Weird cpu line doesn't start with cpu?!: %q", s.Text())
			}
			if !fs.Scan() {
				return point{}, fmt.Errorf("cpu line ended before user data seen: %q", s.Text())
			}
			utime, err := strconv.ParseUint(fs.Text(), 10, 64)
			if err != nil {
				return point{}, err
			}
			if !fs.Scan() {
				return point{}, fmt.Errorf("cpu line ended before user niced data seen: %q", s.Text())
			}
			unicetime, err := strconv.ParseUint(fs.Text(), 10, 64)
			if err != nil {
				return point{}, err
			}
			utime += unicetime
			if !fs.Scan() {
				return point{}, fmt.Errorf("cpu line ended before system data seen: %q", s.Text())
			}
			stime, err := strconv.ParseUint(fs.Text(), 10, 64)
			if err != nil {
				return point{}, err
			}
			if !fs.Scan() {
				return point{}, fmt.Errorf("cpu line ended before idle data seen: %q", s.Text())
			}
			itime, err := strconv.ParseUint(fs.Text(), 10, 64)
			if err != nil {
				return point{}, err
			}
			log.WithFields(log.Fields{
				"raw":  s.Text(),
				"user": utime,
				"sys":  stime,
				"idle": itime,
			}).Debug("reading /proc/stat")
			return point{utime, stime, itime}, nil
		}
	}
	return point{}, fmt.Errorf("No line starting with 'cpu' seen")
}

func (m *Monitor) fetchProcessUsage() (point, error) {
	file, err := os.Open(fmt.Sprintf("/proc/%d/stat", m.process))
	if err != nil {
		return point{}, err
	}
	defer file.Close()
	return parseProcStat(file)
}

func (m *Monitor) fetchProcessMemory() (uint64, error) {
	file, err := os.Open(fmt.Sprintf("/proc/%d/statm", m.process))
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return parseMemStat(file)
}

func (m *Monitor) fetchTotalUsage() (point, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return point{}, err
	}
	defer file.Close()
	return parseGlobalStat(file)
}
