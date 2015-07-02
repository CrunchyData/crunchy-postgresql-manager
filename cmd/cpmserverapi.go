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
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
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

	logit.Info.Println("serveragent starting")

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		&rest.Route{"POST", "/metrics/iostat", cpmserverapi.MetricIostat},
		&rest.Route{"POST", "/metrics/df", cpmserverapi.MetricDf},
		&rest.Route{"POST", "/metrics/mem", cpmserverapi.MetricMEM},
		&rest.Route{"POST", "/metrics/cpu", cpmserverapi.MetricCPU},
		&rest.Route{"POST", "/docker/inspect", cpmserverapi.DockerInspect},
		&rest.Route{"POST", "/docker/remove", cpmserverapi.DockerRemove},
		&rest.Route{"POST", "/docker/start", cpmserverapi.DockerStart},
		&rest.Route{"POST", "/docker/stop", cpmserverapi.DockerStop},
		&rest.Route{"POST", "/docker/run", cpmserverapi.DockerRun},
		&rest.Route{"POST", "/disk/provision", cpmserverapi.DiskProvision},
		&rest.Route{"POST", "/disk/delete", cpmserverapi.DiskDelete},
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":10001", api.MakeHandler()))
	log.Fatal(http.ListenAndServeTLS(":10000", "/var/cpm/keys/cert.pem", "/var/cpm/keys/key.pem", api.MakeHandler()))
}
