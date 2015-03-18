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

package util

import (
	"database/sql"
	"github.com/golang/glog"
	"os"
)

func GetConnection(database string) (*sql.DB, error) {
	var dbHost, dbUser, dbPort, dbPassword string
	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USER")
	dbPort = os.Getenv("DB_PORT")
	dbPassword = os.Getenv("DB_PASSWORD")

	if dbHost == "" || dbUser == "" || dbPort == "" {
		glog.Errorln("DB_HOST [" + dbHost + "]")
		glog.Errorln("DB_USER [" + dbUser + "]")
		glog.Errorln("DB_PORT [" + dbPort + "]")
		glog.Errorln("error in getting required env vars")
		glog.Flush()
		panic("could not get required env vars")
	}

	var dbConn *sql.DB
	var err error
	if dbPassword != "" {
		glog.Infoln("connecting to database " + database + " host=" + dbHost + " user=" + dbUser + " port=" + dbPort + " password=" + dbPassword)
		dbConn, err = sql.Open("postgres", "sslmode=disable user="+dbUser+" host="+dbHost+" port="+dbPort+" dbname="+database+" password="+dbPassword)
	} else {
		glog.Infoln("connecting to database " + database + " host=" + dbHost + " user=" + dbUser + " port=" + dbPort)
		dbConn, err = sql.Open("postgres", "sslmode=disable user="+dbUser+" host="+dbHost+" port="+dbPort+" dbname="+database)
	}
	if err != nil {
		glog.Errorln(err.Error())
	}
	return dbConn, err
}

//used by monitoring
func GetMonitoringConnection(dbHost string, dbUser string, dbPort string, database string, dbPassword string) (*sql.DB, error) {

	var dbConn *sql.DB
	var err error

	if dbPassword == "" {
		glog.Infoln("open db with dbHost=[" + dbHost + "] dbUser=[" + dbUser + "] dbPort=[" + dbPort + "] database=[" + database + "]")
		dbConn, err = sql.Open("postgres", "sslmode=disable user="+dbUser+" host="+dbHost+" port="+dbPort+" dbname="+database)
	} else {
		glog.Infoln("open db with dbHost=[" + dbHost + "] dbUser=[" + dbUser + "] dbPort=[" + dbPort + "] database=[" + database + "] password=[" + dbPassword + "]")
		dbConn, err = sql.Open("postgres", "sslmode=disable user="+dbUser+" host="+dbHost+" port="+dbPort+" dbname="+database+" password="+dbPassword)
	}
	if err != nil {
		glog.Errorln("error in getting connection :" + err.Error())
	}
	return dbConn, err
}
