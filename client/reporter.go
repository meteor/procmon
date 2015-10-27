package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/meteor/procmon"
	"github.com/meteor/procmon/ecu"
	"math"
	"strconv"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Input process is missing")
	}

	process, err := strconv.ParseInt(flag.Arg(0), 10, 32)
	if err != nil {
		log.WithField("input", flag.Arg(0)).WithError(err).Fatal("Couldn't parse input process")
	}

	instance, err := ecu.Mine()
	if err != nil {
		log.WithError(err).Error("Couldn't find instance metadata")
	}

	output := make(chan procmon.Measure, 1)
	monitor, err := procmon.New(output, int(process))
	if err != nil {
		log.WithError(err).Fatal("Couldn't parse input process")
	}
	defer monitor.Stop()

outerloop:
	for {
		select {
		case point, ok := <-output:
			if !ok {
				log.Warn("Not ok, breaking")
				break outerloop
			}
			userSysTotal := point.UserTotal + point.SystemTotal
			var userPerc, sysPerc, userECU, sysECU float64
			if point.UserTotal == 0 {
				userPerc = 0.0
			} else {
				userPerc = 100.0 * float64(point.User) / float64(point.UserTotal)
			}
			if point.SystemTotal == 0 {
				sysPerc = 0.0
			} else {
				sysPerc = 100.0 * float64(point.System) / float64(point.SystemTotal)
			}
			if instance == nil {
				userECU = math.NaN()
				sysECU = math.NaN()
			} else {
				log.WithFields(log.Fields{
					"instance":     instance,
					"instance CPU": instance.ComputeUnitsx10,
					"userSysTotal": userSysTotal,
					"numerator":    float64(instance.ComputeUnitsx10) * float64(point.User),
					"denominator":  (float64(userSysTotal) * 10.0),
				}).Debug("computing")
				userECU = float64(instance.ComputeUnitsx10) * float64(point.User) / (float64(userSysTotal) * 10.0)
				sysECU = float64(instance.ComputeUnitsx10) * float64(point.System) / (float64(userSysTotal) * 10.0)
			}
			log.WithFields(log.Fields{
				"point":      point,
				"user":       userPerc,
				"system":     sysPerc,
				"userInECU":  userECU,
				"sysInECU":   sysECU,
				"memoryInKB": point.Memory,
			}).Info("Got point")
		}
	}
}
