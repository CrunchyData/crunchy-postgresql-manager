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
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"strconv"
	"strings"
)

func GetServer(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetServer: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")
	logit.Info.Println("in GetServer with ID=" + ID)

	results, err := admindb.GetServer(dbConn, ID)
	if err != nil {
		logit.Error.Println("GetServer:" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server := Server{results.ID, results.Name, results.IPAddress,
		results.DockerBridgeIP, results.PGDataPath, results.ServerClass, results.CreateDate}
	logit.Info.Println("GetServer: results=" + results.ID)

	w.WriteJson(&server)
}

//we use AddServer for both updating and inserting based on the ID passed in
func AddServer(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-server")
	if err != nil {
		logit.Error.Println("AddServer: authorize token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	CreateDate := ""
	server := Server{r.PathParam("ID"), r.PathParam("Name"), r.PathParam("IPAddress"), r.PathParam("DockerBridgeIP"), r.PathParam("PGDataPath"), r.PathParam("ServerClass"), CreateDate}

	server.IPAddress = strings.Replace(server.IPAddress, "_", ".", -1)
	server.DockerBridgeIP = strings.Replace(server.DockerBridgeIP, "_", ".", -1)
	server.PGDataPath = strings.Replace(server.PGDataPath, "_", "/", -1)

	if server.Name == "" {
		logit.Error.Println("AddServer: error server name required")
		rest.Error(w, "server name required", http.StatusBadRequest)
		return
	}
	if server.IPAddress == "" {
		logit.Error.Println("AddServer: error ipaddress required")
		rest.Error(w, "server IPAddress required", http.StatusBadRequest)
		return
	}
	if server.PGDataPath == "" {
		logit.Error.Println("AddServer: error pgdatapath required")
		rest.Error(w, "server PGDataPath required", http.StatusBadRequest)
		return
	}

	dbserver := admindb.Server{server.ID, server.Name, server.IPAddress,
		server.DockerBridgeIP, server.PGDataPath, server.ServerClass, CreateDate, ""}
	if dbserver.ID == "0" {
		strid, err := admindb.InsertServer(dbConn, dbserver)
		newid := strconv.Itoa(strid)
		if err != nil {
			logit.Error.Println("AddServer:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		server.ID = newid
	} else {
		admindb.UpdateServer(dbConn, dbserver)
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&server)
}

//monitor server - get info
func MonitorServerGetInfo(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("MonitorServerGetInfo: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ServerID := r.PathParam("ServerID")
	Metric := r.PathParam("Metric")

	if ServerID == "" {
		logit.Error.Println("MonitorServerGetInfo: error ServerID required")
		rest.Error(w, "ServerID required", http.StatusBadRequest)
		return
	}
	if Metric == "" {
		logit.Error.Println("MonitorServerGetInfo: error metric required")
		rest.Error(w, "Metric required", http.StatusBadRequest)
		return
	}

	//go get the IPAddress
	server, err := admindb.GetServer(dbConn, ServerID)
	if err != nil {
		logit.Error.Println("MonitorServerGetInfo:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output string
	url := "http://" + server.IPAddress + ":10001"
	if Metric == "cpmiostat" {
		iostatreq := cpmserverapi.MetricIostatRequest{}
		var iostatResp cpmserverapi.MetricIostatResponse
		iostatResp, err = cpmserverapi.MetricIostatClient(url, &iostatreq)
		if err != nil {
			logit.Error.Println("MonitorServerGetInfo:" + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}
		output = iostatResp.Output
	} else if Metric == "cpmdf" {
		dfreq := cpmserverapi.MetricDfRequest{}
		var dfResp cpmserverapi.MetricDfResponse
		dfResp, err = cpmserverapi.MetricDfClient(url, &dfreq)
		if err != nil {
			logit.Error.Println("MonitorServerGetInfo:" + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}
		output = dfResp.Output
	} else {
		logit.Error.Println("unknown Metric received")
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.(http.ResponseWriter).Write([]byte(output))
	//w.Header().set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
}

func GetAllServers(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllServers: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err := admindb.GetAllServers(dbConn)
	if err != nil {
		logit.Error.Println("GetAllServers: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	servers := make([]Server, len(results))
	i := 0
	for i = range results {
		servers[i].ID = results[i].ID
		servers[i].Name = results[i].Name
		servers[i].IPAddress = results[i].IPAddress
		servers[i].DockerBridgeIP = results[i].DockerBridgeIP
		servers[i].PGDataPath = results[i].PGDataPath
		servers[i].ServerClass = results[i].ServerClass
		servers[i].CreateDate = results[i].CreateDate
		i++
	}

	w.WriteJson(&servers)
}

func DeleteServer(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-server")
	if err != nil {
		logit.Error.Println("DeleteServer: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("DeleteServer: error server id required")
		rest.Error(w, "Server ID required", http.StatusBadRequest)
		return
	}

	err = admindb.DeleteServer(dbConn, ID)
	if err != nil {
		logit.Error.Println("DeleteServer: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}
