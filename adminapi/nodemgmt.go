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
	"database/sql"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserveragent"
	"github.com/crunchydata/crunchy-postgresql-manager/kubeclient"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
)

const CONTAINER_NOT_FOUND = "CONTAINER NOT FOUND"

func GetNode(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetNode: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("GetNode: error node ID required")
		rest.Error(w, "node ID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetDBNode(ID)

	if results.ID == "" {
		rest.NotFound(w, r)
		return
	}
	if err != nil {
		logit.Error.Println("GetNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var currentStatus = "UNKNOWN"

	//go get the docker server IPAddress
	server := admindb.DBServer{}
	server, err = admindb.GetDBServer(results.ServerID)
	if err != nil {
		logit.Error.Println("GetNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var domain string
	domain, err = admindb.GetDomain()
	if err != nil {
		logit.Error.Println("GetNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if KubeEnv {
		var podInfo kubeclient.MyPod
		podInfo, err = kubeclient.GetPod(KubeURL, results.Name)
		if err != nil {
			currentStatus = CONTAINER_NOT_FOUND
		}
		logit.Info.Println("pod info status is " + podInfo.CurrentState.Status)
		if podInfo.CurrentState.Status != "Running" {
			currentStatus = CONTAINER_NOT_FOUND
		}
	} else {
		_, err = cpmserveragent.DockerInspect2Command(results.Name, server.IPAddress)
		if err != nil {
			logit.Error.Println("GetNode: " + err.Error())
			currentStatus = CONTAINER_NOT_FOUND
		}

	}

	if currentStatus != "CONTAINER NOT FOUND" {
		//ping the db on that node to get current status
		var pinghost = results.Name
		if KubeEnv {
			pinghost = results.Name + "-db"
		}
		logit.Info.Println("pinging db on " + pinghost + "." + domain)
		currentStatus, err = GetPGStatus2(pinghost + "." + domain)
		if err != nil {
			logit.Error.Println("GetNode:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logit.Info.Println("pinging db finished")
	}

	node := ClusterNode{results.ID, results.ClusterID, results.ServerID,
		results.Name, results.Role, results.Image, results.CreateDate, currentStatus}

	w.WriteJson(node)
}

func GetAllNodes(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllNodes: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err := admindb.GetAllDBNodes()
	if err != nil {
		logit.Error.Println("GetAllNodes: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	nodes := make([]ClusterNode, len(results))
	i := 0
	for i = range results {
		nodes[i].ID = results[i].ID
		nodes[i].Name = results[i].Name
		nodes[i].ClusterID = results[i].ClusterID
		nodes[i].ServerID = results[i].ServerID
		nodes[i].Role = results[i].Role
		nodes[i].Image = results[i].Image
		nodes[i].CreateDate = results[i].CreateDate
		//nodes[i].Status = "UNKNOWN"
		i++
	}

	w.WriteJson(&nodes)

}

func GetAllNodesNotInCluster(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllNodesNotInCluster: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err := admindb.GetAllDBNodesNotInCluster()
	if err != nil {
		logit.Error.Println("GetAllNodesNotInCluster: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	nodes := make([]ClusterNode, len(results))
	i := 0
	for i = range results {
		nodes[i].ID = results[i].ID
		nodes[i].Name = results[i].Name
		nodes[i].ClusterID = results[i].ClusterID
		nodes[i].ServerID = results[i].ServerID
		nodes[i].Role = results[i].Role
		nodes[i].Image = results[i].Image
		nodes[i].CreateDate = results[i].CreateDate
		//nodes[i].Status = "UNKNOWN"
		i++
	}

	w.WriteJson(&nodes)

}

func GetAllNodesForCluster(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllForCluster: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ClusterID := r.PathParam("ClusterID")
	if ClusterID == "" {
		logit.Error.Println("GetAllNodesForCluster: node ClusterID required")
		rest.Error(w, "node ClusterID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetAllDBNodesForCluster(ClusterID)
	if err != nil {
		logit.Error.Println("GetAllNodesForCluster:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	nodes := make([]ClusterNode, len(results))
	i := 0
	for i = range results {
		nodes[i].ID = results[i].ID
		nodes[i].Name = results[i].Name
		nodes[i].ClusterID = results[i].ClusterID
		nodes[i].ServerID = results[i].ServerID
		nodes[i].Role = results[i].Role
		nodes[i].Image = results[i].Image
		nodes[i].CreateDate = results[i].CreateDate
		//nodes[i].Status = "UNKNOWN"
		i++
	}

	w.WriteJson(&nodes)

}

/*
 TODO refactor this to share code with DeleteCluster!!!!!
*/
func DeleteNode(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-container")
	if err != nil {
		logit.Error.Println("DeleteNode: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("DeleteNode: error node ID required")
		rest.Error(w, "node ID required", http.StatusBadRequest)
		return
	}

	//go get the node we intend to delete
	var dbNode admindb.DBClusterNode
	dbNode, err = admindb.GetDBNode(ID)
	if err != nil {
		logit.Error.Println("DeleteNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//go get the docker server IPAddress
	server := admindb.DBServer{}
	server, err = admindb.GetDBServer(dbNode.ServerID)
	if err != nil {
		logit.Error.Println("DeleteNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = admindb.DeleteDBNode(ID)
	if err != nil {
		logit.Error.Println("DeleteNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("got server IP " + server.IPAddress)

	//it is possible that someone can remove a container
	//outside of us, so we let it pass that we can't remove
	//it

	var output string

	if KubeEnv {
		//delete the kube pod with this name
		err = kubeclient.DeletePod(KubeURL, dbNode.Name)
		if err != nil {
			logit.Error.Println("DeleteNode:" + err.Error())
			rest.Error(w, "error in deleting pod", http.StatusBadRequest)
			return
		}
		//delete the kube service with this name
		err = kubeclient.DeleteService(KubeURL, dbNode.Name)
		if err != nil {
			logit.Error.Println("DeleteNode:" + err.Error())
			rest.Error(w, "error in deleting service 1", http.StatusBadRequest)
			return
		}
		//delete the kube service with this name 5432
		err = kubeclient.DeleteService(KubeURL, dbNode.Name+"-db")
		if err != nil {
			logit.Error.Println("DeleteNode:" + err.Error())
			rest.Error(w, "error in deleting service 2", http.StatusBadRequest)
			return
		}
	} else {
		output, err = cpmserveragent.DockerRemoveContainer(dbNode.Name, server.IPAddress)
		if err != nil {
			logit.Error.Println("DeleteNode: error when trying to remove container " + err.Error())
		}
	}

	//send the server a deletevolume command
	output, err = cpmserveragent.AgentCommand("deletevolume", server.PGDataPath+"/"+dbNode.Name, server.IPAddress)
	logit.Info.Println(output)

	//we should not have to delete the DNS entries because
	//of the dnsbridge, it should remove them when we remove
	//the containers via the docker api

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetAllNodesForServer(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllNodesForServer: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	serverID := r.PathParam("ServerID")
	if serverID == "" {
		logit.Error.Println("GetAllNodesForServer: error serverID required")
		rest.Error(w, "serverID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetAllDBNodesForServer(serverID)
	if err != nil {
		logit.Error.Println("GetAllNodesForServer:" + err.Error())
		logit.Error.Println("error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server, err2 := admindb.GetDBServer(serverID)
	if err2 != nil {
		logit.Error.Println("GetAllNodesForServer:" + err2.Error())
		logit.Error.Println("error " + err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	var output cpmserveragent.InspectOutput
	var e error
	nodes := make([]ClusterNode, len(results))
	i := 0
	for i = range results {
		nodes[i].ID = results[i].ID
		nodes[i].Name = results[i].Name
		nodes[i].ClusterID = results[i].ClusterID
		nodes[i].ServerID = results[i].ServerID
		nodes[i].Role = results[i].Role
		nodes[i].Image = results[i].Image
		nodes[i].CreateDate = results[i].CreateDate
		nodes[i].Status = "down"

		output, e = cpmserveragent.DockerInspect2Command(results[i].Name, server.IPAddress)
		logit.Info.Println("GetAllNodesForServer:" + results[i].Name + " " + output.IPAddress + " " + output.RunningState)
		if e != nil {
			logit.Error.Println("GetAllNodesForServer:" + e.Error())
			logit.Error.Println(e.Error())
			nodes[i].Status = "notfound"
		} else {
			logit.Info.Println("GetAllNodesForServer: setting " + results[i].Name + " to " + output.RunningState)
			nodes[i].Status = output.RunningState
		}

		i++
	}

	w.WriteJson(&nodes)

}

func AdminStartNode(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("AdminStartNode: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("AdminStartNode: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		logit.Error.Println("AdminStartNode:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server := admindb.DBServer{}
	server, err = admindb.GetDBServer(node.ServerID)
	if err != nil {
		logit.Error.Println("AdminStartNode:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output string
	output, err = cpmserveragent.DockerStartContainer(node.Name,
		server.IPAddress)
	if err != nil {
		logit.Error.Println("AdminStartNode: error when trying to start container " + err.Error())
	}
	logit.Info.Println(output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func AdminStopNode(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("AdminStopNode: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("AdminStopNode: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		logit.Error.Println("AdminStopNode:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server := admindb.DBServer{}
	server, err = admindb.GetDBServer(node.ServerID)
	if err != nil {
		logit.Error.Println("AdminStopNode:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output string
	output, err = cpmserveragent.DockerStopContainer(node.Name,
		server.IPAddress)
	if err != nil {
		logit.Error.Println("AdminStopNode error when trying to stop container " + err.Error())
	}
	logit.Info.Println(output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func GetPGStatus2(hostname string) (string, error) {

	dbConn, err := util.GetMonitoringConnection(hostname, "cpmtest", "5432", "cpmtest", "cpmtest")
	defer dbConn.Close()
	var value string

	err = dbConn.QueryRow(fmt.Sprintf("select now()::text")).Scan(&value)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("getpgstatus 2 no rows returned")
		return "OFFLINE", nil
	case err != nil:
		logit.Info.Println("getpgstatus2 error " + err.Error())
		return "OFFLINE", nil
	default:
		logit.Info.Println("getpgstatus2 returned " + value)
	}

	return "RUNNING", nil
}
