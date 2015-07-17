/*
 Copyright 2015 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific langauge governing permissions and
 limitations under the License.
*/

package collect

import (
	"database/sql"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type DBMetric struct {
	MetricType string
	Name       string
	Value      float64
	Timestamp  time.Time
}

func CollectDBSize(gauge *prometheus.GaugeVec) error {
	var dbConn *sql.DB
	var err error

	dbConn, err = util.GetConnection("clusteradmin")
	if err != nil {
		logit.Error.Println(err.Error())
	}
	defer dbConn.Close()

	var domain string
	domain, err = getDomain(dbConn)
	var pgport string
	pgport, err = getPort(dbConn)

	//get all containers
	var containers []admindb.Container
	containers, err = admindb.GetAllContainers(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
	}

	//for each container, collect db size metrics
	i := 0

	for i = range containers {
		//containers[i].ProjectID
		//containers[i].ProjectName
		//containers[i].Name
		//containers[i].ID
		//containers[i].Role
		//containers[i].Image
		logit.Info.Println("dbsize processing " + containers[i].Name)
		err = process(gauge, dbConn, pgport, containers[i].Name, domain, containers[i].Role)
		if err != nil {
			logit.Error.Println(err.Error())
		}

		i++
	}

	return nil
}

func process(gauge *prometheus.GaugeVec, dbConn *sql.DB, port string, containerName string, domain string, containerRole string) error {
	var err error
	var userid, password, database string

	//get node credentials
	userid, password, database, err = getCredential(dbConn, containerName, containerRole)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("dbsize credentials userid:" + userid + " password:" + password + " database:" + database)
	var db *sql.DB
	db, err = util.GetMonitoringConnection(containerName,
		userid, port, database, password)
	defer db.Close()
	if err != nil {
		logit.Error.Println("error in getting connectionto " + containerName)
		return err
	}

	var metrics []DBMetric
	logit.Info.Println("dbsize running pg2 on " + containerName)
	metrics, err = pg2(db)

	//write metrcs to prometheus

	i := 0
	for i = range metrics {
		logit.Info.Println("dbsize setting dbsize metric")
		gauge.WithLabelValues(containerName, metrics[i].Name).Set(metrics[i].Value)
		i++
	}

	return nil

}

//database size in megabytes
func pg2(databaseConn *sql.DB) ([]DBMetric, error) {
	values := []DBMetric{}

	//thisTime := time.Now()
	var intValue int
	var databaseName string

	rows, err := databaseConn.Query("select datname, pg_database_size(d.oid)/1024/1024 from pg_database d")
	if err != nil {
		logit.Error.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		logit.Info.Println("dbsize pg2 got row")
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
