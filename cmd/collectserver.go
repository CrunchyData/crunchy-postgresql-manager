package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/collect"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	uniformDomain     = 200
	normDomain        = 200
	normMean          = 10
	oscillationPeriod = 10 * time.Minute
)

//var guages = make([]prometheus.Gauge, 3)

func main() {

	go func() {
		//register a guage vector
		guage := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cpm_server_cpu",
				Help: "CPU Utilization.",
			}, []string{
				"server",
			})
		prometheus.MustRegister(guage)

		//get servers
		var dbConn *sql.DB
		var err error
		dbConn, err = util.GetConnection("clusteradmin")
		if err != nil {
			logit.Error.Println(err.Error())
		}
		var servers []admindb.Server
		admindb.SetConnection(dbConn)
		servers, err = admindb.GetAllServers()
		if err != nil {
			logit.Error.Println(err.Error())
		}
		dbConn.Close()

		var metric collect.Metric
		for {
			//get metrics for each server
			i := 0
			for i = range servers {
				//v := rand.Float64() * 100.00
				metric, err = collect.Collectcpu(servers[i].IPAddress)
				guage.WithLabelValues(servers[i].Name).Set(metric.Value)
				logit.Info.Println("setting metric for " + servers[i].Name + " to " + strconv.FormatFloat(metric.Value, 'f', -1, 64))
				i++
			}

			time.Sleep(time.Duration(10000 * time.Millisecond))
		}
	}()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)
}
