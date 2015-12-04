package main

import (
	"github.com/crunchydata/crunchy-postgresql-manager/collect"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/swarmapi"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//var gauges = make([]prometheus.Gauge, 3)
var (
	HC_POLL_INT     int64
	SERVER_POLL_INT int64
	CONT_POLL_INT   int64
)

func main() {

	var err error
	var tempVal string

	tempVal = os.Getenv("HC_POLL_INT")
	if tempVal == "" {
		HC_POLL_INT = 3
	}

	tempVal = os.Getenv("SERVER_POLL_INT")
	if tempVal == "" {
		SERVER_POLL_INT = 3
	}
	tempVal = os.Getenv("SERVER_POLL_INT")
	if tempVal == "" {
		CONT_POLL_INT = 3
	}
	HC_POLL_INT, err = strconv.ParseInt(tempVal, 10, 64)
	if err != nil {
		logit.Error.Println(err.Error())
		return
	}
	SERVER_POLL_INT, err = strconv.ParseInt(tempVal, 10, 64)
	if err != nil {
		logit.Error.Println(err.Error())
		return
	}
	CONT_POLL_INT, err = strconv.ParseInt(tempVal, 10, 64)
	if err != nil {
		logit.Error.Println(err.Error())
		return
	}
	logit.Info.Printf("HealthCheck Polling Interval: %d\n", HC_POLL_INT)
	logit.Info.Printf("Server Polling Interval: %d\n", SERVER_POLL_INT)
	logit.Info.Printf("Container Polling Interval: %d\n", CONT_POLL_INT)

	go func() {
		//register a gauge vector
		gauge := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cpm_server_cpu",
				Help: "CPU Utilization.",
			}, []string{
				"server",
			})
		prometheus.MustRegister(gauge)
		gaugeMem := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cpm_server_mem",
				Help: "Memory Utilization.",
			}, []string{
				"server",
			})
		prometheus.MustRegister(gaugeMem)

		var serverinfo swarmapi.DockerInfoResponse
		serverinfo, err = swarmapi.DockerInfo()
		if err != nil {
			logit.Error.Println(err.Error())
		}
		logit.Info.Printf("got %d servers\n", len(serverinfo.Output))

		var metric collect.Metric

		for {
			//get metrics for each server
			for x := 0; x < len(serverinfo.Output); {
				//v := rand.Float64() * 100.00
				parts := strings.Split(serverinfo.Output[x], ":")
				ipaddress := parts[0]
				metric, err = collect.Collectcpu(ipaddress)
				gauge.WithLabelValues(ipaddress).Set(metric.Value)
				logit.Info.Println("setting cpu metric for " + ipaddress + " to " + strconv.FormatFloat(metric.Value, 'f', -1, 64))
				metric, err = collect.Collectmem(ipaddress)
				gaugeMem.WithLabelValues(ipaddress).Set(metric.Value)
				logit.Info.Println("setting mem metric for " + ipaddress + " to " + strconv.FormatFloat(metric.Value, 'f', -1, 64))
				x++
			}

			logit.Info.Println("sleeping in the server metrics ....")

			time.Sleep(time.Duration(SERVER_POLL_INT) * time.Minute)
		}
	}()

	go func() {
		for true {
			collect.Collecthc()
			time.Sleep(time.Duration(HC_POLL_INT) * time.Minute)
		}
	}()

	go func() {
		//register a gauge vector
		dbsizegauge := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cpm_container_dbsize",
				Help: "Database Size.",
			}, []string{
				"containername",
				"databasename",
			})
		prometheus.MustRegister(dbsizegauge)

		for true {
			//dbsizegauge.WithLabelValues("node1", "db1").Set(v)
			logit.Info.Println("collecting dbsize")
			collect.CollectDBSize(dbsizegauge)
			time.Sleep(time.Duration(CONT_POLL_INT) * time.Minute)
		}
	}()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)
}
