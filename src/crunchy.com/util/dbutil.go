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
	"crunchy.com/logutil"
	"database/sql"
	"os"
)

func GetConnection(database string) (*sql.DB, error) {
	var dbHost, dbUser, dbPort string
	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USER")
	dbPort = os.Getenv("DB_PORT")

	if dbHost == "" || dbUser == "" || dbPort == "" {
		logutil.Log("DB_HOST [" + dbHost + "]")
		logutil.Log("DB_USER [" + dbUser + "]")
		logutil.Log("DB_PORT [" + dbPort + "]")
		logutil.Log("error in getting required env vars")
		panic("could not get required env vars")
	}

	var dbConn *sql.DB
	var err error
	dbConn, err = sql.Open("postgres", "sslmode=disable user="+dbUser+" host="+dbHost+" port="+dbPort+" dbname="+database)
	if err != nil {
		logutil.Log(err.Error())
	}
	return dbConn, err
}

//used by monitoring
func GetMonitoringConnection(dbHost string, dbUser string, dbPort string, database string) (*sql.DB, error) {

	var dbConn *sql.DB
	var err error
	logutil.Log("open db with dbHost=[" + dbHost + "] dbUser=[" + dbUser + "] dbPort=[" + dbPort + "] database=[" + database + "]")
	dbConn, err = sql.Open("postgres", "sslmode=disable user="+dbUser+" host="+dbHost+" port="+dbPort+" dbname="+database)
	if err != nil {
		logutil.Log("error in getting connection :" + err.Error())
	}
	return dbConn, err
}
