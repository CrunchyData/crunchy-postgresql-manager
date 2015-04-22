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
	"database/sql"
	"flag"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/mon"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

func main() {
	flag.Parse()

	logit.Info.Println("sleeping during startup to give DNS a chance")
	time.Sleep(time.Millisecond * 7000)

	//verify cpm db exists in influxdb
	mon.Bootdb()

	found := false
	var dbConn *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		dbConn, err = util.GetConnection("clusteradmin")
		if err != nil {
			logit.Error.Println(err.Error())
			logit.Error.Println("could not get initial database connection, will retry in 5 seconds")
			time.Sleep(time.Millisecond * 5000)
		} else {
			logit.Info.Println("got db connection")
			found = true
			break
		}
	}

	if !found {
		logit.Error.Println("could not connect to clusteradmin db")
		panic("could not connect to clusteradmin db")
	}

	mon.SetConnection(dbConn)
	admindb.SetConnection(dbConn)

	mon.LoadSchedules()

	dbConn.Close()

	logit.Info.Println("starting\n")
	command := new(mon.Command)
	rpc.Register(command)
	logit.Info.Println("Command registered\n")
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":13000")
	logit.Info.Println("listening\n")
	if e != nil {
		logit.Error.Println(e.Error())
		panic("could not listen on rpc socker")
	}
	logit.Info.Println("about to serve\n")
	http.Serve(l, nil)
	logit.Info.Println("after serve\n")
}

func bootstrapInflux() {
	//connect
	//verify
}
