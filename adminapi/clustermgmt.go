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
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/template"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	STANDBY = "standby"
)

type AutoClusterInfo struct {
	Name           string
	ProjectID      string
	ClusterType    string
	ClusterProfile string
	Token          string
}

func ScaleUpCluster(w rest.ResponseWriter, r *rest.Request) {

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	cluster, err := admindb.GetCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var containers []admindb.Container
	containers, err = admindb.GetAllContainersForCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//determine number of standby nodes currently
	standbyCnt := 0
	for i := range containers {
		if containers[i].Role == STANDBY {
			standbyCnt++
		}
	}

	logit.Info.Printf("standbyCnt ends at %d\n", standbyCnt)

	//provision new container
	params := new(cpmserverapi.DockerRunRequest)
	params.Image = "cpm-node"
	//TODO make the server choice smart
	params.ServerID = containers[0].ServerID
	params.ProjectID = cluster.ProjectID
	params.ContainerName = cluster.Name + "-" + STANDBY + "-" + fmt.Sprintf("%d", standbyCnt)
	params.Standalone = "false"
	var standby = true
	var PROFILE = "LG"

	logit.Info.Printf("here with ProjectID %s\n", cluster.ProjectID)

	err = provisionImpl(dbConn, params, PROFILE, standby)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = provisionImplInit(dbConn, params, PROFILE, false)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//need to update the new container's ClusterID
	var node admindb.Container
	node, err = admindb.GetContainerByName(dbConn, params.ContainerName)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, "error"+err.Error(), http.StatusBadRequest)
		return
	}

	node.ClusterID = cluster.ID
	node.Role = STANDBY
	err = admindb.UpdateContainer(dbConn, node)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, "error"+err.Error(), http.StatusBadRequest)
		return
	}

	err = configureCluster(dbConn, cluster, false)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	results, err := admindb.GetCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println("GetCluster:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	cluster := Cluster{results.ID, results.ProjectID, results.Name, results.ClusterType,
		results.Status, results.CreateDate, "", results.Containers}
	logit.Info.Println("GetCluser:db call results=" + results.ID)

	w.WriteJson(&cluster)
}

