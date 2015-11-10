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
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/swarmapi"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"strings"
)

// GetServer return a server definition
func GetServer(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")
	logit.Info.Println("in GetServer with ID=" + ID)

	//currently no state about a server is maintained other than IP and port number
	//which we use for the ID, Name, and IPAddress values

	server := types.Server{ID, ID, ID, "", "", ""}

	w.WriteJson(&server)
}

// MonitorServerGetInfo return server information for a given server
func MonitorServerGetInfo(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	Metric := r.PathParam("Metric")

	ServerID := r.PathParam("ServerID")
	if ServerID == "" {
		logit.Error.Println("error ServerID required")
		rest.Error(w, "ServerID required", http.StatusBadRequest)
		return
	}
	cleanIP := strings.Replace(ServerID, "_", ".", -1)
	logit.Info.Println("ServerID=" + ServerID)
	logit.Info.Println("cleanIP=" + cleanIP)
	if Metric == "" {
		logit.Error.Println("MonitorServerGetInfo: error metric required")
		rest.Error(w, "Metric required", http.StatusBadRequest)
		return
	}

	var output string
	if Metric == "cpmiostat" {
		iostatreq := cpmserverapi.MetricIostatRequest{}
		var iostatResp cpmserverapi.MetricIostatResponse
		iostatResp, err = cpmserverapi.MetricIostatClient(cleanIP, &iostatreq)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}
		output = iostatResp.Output
	} else if Metric == "cpmdf" {
		dfreq := cpmserverapi.MetricDfRequest{}
		var dfResp cpmserverapi.MetricDfResponse
		dfResp, err = cpmserverapi.MetricDfClient(cleanIP, &dfreq)
		if err != nil {
			logit.Error.Println(err.Error())
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

// GetAllServers return a list of servers
func GetAllServers(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//use swarm to get the list of servers

	var infoResponse swarmapi.DockerInfoResponse
	infoResponse, err = swarmapi.DockerInfo()
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for x := 0; x < len(infoResponse.Output); {
		logit.Info.Println("got back " + infoResponse.Output[x])
		x++
	}

	servers := make([]types.Server, len(infoResponse.Output))
	i := 0
	for i = range infoResponse.Output {
		servers[i].ID = infoResponse.Output[i]
		servers[i].Name = infoResponse.Output[i]
		servers[i].IPAddress = infoResponse.Output[i]
		i++
	}

	w.WriteJson(&servers)
}
