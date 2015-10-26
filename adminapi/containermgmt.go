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
	"bytes"
	"database/sql"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"os/exec"
	"time"
)

const CONTAINER_NOT_FOUND = "CONTAINER NOT FOUND"

// GetNode returns the container node definition
func GetNode(w rest.ResponseWriter, r *rest.Request) {
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
	if ID == "" {
		logit.Error.Println("error node ID required")
		rest.Error(w, "node ID required", http.StatusBadRequest)
		return
	}

	node, err2 := admindb.GetContainer(dbConn, ID)

	if node.ID == "" {
		rest.NotFound(w, r)
		return
	}
	if err2 != nil {
		logit.Error.Println(err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	var currentStatus = "UNKNOWN"

	//go get the docker server IPAddress
	server := types.Server{}
	server, err = admindb.GetServer(dbConn, node.ServerID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request := &cpmserverapi.DockerInspectRequest{}
	request.ContainerName = node.Name
	_, err = cpmserverapi.DockerInspectClient(server.Name, request)
	if err != nil {
		logit.Error.Println(err.Error())
		currentStatus = CONTAINER_NOT_FOUND
	}

	if currentStatus != "CONTAINER NOT FOUND" {
		currentStatus, err = NewPingPG(dbConn, &node)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logit.Info.Println("pinging db finished")
	}

	clusternode := types.ClusterNode{node.ID, node.ClusterID, node.ServerID,
		node.Name, node.Role, node.Image, node.CreateDate, currentStatus, node.ProjectID, node.ProjectName, node.ServerName, node.ClusterName}

	w.WriteJson(clusternode)
}

// GetAllNodesForProject returns all node definitions for a given project
func GetAllNodesForProject(w rest.ResponseWriter, r *rest.Request) {
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
	if ID == "" {
		logit.Error.Println("project ID required")
		rest.Error(w, "project ID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetAllContainersForProject(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	nodes := make([]types.ClusterNode, len(results))
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

// TODO
func GetAllNodes(w rest.ResponseWriter, r *rest.Request) {
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

	results, err := admindb.GetAllContainers(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	nodes := make([]types.ClusterNode, len(results))
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

// TODO
func GetAllNodesNotInCluster(w rest.ResponseWriter, r *rest.Request) {
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

	results, err := admindb.GetAllContainersNotInCluster(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	nodes := make([]types.ClusterNode, len(results))
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

// GetAllNodesForCluster returns a list of nodes for a given cluster
func GetAllNodesForCluster(w rest.ResponseWriter, r *rest.Request) {
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

	ClusterID := r.PathParam("ClusterID")
	if ClusterID == "" {
		logit.Error.Println("ClusterID required")
		rest.Error(w, "node ClusterID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetAllContainersForCluster(dbConn, ClusterID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	nodes := make([]types.ClusterNode, len(results))
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
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-container")
	if err != nil {
		logit.Error.Println(err.Error())
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
	var dbNode types.Container
	dbNode, err = admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//go get the docker server IPAddress
	server := types.Server{}
	server, err = admindb.GetServer(dbConn, dbNode.ServerID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = admindb.DeleteContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("got server IP " + server.IPAddress)

	//it is possible that someone can remove a container
	//outside of us, so we let it pass that we can't remove
	//it

	request := &cpmserverapi.DockerRemoveRequest{}
	request.ContainerName = dbNode.Name
	_, err = cpmserverapi.DockerRemoveClient(server.Name, request)
	if err != nil {
		logit.Error.Println(err.Error())
	}

	//send the server a deletevolume command
	request2 := &cpmserverapi.DiskDeleteRequest{}
	request2.Path = server.PGDataPath + "/" + dbNode.Name
	_, err = cpmserverapi.DiskDeleteClient(server.Name, request2)
	if err != nil {
		fmt.Println(err.Error())
	}

	//we should not have to delete the DNS entries because
	//of the dnsbridge, it should remove them when we remove
	//the containers via the docker api

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

// GetAllNodesForServer returns a list of all nodes on a given server
func GetAllNodesForServer(w rest.ResponseWriter, r *rest.Request) {
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

	serverID := r.PathParam("ServerID")
	if serverID == "" {
		logit.Error.Println("GetAllNodesForServer: error serverID required")
		rest.Error(w, "serverID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetAllContainersForServer(dbConn, serverID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server, err2 := admindb.GetServer(dbConn, serverID)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	var response cpmserverapi.DockerInspectResponse
	var e error
	nodes := make([]types.ClusterNode, len(results))
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
		response, e = cpmserverapi.DockerInspectClient(server.Name, request)
		logit.Info.Println("GetAllNodesForServer:" + results[i].Name + " " + response.IPAddress + " " + response.RunningState)
		if e != nil {
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

// AdminStartNode starts a container
func AdminStartNode(w rest.ResponseWriter, r *rest.Request) {
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
	if ID == "" {
		logit.Error.Println("AdminStartNode: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server := types.Server{}
	server, err = admindb.GetServer(dbConn, node.ServerID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response cpmserverapi.DockerStartResponse
	request := &cpmserverapi.DockerStartRequest{}
	request.ContainerName = node.Name
	response, err = cpmserverapi.DockerStartClient(server.Name, request)
	if err != nil {
		logit.Error.Println(err.Error())
	}
	logit.Info.Println(response.Output)

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

// AdminStopNode stops a container
func AdminStopNode(w rest.ResponseWriter, r *rest.Request) {
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
	if ID == "" {
		logit.Error.Println("AdminStopNode: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server := types.Server{}
	server, err = admindb.GetServer(dbConn, node.ServerID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request := &cpmserverapi.DockerStopRequest{}
	request.ContainerName = node.Name
	_, err = cpmserverapi.DockerStopClient(server.Name, request)
	if err != nil {
		logit.Error.Println(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

// NewPingPG performs a simple network (nc) ping on a given container's postgres database port
func NewPingPG(dbConn *sql.DB, node *types.Container) (string, error) {

	var err error

	var pgport types.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")
	if err != nil {
		logit.Error.Println(err.Error())
		return "OFFLINE", err
	}

	var cmd *exec.Cmd
	logit.Info.Println("about to newping " + node.Name + " at port " + pgport.Value)
	cmd = exec.Command("ping-wrapper.sh", node.Name, pgport.Value)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
		return "OFFLINE", err
	}
	var rc = out.String()
	logit.Info.Println("NewPingPG rc is " + rc)
	if rc != "connected" {
		return "OFFLINE", nil
	}

	return "RUNNING", nil
}

// AdminStartServerContainers starts all containers on a given server
func AdminStartServerContainers(w rest.ResponseWriter, r *rest.Request) {
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

	//serverID
	serverid := r.PathParam("ID")
	if serverid == "" {
		logit.Error.Println(" error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	containers, err := admindb.GetAllContainersForServer(dbConn, serverid)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//for each, get server, start container
	//use a 'best effort' approach here since containers
	//can be removed outside of CPM's control

	for i := range containers {
		//fetch the server
		server := types.Server{}
		server, err = admindb.GetServer(dbConn, containers[i].ServerID)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//start the container
		var response cpmserverapi.DockerStartResponse
		var err error
		request := &cpmserverapi.DockerStartRequest{}
		request.ContainerName = containers[i].Name
		response, err = cpmserverapi.DockerStartClient(server.Name, request)
		if err != nil {
			logit.Error.Println("AdminStartServerContainers: error when trying to start container " + err.Error())
		}
		logit.Info.Println(response.Output)

	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

// AdminStopServerContainers stops all containers on a given server
func AdminStopServerContainers(w rest.ResponseWriter, r *rest.Request) {
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
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//for each, get server, stop container
	for i := range containers {
		server := types.Server{}
		server, err = admindb.GetServer(dbConn, containers[i].ServerID)
		if err != nil {
			logit.Error.Println(err.Error())
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
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sleepTime, _ := time.ParseDuration("2s")
		time.Sleep(sleepTime)

		//stop container
		request := &cpmserverapi.DockerStopRequest{}
		request.ContainerName = containers[i].Name
		_, err = cpmserverapi.DockerStopClient(server.Name, request)
		if err != nil {
			logit.Error.Println("AdminStopServerContainers: error when trying to start container " + err.Error())
		}
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}
