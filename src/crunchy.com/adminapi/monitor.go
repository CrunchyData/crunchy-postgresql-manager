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

package main

import (
	"crunchy.com/admindb"
	"crunchy.com/util"
	"database/sql"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
	"github.com/influxdb/influxdb/client"
	_ "github.com/lib/pq"
	"net/http"
)

type MyPoint struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}
type PG2Data struct {
	Color string
	Data  []MyPoint
	Name  string
}

func GetServerMetrics(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetServerMetrics: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	serverID := r.PathParam("ServerID")
	if serverID == "" {
		glog.Errorln("GetServerMetrics: error serverID required")
		rest.Error(w, "serverID required", http.StatusBadRequest)
		return
	}

	interval := r.PathParam("Interval")
	if interval == "" {
		glog.Errorln("GetServerMetrics: error Interval required")
		rest.Error(w, "Interval required", http.StatusBadRequest)
		return
	}

	metric := r.PathParam("Metric")
	if interval == "" {
		glog.Errorln("GetServerMetrics: error Metric required")
		rest.Error(w, "Metric required", http.StatusBadRequest)
		return
	}

	server := admindb.DBServer{}
	server, err = admindb.GetDBServer(serverID)
	if err != nil {
		glog.Errorln("GetServerCPUMetrics: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var domain string
	domain, err = admindb.GetDomain()
	if err != nil {
		glog.Errorln("GetServerCPUMetrics: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var hostname = "cpm-mon"
	if kubeEnv {
		hostname = hostname + "-api"
	}
	c, err := client.NewClient(&client.ClientConfig{
		Host:     hostname + "." + domain + ":8086",
		Username: "root",
		Password: "root",
		Database: "cpm",
	})

	if err != nil {
		glog.Errorln("GetServerCPUMetrics: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var results []*client.Series

	var query = "select value, server from " + metric + " where server = '" + server.Name + "' and time > now() - " + interval + " order asc limit 1000"
	glog.Infoln(query)

	results, err = c.Query(query)
	if err != nil {
		glog.Errorln("GetServerCPUMetrics: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(results) > 0 {
		glog.Infof("results len = %d\n", len(results[0].Points))
		glog.Infof("results x=%f y=%f\n", results[0].Points[0][0], results[0].Points[0][3])
	}
	w.WriteJson(&results)

}

//get database sizes in a container
func GetPG2(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetPG2: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	Name := r.PathParam("Name")
	if Name == "" {
		glog.Errorln("GetPG2: error Name required")
		rest.Error(w, "Name required", http.StatusBadRequest)
		return
	}

	interval := r.PathParam("Interval")
	if interval == "" {
		glog.Errorln("GetPG2: error Interval required")
		rest.Error(w, "Interval required", http.StatusBadRequest)
		return
	}

	//get domain name
	var domain string
	var hostname = "cpm-mon"

	domain, err = admindb.GetDomain()

	if kubeEnv {
		hostname = hostname + "-api"
	}
	c, err := client.NewClient(&client.ClientConfig{
		Host:     hostname + "." + domain + ":8086",
		Username: "root",
		Password: "root",
		Database: "cpm",
	})

	if err != nil {
		glog.Errorln("GetPG2: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get list of databases on node
	var databaseConn *sql.DB
	databaseConn, err = util.GetMonitoringConnection(Name+"."+domain, "postgres", "5432", "postgres", "")
	defer databaseConn.Close()

	var databases []string
	databases, err = GetAllDatabases(databaseConn)
	if err != nil {
		glog.Errorln("GetPG2: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	glog.Infof("databases len = %d\n", len(databases))

	bigresults := make([]PG2Data, 0)
	var results []*client.Series

	y := 0
	for y = range databases {
		var query = "select value, database from pg2 where database = '" + databases[y] + "' and container = '" + Name + "' and time > now() - " + interval + " order asc limit 1000"
		pgdata := PG2Data{}

		glog.Infoln(query)

		results, err = c.Query(query)

		points := make([]MyPoint, 0)

		pgdata.Color = "#c05020"
		pgdata.Name = databases[y]

		if len(results) > 0 {
			glog.Infof("results len = %d\n", len(results[0].Points))
			glog.Infof("results x=%f y=%f\n", results[0].Points[0][0], results[0].Points[0][3])
			i := 0
			for i = range results[0].Points {
				pt := MyPoint{}
				pt.X = int64(results[0].Points[i][0].(float64) / 1000)
				pt.Y = results[0].Points[i][3].(float64)

				points = append(points, pt)
			}
		} else {
			glog.Infof("results len = 0 for database %s\n", databases[y])
		}
		pgdata.Data = points
		bigresults = append(bigresults, pgdata)

	}

	if len(bigresults) == 0 {
		glog.Errorln("GetPG2: no data found")
		rest.Error(w, "GetPG2:  no data found", http.StatusBadRequest)
		return
	} else {
		glog.Infof("bigresults len = %d\n", len(bigresults))
	}

	w.WriteJson(&bigresults)

}

func GetAllDatabases(conn *sql.DB) ([]string, error) {
	glog.Infoln("GetAllDatabases: called")
	m := make([]string, 0)

	var rows *sql.Rows
	var err error
	rows, err = conn.Query("select datname from pg_database")
	if err != nil {
		return m, err
	}
	defer rows.Close()
	var value string
	for rows.Next() {
		if err = rows.Scan(
			&value); err != nil {
			return m, err
		}
		m = append(m, value)
	}
	if err = rows.Err(); err != nil {
		return m, err
	}
	return m, nil
}
