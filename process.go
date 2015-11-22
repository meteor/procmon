// Package procmon is a simple library for recording samples of a
// process and sending the results to datadog.
package procmon

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

// Measure is a point in time measure of a process's resource consumption.
type Measure struct {
	// User is the number of usermode CPU jiffies used
	User uint64
	// System is the number of kernelmode CPU jiffies used
	System uint64
	// UserTotal is the number of usermode CPU jiffies used across all processes
	UserTotal uint64
	// SystemTotal is the number of kernelmode CPU jiffies used across all processes
	SystemTotal uint64
	// IdleTotal is the number of CPU jiffies spent in the idle task.
	IdleTotal uint64
	// Memory is the amount of memory used in kB
	Memory uint64
}

// Point in time measure of a process's state
type point struct {
	user   uint64
	system uint64
	idle   uint64
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

			log.WithFields(log.Fields{
				"new total":  newtotal,
				"new target": newtarget,
				"old total":  m.total,
				"old target": m.stats,
			}).Debug("tick")
			select {
			case m.Output <- Measure{
				newtarget.user - m.stats.user,
				newtarget.system - m.stats.system,
				newtotal.user - m.total.user,
				newtotal.system - m.total.system,
				newtotal.idle - m.total.idle,
				memory,
			}:
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
