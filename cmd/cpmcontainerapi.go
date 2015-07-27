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
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"log"
	"net/http"
)

func init() {
	fmt.Println("before parsing in init")
	flag.Parse()
}

var CPMDIR = "/var/cpm/"
var CPMBIN = CPMDIR + "bin/"

func main() {

	logit.Info.Println("containeragent starting")

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		&rest.Route{"POST", "/remotewritefile", cpmcontainerapi.RemoteWritefile},
		&rest.Route{"POST", "/seed", cpmcontainerapi.Seed},
		&rest.Route{"POST", "/startpg", cpmcontainerapi.StartPG},
		&rest.Route{"POST", "/stoppg", cpmcontainerapi.StopPG},
		&rest.Route{"POST", "/startpgonstandby", cpmcontainerapi.StartPGOnStandby},
		&rest.Route{"POST", "/initdb", cpmcontainerapi.Initdb},
		&rest.Route{"POST", "/startpgpool", cpmcontainerapi.StartPgpool},
		&rest.Route{"POST", "/stoppgpool", cpmcontainerapi.StopPgpool},
		&rest.Route{"POST", "/basebackup", cpmcontainerapi.Basebackup},
		&rest.Route{"POST", "/failover", cpmcontainerapi.Failover},
		&rest.Route{"POST", "/controldata", cpmcontainerapi.Controldata},
		&rest.Route{"POST", "/badgergenerate", cpmcontainerapi.BadgerGenerate},
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("/pgdata/pg_log"))))

	log.Fatal(http.ListenAndServe(":10001", nil))
	log.Fatal(http.ListenAndServeTLS(":10000", "/var/cpm/keys/cert.pem", "/var/cpm/keys/key.pem", nil))
}
