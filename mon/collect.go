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
	"database/sql"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmagent"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/myinfluxdb/client"
	"strconv"
	"strings"
	"time"
)

var CPMBIN = "/opt/cpm/bin/"

type DBMetric struct {
	MetricType string
	Name       string
	Value      float64
	Timestamp  time.Time
}

func collectServerMetrics(metricName string, server string) (DBMetric, error) {
	var values DBMetric
	var err error

	logit.Info.Println("collecting metric " + metricName + " for server " + server)
	switch metricName {
	case "cpu":
		logit.Info.Println("cpu collecting ")
		values, err = cpu(server)
	case "mem":
		logit.Info.Println("mem collecting ")
		values, err = mem(server)
	default:
		logit.Info.Println(metricName + " not implemented yet ")
	}
	if err != nil {
		logit.Error.Println("error in collecting " + metricName + " " + err.Error())
		return values, err
	}

	return values, nil
}

func collectContainerMetrics(metricName string, databaseConn *sql.DB) ([]DBMetric, error) {
	var values []DBMetric
	var err error

	logit.Info.Println("collecting metric..." + metricName)
	switch metricName {
	case "pg1":
		values, err = pg1(databaseConn)
	case "pg2":
		values, err = pg2(databaseConn)
	default:
		logit.Info.Println(metricName + " not implemented yet ")
	}

	if err != nil {
		logit.Error.Println("error in collecting " + metricName + " " + err.Error())
		return values, err
	}

	return values, nil
}

//dummy random value
func pg1(databaseConn *sql.DB) ([]DBMetric, error) {
	values := make([]DBMetric, 1)
	var err error
	values[0].Timestamp = time.Now()
	values[0].Name = "a"
	values[0].MetricType = "pg1"
	var intValue int

	err = databaseConn.QueryRow(fmt.Sprintf("select trunc(random() * 10 + 1) from  generate_series(1,1)")).Scan(&intValue)
	if err != nil {
		logit.Error.Println("pg1:error:" + err.Error())
		return values, err
	}
	values[0].Value = float64(intValue)

	return values, err
}

//database size in megabytes
func pg2(databaseConn *sql.DB) ([]DBMetric, error) {
	values := []DBMetric{}

	//thisTime := time.Now()
	var intValue int
	var databaseName string

	rows, err := databaseConn.Query("select datname, pg_database_size(d.oid)/1024/1024 from pg_database d")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := DBMetric{}
		if err = rows.Scan(
			&databaseName,
			&intValue); err != nil {
			return nil, err
		}
		m.Name = databaseName
		//m.Timestamp = thisTime
		m.Timestamp = time.Now()
		m.MetricType = "pg2"
		m.Value = float64(intValue)
		values = append(values, m)
	}
	if err = rows.Err(); err != nil {
		logit.Error.Println("pg2:error:" + err.Error())
		return nil, err
	}

	return values, err
}

//last minute load average of all cpu(s)
//range returned is 0.00 - 0.99
func cpu(server string) (DBMetric, error) {
	values := DBMetric{}
	var err error
	values.Timestamp = time.Now()
	values.MetricType = "cpu"
	var output string

	output, err = cpmagent.AgentCommand(CPMBIN+"monitor-load", "", server)
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
func mem(server string) (DBMetric, error) {
	values := DBMetric{}
	var err error
	values.Timestamp = time.Now()
	values.MetricType = "mem"
	var output string

	output, err = cpmagent.AgentCommand(CPMBIN+"monitor-mem", "", server)
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

func hc1(scheduleTS int64, nodeName string, databaseConn *sql.DB, c *client.Client) {
	var err error
	var strValue string

	err = databaseConn.QueryRow(fmt.Sprintf("select now()::text")).Scan(&strValue)
	if err != nil {
		logit.Error.Println(err.Error())
		//hc1 - database down condition
		series := &client.Series{
			Name:    "hc1",
			Columns: []string{"seconds", "service", "servicetype", "status"},
			Points: [][]interface{}{
				{scheduleTS, nodeName, "db", "down"},
			},
		}
		if err = c.WriteSeries([]*client.Series{series}); err != nil {
			logit.Error.Println("hc1 error writing to influxdb " + err.Error())
		}

	} else {
		//hc1 - database up condition
		series := &client.Series{
			Name:    "hc1",
			Columns: []string{"seconds", "service", "servicetype", "status"},
			Points: [][]interface{}{
				{scheduleTS, nodeName, "db", "up"},
			},
		}
		if err = c.WriteSeries([]*client.Series{series}); err != nil {
			logit.Error.Println("hc1 error writing to influxdb " + err.Error())
		}
	}

}
