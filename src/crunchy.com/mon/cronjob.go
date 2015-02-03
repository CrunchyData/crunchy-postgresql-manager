/*
 Copyright 2015 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package mon

import (
	"crunchy.com/admindb"
	"crunchy.com/logutil"
	"crunchy.com/util"
	"database/sql"
	"github.com/myinfluxdb/client"
	"strconv"
)

type DefaultJob struct {
	request MonRequest
}

//this is the func that implements the cron Job interface
func (t DefaultJob) Run() {
	logutil.Log("running Schedule:" + t.request.Schedule.Name)
	RunMonJob(&t.request)
}

func RunMonJob(args *MonRequest) error {
	logutil.Log("mon.RunMonJob called")
	logutil.Log("with Schedule Name=" + args.Schedule.Name)

	c, err := client.NewClient(&client.ClientConfig{
		Username: "root",
		Password: "root",
		Database: "cpm",
	})
	if err != nil {
		logutil.Log("error in connection to fluxdb " + err.Error())
		return err
	}

	servers, nodes, metrics, err := getData()
	if err != nil {
		logutil.Log("error: RunMonJob " + err.Error())
		return err
	}

	var value DBMetric
	x := 0
	for x = range servers {
		//get connection to server
		logutil.Log("collecting for server " + servers[x].Name)
		//collect all the server metrics
		i := 0
		for i = range metrics {
			if metrics[i].ScheduleName == args.Schedule.Name {
				if metrics[i].MetricType == "server" {
					value, err = collectServerMetrics(metrics[i].Name, servers[x].IPAddress)
					//add value to influxdb here
					logutil.Log(metrics[i].Name + " value " + strconv.FormatFloat(value.Value, 'f', 3, 64))
					series := &client.Series{
						Name:    metrics[i].Name,
						Columns: []string{"value", "server"},
						Points: [][]interface{}{
							{value.Value, servers[x].Name},
						},
					}
					if err = c.WriteSeries([]*client.Series{series}); err != nil {
						logutil.Log("error writing to influxdb " + err.Error())
					}

				} else if metrics[i].MetricType == "healthck" {
					logutil.Log("healthck server metrics callled on " + args.Schedule.Name)
				}
			}
		}
	}

	y := 0
	for y = range nodes {
		//get connection to database
		logutil.Log("collecting for node " + nodes[y].Name)
		var databaseConn *sql.DB
		databaseConn, err = util.GetMonitoringConnection(nodes[y].Name+".crunchy.lab", "postgres", "5432", "postgres")
		if err != nil {
			logutil.Log("error in getting connection to " + nodes[y].Name)
		} else {
			//collect all the database metrics
			i := 0
			for i = range metrics {
				if metrics[i].ScheduleName == args.Schedule.Name {
					if metrics[i].MetricType == "database" {
						value, err = collectContainerMetrics(metrics[i].Name, databaseConn)
						logutil.Log(metrics[i].Name + " value " + strconv.FormatFloat(value.Value, 'f', 3, 64))
						//add value to influxdb here
						series := &client.Series{
							Name:    metrics[i].Name,
							Columns: []string{"value", "database"},
							Points: [][]interface{}{
								{value.Value, nodes[y].Name},
							},
						}
						if err = c.WriteSeries([]*client.Series{series}); err != nil {
							logutil.Log("error writing to influxdb " + err.Error())
						}
					} else if metrics[i].MetricType == "healthck" {
						logutil.Log("healthck metric database run on schedule " + args.Schedule.Name)
					}
				}
			}
		}
		databaseConn.Close()
	}

	return err

}

func getData() ([]admindb.DBServer, []admindb.DBClusterNode, []MonMetric, error) {
	var servers []admindb.DBServer
	var nodes []admindb.DBClusterNode
	var metrics []MonMetric
	var dbConn *sql.DB
	var err error

	dbConn, err = util.GetConnection("clusteradmin")
	if err != nil {
		logutil.Log(err.Error())
		return servers, nodes, metrics, err
	}
	defer dbConn.Close()

	admindb.SetConnection(dbConn)

	servers, err = admindb.GetAllDBServers()
	if err != nil {
		logutil.Log(err.Error())
		return servers, nodes, metrics, err
	}
	nodes, err = admindb.GetAllDBNodes()
	if err != nil {
		logutil.Log(err.Error())
		return servers, nodes, metrics, err
	}

	SetConnection(dbConn)
	metrics, err = DBGetMetrics()
	if err != nil {
		logutil.Log(err.Error())
		return servers, nodes, metrics, err
	}

	logutil.Log("got this many metrics " + strconv.Itoa(len(metrics)))
	return servers, nodes, metrics, err
}
