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
	"crunchy.com/admindb"
	"crunchy.com/cpmagent"
	"crunchy.com/template"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AutoClusterInfo struct {
	Name           string
	ClusterType    string
	ClusterProfile string
	Token          string
}

func GetCluster(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	results, err := admindb.GetDBCluster(ID)
	if err != nil {
		glog.Errorln("GetCluster:" + err.Error())
		rest.Error(w, err.Error(), 400)
	}
	cluster := Cluster{results.ID, results.Name, results.ClusterType,
		results.Status, results.CreateDate, ""}
	glog.Infoln("GetCluser:db call results=" + results.ID)

	w.WriteJson(&cluster)
}

func ConfigureCluster(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-cluster")
	if err != nil {
		glog.Errorln("ConfigureCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	cluster, err := admindb.GetDBCluster(ID)
	if err != nil {
		glog.Errorln("ConfigureCluster:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	err = configureCluster(cluster, false)
	if err != nil {
		glog.Errorln("ConfigureCluster:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func configureCluster(cluster admindb.DBCluster, autocluster bool) error {
	glog.Infoln("configureCluster:GetDBCluster")

	//get master node for this cluster
	master, err := admindb.GetDBNodeMaster(cluster.ID)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}

	glog.Infoln("configureCluster:GetDBNodeMaster")

	//configure master postgresql.conf file
	var data string
	if cluster.ClusterType == "synchronous" {
		data, err = template.Postgresql("master", "5432", "*")
	} else {
		data, err = template.Postgresql("master", "5432", "")
	}
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}

	glog.Infoln("configureCluster:master postgresql.conf generated")

	//write master postgresql.conf file remotely
	err = RemoteWritefile("/pgdata/postgresql.conf", data, master.Name)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}

	glog.Infoln("configureCluster:master postgresql.conf copied to remote")

	//get domain name
	var domainname admindb.DBSetting
	domainname, err = admindb.GetDBSetting("DOMAIN-NAME")
	if err != nil {
		glog.Errorln("configureCluster: DOMAIN-NAME err " + err.Error())
		return err
	}

	//configure master pg_hba.conf file
	data, err = template.Hba("master", master.Name, "5432", cluster.ID, domainname.Value)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}

	glog.Infoln("configureCluster:master pg_hba.conf generated")

	//write master pg_hba.conf file remotely
	err = RemoteWritefile("/pgdata/pg_hba.conf", data, master.Name)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}

	glog.Infoln("configureCluster:master pg_hba.conf copied remotely")

	//restart postgres after the config file changes
	var commandoutput string
	commandoutput, err = PGCommand("/cluster/bin/stoppg.sh", master.Name)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}
	glog.Infoln("configureCluster: master stoppg output was" + commandoutput)

	commandoutput, err = PGCommand("/cluster/bin/startpg.sh", master.Name)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}
	glog.Infoln("configureCluster:master startpg output was" + commandoutput)

	//sleep loop until the master's PG can respond
	var found = false
	var currentStatus string
	for i := 0; i < 20; i++ {
		currentStatus, err = GetPGStatus(master.Name)
		if currentStatus == "RUNNING" {
			glog.Infoln("master is running...continuing")
			found = true
			break
		} else {
			glog.Infoln("sleeping 1 sec waiting on master..")
			time.Sleep(1000 * time.Millisecond)
		}
	}
	if !found {
		glog.Infoln("configureCluster: timed out waiting on master pg to start")
		return errors.New("timeout waiting for master pg to respond")
	}

	standbynodes, err2 := admindb.GetAllDBStandbyNodes(cluster.ID)
	if err2 != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}
	//configure all standby nodes
	i := 0
	for i = range standbynodes {
		if standbynodes[i].Role == "standby" {

			//stop standby
			if !autocluster {
				commandoutput, err = PGCommand("/cluster/bin/stoppg.sh", standbynodes[i].Name)
				if err != nil {
					glog.Errorln("configureCluster:" + err.Error())
					return err
				}
				glog.Infoln("configureCluster:stop output was" + commandoutput)
			}

			//create base backup from master
			commandoutput, err = cpmagent.Command1("/cluster/bin/basebackup.sh", master.Name+"."+domainname.Value, standbynodes[i].Name)
			if err != nil {
				glog.Errorln("configureCluster:" + err.Error())
				return err
			}
			glog.Infoln("configureCluster:basebackup output was" + commandoutput)

			data, err = template.Recovery(master.Name, "5432", "postgres")
			if err != nil {
				glog.Errorln("configureCluster:" + err.Error())
				return err
			}
			glog.Infoln("configureCluster:standby recovery.conf generated")

			//write standby recovery.conf file remotely
			err = RemoteWritefile("/pgdata/recovery.conf", data, standbynodes[i].Name)
			if err != nil {
				glog.Errorln("configureCluster:" + err.Error())
				return err
			}
			glog.Infoln("configureCluster:standby recovery.conf copied remotely")

			data, err = template.Postgresql("standby", "5432", "")
			if err != nil {
				glog.Errorln("configureCluster:" + err.Error())
				return err
			}

			//write standby postgresql.conf file remotely
			err = RemoteWritefile("/pgdata/postgresql.conf", data, standbynodes[i].Name)
			if err != nil {
				glog.Errorln("configureCluster:" + err.Error())
				return err
			}
			glog.Infoln("configureCluster:standby postgresql.conf copied remotely")

			//configure standby pg_hba.conf file
			data, err = template.Hba("standby", standbynodes[i].Name, "5432", cluster.ID, domainname.Value)
			if err != nil {
				glog.Errorln("configureCluster:" + err.Error())
				return err
			}

			glog.Infoln("configureCluster:standby pg_hba.conf generated")

			//write standby pg_hba.conf file remotely
			err = RemoteWritefile("/pgdata/pg_hba.conf", data, standbynodes[i].Name)
			if err != nil {
				glog.Errorln("configureCluster:" + err.Error())
				return err
			}
			glog.Infoln("configureCluster:standby pg_hba.conf copied remotely")

			//start standby

			commandoutput, err = PGCommand("/cluster/bin/startpgonstandby.sh", standbynodes[i].Name)
			if err != nil {
				glog.Errorln("configureCluster:" + err.Error())
				return err
			}
			glog.Infoln("configureCluster:standby startpg output was" + commandoutput)
		}
		i++
	}

	glog.Infoln("configureCluster: sleeping 5 seconds before configuring pgpool...")
	time.Sleep(5000 * time.Millisecond)

	pgpoolNode, err4 := admindb.GetDBNodePgpool(cluster.ID)
	if err4 != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}
	glog.Infoln("configureCluster:" + pgpoolNode.Name)

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
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}

	glog.Infoln("configureCluster:pgpool pgpool.conf generated")

	//write pgpool.conf to remote pool node
	err = RemoteWritefile("/cluster/bin/pgpool.conf", data, pgpoolNode.Name)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}
	glog.Infoln("configureCluster:pgpool pgpool.conf copied remotely")

	//generate pool_passwd
	data, err = template.Poolpasswd()
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}

	glog.Infoln("configureCluster:pgpool pool_passwd generated")

	//write pgpool.conf to remote pool node
	err = RemoteWritefile("/cluster/bin/pool_passwd", data, pgpoolNode.Name)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}
	glog.Infoln("configureCluster:pgpool pool_passwd copied remotely")

	//generate pool_hba.conf
	data, err = template.Poolhba()
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}

	glog.Infoln("configureCluster:pgpool pool_hba generated")

	//write pgpool.conf to remote pool node
	err = RemoteWritefile("/cluster/bin/pool_hba.conf", data, pgpoolNode.Name)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}
	glog.Infoln("configureCluster:pgpool pool_hba copied remotely")

	//start pgpool
	commandoutput, err = PGCommand("/cluster/bin/startpgpool.sh", pgpoolNode.Name)
	if err != nil {
		glog.Errorln("configureCluster: " + err.Error())
		return err
	}
	glog.Infoln("configureCluster: pgpool startpgpool output was" + commandoutput)

	//finally, update the cluster to show that it is
	//initialized!
	cluster.Status = "initialized"
	err = admindb.UpdateDBCluster(cluster)
	if err != nil {
		glog.Errorln("configureCluster:" + err.Error())
		return err
	}

	return nil

}

