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
	"crunchy.com/cpmagent"
	"database/sql"
	"fmt"
	"github.com/golang/glog"
	"strconv"
	"strings"
	"time"
)

type DBMetric struct {
	MetricType string
	Value      float64
	Timestamp  time.Time
}

func collectServerMetrics(metricName string, server string) (DBMetric, error) {
	var values DBMetric
	var err error

	glog.Infoln("collecting metric " + metricName + " for server " + server)
	switch metricName {
	case "cpu":
		glog.Infoln("cpu collecting ")
		values, err = cpu(server)
	case "mem":
		glog.Infoln("mem collecting ")
		values, err = mem(server)
	default:
		glog.Infoln(metricName + " not implemented yet ")
	}
	if err != nil {
		glog.Errorln("error in collecting " + metricName + " " + err.Error())
		return values, err
	}

	return values, nil
}

func collectContainerMetrics(metricName string, databaseConn *sql.DB) (DBMetric, error) {
	var values DBMetric
	var err error

	glog.Infoln("collecting metric..." + metricName)
	switch metricName {
	case "pg1":
		values, err = pg1(databaseConn)
	case "pg2":
		values, err = pg2(databaseConn)
	default:
		glog.Infoln(metricName + " not implemented yet ")
	}

	if err != nil {
		glog.Errorln("error in collecting " + metricName + " " + err.Error())
		return values, err
	}

	return values, nil
}

func pg1(databaseConn *sql.DB) (DBMetric, error) {
	values := DBMetric{}
	var err error
	values.Timestamp = time.Now()
	values.MetricType = "pg1"
	var intValue int

	err = databaseConn.QueryRow(fmt.Sprintf("select trunc(random() * 10 + 1) from  generate_series(1,1)")).Scan(&intValue)
	if err != nil {
		glog.Errorln("pg1:error:" + err.Error())
		return values, err
	}
	values.Value = float64(intValue)

	return values, err
}
func pg2(databaseConn *sql.DB) (DBMetric, error) {
	values := DBMetric{}
	var err error
	values.Timestamp = time.Now()
	values.MetricType = "pg2"
	var intValue int

	err = databaseConn.QueryRow(fmt.Sprintf("select trunc(random() * 10 + 1) from  generate_series(1,1)")).Scan(&intValue)
	if err != nil {
		glog.Errorln("pg1:error:" + err.Error())
		return values, err
	}
	values.Value = float64(intValue)

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

	output, err = cpmagent.AgentCommand("/cluster/bin/monitor-load", "", server)
	if err != nil {
		glog.Errorln("cpu metric error:" + err.Error())
		return values, err
	}

	output = strings.TrimSpace(output)

	values.Value, err = strconv.ParseFloat(output, 64)
	if err != nil {
		glog.Errorln("parseFloat error in cpu metric " + err.Error())
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

	output, err = cpmagent.AgentCommand("/cluster/bin/monitor-mem", "", server)
	if err != nil {
		glog.Errorln("mem metric error:" + err.Error())
		return values, err
	}

	output = strings.TrimSpace(output)

	values.Value, err = strconv.ParseFloat(output, 64)
	if err != nil {
		glog.Errorln("parseFloat error in mem metric " + err.Error())
	}

	return values, err
}

func hc1(databaseConn *sql.DB) error {
	var err error
	var strValue string

	err = databaseConn.QueryRow(fmt.Sprintf("select now()::text")).Scan(&strValue)
	return err
}
