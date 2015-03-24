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
	"crunchy.com/backup"
	"crunchy.com/mylog"
	"crunchy.com/util"
	"database/sql"
	"flag"
	"log"
	//"github.com/golang/glog"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

var logger mylog.MyLogger

var LOGPATH = flag.String("logpath", "backupserver.log", "log directory to write to")

func main() {
	flag.Parse()

	log.Println("logging to " + *LOGPATH)

	logger = mylog.Init(*LOGPATH)
	backup.SetLogger(logger)
	admindb.SetLogger(logger)

	logger.Info.Println("sleeping during startup to give DNS a chance")
	time.Sleep(time.Millisecond * 7000)

	found := false
	var dbConn *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		dbConn, err = util.GetConnection("clusteradmin")
		if err != nil {
			logger.Error.Println(err.Error())
			logger.Error.Println("could not get initial database connection, will retry in 5 seconds")
			time.Sleep(time.Millisecond * 5000)
		} else {
			logger.Info.Println("got db connection")
			found = true
			break
		}
	}

	if !found {
		panic("could not connect to clusteradmin db")
	}

	backup.SetConnection(dbConn)
	admindb.SetConnection(dbConn)

	backup.LoadSchedules()

	logger.Info.Println("starting\n")
	command := new(backup.Command)
	rpc.Register(command)
	logger.Info.Println("Command registered\n")
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":13000")
	logger.Info.Println("listening\n")
	if e != nil {
		logger.Error.Println(e.Error())
		panic("could not listen on rpc socker")
	}
	logger.Info.Println("about to serve\n")
	http.Serve(l, nil)
	logger.Info.Println("after serve\n")
}
