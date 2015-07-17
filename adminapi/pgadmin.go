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

package adminapi

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"time"
)

func AdminStartpg(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println("AdminStartpg: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("AdminStartpg: error node ID required")
		rest.Error(w, "node ID required", http.StatusBadRequest)
		return
	}

	dbNode, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("AdminStartpg: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dbNode.Role == "pgpool" {
		var spgresp cpmcontainerapi.StartPgpoolResponse
		spgresp, err = cpmcontainerapi.StartPgpoolClient(dbNode.Name)
		logit.Info.Println("AdminStartpg:" + spgresp.Output)
	} else {
		var srep cpmcontainerapi.StartPGResponse
		srep, err = cpmcontainerapi.StartPGClient(dbNode.Name)
		logit.Info.Println("AdminStartpg:" + srep.Output)
	}

	if err != nil {
		logit.Error.Println("AdminStartpg:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//give the UI a chance to see the start
	time.Sleep(3000 * time.Millisecond)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AdminStoppg(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println("AdminStoppg: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logit.Info.Println("AdminStoppg:called")
	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("AdminStoppg:ID not found error")
		rest.Error(w, "node ID required", http.StatusBadRequest)
		return
	}

	var dbNode admindb.Container
	dbNode, err = admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("AdminStartpg: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logit.Info.Println("AdminStoppg: in stop with dbnode")

	if dbNode.Role == "pgpool" {
		var stoppoolResp cpmcontainerapi.StopPgpoolResponse
		stoppoolResp, err = cpmcontainerapi.StopPgpoolClient(dbNode.Name)
		logit.Info.Println("AdminStoppg:" + stoppoolResp.Output)
	} else {
		var stoppgResp cpmcontainerapi.StopPGResponse
		stoppgResp, err = cpmcontainerapi.StopPGClient(dbNode.Name)
		logit.Info.Println("AdminStoppg:" + stoppgResp.Output)
	}
	if err != nil {
		logit.Error.Println("AdminStoppg:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//give the UI a chance to see the stop
	time.Sleep(3000 * time.Millisecond)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}
