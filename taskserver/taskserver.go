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
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/task"
	"log"
	"net/http"
	"time"
)

func main() {

	logit.Info.Println("sleeping during startup to give DNS a chance")
	sleepTime, _ := time.ParseDuration("7s")
	time.Sleep(sleepTime)

	task.LoadSchedules()

	logit.Info.Println("taskserver starting")

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		&rest.Route{"POST", "/status/add", task.StatusAdd},
		&rest.Route{"POST", "/status/update", task.StatusUpdate},
		&rest.Route{"POST", "/executenow", task.ExecuteNow},
		&rest.Route{"POST", "/reload", task.Reload},
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	log.Fatal(http.ListenAndServe(":13001", nil))

}
