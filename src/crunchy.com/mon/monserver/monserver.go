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
	"crunchy.com/logutil"
	"crunchy.com/mon"
	"crunchy.com/util"
	"database/sql"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

func main() {

	logutil.Log("sleeping during startup to give DNS a chance")
	time.Sleep(time.Millisecond * 7000)

	//verify cpm db exists in influxdb
	mon.Bootdb()

	found := false
	var dbConn *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		dbConn, err = util.GetConnection("clusteradmin")
		if err != nil {
			logutil.Log(err.Error())
			logutil.Log("could not get initial database connection, will retry in 5 seconds")
			time.Sleep(time.Millisecond * 5000)
		} else {
			logutil.Log("got db connection")
			found = true
			break
		}
	}

	if !found {
		panic("could not connect to clusteradmin db")
	}

	mon.SetConnection(dbConn)

	mon.LoadSchedules()

	dbConn.Close()

	logutil.Log("starting\n")
	command := new(mon.Command)
	rpc.Register(command)
	logutil.Log("Command registered\n")
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":13000")
	logutil.Log("listening\n")
	if e != nil {
		logutil.Log(e.Error())
		panic("could not listen on rpc socker")
	}
	logutil.Log("about to serve\n")
	http.Serve(l, nil)
	logutil.Log("after serve\n")
}

func bootstrapInflux() {
	//connect
	//verify
}
