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
	"crunchy.com/cpmagent"
	"crunchy.com/logutil"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"time"
)

func AdminStartpg(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logutil.Log("AdminStartpg: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logutil.Log("AdminStartpg: error node ID required")
		rest.Error(w, "node ID required", 400)
		return
	}

	dbNode, err := admindb.GetDBNode(ID)
	if err != nil {
		logutil.Log("AdminStartpg: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	var output string
	output, err = cpmagent.AgentCommand("/cluster/bin/startpg.sh", "", dbNode.Name)
	if err != nil {
		logutil.Log("AdminStartpg:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	logutil.Log("AdminStartpg:" + output)

	//give the UI a chance to see the start
	time.Sleep(3000 * time.Millisecond)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AdminStoppg(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logutil.Log("AdminStoppg: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	logutil.Log("AdminStoppg:called")
	ID := r.PathParam("ID")
	if ID == "" {
		logutil.Log("AdminStoppg:ID not found error")
		rest.Error(w, "node ID required", 400)
		return
	}

	dbNode, err := admindb.GetDBNode(ID)
	if err != nil {
		logutil.Log("AdminStartpg: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	logutil.Log("AdminStoppg: in stop with dbnode")

	var output string
	output, err = cpmagent.AgentCommand("/cluster/bin/stoppg.sh", "", dbNode.Name)
	if err != nil {
		logutil.Log("AdminStoppg:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	logutil.Log("AdminStoppg:" + output)

	//give the UI a chance to see the stop
	time.Sleep(3000 * time.Millisecond)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}
