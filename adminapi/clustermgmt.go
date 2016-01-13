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
	"github.com/crunchydata/crunchy-postgresql-manager/swarmapi"
	"github.com/crunchydata/crunchy-postgresql-manager/template"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
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

// ScaleUpCluster increases the count of standby containers in a cluster
func ScaleUpCluster(w rest.ResponseWriter, r *rest.Request) {

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
	cluster, err := admindb.GetCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var containers []types.Container
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

	//logit.Info.Printf("standbyCnt ends at %d\n", standbyCnt)

	//provision new container
	params := new(swarmapi.DockerRunRequest)
	params.Image = "cpm-node"
	//TODO make the server choice smart
	params.ProjectID = cluster.ProjectID
	params.ContainerName = cluster.Name + "-" + STANDBY + "-" + fmt.Sprintf("%d", standbyCnt)
	params.Standalone = "false"
	var standby = true
	params.Profile = "LG"

	//logit.Info.Printf("here with ProjectID %s\n", cluster.ProjectID)

	_, err = provisionImpl(dbConn, params, standby)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = provisionImplInit(dbConn, params, false)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//need to update the new container's ClusterID
	var node types.Container
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

	err = configureCluster(params.Profile, dbConn, cluster, false)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

// GetCluster returns a given cluster definition
func GetCluster(w rest.ResponseWriter, r *rest.Request) {
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
	results, err := admindb.GetCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	cluster := types.Cluster{}
	cluster.ID = results.ID
	cluster.ProjectID = results.ProjectID
	cluster.Name = results.Name
	cluster.ClusterType = results.ClusterType
	cluster.Status = results.Status
	cluster.CreateDate = results.CreateDate
	cluster.Containers = results.Containers
	//logit.Info.Println("GetCluser:db call results=" + results.ID)

	w.WriteJson(&cluster)
}

// TODO
func ConfigureCluster(w rest.ResponseWriter, r *rest.Request) {

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println(err.Error())
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

	err = configureCluster("SM", dbConn, cluster, false)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func configureCluster(profile string, dbConn *sql.DB, cluster types.Cluster, autocluster bool) error {
	logit.Info.Println("configureCluster:GetCluster")

	//get master node for this cluster
	master, err := admindb.GetContainerMaster(dbConn, cluster.ID)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	var pgport types.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	var sleepSetting types.Setting
	sleepSetting, err = admindb.GetSetting(dbConn, "SLEEP-PROV")
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	//logit.Info.Println("configureCluster:GetContainerMaster")

	//configure master postgresql.conf file
	var data string
	info := new(template.PostgresqlParameters)
	info.PG_PORT = pgport.Value
	err = template.GetTuningParms(dbConn, profile, info)
	if err != nil {
		logit.Error.Println("configureCluster:" + err.Error())
		return err
	}

	if cluster.ClusterType == "synchronous" {
		info.CLUSTER_TYPE = "*"
		data, err = template.Postgresql("master", info)
	} else {
		//TODO verify these next 2 lines look erroneous
		info.CLUSTER_TYPE = "*"
		data, err = template.Postgresql("master", info)

		//info.CLUSTER_TYPE = ""
		//data, err = template.Postgresql("master", info)
	}
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("configureCluster:master postgresql.conf generated")

	//write master postgresql.conf file remotely
	_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/postgresql.conf", data, master.Name)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("configureCluster:master postgresql.conf copied to remote")

	//get domain name
	var domainname types.Setting
	domainname, err = admindb.GetSetting(dbConn, "DOMAIN-NAME")
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//configure master pg_hba.conf file
	rules := make([]template.Rule, 0)
	data, err = template.Hba(dbConn, "master", master.Name, pgport.Value, cluster.ID, domainname.Value, rules)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("configureCluster:master pg_hba.conf generated")

	//write master pg_hba.conf file remotely
	_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/pg_hba.conf", data, master.Name)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("configureCluster:master pg_hba.conf copied remotely")

	//restart postgres after the config file changes
	var stopResp cpmcontainerapi.StopPGResponse
	stopResp, err = cpmcontainerapi.StopPGClient(master.Name)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	logit.Info.Println("configureCluster: master stoppg output was" + stopResp.Output)

	var startResp cpmcontainerapi.StartPGResponse
	startResp, err = cpmcontainerapi.StartPGClient(master.Name)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	logit.Info.Println("configureCluster:master startpg output was" + startResp.Output)

	//sleep loop until the master's PG can respond
	var sleepTime time.Duration
	sleepTime, err = time.ParseDuration(sleepSetting.Value)
	var found = false
	var currentStatus string
	var masterhost = master.Name
	for i := 0; i < 20; i++ {
		//currentStatus, err = GetPGStatus2(dbConn, master.Name, masterhost)
		currentStatus, err = util.FastPing(pgport.Value, master.Name)
		if currentStatus == "RUNNING" {
			logit.Info.Println("master is running...continuing")
			found = true
			break
		} else {
			logit.Info.Println("sleeping waiting on master..")
			time.Sleep(sleepTime)
		}
	}
	if !found {
		logit.Info.Println("configureCluster: timed out waiting on master pg to start")
		return errors.New("timeout waiting for master pg to respond")
	}

	standbynodes, err2 := admindb.GetAllStandbyContainers(dbConn, cluster.ID)
	if err2 != nil {
		logit.Error.Println(err.Error())
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
					logit.Error.Println(err.Error())
					logit.Error.Println("configureCluster:stop output was" + stopPGResp.Output)
					return err
				}
				//logit.Info.Println("configureCluster:stop output was" + stopPGResp.Output)
			}

			//
			var username = "cpmtest"
			var password = "cpmtest"

			//create base backup from master
			var backupresp cpmcontainerapi.BasebackupResponse
			backupresp, err = cpmcontainerapi.BasebackupClient(masterhost+"."+domainname.Value, standbynodes[i].Name, username, password)
			if err != nil {
				logit.Error.Println(err.Error())
				logit.Error.Println("configureCluster:basebackup output was" + backupresp.Output)
				return err
			}
			//logit.Info.Println("configureCluster:basebackup output was" + backupresp.Output)

			data, err = template.Recovery(masterhost, pgport.Value, "postgres")
			if err != nil {
				logit.Error.Println(err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby recovery.conf generated")

			//write standby recovery.conf file remotely
			_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/recovery.conf", data, standbynodes[i].Name)
			if err != nil {
				logit.Error.Println(err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby recovery.conf copied remotely")

			info := new(template.PostgresqlParameters)
			info.PG_PORT = pgport.Value
			info.CLUSTER_TYPE = ""
			err = template.GetTuningParms(dbConn, profile, info)
			if err != nil {
				logit.Error.Println(err.Error())
				return err
			}

			data, err = template.Postgresql(STANDBY, info)
			if err != nil {
				logit.Error.Println(err.Error())
				return err
			}

			//write standby postgresql.conf file remotely
			_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/postgresql.conf", data, standbynodes[i].Name)
			if err != nil {
				logit.Error.Println(err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby postgresql.conf copied remotely")

			//configure standby pg_hba.conf file
			data, err = template.Hba(dbConn, STANDBY, standbynodes[i].Name, pgport.Value, cluster.ID, domainname.Value, rules)
			if err != nil {
				logit.Error.Println(err.Error())
				return err
			}

			logit.Info.Println("configureCluster:standby pg_hba.conf generated")

			//write standby pg_hba.conf file remotely
			_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/pg_hba.conf", data, standbynodes[i].Name)
			if err != nil {
				logit.Error.Println(err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby pg_hba.conf copied remotely")

			//start standby

			var stResp cpmcontainerapi.StartPGOnStandbyResponse
			stResp, err = cpmcontainerapi.StartPGOnStandbyClient(standbynodes[i].Name)
			if err != nil {
				logit.Error.Println(err.Error())
				return err
			}
			logit.Info.Println("configureCluster:standby startpg output was" + stResp.Output)
		}
		i++
	}

	logit.Info.Println("configureCluster: sleeping 5 seconds before configuring pgpool...")
	clustersleepTime, _ := time.ParseDuration("5s")
	time.Sleep(clustersleepTime)

	pgpoolNode, err4 := admindb.GetContainerPgpool(dbConn, cluster.ID)
	//logit.Info.Println("configureCluster: lookup pgpool node")
	if err4 != nil {
		logit.Error.Println(err.Error())
		return err
	}
	//logit.Info.Println("configureCluster:" + pgpoolNode.Name)

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
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("configureCluster:pgpool pgpool.conf generated")

	//write pgpool.conf to remote pool node
	_, err = cpmcontainerapi.RemoteWritefileClient(util.GetBase()+"/bin/"+"pgpool.conf", data, pgpoolNode.Name)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	logit.Info.Println("configureCluster:pgpool pgpool.conf copied remotely")

	//generate pool_passwd
	data, err = template.Poolpasswd()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("configureCluster:pgpool pool_passwd generated")

	//write pgpool.conf to remote pool node
	_, err = cpmcontainerapi.RemoteWritefileClient(util.GetBase()+"/bin/"+"pool_passwd", data, pgpoolNode.Name)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	logit.Info.Println("configureCluster:pgpool pool_passwd copied remotely")

	//g enerate pool_hba.conf
	cars := make([]template.Rule, 0)
	data, err = template.Poolhba(cars)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("configureCluster:pgpool pool_hba generated")

	//write pgpool.conf to remote pool node
	_, err = cpmcontainerapi.RemoteWritefileClient(util.GetBase()+"/bin/"+"pool_hba.conf", data, pgpoolNode.Name)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	logit.Info.Println("configureCluster:pgpool pool_hba copied remotely")

	//start pgpool
	var startPoolResp cpmcontainerapi.StartPgpoolResponse
	startPoolResp, err = cpmcontainerapi.StartPgpoolClient(pgpoolNode.Name)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	logit.Info.Println("configureCluster: pgpool startpgpool output was" + startPoolResp.Output)

	//finally, update the cluster to show that it is
	//initialized!
	cluster.Status = "initialized"
	err = admindb.UpdateCluster(dbConn, cluster)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	return nil

}

// GetAllClustersForProject returns a list of cluster definitions for a given project
func GetAllClustersForProject(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("project ID required")
		rest.Error(w, "project ID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetAllClustersForProject(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	clusters := make([]types.Cluster, len(results))
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

// PostCluster updates or inserts a new cluster definition
func PostCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	//logit.Info.Println("PostCluster: in PostCluster")
	cluster := types.Cluster{}
	err = r.DecodeJsonPayload(&cluster)
	if err != nil {
		logit.Error.Println("error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, cluster.Token, "perm-cluster")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if cluster.Name == "" {
		logit.Error.Println("PostCluster: error in Name")
		rest.Error(w, "cluster name required", http.StatusBadRequest)
		return
	}

	//logit.Info.Println("PostCluster: have ID=" + cluster.ID + " Name=" + cluster.Name + " type=" + cluster.ClusterType + " status=" + cluster.Status)
	dbcluster := types.Cluster{}
	dbcluster.ID = cluster.ID
	dbcluster.ProjectID = cluster.ProjectID
	dbcluster.Name = cluster.Name
	dbcluster.ClusterType = cluster.ClusterType
	dbcluster.Status = cluster.Status
	dbcluster.Containers = cluster.Containers
	if cluster.ID == "" {
		strid, err := admindb.InsertCluster(dbConn, dbcluster)
		newid := strconv.Itoa(strid)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cluster.ID = newid
	} else {
		//logit.Info.Println("PostCluster: about to call UpdateCluster")
		err2 := admindb.UpdateCluster(dbConn, dbcluster)
		if err2 != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteJson(&cluster)
}

// DeleteCluster deletes an existing cluster definition
func DeleteCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("cluster ID required")
		rest.Error(w, "cluster ID required", http.StatusBadRequest)
		return
	}

	cluster, err := admindb.GetCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	//delete docker containers
	containers, err := admindb.GetAllContainersForCluster(dbConn, cluster.ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	var infoResponse swarmapi.DockerInfoResponse

	infoResponse, err = swarmapi.DockerInfo()
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	i := 0
	servers := make([]types.Server, len(infoResponse.Output))
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
	i = 0
	//server := types.Server{}
	for i = range containers {

		//logit.Info.Println("DeleteCluster: got server IP " + server.IPAddress)

		//it is possible that someone can remove a container
		//outside of us, so we let it pass that we can't remove
		//it
		//err = removeContainer(server.IPAddress, containers[i].Name)
		dremreq := &swarmapi.DockerRemoveRequest{}
		dremreq.ContainerName = containers[i].Name
		//logit.Info.Println("will attempt to delete container " + dremreq.ContainerName)
		_, err = swarmapi.DockerRemove(dremreq)
		if err != nil {
			logit.Error.Println("error when trying to remove container" + err.Error())
		}

		//send all the servers a deletevolume command
		ddreq := &cpmserverapi.DiskDeleteRequest{}
		ddreq.Path = pgdatapath.Value + "/" + containers[i].Name
		for _, each := range servers {
			_, err = cpmserverapi.DiskDeleteClient(each.Name, ddreq)
			if err != nil {
				logit.Error.Println("error when trying to remove disk volume" + err.Error())
			}
		}

		i++
	}

	//delete the container entries
	//delete the cluster entry
	admindb.DeleteCluster(dbConn, ID)

	for i = range containers {

		err = admindb.DeleteContainer(dbConn, containers[i].ID)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

// AdminFailover causes a cluster failorver to be performed for a given cluster
func AdminFailover(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println("authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("node ID required error")
		rest.Error(w, "node ID required", http.StatusBadRequest)
		return
	}

	//dbNode is the standby node we are going to fail over and
	//make the new master in the cluster
	var dbNode types.Container
	dbNode, err = admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cluster, err := admindb.GetCluster(dbConn, dbNode.ClusterID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var failoverResp cpmcontainerapi.FailoverResponse
	failoverResp, err = cpmcontainerapi.FailoverClient(dbNode.Name)
	if err != nil {
		logit.Error.Println("fail-over error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logit.Info.Println("AdminFailover: fail-over output " + failoverResp.Output)

	//update the old master to standalone role
	oldMaster := types.Container{}
	oldMaster, err = admindb.GetContainerMaster(dbConn, dbNode.ClusterID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	oldMaster.Role = "standalone"
	oldMaster.ClusterID = "-1"
	err = admindb.UpdateContainer(dbConn, oldMaster)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//update the failover node to master role
	dbNode.Role = "master"
	err = admindb.UpdateContainer(dbConn, dbNode)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//stop pg on the old master
	//params.IPAddress1 = oldMaster.IPAddress
	var stopPGResp cpmcontainerapi.StopPGResponse
	stopPGResp, err = cpmcontainerapi.StopPGClient(oldMaster.Name)
	if err != nil {
		logit.Error.Println(err.Error() + stopPGResp.Output)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = configureCluster("SM", dbConn, cluster, false)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

	return
}

// TODO
func EventJoinCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	IDList := r.PathParam("IDList")
	if IDList == "" {
		logit.Error.Println("IDList required")
		rest.Error(w, "IDList required", http.StatusBadRequest)
		return
	} else {
		logit.Info.Println("EventJoinCluster: IDList=[" + IDList + "]")
	}

	MasterID := r.PathParam("MasterID")
	if MasterID == "" {
		logit.Error.Println("MasterID required")
		rest.Error(w, "MasterID required", http.StatusBadRequest)
		return
	} else {
		logit.Info.Println("EventJoinCluster: MasterID=[" + MasterID + "]")
	}
	ClusterID := r.PathParam("ClusterID")
	if ClusterID == "" {
		logit.Error.Println("ClusterID required")
		rest.Error(w, "node ClusterID required", http.StatusBadRequest)
		return
	} else {
		logit.Info.Println("EventJoinCluster: ClusterID=[" + ClusterID + "]")
	}

	var idList = strings.Split(IDList, "_")
	i := 0
	pgpoolCount := 0

	origDBNode := types.Container{}
	for i = range idList {
		if idList[i] != "" {
			logit.Info.Println("EventJoinCluster: idList[" + strconv.Itoa(i) + "]=" + idList[i])
			origDBNode, err = admindb.GetContainer(dbConn, idList[i])
			if err != nil {
				logit.Error.Println(err.Error())
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
				logit.Error.Println(err.Error())
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
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		origDBNode.ClusterID = ClusterID
		origDBNode.Role = "master"
		err = admindb.UpdateContainer(dbConn, origDBNode)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

// AutoCluster creates a new cluster
func AutoCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	logit.Info.Println("AUTO CLUSTER PROFILE starts")
	params := AutoClusterInfo{}
	err = r.DecodeJsonPayload(&params)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, params.Token, "perm-cluster")
	if err != nil {
		logit.Error.Println(err.Error())
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
	dbcluster := types.Cluster{}
	dbcluster.ID = ""
	dbcluster.ProjectID = params.ProjectID
	dbcluster.Name = util.CleanName(params.Name)
	dbcluster.ClusterType = params.ClusterType
	dbcluster.Status = "uninitialized"
	dbcluster.Containers = make(map[string]string)

	var ival int
	ival, err = admindb.InsertCluster(dbConn, dbcluster)
	clusterID := strconv.Itoa(ival)
	dbcluster.ID = clusterID
	//logit.Info.Println(clusterID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, "Insert Cluster error:"+err.Error(), http.StatusBadRequest)
		return
	}

	//lookup profile
	profile, err2 := getClusterProfileInfo(dbConn, params.ClusterProfile)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	//var masterServer types.Server
	//var chosenServers []types.Server
	if profile.Algo == "round-robin" {
		//masterServer, chosenServers, err2 = roundRobin(dbConn, profile)
	} else {
		logit.Error.Println("AutoCluster: error-unsupported algorithm request")
		rest.Error(w, "AutoCluster error: unsupported algorithm", http.StatusBadRequest)
		return
	}

	//create master container
	dockermaster := swarmapi.DockerRunRequest{}
	dockermaster.Image = "cpm-node"
	dockermaster.ContainerName = params.Name + "-master"
	dockermaster.ProjectID = params.ProjectID
	dockermaster.Standalone = "false"
	dockermaster.Profile = profile.MasterProfile
	if err != nil {
		logit.Error.Println("AutoCluster: error-create master node " + err.Error())
		rest.Error(w, "AutoCluster error"+err.Error(), http.StatusBadRequest)
		return
	}

	//	provision the master
	logit.Info.Println("dockermaster profile is " + dockermaster.Profile)

	_, err2 = provisionImpl(dbConn, &dockermaster, false)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-provision master " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("AUTO CLUSTER PROFILE master container created")
	var node types.Container
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

	var sleepSetting types.Setting
	sleepSetting, err2 = admindb.GetSetting(dbConn, "SLEEP-PROV")
	if err2 != nil {
		logit.Error.Println("SLEEP-PROV setting error " + err2.Error())
		rest.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	var sleepTime time.Duration
	sleepTime, err2 = time.ParseDuration(sleepSetting.Value)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		rest.Error(w, err2.Error(), http.StatusInternalServerError)
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

	dockerstandby := make([]swarmapi.DockerRunRequest, count)
	for i := 0; i < count; i++ {
		logit.Info.Println("working on standby ....")
		//	loop - provision standby
		dockerstandby[i].ProjectID = params.ProjectID
		dockerstandby[i].Image = "cpm-node"
		dockerstandby[i].ContainerName = params.Name + "-" + STANDBY + "-" + strconv.Itoa(i)
		dockerstandby[i].Standalone = "false"
		dockerstandby[i].Profile = profile.StandbyProfile

		_, err2 = provisionImpl(dbConn, &dockerstandby[i], true)
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
	dockerpgpool := swarmapi.DockerRunRequest{}
	dockerpgpool.ContainerName = params.Name + "-pgpool"
	dockerpgpool.Image = "cpm-pgpool"
	dockerpgpool.ProjectID = params.ProjectID
	dockerpgpool.Standalone = "false"
	dockerpgpool.Profile = profile.StandbyProfile

	_, err2 = provisionImpl(dbConn, &dockerpgpool, true)
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
	dockermaster.Profile = profile.MasterProfile
	err2 = provisionImplInit(dbConn, &dockermaster, false)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-provisionInit master " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	//make sure every node is ready
	err2 = waitTillAllReady(dockermaster, dockerpgpool, dockerstandby, sleepTime)
	if err2 != nil {
		logit.Error.Println("cluster members not responding in time")
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	//configure cluster
	//	ConfigureCluster
	logit.Info.Println("AUTO CLUSTER PROFILE configure cluster ")
	err2 = configureCluster(profile.MasterProfile, dbConn, dbcluster, true)
	if err2 != nil {
		logit.Error.Println("AutoCluster: error-configure cluster " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("AUTO CLUSTER PROFILE done")
	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

// round-robin provisioning algorithm -
//  to promote least used servers, incoming servers list
//  should be sorted by class and least used order
//  returns the master server and the list of standby servers
/**
func roundRobin(dbConn *sql.DB, profile types.ClusterProfiles) (types.Server, []types.Server, error) {
	var masterServer types.Server
	count, err := strconv.Atoi(profile.Count)

	//add 1 to the standby count to make room for the pgpool node
	count++

	//create a slice to hold servers for standby and pgpool nodes
	//assumes 1 pgpool node per cluster which is enforced by auto-cluster
	chosen := make([]types.Server, count)

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
*/

func waitTillAllReady(dockermaster swarmapi.DockerRunRequest, dockerpgpool swarmapi.DockerRunRequest, dockerstandby []swarmapi.DockerRunRequest, sleepTime time.Duration) error {
	err := waitTillReady(dockermaster.ContainerName, sleepTime)
	if err != nil {
		logit.Error.Println("time out waiting for " + dockermaster.ContainerName)
		return err
	}
	err = waitTillReady(dockerpgpool.ContainerName, sleepTime)
	if err != nil {
		logit.Error.Println("time out waiting for " + dockerpgpool.ContainerName)
		return err
	}
	for x := 0; x < len(dockerstandby); x++ {
		err = waitTillReady(dockerstandby[x].ContainerName, sleepTime)
		if err != nil {
			logit.Error.Println("time out waiting for " + dockerstandby[x].ContainerName)
			return err
		}
	}
	return nil

}

// StartCluster starts all nodes in a cluster
func StartCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("StartCluster: error cluster ID required")
		rest.Error(w, "cluster ID required", http.StatusBadRequest)
		return
	}

	cluster, err := admindb.GetCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	//start docker containers
	containers, err := admindb.GetAllContainersForCluster(dbConn, cluster.ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	i := 0

	i = 0
	var response swarmapi.DockerStartResponse
	for i = range containers {

		req := &swarmapi.DockerStartRequest{}
		req.ContainerName = containers[i].Name
		logit.Info.Println("will attempt to start container " + req.ContainerName)
		response, err = swarmapi.DockerStart(req)
		if err != nil {
			logit.Error.Println("StartCluster: error when trying to start container" + err.Error())
		}
		logit.Info.Println("StartCluster: started " + response.Output)

		i++
	}

	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

// StopCluster stops all nodes in a cluster
func StopCluster(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-cluster")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("StopCluster: error cluster ID required")
		rest.Error(w, "cluster ID required", http.StatusBadRequest)
		return
	}

	cluster, err := admindb.GetCluster(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	//stop docker containers
	containers, err := admindb.GetAllContainersForCluster(dbConn, cluster.ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}

	i := 0

	i = 0
	var response swarmapi.DockerStopResponse
	//server := types.Server{}
	for i = range containers {

		req := &swarmapi.DockerStopRequest{}
		req.ContainerName = containers[i].Name
		logit.Info.Println("will attempt to stop container " + req.ContainerName)
		response, err = swarmapi.DockerStop(req)
		if err != nil {
			logit.Error.Println("StopCluster: error when trying to stop container" + err.Error())
		}
		logit.Info.Println("StopCluster: started " + response.Output)

		i++
	}

	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}
