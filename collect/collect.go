package collect

import (
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserveragent"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"strconv"
	"strings"
	"time"
)

type Metric struct {
	MetricType string
	Name       string
	Value      float64
	Timestamp  time.Time
}

//last minute load average of all cpu(s)
//range returned is 0.00 - 0.99
func Collectcpu(serverName string) (Metric, error) {
	values := Metric{}
	var err error
	values.Timestamp = time.Now()
	values.MetricType = "cpu"
	var output string

	output, err = cpmserveragent.AgentCommand("monitor-load", "", serverName)
	if err != nil {
		logit.Error.Println("cpu metric error:" + err.Error())
		return values, err
	}

	output = strings.TrimSpace(output)

	values.Value, err = strconv.ParseFloat(output, 64)
	if err != nil {
		logit.Error.Println("parseFloat error in cpu metric " + err.Error())
	}

	return values, err
}

//percent memory utilization
//range returned is 1-100
func Collectmem(serverName string) (Metric, error) {
	values := Metric{}
	var err error
	values.Timestamp = time.Now()
	values.MetricType = "mem"
	var output string

	output, err = cpmserveragent.AgentCommand("monitor-mem", "", serverName)
	if err != nil {
		logit.Error.Println("mem metric error:" + err.Error())
		return values, err
	}

	output = strings.TrimSpace(output)

	values.Value, err = strconv.ParseFloat(output, 64)
	if err != nil {
		logit.Error.Println("parseFloat error in mem metric " + err.Error())
	}

	return values, err
}
