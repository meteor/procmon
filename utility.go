package procmon

import "github.com/meteor/procmon/ecu"
import "math"

// UserPerc calculates the percentage of CPU time spent in userland by
// the monitored process.
func (m *Measure) UserPerc() float64 {
	return 100.0 * float64(m.User) / float64(m.UserTotal+m.SystemTotal)
}

// SysPerc calculates the percentage of CPU time spent in the kernel by
// the monitored process.
func (m *Measure) SysPerc() float64 {
	return 100.0 * float64(m.System) / float64(m.UserTotal+m.SystemTotal)
}

// UserInECU calculates the amount of CPU time spent in userland, as
// measured in ECUs, for the monitored process assuming it is running
// on a machine of type `instance`.
func (m *Measure) UserInECU(instance *ecu.Instance) float64 {
	if instance == nil {
		return math.NaN()
	} else {
		return float64(instance.ComputeUnitsx10) * float64(m.User) / (float64(m.UserTotal+m.SystemTotal) * 10.0)
	}
}

// SysInECU calculates the amount of CPU time spent in the kernel, as
// measured in ECUs, for the monitored process assuming it is running
// on a machine of type `instance`.
func (m *Measure) SysInECU(instance *ecu.Instance) float64 {
	if instance == nil {
		return math.NaN()
	} else {
		return float64(instance.ComputeUnitsx10) * float64(m.System) / (float64(m.UserTotal+m.SystemTotal) * 10.0)
	}
}
