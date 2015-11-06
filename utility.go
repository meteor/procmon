package procmon

import "github.com/meteor/procmon/ecu"
import "math"

// Total calculates the total jiffies spent between ticks.
func (m *Measure) Total() uint64 {
	return m.UserTotal + m.SystemTotal + m.IdleTotal
}

// UserPerc calculates the percentage of CPU time spent in userland by
// the monitored process.
func (m *Measure) UserPerc() float64 {
	return 100.0 * float64(m.User) / float64(m.Total())
}

// IdlePerc calculates the percentage of CPU time spent idling.
func (m *Measure) IdlePerc() float64 {
	return 100.0 * float64(m.IdleTotal) / float64(m.Total())
}

// SysPerc calculates the percentage of CPU time spent in the kernel by
// the monitored process.
func (m *Measure) SysPerc() float64 {
	return 100.0 * float64(m.System) / float64(m.Total())
}

func (m *Measure) scaleBy(datum uint64, instance *ecu.Instance) float64 {
	return float64(instance.ComputeUnitsx10) * float64(datum) /
		(float64(m.Total()) * 10.0)
}

// UserInECU calculates the amount of CPU time spent in userland, as
// measured in ECUs, for the monitored process assuming it is running
// on a machine of type `instance`.
func (m *Measure) UserInECU(instance *ecu.Instance) float64 {
	if instance == nil {
		return math.NaN()
	} else {
		return m.scaleBy(m.User, instance)
	}
}

// SysInECU calculates the amount of CPU time spent in the kernel, as
// measured in ECUs, for the monitored process assuming it is running
// on a machine of type `instance`.
func (m *Measure) SysInECU(instance *ecu.Instance) float64 {
	if instance == nil {
		return math.NaN()
	} else {
		return m.scaleBy(m.System, instance)
	}
}
