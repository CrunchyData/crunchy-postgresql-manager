package collect

import (
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
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

	var response cpmserverapi.MetricCPUResponse
	request := &cpmserverapi.MetricCPURequest{}
	var url = "http://" + serverName + ":10001"
	response, err = cpmserverapi.MetricCPUClient(url, request)
	if err != nil {
		logit.Error.Println("cpu metric error:" + err.Error())
		return values, err
	}

	var output = strings.TrimSpace(response.Output)

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

	var response cpmserverapi.MetricMEMResponse
	request := &cpmserverapi.MetricMEMRequest{}
	var url = "http://" + serverName + ":10001"
	response, err = cpmserverapi.MetricMEMClient(url, request)
	if err != nil {
		logit.Error.Println("mem metric error:" + err.Error())
		return values, err
	}

	output = strings.TrimSpace(response.Output)

	values.Value, err = strconv.ParseFloat(output, 64)
	if err != nil {
		logit.Error.Println("parseFloat error in mem metric " + err.Error())
	}

	return values, err
}
