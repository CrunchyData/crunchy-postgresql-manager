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
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	_ "github.com/lib/pq"
)

type MonSchedule struct {
	Name    string
	Cronexp string
}

type MonMetric struct {
	Name         string
	MetricType   string
	ScheduleName string
}

var dbConn *sql.DB

func SetConnection(conn *sql.DB) {
	logit.Info.Println("mondb:SetConnection: called to open dbConn")
	dbConn = conn
}

func DBGetSchedules() ([]MonSchedule, error) {
	logit.Info.Println("DBGetSchedules called")
	var rows *sql.Rows
	var err error

	rows, err = dbConn.Query(fmt.Sprintf("select name, cronexp from monschedule"))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	schedules := make([]MonSchedule, 0)
	for rows.Next() {
		s := MonSchedule{}
		if err = rows.Scan(
			&s.Name,
			&s.Cronexp); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func DBGetMetrics() ([]MonMetric, error) {
	logit.Info.Println("DBGetMetrics called")
	var rows *sql.Rows
	var err error

	rows, err = dbConn.Query(fmt.Sprintf("select name, metrictype, schedule from monmetric"))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	metrics := make([]MonMetric, 0)
	for rows.Next() {
		s := MonMetric{}
		if err = rows.Scan(
			&s.Name,
			&s.MetricType,
			&s.ScheduleName); err != nil {
			return nil, err
		}
		metrics = append(metrics, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}