func GetAllClusters(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetAllClusters: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	results, err := admindb.GetAllDBClusters()
	if err != nil {
		glog.Errorln("GetAllClusters: error-" + err.Error())
		rest.Error(w, err.Error(), 400)
	}
	clusters := make([]Cluster, len(results))
	i := 0
	for i = range results {
		clusters[i].ID = results[i].ID
		clusters[i].Name = results[i].Name
		clusters[i].ClusterType = results[i].ClusterType
		clusters[i].Status = results[i].Status
		clusters[i].CreateDate = results[i].CreateDate
		i++
	}

	w.WriteJson(&clusters)
}

//we use POST for both updating and inserting based on the ID passed in
func PostCluster(w rest.ResponseWriter, r *rest.Request) {
	glog.Infoln("PostCluster: in PostCluster")
	cluster := Cluster{}
	err := r.DecodeJsonPayload(&cluster)
	if err != nil {
		glog.Errorln("PostCluster: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(cluster.Token, "perm-cluster")
	if err != nil {
		glog.Errorln("PostCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	if cluster.Name == "" {
		glog.Errorln("PostCluster: error in Name")
		rest.Error(w, "cluster name required", 400)
		return
	}

	glog.Infoln("PostCluster: have ID=" + cluster.ID + " Name=" + cluster.Name + " type=" + cluster.ClusterType + " status=" + cluster.Status)
	dbcluster := admindb.DBCluster{cluster.ID, cluster.Name, cluster.ClusterType, cluster.Status, ""}
	if cluster.ID == "" {
		strid, err := admindb.InsertDBCluster(dbcluster)
		newid := strconv.Itoa(strid)
		if err != nil {
			glog.Errorln("PostCluster:" + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}
		cluster.ID = newid
	} else {
		glog.Infoln("PostCluster: about to call UpdateDBCluster")
		err2 := admindb.UpdateDBCluster(dbcluster)
		if err2 != nil {
			glog.Errorln("PostCluster: error in UpdateDBCluster " + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}
	}

	w.WriteJson(&cluster)
}

func DeleteCluster(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-cluster")
	if err != nil {
		glog.Errorln("DeleteCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		glog.Errorln("DeleteCluster: error cluster ID required")
		rest.Error(w, "cluster ID required", 400)
		return
	}

	cluster, err := admindb.GetDBCluster(ID)
	if err != nil {
		glog.Errorln("DeleteCluster:" + err.Error())
		rest.Error(w, err.Error(), 400)
	}

	//delete docker containers
	containers, err := admindb.GetAllDBNodesForCluster(ID)
	if err != nil {
		glog.Errorln("DeleteCluster:" + err.Error())
		rest.Error(w, err.Error(), 400)
	}

	i := 0

	//handle the case where we want to delete a cluster but
	//it is not initialized, we can reuse the containers
	if cluster.Status == "uninitialized" {
		glog.Infoln("DeleteCluster: delete cluster but not the nodes")
		for i = range containers {
			containers[i].ClusterID = "-1"
			err = admindb.UpdateDBNode(containers[i])
			if err != nil {
				glog.Errorln("DeleteCluster:" + err.Error())
				rest.Error(w, err.Error(), 400)
				return
			}
		}

		err = admindb.DeleteDBCluster(ID)
		if err != nil {
			glog.Errorln("DeleteCluster:" + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}

		status := SimpleStatus{}
		status.Status = "OK"
		w.WriteHeader(http.StatusOK)
		w.WriteJson(&status)
		return
	}

	i = 0
	var output string
	server := admindb.DBServer{}
	for i = range containers {

		//go get the docker server IPAddress
		server, err = admindb.GetDBServer(containers[i].ServerID)
		if err != nil {
			glog.Errorln("DeleteCluster:" + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}

		glog.Infoln("DeleteCluster: got server IP " + server.IPAddress)

		//it is possible that someone can remove a container
		//outside of us, so we let it pass that we can't remove
		//it
		//err = removeContainer(server.IPAddress, containers[i].Name)
		output, err = cpmagent.DockerRemoveContainer(containers[i].Name,
			server.IPAddress)
		if err != nil {
			glog.Errorln("DeleteCluster: error when trying to remove container" + err.Error())
		}

		//send the server a deletevolume command
		output, err = cpmagent.AgentCommand(CPMBIN+"deletevolume", server.PGDataPath+"/"+containers[i].Name, server.IPAddress)
		glog.Infoln("DeleteCluster:" + output)

		i++
	}

	//delete the container entries
	//delete the cluster entry
	admindb.DeleteDBCluster(ID)

	for i = range containers {

		err = admindb.DeleteDBNode(containers[i].ID)
		if err != nil {
			glog.Errorln("DeleteCluster:" + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

func AdminFailover(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-cluster")
	if err != nil {
		glog.Errorln("AdminFailover: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	ID := r.PathParam("ID")
	if ID == "" {
		glog.Errorln("AdminFailover: node ID required error")
		rest.Error(w, "node ID required", 400)
		return
	}

	//dbNode is the standby node we are going to fail over and
	//make the new master in the cluster
	var dbNode admindb.DBClusterNode
	dbNode, err = admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	var output string

	cluster, err := admindb.GetDBCluster(dbNode.ClusterID)
	if err != nil {
		glog.Errorln("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	output, err = cpmagent.AgentCommand("/cluster/bin/fail-over.sh", dbNode.Name, dbNode.Name)
	if err != nil {
		glog.Errorln("AdminFailover: fail-over error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	glog.Infoln("AdminFailover: fail-over output " + output)

	//update the old master to standalone role
	oldMaster := admindb.DBClusterNode{}
	oldMaster, err = admindb.GetDBNodeMaster(dbNode.ClusterID)
	if err != nil {
		glog.Errorln("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	oldMaster.Role = "standalone"
	oldMaster.ClusterID = "-1"
	err = admindb.UpdateDBNode(oldMaster)
	if err != nil {
		glog.Errorln("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//update the failover node to master role
	dbNode.Role = "master"
	err = admindb.UpdateDBNode(dbNode)
	if err != nil {
		glog.Errorln("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//stop pg on the old master
	//params.IPAddress1 = oldMaster.IPAddress

	output, err = cpmagent.AgentCommand("/cluster/bin/stoppg.sh", oldMaster.Name, oldMaster.Name)
	if err != nil {
		glog.Errorln("AdminFailover: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	clusterNodes, err := admindb.GetAllDBNodesForCluster(dbNode.ClusterID)
	if err != nil {
		glog.Errorln("AdminFailover:" + err.Error())
		rest.Error(w, err.Error(), 400)
	}

	i := 0
	for i = range clusterNodes {

		if clusterNodes[i].Name == oldMaster.Name {
			glog.Infoln("AdminFailover: fail-over is skipping previous master")
		} else if clusterNodes[i].Name == dbNode.Name {
			glog.Infoln("fail-over is skipping new master " + clusterNodes[i].Name)
		} else {
			if clusterNodes[i].Image == "cpm-pgpool" {
				glog.Infoln("AdminFailover: fail-over is reconfiguring pgpool  " + clusterNodes[i].Name)
				//reconfigure pgpool node
			} else {
				//reconfigure other standby nodes
				glog.Infoln("AdminFailover: fail-over is reconfiguring standby  " + clusterNodes[i].Name)
				//stop standby
				var commandoutput string
				commandoutput, err = PGCommand("/cluster/bin/stoppg.sh", clusterNodes[i].Name)
				if err != nil {
					glog.Errorln("AdminFailover:" + err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}
				glog.Infoln("AdminFailover: fail-over stop output was" + commandoutput)

				var domainname admindb.DBSetting
				domainname, err = admindb.GetDBSetting("DOMAIN-NAME")
				if err != nil {
					glog.Errorln("configureCluster: DOMAIN-NAME err " + err.Error())
					rest.Error(w, err.Error(), 400)
				}
				//create base backup from master
				commandoutput, err = cpmagent.Command1("/cluster/bin/basebackup.sh", dbNode.Name+"."+domainname.Value, clusterNodes[i].Name)
				if err != nil {
					glog.Errorln("AdminFailover:" + err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}
				glog.Infoln("AdminFailover: fail-over basebackup output was" + commandoutput)

				var data string
				data, err = template.Recovery(dbNode.Name, "5432", "postgres")
				if err != nil {
					glog.Errorln("AdminFailover:" + err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}
				glog.Infoln("AdminFailover:fail-over\t standby recovery.conf generated")

				//write standby recovery.conf file remotely
				err = RemoteWritefile("/pgdata/recovery.conf", data, clusterNodes[i].Name)
				if err != nil {
					glog.Errorln("AdminFailover:" + err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}
				glog.Infoln("AdminFailover: fail-over standby recovery.conf copied remotely")

				if cluster.ClusterType == "synchronous" {
					data, err = template.Postgresql("standby", "5432", "*")
				} else {
					data, err = template.Postgresql("standby", "5432", "")
				}
				if err != nil {
					glog.Errorln("AdminFailover: " + err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}

				//write standby postgresql.conf file remotely
				err = RemoteWritefile("/pgdata/postgresql.conf", data, clusterNodes[i].Name)
				if err != nil {
					glog.Errorln("AdminFailover: " + err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}
				glog.Infoln("AdminFailover: standby postgresql.conf copied remotely")

				//configure standby pg_hba.conf file
				data, err = template.Hba("standby", clusterNodes[i].Name, "5432", dbNode.ClusterID, domainname.Value)
				if err != nil {
					glog.Errorln("AdminFailover:" + err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}

				glog.Infoln("AdminFailover: fail-over\t standby pg_hba.conf generated")

				//write standby pg_hba.conf file remotely
				err = RemoteWritefile("/pgdata/pg_hba.conf", data, clusterNodes[i].Name)
				if err != nil {
					glog.Errorln("AdminFailover: " + err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}
				glog.Infoln("AdminFailover:  standby pg_hba.conf copied remotely")

				//start standby

				commandoutput, err = PGCommand("/cluster/bin/startpgonstandby.sh", clusterNodes[i].Name)
				if err != nil {
					glog.Errorln("AdminFailover:" + err.Error())
					rest.Error(w, err.Error(), 400)
					return
				}
				glog.Infoln("AdminFailover: standby startpg output was" + commandoutput)
			}
		}

		i++
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func EventJoinCluster(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-cluster")
	if err != nil {
		glog.Errorln("EventJoinCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	IDList := r.PathParam("IDList")
	if IDList == "" {
		glog.Errorln("EventJoinCluster: error IDList required")
		rest.Error(w, "IDList required", 400)
		return
	} else {
		glog.Infoln("EventJoinCluster: IDList=[" + IDList + "]")
	}

	MasterID := r.PathParam("MasterID")
	if MasterID == "" {
		glog.Errorln("EventJoinCluster: error MasterID required")
		rest.Error(w, "MasterID required", 400)
		return
	} else {
		glog.Infoln("EventJoinCluster: MasterID=[" + MasterID + "]")
	}
	ClusterID := r.PathParam("ClusterID")
	if ClusterID == "" {
		glog.Errorln("EventJoinCluster: error ClusterID required")
		rest.Error(w, "node ClusterID required", 400)
		return
	} else {
		glog.Infoln("EventJoinCluster: ClusterID=[" + ClusterID + "]")
	}

	var idList = strings.Split(IDList, "_")
	i := 0
	origDBNode := admindb.DBClusterNode{}
	for i = range idList {
		if idList[i] != "" {
			glog.Infoln("EventJoinCluster: idList[" + strconv.Itoa(i) + "]=" + idList[i])
			origDBNode, err = admindb.GetDBNode(idList[i])
			if err != nil {
				glog.Errorln("EventJoinCluster:" + err.Error())
				rest.Error(w, err.Error(), 400)
				return
			}

			//update the node to be in the cluster
			origDBNode.ClusterID = ClusterID
			if origDBNode.Image == "cpm-node" {
				origDBNode.Role = "standby"
			} else {
				origDBNode.Role = "pgpool"
			}
			err = admindb.UpdateDBNode(origDBNode)
			if err != nil {
				glog.Errorln("EventJoinCluster:" + err.Error())
				rest.Error(w, err.Error(), 400)
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
		origDBNode, err = admindb.GetDBNode(MasterID)
		if err != nil {
			glog.Errorln("EventJoinCluster:" + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}

		origDBNode.ClusterID = ClusterID
		origDBNode.Role = "master"
		err = admindb.UpdateDBNode(origDBNode)
		if err != nil {
			glog.Errorln("EventJoinCluster:" + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AutoCluster(w rest.ResponseWriter, r *rest.Request) {
	glog.Infoln("AUTO CLUSTER PROFILE starts")
	glog.Infoln("AutoCluster: start AutoCluster")
	params := AutoClusterInfo{}
	err := r.DecodeJsonPayload(&params)
	if err != nil {
		glog.Errorln("AutoCluster: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(params.Token, "perm-cluster")
	if err != nil {
		glog.Errorln("AutoCluster: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	if params.Name == "" {
		glog.Errorln("AutoCluster: error in Name")
		rest.Error(w, "cluster name required", 400)
		return
	}
	if params.ClusterType == "" {
		glog.Errorln("AutoCluster: error in ClusterType")
		rest.Error(w, "ClusterType name required", 400)
		return
	}
	if params.ClusterProfile == "" {
		glog.Errorln("AutoCluster: error in ClusterProfile")
		rest.Error(w, "ClusterProfile name required", 400)
		return
	}

	glog.Infoln("AutoCluster: Name=" + params.Name + " ClusterType=" + params.ClusterType + " Profile=" + params.ClusterProfile)

	//create cluster definition
	dbcluster := admindb.DBCluster{"", params.Name, params.ClusterType, "uninitialized", ""}
	var ival int
	ival, err = admindb.InsertDBCluster(dbcluster)
	clusterID := strconv.Itoa(ival)
	dbcluster.ID = clusterID
	glog.Infoln(clusterID)
	if err != nil {
		glog.Errorln("AutoCluster:" + err.Error())
		rest.Error(w, "Insert Cluster error:"+err.Error(), 400)
		return
	}

	//lookup profile
	profile, err2 := getClusterProfileInfo(params.ClusterProfile)
	if err2 != nil {
		glog.Errorln("AutoCluster: error-" + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), 400)
		return
	}

	var masterServer admindb.DBServer
	var chosenServers []admindb.DBServer
	if profile.Algo == "round-robin" {
		masterServer, chosenServers, err2 = roundRobin(profile)
	} else {
		glog.Errorln("AutoCluster: error-unsupported algorithm request")
		rest.Error(w, "AutoCluster error: unsupported algorithm", 400)
		return
	}

	//create master container
	docker := new(cpmagent.DockerRunArgs)
	docker.Image = "cpm-node"
	docker.ContainerName = params.Name + "-master"
	docker.ServerID = masterServer.ID
	docker.Standalone = "false"
	if err != nil {
		glog.Errorln("AutoCluster: error-create master node " + err.Error())
		rest.Error(w, "AutoCluster error"+err.Error(), 400)
		return
	}

	//	provision the master
	err2 = provisionImpl(docker, profile.MasterProfile, false)
	if err2 != nil {
		glog.Errorln("AutoCluster: error-provision master " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), 400)
		return
	}
	glog.Infoln("AUTO CLUSTER PROFILE master container created")
	var node admindb.DBClusterNode
	//update node with cluster iD
	node, err2 = admindb.GetDBNodeByName(docker.ContainerName)
	if err2 != nil {
		glog.Errorln("AutoCluster: error-get node by name " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), 400)
		return
	}

	node.ClusterID = clusterID
	node.Role = "master"
	err2 = admindb.UpdateDBNode(node)
	if err2 != nil {
		glog.Errorln("AutoCluster: error-update standby node " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), 400)
		return
	}

	//create standby containers
	var count int
	count, err2 = strconv.Atoi(profile.Count)
	for i := 0; i < count; i++ {
		glog.Infoln("working on standby ....")
		//	loop - provision standby
		docker.ServerID = chosenServers[i].ID
		docker.ContainerName = params.Name + "-standby-" + strconv.Itoa(i)
		err2 = provisionImpl(docker, profile.StandbyProfile, true)
		if err2 != nil {
			glog.Errorln("AutoCluster: error-provision master " + err2.Error())
			rest.Error(w, "AutoCluster error"+err2.Error(), 400)
			return
		}

		//update node with cluster iD
		node, err2 = admindb.GetDBNodeByName(docker.ContainerName)
		if err2 != nil {
			glog.Errorln("AutoCluster: error-get node by name " + err2.Error())
			rest.Error(w, "AutoCluster error"+err2.Error(), 400)
			return
		}

		node.ClusterID = clusterID
		node.Role = "standby"
		err2 = admindb.UpdateDBNode(node)
		if err2 != nil {
			glog.Errorln("AutoCluster: error-update standby node " + err2.Error())
			rest.Error(w, "AutoCluster error"+err2.Error(), 400)
			return
		}
	}
	glog.Infoln("AUTO CLUSTER PROFILE standbys created")
	//create pgpool container
	//	provision
	docker.ContainerName = params.Name + "-pgpool"
	docker.Image = "cpm-pgpool"
	docker.ServerID = chosenServers[count].ID
	err2 = provisionImpl(docker, profile.StandbyProfile, true)
	if err2 != nil {
		glog.Errorln("AutoCluster: error-provision pgpool " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), 400)
		return
	}
	glog.Infoln("AUTO CLUSTER PROFILE pgpool created")
	//update node with cluster ID
	node, err2 = admindb.GetDBNodeByName(docker.ContainerName)
	if err2 != nil {
		glog.Errorln("AutoCluster: error-get pgpool node by name " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), 400)
		return
	}

	node.ClusterID = clusterID
	node.Role = "pgpool"
	err2 = admindb.UpdateDBNode(node)
	if err2 != nil {
		glog.Errorln("AutoCluster: error-update pgpool node " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), 400)
		return
	}

	//configure cluster
	//	ConfigureCluster
	glog.Infoln("AUTO CLUSTER PROFILE configure cluster ")
	err2 = configureCluster(dbcluster, true)
	if err2 != nil {
		glog.Errorln("AutoCluster: error-configure cluster " + err2.Error())
		rest.Error(w, "AutoCluster error"+err2.Error(), 400)
		return
	}

	glog.Infoln("AUTO CLUSTER PROFILE done")
	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

// round-robin provisioning algorithm -
//  to promote least used servers, incoming servers list
//  should be sorted by class and least used order
//  returns the master server and the list of standby servers
func roundRobin(profile ClusterProfiles) (admindb.DBServer, []admindb.DBServer, error) {
	var masterServer admindb.DBServer
	count, err := strconv.Atoi(profile.Count)

	//add 1 to the standby count to make room for the pgpool node
	count++

	//create a slice to hold servers for standby and pgpool nodes
	//assumes 1 pgpool node per cluster which is enforced by auto-cluster
	chosen := make([]admindb.DBServer, count)

	//get all the servers available
	servers, err := admindb.GetAllDBServersByClassByCount()
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

	glog.Infoln("round-robin: master " + masterServer.Name + " class=" + masterServer.ServerClass)
	for x := 0; x < len(chosen); x++ {
		glog.Infoln("round-robin: other servers " + chosen[x].Name + " class=" + chosen[x].ServerClass)
	}
	return masterServer, chosen, nil
}
