package ecu

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Mine fetches the instance type for the current running instance and
// retrieves data for it.
func Mine() (*Instance, error) {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-type")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"body":      body,
		"as string": string(body),
		"bytes":     []byte("m1.medium"),
	}).Info("body is")
	instance, ok := LookupName(string(body))
	if !ok {
		return nil, fmt.Errorf("Couldn't find instance type %q %q", body, instance)
	}
	return instance, nil
}
