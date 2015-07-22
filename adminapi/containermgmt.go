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
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"time"
)

const CONTAINER_NOT_FOUND = "CONTAINER NOT FOUND"

func GetNode(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
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

	results, err2 := admindb.GetContainer(dbConn, ID)

	if results.ID == "" {
		rest.NotFound(w, r)
		return
	}
	if err2 != nil {
		logit.Error.Println("GetNode: " + err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	var currentStatus = "UNKNOWN"

	//go get the docker server IPAddress
	server := admindb.Server{}
	server, err = admindb.GetServer(dbConn, results.ServerID)
	if err != nil {
		logit.Error.Println("GetNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var domain string
	domain, err = admindb.GetDomain(dbConn)
	if err != nil {
		logit.Error.Println("GetNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request := &cpmserverapi.DockerInspectRequest{}
	request.ContainerName = results.Name
	var url = "http://" + server.IPAddress + ":10001"
	_, err = cpmserverapi.DockerInspectClient(url, request)
	if err != nil {
		logit.Error.Println("GetNode: " + err.Error())
		currentStatus = CONTAINER_NOT_FOUND
	}

	if currentStatus != "CONTAINER NOT FOUND" {
		//ping the db on that node to get current status
		var pinghost = results.Name
		if KubeEnv {
			pinghost = results.Name + "-db"
		}
		logit.Info.Println("pinging db on " + pinghost + "." + domain)
		currentStatus, err = GetPGStatus2(dbConn, results.Name, pinghost+"."+domain)
		if err != nil {
			logit.Error.Println("GetNode:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logit.Info.Println("pinging db finished")
	}

	node := ClusterNode{results.ID, results.ClusterID, results.ServerID,
		results.Name, results.Role, results.Image, results.CreateDate, currentStatus, results.ProjectID, results.ProjectName, results.ServerName, results.ClusterName}

	w.WriteJson(node)
}

func GetAllNodesForProject(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllNodes: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("GetAllNodesForProject: error project ID required")
		rest.Error(w, "project ID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetAllContainersForProject(dbConn, ID)
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
		nodes[i].ProjectID = results[i].ProjectID
		nodes[i].ProjectName = results[i].ProjectName
		nodes[i].ServerName = results[i].ServerName
		nodes[i].ClusterName = results[i].ClusterName
		//nodes[i].Status = "UNKNOWN"
		i++
	}

	w.WriteJson(&nodes)

}
func GetAllNodes(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllNodes: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err := admindb.GetAllContainers(dbConn)
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
		nodes[i].ProjectID = results[i].ProjectID
		nodes[i].ProjectName = results[i].ProjectName
		nodes[i].ServerName = results[i].ServerName
		nodes[i].ClusterName = results[i].ClusterName
		//nodes[i].Status = "UNKNOWN"
		i++
	}

	w.WriteJson(&nodes)

}

func GetAllNodesNotInCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllNodesNotInCluster: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err := admindb.GetAllContainersNotInCluster(dbConn)
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
		nodes[i].ProjectID = results[i].ProjectID
		nodes[i].ProjectName = results[i].ProjectName
		nodes[i].ServerName = results[i].ServerName
		nodes[i].ClusterName = results[i].ClusterName
		//nodes[i].Status = "UNKNOWN"
		i++
	}

	w.WriteJson(&nodes)

}

func GetAllNodesForCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
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

	results, err := admindb.GetAllContainersForCluster(dbConn, ClusterID)
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
		nodes[i].ProjectID = results[i].ProjectID
		nodes[i].ProjectName = results[i].ProjectName
		nodes[i].ServerName = results[i].ServerName
		nodes[i].ClusterName = results[i].ClusterName
		//nodes[i].Status = "UNKNOWN"
		i++
	}

	w.WriteJson(&nodes)

}

/*
 TODO refactor this to share code with DeleteCluster!!!!!
*/
func DeleteNode(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-container")
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
	var dbNode admindb.Container
	dbNode, err = admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("DeleteNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//go get the docker server IPAddress
	server := admindb.Server{}
	server, err = admindb.GetServer(dbConn, dbNode.ServerID)
	if err != nil {
		logit.Error.Println("DeleteNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var url = "http://" + server.IPAddress + ":10001"

	err = admindb.DeleteContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("DeleteNode: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("got server IP " + server.IPAddress)

	//it is possible that someone can remove a container
	//outside of us, so we let it pass that we can't remove
	//it

	request := &cpmserverapi.DockerRemoveRequest{}
	request.ContainerName = dbNode.Name
	_, err = cpmserverapi.DockerRemoveClient(url, request)
	if err != nil {
		logit.Error.Println("DeleteNode: error when trying to remove container " + err.Error())
	}

	//send the server a deletevolume command
	request2 := &cpmserverapi.DiskDeleteRequest{}
	request2.Path = server.PGDataPath + "/" + dbNode.Name
	_, err = cpmserverapi.DiskDeleteClient(url, request2)
	if err != nil {
		fmt.Println(err.Error())
	}

	//we should not have to delete the DNS entries because
	//of the dnsbridge, it should remove them when we remove
	//the containers via the docker api

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetAllNodesForServer(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
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

	results, err := admindb.GetAllContainersForServer(dbConn, serverID)
	if err != nil {
		logit.Error.Println("GetAllNodesForServer:" + err.Error())
		logit.Error.Println("error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server, err2 := admindb.GetServer(dbConn, serverID)
	if err2 != nil {
		logit.Error.Println("GetAllNodesForServer:" + err2.Error())
		logit.Error.Println("error " + err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	var response cpmserverapi.DockerInspectResponse
	var e error
	var url string
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
		nodes[i].ProjectID = results[i].ProjectID
		nodes[i].ProjectName = results[i].ProjectName
		nodes[i].ServerName = results[i].ServerName
		nodes[i].Status = "down"

		request := &cpmserverapi.DockerInspectRequest{}
		request.ContainerName = results[i].Name
		url = "http://" + server.IPAddress + ":10001"
		response, e = cpmserverapi.DockerInspectClient(url, request)
		logit.Info.Println("GetAllNodesForServer:" + results[i].Name + " " + response.IPAddress + " " + response.RunningState)
		if e != nil {
			logit.Error.Println("GetAllNodesForServer:" + e.Error())
			logit.Error.Println(e.Error())
			nodes[i].Status = "notfound"
		} else {
			logit.Info.Println("GetAllNodesForServer: setting " + results[i].Name + " to " + response.RunningState)
			nodes[i].Status = response.RunningState
		}

		i++
	}

	w.WriteJson(&nodes)

}

func AdminStartNode(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
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

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("AdminStartNode:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server := admindb.Server{}
	server, err = admindb.GetServer(dbConn, node.ServerID)
	if err != nil {
		logit.Error.Println("AdminStartNode:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var url = "http://" + server.IPAddress + ":10001"
	var response cpmserverapi.DockerStartResponse
	request := &cpmserverapi.DockerStartRequest{}
	request.ContainerName = node.Name
	response, err = cpmserverapi.DockerStartClient(url, request)
	if err != nil {
		logit.Error.Println("AdminStartNode: error when trying to start container " + err.Error())
	}
	logit.Info.Println(response.Output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func AdminStopNode(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
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

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("AdminStopNode:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server := admindb.Server{}
	server, err = admindb.GetServer(dbConn, node.ServerID)
	if err != nil {
		logit.Error.Println("AdminStopNode:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request := &cpmserverapi.DockerStopRequest{}
	request.ContainerName = node.Name
	var url = "http://" + server.IPAddress + ":10001"
	_, err = cpmserverapi.DockerStopClient(url, request)
	if err != nil {
		logit.Error.Println("AdminStopNode error when trying to stop container " + err.Error())
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func GetPGStatus2(dbConn *sql.DB, nodename string, hostname string) (string, error) {

	//fetch cpmtest user credentials
	nodeuser, err := admindb.GetContainerUser(dbConn, nodename, "cpmtest")
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}

	logit.Info.Println("cpmtest password is " + nodeuser.Passwd)

	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")

	dbConn2, err := util.GetMonitoringConnection(hostname, "cpmtest", pgport.Value, "cpmtest", nodeuser.Passwd)
	defer dbConn2.Close()
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}

	var value string

	err = dbConn2.QueryRow(fmt.Sprintf("select now()::text")).Scan(&value)
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

func AdminStartServerContainers(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("AdminStartServerContainers: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//serverID
	serverid := r.PathParam("ID")
	if serverid == "" {
		logit.Error.Println("AdminStartServerContainers: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	containers, err := admindb.GetAllContainersForServer(dbConn, serverid)
	if err != nil {
		logit.Error.Println("AdminStartServerContainers:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//for each, get server, start container
	//use a 'best effort' approach here since containers
	//can be removed outside of CPM's control

	var url string

	for i := range containers {
		//fetch the server
		server := admindb.Server{}
		server, err = admindb.GetServer(dbConn, containers[i].ServerID)
		if err != nil {
			logit.Error.Println("AdminStartServerContainers:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//start the container
		var response cpmserverapi.DockerStartResponse
		var err error
		request := &cpmserverapi.DockerStartRequest{}
		request.ContainerName = containers[i].Name
		url = "http://" + server.IPAddress + ":10001"
		response, err = cpmserverapi.DockerStartClient(url, request)
		if err != nil {
			logit.Error.Println("AdminStartServerContainers: error when trying to start container " + err.Error())
		}
		logit.Info.Println(response.Output)

	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}
func AdminStopServerContainers(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("AdminStopServerContainers: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//serverID
	serverid := r.PathParam("ID")
	if serverid == "" {
		logit.Error.Println("AdminStopoServerContainers: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	//fetch the server
	containers, err := admindb.GetAllContainersForServer(dbConn, serverid)
	if err != nil {
		logit.Error.Println("AdminStopServerContainers:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var url string
	//for each, get server, stop container
	for i := range containers {
		server := admindb.Server{}
		server, err = admindb.GetServer(dbConn, containers[i].ServerID)
		if err != nil {
			logit.Error.Println("AdminStopServerContainers:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//send stop command before stopping container
		if containers[i].Role == "pgpool" {
			var stoppoolResp cpmcontainerapi.StopPgpoolResponse
			stoppoolResp, err = cpmcontainerapi.StopPgpoolClient(containers[i].Name)
			logit.Info.Println("AdminStoppg:" + stoppoolResp.Output)
		} else {
			var stopResp cpmcontainerapi.StopPGResponse
			stopResp, err = cpmcontainerapi.StopPGClient(containers[i].Name)
			logit.Info.Println("AdminStoppg:" + stopResp.Output)
		}
		if err != nil {
			logit.Error.Println("AdminStopServerContainers:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		time.Sleep(2000 * time.Millisecond)
		//stop container
		request := &cpmserverapi.DockerStopRequest{}
		request.ContainerName = containers[i].Name
		url = "http://" + server.IPAddress + ":10001"
		_, err = cpmserverapi.DockerStopClient(url, request)
		if err != nil {
			logit.Error.Println("AdminStopServerContainers: error when trying to start container " + err.Error())
		}
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}