func ConfigureCluster(w rest.ResponseWriter, r *rest.Request) {

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println("ConfigureCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	cluster, err := admindb.GetCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println("ConfigureCluster:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = configureCluster(dbConn, cluster, false)
	if err != nil {
		logit.Error.Println("ConfigureCluster:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func configureCluster(dbConn *sql.DB, cluster admindb.Cluster, autocluster bool) error {
	logit.Info.Println("configureCluster:GetCluster")

	//get master node for this cluster
	master, err := admindb.GetContainerMaster(dbConn, cluster.ID)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	logit.Info.Println("configureCluster:GetContainerMaster")

	//configure master postgresql.conf file
	var data string
	if cluster.ClusterType == "synchronous" {
		data, err = template.Postgresql("master", pgport.Value, "*")
	} else {
		data, err = template.Postgresql("master", pgport.Value, "")
	}
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	logit.Info.Println("configureCluster:master postgresql.conf generated")

	//write master postgresql.conf file remotely
	_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/postgresql.conf", data, master.Name)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	logit.Info.Println("configureCluster:master postgresql.conf copied to remote")

	//get domain name
	var domainname admindb.Setting
	domainname, err = admindb.GetSetting(dbConn, "DOMAIN-NAME")
	if err != nil {
		logit.Error.Println("configureCluster: DOMAIN-NAME err " + err.Error())
		return err
	}

	//configure master pg_hba.conf file
	rules := make([]template.Rule, 0)
	data, err = template.Hba(dbConn, "master", master.Name, pgport.Value, cluster.ID, domainname.Value, rules)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	logit.Info.Println("configureCluster:master pg_hba.conf generated")

	//write master pg_hba.conf file remotely
	_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/pg_hba.conf", data, master.Name)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	logit.Info.Println("configureCluster:master pg_hba.conf copied remotely")

	//restart postgres after the config file changes
	var stopResp cpmcontainerapi.StopPGResponse
	stopResp, err = cpmcontainerapi.StopPGClient(master.Name)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}
	logit.Info.Println("configureCluster: master stoppg output was" + stopResp.Output)

	var startResp cpmcontainerapi.StartPGResponse
	startResp, err = cpmcontainerapi.StartPGClient(master.Name)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}
	logit.Info.Println("configureCluster:master startpg output was" + startResp.Output)

	//sleep loop until the master's PG can respond
	var found = false
	var currentStatus string
	var masterhost = master.Name
	for i := 0; i < 20; i++ {
		currentStatus, err = GetPGStatus2(dbConn, master.Name, masterhost)
		if currentStatus == "RUNNING" {
			logit.Info.Println("master is running...continuing")
			found = true
			break
		} else {
			logit.Info.Println("sleeping 1 sec waiting on master..")
			time.Sleep(1000 * time.Millisecond)
		}
	}
	if !found {
		logit.Info.Println("configureCluster: timed out waiting on master pg to start")
		return errors.New("timeout waiting for master pg to respond")
	}

	standbynodes, err2 := admindb.GetAllStandbyContainers(dbConn, cluster.ID)
	if err2 != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}
	//configure all standby nodes
	var stopPGResp cpmcontainerapi.StopPGResponse
	i := 0
	for i = range standbynodes {
		if standbynodes[i].Role == STANDBY {

			//stop standby
			if !autocluster {
				stopPGResp, err = cpmcontainerapi.StopPGClient(standbynodes[i].Name)
				if err != nil {
					logit.Error.Println("configureCluster:" + err.Error())
					return err
				}
				logit.Info.Println("configureCluster:stop output was" + stopPGResp.Output)
			}

			//create base backup from master
			var backupresp cpmcontainerapi.BasebackupResponse
			backupresp, err = cpmcontainerapi.BasebackupClient(masterhost+"."+domainname.Value, standbynodes[i].Name)
			if err != nil {
				logit.Error.Println("configureCluster:" + err.Error())
				return err
			}
			logit.Info.Println("configureCluster:basebackup output was" + backupresp.Output)

			data, err = template.Recovery(masterhost, pgport.Value, "postgres")
			if err != nil {
				logit.Error.Println("configureCluster:" + err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby recovery.conf generated")

			//write standby recovery.conf file remotely
			_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/recovery.conf", data, standbynodes[i].Name)
			if err != nil {
				logit.Error.Println("configureCluster:" + err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby recovery.conf copied remotely")

			data, err = template.Postgresql(STANDBY, pgport.Value, "")
			if err != nil {
				logit.Error.Println("configureCluster:" + err.Error())
				return err
			}

			//write standby postgresql.conf file remotely
			_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/postgresql.conf", data, standbynodes[i].Name)
			if err != nil {
				logit.Error.Println("configureCluster:" + err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby postgresql.conf copied remotely")

			//configure standby pg_hba.conf file
			data, err = template.Hba(dbConn, STANDBY, standbynodes[i].Name, pgport.Value, cluster.ID, domainname.Value, rules)
			if err != nil {
				logit.Error.Println("configureCluster:" + err.Error())
				return err
			}

			logit.Info.Println("configureCluster:standby pg_hba.conf generated")

			//write standby pg_hba.conf file remotely
			_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/pg_hba.conf", data, standbynodes[i].Name)
			if err != nil {
				logit.Error.Println("configureCluster:" + err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby pg_hba.conf copied remotely")

			//start standby

			var stResp cpmcontainerapi.StartPGOnStandbyResponse
			stResp, err = cpmcontainerapi.StartPGOnStandbyClient(standbynodes[i].Name)
			if err != nil {
				logit.Error.Println("configureCluster:" + err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby startpg output was" + stResp.Output)
		}
		i++
	}

	logit.Info.Println("configureCluster: sleeping 5 seconds before configuring pgpool...")
	time.Sleep(5000 * time.Millisecond)

	pgpoolNode, err4 := admindb.GetContainerPgpool(dbConn, cluster.ID)
	logit.Info.Println("configureCluster: lookup pgpool node")
	if err4 != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}
	logit.Info.Println("configureCluster:" + pgpoolNode.Name)

	//configure the pgpool includes all standby nodes AND the master node
	poolnames := make([]string, len(standbynodes)+1)

	i = 0
	for i = range standbynodes {
		poolnames[i] = standbynodes[i].Name + "." + domainname.Value
		i++
	}
	poolnames[i] = master.Name + "." + domainname.Value

	//generate pgpool.conf HOST_LIST
	data, err = template.Poolconf(poolnames)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	logit.Info.Println("configureCluster:pgpool pgpool.conf generated")

	//write pgpool.conf to remote pool node
	_, err = cpmcontainerapi.RemoteWritefileClient(util.GetBase()+"/bin/"+"pgpool.conf", data, pgpoolNode.Name)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}
	logit.Info.Println("configureCluster:pgpool pgpool.conf copied remotely")

	//generate pool_passwd
	data, err = template.Poolpasswd()
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	logit.Info.Println("configureCluster:pgpool pool_passwd generated")

	//write pgpool.conf to remote pool node
	_, err = cpmcontainerapi.RemoteWritefileClient(util.GetBase()+"/bin/"+"pool_passwd", data, pgpoolNode.Name)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}
	logit.Info.Println("configureCluster:pgpool pool_passwd copied remotely")

	//generate pool_hba.conf
	data, err = template.Poolhba()
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	logit.Info.Println("configureCluster:pgpool pool_hba generated")

	//write pgpool.conf to remote pool node
	_, err = cpmcontainerapi.RemoteWritefileClient(util.GetBase()+"/bin/"+"pool_hba.conf", data, pgpoolNode.Name)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}
	logit.Info.Println("configureCluster:pgpool pool_hba copied remotely")

	//start pgpool
	var startPoolResp cpmcontainerapi.StartPgpoolResponse
	startPoolResp, err = cpmcontainerapi.StartPgpoolClient(pgpoolNode.Name)
	if err != nil {
		logit.Error.Println("configureCluster: " + err.Error())
		return err
	}
	logit.Info.Println("configureCluster: pgpool startpgpool output was" + startPoolResp.Output)

	//finally, update the cluster to show that it is
	//initialized!
	cluster.Status = "initialized"
	err = admindb.UpdateCluster(dbConn, cluster)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	return nil

}

func GetAllClustersForProject(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllClustersForProject: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("GetAllClustersForProject: error project ID required")
		rest.Error(w, "project ID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetAllClustersForProject(dbConn, ID)
	if err != nil {
		logit.Error.Println("GetAllClustersForProject: error-" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	clusters := make([]Cluster, len(results))
	i := 0
	for i = range results {
		clusters[i].ID = results[i].ID
		clusters[i].ProjectID = results[i].ProjectID
		clusters[i].Name = results[i].Name
		clusters[i].ClusterType = results[i].ClusterType
		clusters[i].Status = results[i].Status
		clusters[i].CreateDate = results[i].CreateDate
		clusters[i].Containers = results[i].Containers
		i++
	}

	w.WriteJson(&clusters)
}

func GetAllClusters(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllClusters: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err := admindb.GetAllClusters(dbConn)
	if err != nil {
		logit.Error.Println("GetAllClusters: error-" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	clusters := make([]Cluster, len(results))
	i := 0
	for i = range results {
		clusters[i].ID = results[i].ID
		clusters[i].ProjectID = results[i].ProjectID
		clusters[i].Name = results[i].Name
		clusters[i].ClusterType = results[i].ClusterType
		clusters[i].Status = results[i].Status
		clusters[i].CreateDate = results[i].CreateDate
		clusters[i].Containers = results[i].Containers
		i++
	}

	w.WriteJson(&clusters)
}

//we use POST for both updating and inserting based on the ID passed in
func PostCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	logit.Info.Println("PostCluster: in PostCluster")
	cluster := Cluster{}
	err = r.DecodeJsonPayload(&cluster)
	if err != nil {
		logit.Error.Println("PostCluster: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, cluster.Token, "perm-cluster")
	if err != nil {
		logit.Error.Println("PostCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if cluster.Name == "" {
		logit.Error.Println("PostCluster: error in Name")
		rest.Error(w, "cluster name required", http.StatusBadRequest)
		return
	}

	logit.Info.Println("PostCluster: have ID=" + cluster.ID + " Name=" + cluster.Name + " type=" + cluster.ClusterType + " status=" + cluster.Status)
	dbcluster := admindb.Cluster{cluster.ID, cluster.ProjectID, cluster.Name, cluster.ClusterType, cluster.Status, "", cluster.Containers}
	if cluster.ID == "" {
		strid, err := admindb.InsertCluster(dbConn, dbcluster)
		newid := strconv.Itoa(strid)
		if err != nil {
			logit.Error.Println("PostCluster:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cluster.ID = newid
	} else {
		logit.Info.Println("PostCluster: about to call UpdateCluster")
		err2 := admindb.UpdateCluster(dbConn, dbcluster)
		if err2 != nil {
			logit.Error.Println("PostCluster: error in UpdateCluster " + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteJson(&cluster)
}

func DeleteCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println("DeleteCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("DeleteCluster: error cluster ID required")
		rest.Error(w, "cluster ID required", http.StatusBadRequest)
		return
	}

	cluster, err := admindb.GetCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println("DeleteCluster:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	//delete docker containers
	containers, err := admindb.GetAllContainersForCluster(dbConn, cluster.ID)
	if err != nil {
		logit.Error.Println("DeleteCluster:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	i := 0

	i = 0
	server := admindb.Server{}
	for i = range containers {

		//go get the docker server IPAddress
		server, err = admindb.GetServer(dbConn, containers[i].ServerID)
		if err != nil {
			logit.Error.Println("DeleteCluster:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logit.Info.Println("DeleteCluster: got server IP " + server.IPAddress)

		//it is possible that someone can remove a container
		//outside of us, so we let it pass that we can't remove
		//it
		//err = removeContainer(server.IPAddress, containers[i].Name)
		dremreq := &cpmserverapi.DockerRemoveRequest{}
		dremreq.ContainerName = containers[i].Name
		logit.Info.Println("will attempt to delete container " + dremreq.ContainerName)
		var url = "http://" + server.IPAddress + ":10001"
		_, err = cpmserverapi.DockerRemoveClient(url, dremreq)
		if err != nil {
			logit.Error.Println("DeleteCluster: error when trying to remove container" + err.Error())
		}

		//send the server a deletevolume command
		url = "http://" + server.IPAddress + ":10001"
		ddreq := &cpmserverapi.DiskDeleteRequest{}
		ddreq.Path = server.PGDataPath + "/" + containers[i].Name
		_, err = cpmserverapi.DiskDeleteClient(url, ddreq)
		if err != nil {
			logit.Error.Println("DeleteCluster: error when trying to remove disk volume" + err.Error())
		}

		i++
	}

	//delete the container entries
	//delete the cluster entry
	admindb.DeleteCluster(dbConn, ID)

	for i = range containers {

		err = admindb.DeleteContainer(dbConn, containers[i].ID)
		if err != nil {
			logit.Error.Println("DeleteCluster:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

func AdminFailover(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println("AdminFailover: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("AdminFailover: node ID required error")
		rest.Error(w, "node ID required", http.StatusBadRequest)
		return
	}

	//dbNode is the standby node we are going to fail over and
	//make the new master in the cluster
	var dbNode admindb.Container
	dbNode, err = admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cluster, err := admindb.GetCluster(dbConn, dbNode.ClusterID)
	if err != nil {
		logit.Error.Println("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var failoverResp cpmcontainerapi.FailoverResponse
	failoverResp, err = cpmcontainerapi.FailoverClient(dbNode.Name)
	if err != nil {
		logit.Error.Println("AdminFailover: fail-over error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logit.Info.Println("AdminFailover: fail-over output " + failoverResp.Output)

	//update the old master to standalone role
	oldMaster := admindb.Container{}
	oldMaster, err = admindb.GetContainerMaster(dbConn, dbNode.ClusterID)
	if err != nil {
		logit.Error.Println("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	oldMaster.Role = "standalone"
	oldMaster.ClusterID = "-1"
	err = admindb.UpdateContainer(dbConn, oldMaster)
	if err != nil {
		logit.Error.Println("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//update the failover node to master role
	dbNode.Role = "master"
	err = admindb.UpdateContainer(dbConn, dbNode)
	if err != nil {
		logit.Error.Println("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//stop pg on the old master
	//params.IPAddress1 = oldMaster.IPAddress
	var stopPGResp cpmcontainerapi.StopPGResponse
	stopPGResp, err = cpmcontainerapi.StopPGClient(oldMaster.Name)
	if err != nil {
		logit.Error.Println("AdminFailover: " + err.Error() + stopPGResp.Output)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = configureCluster(dbConn, cluster, false)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

	return
}

func EventJoinCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println("EventJoinCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	IDList := r.PathParam("IDList")
	if IDList == "" {
		logit.Error.Println("EventJoinCluster: error IDList required")
		rest.Error(w, "IDList required", http.StatusBadRequest)
		return
	} else {
		logit.Info.Println("EventJoinCluster: IDList=[" + IDList + "]")
	}

	MasterID := r.PathParam("MasterID")
	if MasterID == "" {
		logit.Error.Println("EventJoinCluster: error MasterID required")
		rest.Error(w, "MasterID required", http.StatusBadRequest)
		return
	} else {
		logit.Info.Println("EventJoinCluster: MasterID=[" + MasterID + "]")
	}
	ClusterID := r.PathParam("ClusterID")
	if ClusterID == "" {
		logit.Error.Println("EventJoinCluster: error ClusterID required")
		rest.Error(w, "node ClusterID required", http.StatusBadRequest)
		return
	} else {
		logit.Info.Println("EventJoinCluster: ClusterID=[" + ClusterID + "]")
	}

	var idList = strings.Split(IDList, "_")
	i := 0
	pgpoolCount := 0

	origDBNode := admindb.Container{}
	for i = range idList {
		if idList[i] != "" {
			logit.Info.Println("EventJoinCluster: idList[" + strconv.Itoa(i) + "]=" + idList[i])
			origDBNode, err = admindb.GetContainer(dbConn, idList[i])
			if err != nil {
				logit.Error.Println("EventJoinCluster:" + err.Error())
				rest.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			//update the node to be in the cluster
			origDBNode.ClusterID = ClusterID
			if origDBNode.Image == "cpm-node" {
				origDBNode.Role = STANDBY
			} else {
				origDBNode.Role = "pgpool"
				pgpoolCount++
			}

			if pgpoolCount > 1 {
				logit.Error.Println("EventJoinCluster: more than 1 pgpool is in the cluster")
				rest.Error(w, "only 1 pgpool is allowed in a cluster", http.StatusBadRequest)
				return
			}

			err = admindb.UpdateContainer(dbConn, origDBNode)
			if err != nil {
				logit.Error.Println("EventJoinCluster:" + err.Error())
				rest.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		i++
	}

	//we use the -1 value to indicate that we are only adding
	//to an existing cluster, the UI doesn't know who the master
	//is at this point
	if MasterID != "-1" {
		//update the master node
		origDBNode, err = admindb.GetContainer(dbConn, MasterID)
		if err != nil {
			logit.Error.Println("EventJoinCluster:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		origDBNode.ClusterID = ClusterID
		origDBNode.Role = "master"
		err = admindb.UpdateContainer(dbConn, origDBNode)
		if err != nil {
			logit.Error.Println("EventJoinCluster:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AutoCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	logit.Info.Println("AUTO CLUSTER PROFILE starts")
	logit.Info.Println("AutoCluster: start AutoCluster")
	params := AutoClusterInfo{}
	err = r.DecodeJsonPayload(&params)
	if err != nil {
		logit.Error.Println("AutoCluster: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, params.Token, "perm-cluster")
	if err != nil {
		logit.Error.Println("AutoCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if params.Name == "" {
		logit.Error.Println("AutoCluster: error in Name")
		rest.Error(w, "cluster name required", http.StatusBadRequest)
		return
	}
	if params.ClusterType == "" {
		logit.Error.Println("AutoCluster: error in ClusterType")
		rest.Error(w, "ClusterType name required", http.StatusBadRequest)
		return
	}
	if params.ProjectID == "" {
		logit.Error.Println("AutoCluster: error in ProjectID")
		rest.Error(w, "ProjectID name required", http.StatusBadRequest)
		return
	}
	if params.ClusterProfile == "" {
		logit.Error.Println("AutoCluster: error in ClusterProfile")
		rest.Error(w, "ClusterProfile name required", http.StatusBadRequest)
		return
	}

	logit.Info.Println("AutoCluster: Name=" + params.Name + " ClusterType=" + params.ClusterType + " Profile=" + params.ClusterProfile + " ProjectID=" + params.ProjectID)

	//create cluster definition
	dbcluster := admindb.Cluster{"", params.ProjectID, params.Name, params.ClusterType, "uninitialized", "", make(map[string]string)}
	var ival int
	ival, err = admindb.InsertCluster(dbConn, dbcluster)
	clusterID := strconv.Itoa(ival)
	dbcluster.ID = clusterID
	logit.Info.Println(clusterID)
	if err != nil {
		logit.Error.Println("AutoCluster:" + err.Error())
		rest.Error(w, "Insert Cluster error:"+err.Error(), http.StatusBadRequest)
		return
	}

	//lookup profile
	profile, err2 := getClusterProfileInfo(dbConn, params.ClusterProfile)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-" + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	var masterServer admindb.Server
	var chosenServers []admindb.Server
	if profile.Algo == "round-robin" {
		masterServer, chosenServers, err2 = roundRobin(dbConn, profile)
	} else {
		logit.Error.Println("AutoCluster: error-unsupported algorithm request")
		rest.Error(w, "AutoCluster error: unsupported algorithm", http.StatusBadRequest)
		return
	}

	//create master container
	dockermaster := cpmserverapi.DockerRunRequest{}
	dockermaster.Image = "cpm-node"
	dockermaster.ContainerName = params.Name + "-master"
	dockermaster.ServerID = masterServer.ID
	dockermaster.ProjectID = params.ProjectID
	dockermaster.Standalone = "false"
	if err != nil {
		logit.Error.Println("AutoCluster: error-create master node " + err.Error())
		rest.Error(w, "AutoCluster error"+err.Error(), http.StatusBadRequest)
		return
	}

	//	provision the master
	err2 = provisionImpl(dbConn, &dockermaster, profile.MasterProfile, false)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-provision master " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("AUTO CLUSTER PROFILE master container created")
	var node admindb.Container
	//update node with cluster iD
	node, err2 = admindb.GetContainerByName(dbConn, dockermaster.ContainerName)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-get node by name " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	node.ClusterID = clusterID
	node.Role = "master"
	err2 = admindb.UpdateContainer(dbConn, node)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-update standby node " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	//create standby containers
	var count int
	count, err2 = strconv.Atoi(profile.Count)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	dockerstandby := make([]cpmserverapi.DockerRunRequest, count)
	for i := 0; i < count; i++ {
		logit.Info.Println("working on standby ....")
		//	loop - provision standby
		dockerstandby[i].ServerID = chosenServers[i].ID
		dockerstandby[i].ProjectID = params.ProjectID
		dockerstandby[i].Image = "cpm-node"
		dockerstandby[i].ContainerName = params.Name + "-" + STANDBY + "-" + strconv.Itoa(i)
		dockerstandby[i].Standalone = "false"
		err2 = provisionImpl(dbConn, &dockerstandby[i], profile.StandbyProfile, true)
		if err2 != nil {
			logit.Error.Println("AutoCluster: error-provision master " + err2.Error())
			rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
			return
		}

		//update node with cluster iD
		node, err2 = admindb.GetContainerByName(dbConn, dockerstandby[i].ContainerName)
		if err2 != nil {
			logit.Error.Println("AutoCluster: error-get node by name " + err2.Error())
			rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
			return
		}

		node.ClusterID = clusterID
		node.Role = STANDBY
		err2 = admindb.UpdateContainer(dbConn, node)
		if err2 != nil {
			logit.Error.Println("AutoCluster: error-update standby node " + err2.Error())
			rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
			return
		}
	}
	logit.Info.Println("AUTO CLUSTER PROFILE standbys created")
	//create pgpool container
	//	provision
	dockerpgpool := cpmserverapi.DockerRunRequest{}
	dockerpgpool.ContainerName = params.Name + "-pgpool"
	dockerpgpool.Image = "cpm-pgpool"
	dockerpgpool.ServerID = chosenServers[count].ID
	dockerpgpool.ProjectID = params.ProjectID
	dockerpgpool.Standalone = "false"

	err2 = provisionImpl(dbConn, &dockerpgpool, profile.StandbyProfile, true)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-provision pgpool " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}
	logit.Info.Println("AUTO CLUSTER PROFILE pgpool created")
	//update node with cluster ID
	node, err2 = admindb.GetContainerByName(dbConn, dockerpgpool.ContainerName)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-get pgpool node by name " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	node.ClusterID = clusterID
	node.Role = "pgpool"
	err2 = admindb.UpdateContainer(dbConn, node)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-update pgpool node " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	//init the master DB
	//	provision the master
	err2 = provisionImplInit(dbConn, &dockermaster, profile.MasterProfile, false)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-provisionInit master " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	//make sure every node is ready
	err2 = waitTillAllReady(dockermaster, dockerpgpool, dockerstandby)
	if err2 != nil {
		logit.Error.Println("cluster members not responding in time")
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	//configure cluster
	//	ConfigureCluster
	logit.Info.Println("AUTO CLUSTER PROFILE configure cluster ")
	err2 = configureCluster(dbConn, dbcluster, true)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-configure cluster " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("AUTO CLUSTER PROFILE done")
	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

// round-robin provisioning algorithm -
//  to promote least used servers, incoming servers list
//  should be sorted by class and least used order
//  returns the master server and the list of standby servers
func roundRobin(dbConn *sql.DB, profile ClusterProfiles) (admindb.Server, []admindb.Server, error) {
	var masterServer admindb.Server
	count, err := strconv.Atoi(profile.Count)

	//add 1 to the standby count to make room for the pgpool node
	count++

	//create a slice to hold servers for standby and pgpool nodes
	//assumes 1 pgpool node per cluster which is enforced by auto-cluster
	chosen := make([]admindb.Server, count)

	//get all the servers available
	servers, err := admindb.GetAllServersByClassByCount(dbConn)
	if err != nil {
		return masterServer, chosen, err
	}
	if len(servers) == 0 {
		return masterServer, chosen, errors.New("no servers defined")
	}

	//find the server for the master
	//search from last used to end of servers list
	found := false
	for j := 0; j < len(servers); j++ {
		if profile.MasterServer == servers[j].ServerClass {
			found = true
			masterServer = servers[j]
			break
		}
	}

	//give up on finding a match and use any server
	if !found {
		for j := 0; j < len(servers); j++ {
			masterServer = servers[j]
			break
		}
	}

	//find the servers for all the other nodes (standby, pgpool)
	//avoiding the use of the masterServer for HA

	lastused := 0

	for i := 0; i < count; i++ {

		found = false

		//search from last used to end of servers list
		for j := lastused; j < len(servers); j++ {
			if servers[j].ID != masterServer.ID &&
				servers[j].ServerClass == profile.StandbyServer {
				chosen[i] = servers[j]
				found = true
				lastused = j
				break
			}
		}

		if !found {
			//search from start of servers list to end
			for j := 0; j < len(servers); j++ {
				if servers[j].ID != masterServer.ID && servers[j].ServerClass == profile.StandbyServer {
					chosen[i] = servers[j]
					found = true
					lastused = j
					break
				}
			}

		}

		//if still not found, use any server
		if !found {
			//search from start of servers list to end
			for j := 0; j < len(servers); j++ {
				chosen[i] = servers[j]
				found = true
				lastused = j
				break
			}

		}

	}

	logit.Info.Println("round-robin: master " + masterServer.Name + " class=" + masterServer.ServerClass)
	for x := 0; x < len(chosen); x++ {
		logit.Info.Println("round-robin: other servers " + chosen[x].Name + " class=" + chosen[x].ServerClass)
	}
	return masterServer, chosen, nil
}

func waitTillAllReady(dockermaster cpmserverapi.DockerRunRequest, dockerpgpool cpmserverapi.DockerRunRequest, dockerstandby []cpmserverapi.DockerRunRequest) error {
	err := waitTillReady(dockermaster.ContainerName)
	if err != nil {
		logit.Error.Println("time out waiting for " + dockermaster.ContainerName)
		return err
	}
	err = waitTillReady(dockerpgpool.ContainerName)
	if err != nil {
		logit.Error.Println("time out waiting for " + dockerpgpool.ContainerName)
		return err
	}
	for x := 0; x < len(dockerstandby); x++ {
		err = waitTillReady(dockerstandby[x].ContainerName)
		if err != nil {
			logit.Error.Println("time out waiting for " + dockerstandby[x].ContainerName)
			return err
		}
	}
	return nil

}
