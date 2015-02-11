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
	"bytes"
	"crunchy.com/admindb"
	"crunchy.com/cpmagent"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
	"net/http"
	"os"
	"os/exec"
)

func GetNode(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetNode: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		glog.Errorln("GetNode: error node ID required")
		rest.Error(w, "node ID required", 400)
		return
	}

	results, err := admindb.GetDBNode(ID)

	if results.ID == "" {
		rest.NotFound(w, r)
		return
	}
	if err != nil {
		glog.Errorln("GetNode: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//ping the db on that node to get current status
	var currentStatus = "UNKNOWN"

	currentStatus, err = GetPGStatus(results.Name)
	if err != nil {
		glog.Errorln("GetNode:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	node := ClusterNode{results.ID, results.ClusterID, results.ServerID,
		results.Name, results.Role, results.Image, results.CreateDate, currentStatus}
	w.WriteJson(node)
}

func GetPGStatus(hostname string) (string, error) {
	var currentStatus = "UNKNOWN"
	var err error

	cmd := exec.Command(CPMBIN+"pgstatus",
		hostname,
		"5432",
		"cpmtest",
		"cpmtest",
		"cpmtest")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		glog.Errorln("GetNode:" + err.Error())
		return "", err
	}
	glog.Infoln("GetPGStatus: command output was " + out.String())

	glog.Infoln("GetPGStatus: output from ping was [" + out.String() + "]")
	currentStatus = "OFFLINE"

	if out.String() == "up" {
		currentStatus = "RUNNING"
	}

	return currentStatus, err
}

func GetAllNodes(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetAllNodes: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	results, err := admindb.GetAllDBNodes()
	if err != nil {
		glog.Errorln("GetAllNodes: " + err.Error())
		rest.Error(w, err.Error(), 400)
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
		nodes[i].Status = "UNKNOWN"
		i++
	}

	w.WriteJson(&nodes)

}

func GetAllNodesNotInCluster(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetAllNodesNotInCluster: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	results, err := admindb.GetAllDBNodesNotInCluster()
	if err != nil {
		glog.Errorln("GetAllNodesNotInCluster: " + err.Error())
		rest.Error(w, err.Error(), 400)
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
		nodes[i].Status = "UNKNOWN"
		i++
	}

	w.WriteJson(&nodes)

}

func GetAllNodesForCluster(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetAllForCluster: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ClusterID := r.PathParam("ClusterID")
	if ClusterID == "" {
		glog.Errorln("GetAllNodesForCluster: node ClusterID required")
		rest.Error(w, "node ClusterID required", 400)
		return
	}

	results, err := admindb.GetAllDBNodesForCluster(ClusterID)
	if err != nil {
		glog.Errorln("GetAllNodesForCluster:" + err.Error())
		rest.Error(w, err.Error(), 400)
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
		nodes[i].Status = "UNKNOWN"
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
		glog.Errorln("DeleteNode: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		glog.Errorln("DeleteNode: error node ID required")
		rest.Error(w, "node ID required", 400)
		return
	}

	//go get the node we intend to delete
	var dbNode admindb.DBClusterNode
	dbNode, err = admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("DeleteNode: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//go get the docker server IPAddress
	server := admindb.DBServer{}
	server, err = admindb.GetDBServer(dbNode.ServerID)
	if err != nil {
		glog.Errorln("DeleteNode: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	err = admindb.DeleteDBNode(ID)
	if err != nil {
		glog.Errorln("DeleteNode: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	glog.Infoln("got server IP " + server.IPAddress)

	//it is possible that someone can remove a container
	//outside of us, so we let it pass that we can't remove
	//it
	kubeEnv := false
	kube := os.Getenv("KUBE_URL")
	glog.Infoln("KUBE_URL=[" + kube + "]")
	if kube != "" {
		glog.Infoln("KUBE_URL value set, assume Kube environment")
		kubeEnv = true
	}

	var output string

	if kubeEnv {
		//delete the kube pod with this name
		err = DeletePod(kube, dbNode.Name)
		if err != nil {
			glog.Errorln("DeleteNode:" + err.Error())
			rest.Error(w, "error in deleting pod", 400)
			return
		}
	} else {
		output, err = cpmagent.DockerRemoveContainer(dbNode.Name, server.IPAddress)
		if err != nil {
			glog.Errorln("DeleteNode: error when trying to remove container " + err.Error())
		}
	}

	//send the server a deletevolume command
	output, err = cpmagent.AgentCommand(CPMBIN+"deletevolume", server.PGDataPath+"/"+dbNode.Name, server.IPAddress)
	glog.Infoln(output)

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
		glog.Errorln("GetAllNodesForServer: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	serverID := r.PathParam("ServerID")
	if serverID == "" {
		glog.Errorln("GetAllNodesForServer: error serverID required")
		rest.Error(w, "serverID required", 400)
		return
	}

	results, err := admindb.GetAllDBNodesForServer(serverID)
	if err != nil {
		glog.Errorln("GetAllNodesForServer:" + err.Error())
		glog.Errorln("error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	server, err2 := admindb.GetDBServer(serverID)
	if err2 != nil {
		glog.Errorln("GetAllNodesForServer:" + err2.Error())
		glog.Errorln("error " + err2.Error())
		rest.Error(w, err2.Error(), 400)
		return
	}

	var output cpmagent.InspectOutput
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
		output, e = cpmagent.DockerInspect2Command(results[i].Name, server.IPAddress)
		glog.Infoln("GetAllNodesForServer:" + results[i].Name + " " + output.IPAddress + " " + output.RunningState)
		if e != nil {
			glog.Errorln("GetAllNodesForServer:" + e.Error())
			glog.Errorln(e.Error())
		} else {
			glog.Infoln("GetAllNodesForServer: setting " + results[i].Name + " to " + output.RunningState)
			nodes[i].Status = output.RunningState
		}
		i++
	}

	w.WriteJson(&nodes)

}

func AdminStartNode(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("AdminStartNode: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		glog.Errorln("AdminStartNode: error ID required")
		rest.Error(w, "ID required", 400)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("AdminStartNode:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	server := admindb.DBServer{}
	server, err = admindb.GetDBServer(node.ServerID)
	if err != nil {
		glog.Errorln("AdminStartNode:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	var output string
	output, err = cpmagent.DockerStartContainer(node.Name,
		server.IPAddress)
	if err != nil {
		glog.Errorln("AdminStartNode: error when trying to start container " + err.Error())
	}
	glog.Infoln(output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func AdminStopNode(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("AdminStopNode: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		glog.Errorln("AdminStopNode: error ID required")
		rest.Error(w, "ID required", 400)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("AdminStopNode:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	server := admindb.DBServer{}
	server, err = admindb.GetDBServer(node.ServerID)
	if err != nil {
		glog.Errorln("AdminStopNode:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	var output string
	output, err = cpmagent.DockerStopContainer(node.Name,
		server.IPAddress)
	if err != nil {
		glog.Errorln("AdminStopNode error when trying to stop container " + err.Error())
	}
	glog.Infoln(output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}
