// Package procmon is a simple library for recording samples of a
// process and sending the results to datadog.
package procmon

import (
	"bufio"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// Measure is a point in time measure of a process's resource consumption.
type Measure struct {
	// User is the % of usermode CPU time used
	User float64
	// System is the % of kernel mode CPU time used
	System float64
	// Memory is the amount of memory used in kB
	Memory uint64
}

// Point in time measure of a process's state
type point struct {
	user   uint64
	system uint64
}

// Monitor represents a continuous monitoring of a given Linux
// process.
type Monitor struct {
	Output  chan<- Measure
	ticker  *time.Ticker
	process int
	done    chan bool
	stats   point
	total   point
}

// New creates a new monitor and starts it.
func New(out chan<- Measure, process int) (*Monitor, error) {
	m := new(Monitor)
	m.ticker = time.NewTicker(5 * time.Second)
	m.done = make(chan bool, 1)
	m.process = process
	m.Output = out
	if err := m.preflight(); err != nil {
		return nil, err
	}
	go m.Monitor()
	return m, nil
}

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
	return point{cutime, cstime}, nil
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
			for i := 0; i < 2; i++ {
				if !fs.Scan() {
					return point{}, fmt.Errorf("cpu line ended before user data seen: %q", s.Text())
				}
			}
			utime, err := strconv.ParseUint(fs.Text(), 10, 64)
			if err != nil {
				return point{}, err
			}
			for i := 0; i < 2; i++ {
				if !fs.Scan() {
					return point{}, fmt.Errorf("cpu line ended before system data seen: %q", s.Text())
				}
			}
			stime, err := strconv.ParseUint(fs.Text(), 10, 64)
			if err != nil {
				return point{}, err
			}
			return point{utime, stime}, nil
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

// Monitor runs in a background goroutine that can be halted with Stop
// and monitors process metrics, submitting results to datadog every 5
// seconds.
func (m *Monitor) Monitor() {
	var err error
	m.stats, err = m.fetchProcessUsage()
	if err != nil {
		log.WithField("process", m.process).WithError(err).
			Error("couldn't read process stats")
	}
	m.total, err = m.fetchTotalUsage()
	if err != nil {
		log.WithField("process", m.process).WithError(err).
			Error("couldn't read total CPU stats")
	}
	for {
		select {
		case <-m.ticker.C:
			newtarget, err := m.fetchProcessUsage()
			if err != nil {
				log.WithField("process", m.process).WithError(err).
					Error("couldn't read process stats")
				m.ticker.Stop()
				close(m.Output)
				return
			}
			memory, err := m.fetchProcessMemory()
			if err != nil {
				log.WithField("process", m.process).WithError(err).
					Error("couldn't read process stats")
				m.ticker.Stop()
				close(m.Output)
				return
			}
			newtotal, err := m.fetchTotalUsage()
			if err != nil {
				// this is a weird one, as it indicates something has
				// gone seriously haywire.  Still, closing as normal.
				log.WithField("process", m.process).WithError(err).
					Error("couldn't read total CPU stats")
				m.ticker.Stop()
				close(m.Output)
				return
			}
			userTotalDiff := newtotal.user - m.total.user
			sysTotalDiff := newtotal.system - m.total.system

			var userPerc, sysPerc float64
			if userTotalDiff == 0 {
				userPerc = 0.0
			} else {
				userPerc = 100.0 * float64(newtarget.user-m.stats.user) /
					float64(userTotalDiff)
			}
			if sysTotalDiff == 0 {
				sysPerc = 0.0
			} else {
				sysPerc = 100.0 * float64(newtarget.system-m.stats.system) /
					float64(sysTotalDiff)
			}

			select {
			case m.Output <- Measure{userPerc, sysPerc, memory}:
			default:
				log.WithField("process", m.process).Warn("Output full, dropping update")
			}
			m.stats = newtarget
			m.total = newtotal
		case <-m.done:
			return
		}
	}
}

// Stop halts the background monitoring task.
func (m *Monitor) Stop() {
	m.done <- true
	m.ticker.Stop()
}
