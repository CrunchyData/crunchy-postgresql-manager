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
	"github.com/crunchydata/crunchy-postgresql-manager/types"
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

// CollectDBSize collect and persist the database size metric into the prometheus db
func CollectDBSize(gauge *prometheus.GaugeVec) error {
	var dbConn *sql.DB
	var err error

	dbConn, err = util.GetConnection("clusteradmin")
	if err != nil {
		logit.Error.Println(err.Error())
	}
	defer dbConn.Close()

	//var domain string
	//domain, err = getDomain(dbConn)

	//get all containers
	var containers []types.Container
	containers, err = admindb.GetAllContainers(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
	}

	//for each container, collect db size metrics
	i := 0
	var credential types.Credential

	for i = range containers {
		//logit.Info.Println("dbsize processing " + containers[i].Name)
		credential, err = admindb.GetUserCredentials(dbConn, &containers[i])
		if err != nil {
			logit.Error.Println(err.Error())
		}
		err = process(&containers[i], &credential, gauge)
		if err != nil {
			logit.Error.Println(err.Error())
		}

		i++
	}

	return nil
}

func process(node *types.Container, credential *types.Credential, gauge *prometheus.GaugeVec) error {
	var err error

	//logit.Info.Println("dbsize node=" + node.Name + " credentials Username:" + credential.Username + " Password:" + credential.Password + " Database:" + credential.Database + " Host:" + credential.Host)
	var db *sql.DB
	db, err = util.GetMonitoringConnection(credential.Host,
		credential.Username, credential.Port, credential.Database, credential.Password)
	defer db.Close()
	if err != nil {
		logit.Error.Println("error in getting connectionto " + credential.Host)
		return err
	}

	var metrics []DBMetric
	//logit.Info.Println("dbsize running pg2 on " + node.Name)
	metrics, err = pg2(db)

	//write metrcs to prometheus

	i := 0
	for i = range metrics {
		//logit.Info.Println("dbsize setting dbsize metric")
		gauge.WithLabelValues(node.Name, metrics[i].Name).Set(metrics[i].Value)
		i++
	}

	return nil

}

// pg2 calculate the database size in megabytes and return the metrics
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
		//logit.Info.Println("dbsize pg2 got row")
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
		logit.Error.Println(err.Error())
		return nil, err
	}

	return values, err
}
