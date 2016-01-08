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
	"github.com/crunchydata/crunchy-postgresql-manager/swarmapi"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"strings"
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

	request := &swarmapi.DockerInspectRequest{}
	var inspectInfo swarmapi.DockerInspectResponse
	request.ContainerName = node.Name
	inspectInfo, err = swarmapi.DockerInspect(request)
	if err != nil {
		logit.Error.Println(err.Error())
		currentStatus = CONTAINER_NOT_FOUND
	}

	if currentStatus != "CONTAINER NOT FOUND" {
		var pgport types.Setting
		pgport, err = admindb.GetSetting(dbConn, "PG-PORT")
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		currentStatus, err = util.FastPing(pgport.Value, node.Name)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//logit.Info.Println("pinging db finished")
	}

	clusternode := new(types.ClusterNode)
	clusternode.ID = node.ID
	clusternode.ClusterID = node.ClusterID
	clusternode.Name = node.Name
	clusternode.Role = node.Role
	clusternode.Image = node.Image
	clusternode.CreateDate = node.CreateDate
	clusternode.Status = currentStatus
	clusternode.ProjectID = node.ProjectID
	clusternode.ProjectName = node.ProjectName
	clusternode.ClusterName = node.ClusterName
	clusternode.ServerID = inspectInfo.ServerID
	clusternode.IPAddress = inspectInfo.IPAddress

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
		nodes[i].Role = results[i].Role
		nodes[i].Image = results[i].Image
		nodes[i].CreateDate = results[i].CreateDate
		nodes[i].ProjectID = results[i].ProjectID
		nodes[i].ProjectName = results[i].ProjectName
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
		nodes[i].Role = results[i].Role
		nodes[i].Image = results[i].Image
		nodes[i].CreateDate = results[i].CreateDate
		nodes[i].ProjectID = results[i].ProjectID
		nodes[i].ProjectName = results[i].ProjectName
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
		nodes[i].Role = results[i].Role
		nodes[i].Image = results[i].Image
		nodes[i].CreateDate = results[i].CreateDate
		nodes[i].ProjectID = results[i].ProjectID
		nodes[i].ProjectName = results[i].ProjectName
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
		nodes[i].Role = results[i].Role
		nodes[i].Image = results[i].Image
		nodes[i].CreateDate = results[i].CreateDate
		nodes[i].ProjectID = results[i].ProjectID
		nodes[i].ProjectName = results[i].ProjectName
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

	var infoResponse swarmapi.DockerInfoResponse
	infoResponse, err = swarmapi.DockerInfo()
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	servers := make([]types.Server, len(infoResponse.Output))
	i := 0
	for i = range infoResponse.Output {
		servers[i].ID = infoResponse.Output[i]
		servers[i].Name = infoResponse.Output[i]
		servers[i].IPAddress = infoResponse.Output[i]
		i++
	}

	var pgdatapath types.Setting
	pgdatapath, err = admindb.GetSetting(dbConn, "PG-DATA-PATH")
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

	logit.Info.Println("remove 1")
	//it is possible that someone can remove a container
	//outside of us, so we let it pass that we can't remove
	//it

	request := &swarmapi.DockerRemoveRequest{}
	request.ContainerName = dbNode.Name
	_, err = swarmapi.DockerRemove(request)
	if err != nil {
		logit.Error.Println(err.Error())
	}

	logit.Info.Println("remove 2")
	//send the server a deletevolume command
	request2 := &cpmserverapi.DiskDeleteRequest{}
	request2.Path = pgdatapath.Value + "/" + dbNode.Name
	for _, each := range servers {
		_, err = cpmserverapi.DiskDeleteClient(each.Name, request2)
		if err != nil {
			logit.Error.Println(err.Error())
		}
	}
	logit.Info.Println("remove 3")

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

	serverIPAddress := strings.Replace(serverID, "_", ".", -1)

	results, err := swarmapi.DockerPs(serverIPAddress)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nodes := make([]types.ClusterNode, len(results.Output))
	i := 0
	var container types.Container

	for _, each := range results.Output {
		logit.Info.Println("got back Name:" + each.Name + " Status:" + each.Status + " Image:" + each.Image)
		nodes[i].Name = each.Name
		container, err = admindb.GetContainerByName(dbConn, each.Name)
		if err != nil {
			logit.Error.Println(err.Error())
			nodes[i].ID = "unknown"
			nodes[i].ProjectID = "unknown"
		} else {
			nodes[i].ID = container.ID
			nodes[i].ProjectID = container.ProjectID
		}

		nodes[i].Status = each.Status
		nodes[i].Image = each.Image
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

	/**
	server := types.Server{}
	server, err = admindb.GetServer(dbConn, node.ServerID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	*/

	var response swarmapi.DockerStartResponse
	request := &swarmapi.DockerStartRequest{}
	request.ContainerName = node.Name
	response, err = swarmapi.DockerStart(request)
	if err != nil {
		logit.Error.Println(err.Error())
		logit.Error.Println(response.Output)
	}
	//logit.Info.Println(response.Output)

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

	/**
	server := types.Server{}
	server, err = admindb.GetServer(dbConn, node.ServerID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	*/

	request := &swarmapi.DockerStopRequest{}
	request.ContainerName = node.Name
	_, err = swarmapi.DockerStop(request)
	if err != nil {
		logit.Error.Println(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

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

	cleanIP := strings.Replace(serverid, "_", ".", -1)

	containers, err := swarmapi.DockerPs(cleanIP)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//for each, get server, start container
	//use a 'best effort' approach here since containers
	//can be removed outside of CPM's control

	for _, each := range containers.Output {

		//start the container
		var response swarmapi.DockerStartResponse
		var err error
		request := &swarmapi.DockerStartRequest{}
		logit.Info.Println("trying to start " + each.Name)
		request.ContainerName = each.Name
		response, err = swarmapi.DockerStart(request)
		if err != nil {
			logit.Error.Println("AdminStartServerContainers: error when trying to start container " + err.Error())
			logit.Error.Println(response.Output)
		}
		//logit.Info.Println(response.Output)

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

	cleanIP := strings.Replace(serverid, "_", ".", -1)

	containers, err := swarmapi.DockerPs(cleanIP)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//for each, get server, stop container
	for _, each := range containers.Output {

		if strings.HasPrefix(each.Status, "Up") {
			//stop container
			request := &swarmapi.DockerStopRequest{}
			request.ContainerName = each.Name
			logit.Info.Println("stopping " + request.ContainerName)
			_, err = swarmapi.DockerStop(request)
			if err != nil {
				logit.Error.Println("AdminStopServerContainers: error when trying to start container " + err.Error())
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}
