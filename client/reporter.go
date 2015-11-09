package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/meteor/procmon"
	"github.com/meteor/procmon/ecu"
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
			log.WithFields(log.Fields{
				"user":       point.UserPerc(),
				"system":     point.SysPerc(),
				"userInECU":  point.UserInECU(instance),
				"sysInECU":   point.SysInECU(instance),
				"memoryInKB": point.Memory,
			}).Debug("Got point")
		}
	}
}
