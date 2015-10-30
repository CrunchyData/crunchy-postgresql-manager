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
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/adminapi"
)

func init() {
	fmt.Println("before parsing in init")
	flag.Parse()

}

var CPMDIR = "/var/cpm/"
var CPMBIN = CPMDIR + "bin/"

func main() {

	fmt.Println("at top of adminapi main")

	var err error

	handler := rest.ResourceHandler{
		PreRoutingMiddlewares: []rest.Middleware{
			&rest.CorsMiddleware{
				RejectNonCorsRequests: false,
				OriginValidator: func(origin string, request *rest.Request) bool {
					return true
				},
				AllowedMethods: []string{"DELETE", "GET", "POST", "PUT"},
				AllowedHeaders: []string{
					"Accept", "Content-Type", "X-Custom-Header", "Origin"},
				AccessControlAllowCredentials: true,
				AccessControlMaxAge:           3600,
			},
		},
		EnableRelaxedContentType: true,
	}

	err = handler.SetRoutes(
		// Cluster Information
		// Get all the clusters for a given project
		&rest.Route{"GET", "/projectclusters/:ID.:Token", adminapi.GetAllClustersForProject},
		//		&rest.Route{"GET", "/clusters/:Token", adminapi.GetAllClusters},

		// Server Information
		// Get all the servers defined in CPM
		&rest.Route{"GET", "/servers/:Token", adminapi.GetAllServers},
		// update a cluster
		&rest.Route{"POST", "/cluster", adminapi.PostCluster},
		// create a new cluster
		&rest.Route{"POST", "/autocluster", adminapi.AutoCluster},
		// update a setting value
		&rest.Route{"POST", "/savesetting", adminapi.SaveSetting},
		//		&rest.Route{"POST", "/savesettings", adminapi.SaveSettings},
		//		&rest.Route{"POST", "/saveprofiles", adminapi.SaveProfiles},
		//		&rest.Route{"POST", "/saveclusterprofiles", adminapi.SaveClusterProfiles},
		// get all the settings
		&rest.Route{"GET", "/settings/:Token", adminapi.GetAllSettings},
		// get all the general settings
		&rest.Route{"GET", "/generalsettings/:Token", adminapi.GetAllGeneralSettings},
		// define a new CPM server
		&rest.Route{"GET", "/addserver/:ID.:Name.:IPAddress.:DockerBridgeIP.:PGDataPath.:ServerClass.:Token", adminapi.AddServer},
		// get a CPM server
		&rest.Route{"GET", "/server/:ID.:Token", adminapi.GetServer},
		// get a cluster
		&rest.Route{"GET", "/cluster/:ID.:Token", adminapi.GetCluster},
		// configure a cluster
		&rest.Route{"GET", "/cluster/configure/:ID.:Token", adminapi.ConfigureCluster},
		// scale up a cluster by adding a new standby container
		&rest.Route{"GET", "/cluster/scale/:ID.:Token", adminapi.ScaleUpCluster},
		// delete a cluster and all of its nodes
		&rest.Route{"GET", "/cluster/delete/:ID.:Token", adminapi.DeleteCluster},
		// perform a docker start on a given clusters set of containers
		&rest.Route{"GET", "/cluster/start/:ID.:Token", adminapi.StartCluster},
		// perform a docker stop on a given clusters set of containers
		&rest.Route{"GET", "/cluster/stop/:ID.:Token", adminapi.StopCluster},
		// remove a CPM server definition
		&rest.Route{"GET", "/deleteserver/:ID.:Token", adminapi.DeleteServer},

		// get all the containers for a given project
		&rest.Route{"GET", "/projectnodes/:ID.:Token", adminapi.GetAllNodesForProject},
		// get all the containers in the whole CPM system
		&rest.Route{"GET", "/nodes/:Token", adminapi.GetAllNodes},
		// get all the containers that are not part of a cluster
		&rest.Route{"GET", "/nodes/nocluster/:Token", adminapi.GetAllNodesNotInCluster},
		// get all the containers within a given cluster
		&rest.Route{"GET", "/clusternodes/:ClusterID.:Token", adminapi.GetAllNodesForCluster},
		// get all the containers running on a given CPM server
		&rest.Route{"GET", "/nodes/forserver/:ServerID.:Token", adminapi.GetAllNodesForServer},
		// perform a docker run to create a new container
		&rest.Route{"POST", "/provision", adminapi.Provision},
		// get a container
		&rest.Route{"GET", "/node/:ID.:Token", adminapi.GetNode},
		// delete a container
		&rest.Route{"GET", "/deletenode/:ID.:Token", adminapi.DeleteNode},

		// get a CPM server monitoring metric
		&rest.Route{"GET", "/monitor/server-getinfo/:ServerID.:Metric.:Token", adminapi.MonitorServerGetInfo},
		// get a container database pg_settings output
		&rest.Route{"GET", "/monitor/container/settings/:ID.:Token", adminapi.MonitorContainerSettings},
		// get a container database pg_stat_statements output
		&rest.Route{"GET", "/monitor/container/statements/:ID.:Token", adminapi.MonitorStatements},
		// get a cluster master pg_stat_Replication output
		&rest.Route{"GET", "/monitor/container/repl/:ID.:Token", adminapi.ContainerInfoStatrepl},
		// get a container database pg_stat_database output
		&rest.Route{"GET", "/monitor/container/database/:ID.:Token", adminapi.ContainerInfoStatdatabase},
		// get a container database bgwriter output
		&rest.Route{"GET", "/monitor/container/bgwriter/:ID.:Token", adminapi.ContainerInfoBgwriter},
		// run a pgbadger on a container to generate a pgbadger report
		&rest.Route{"GET", "/monitor/container/badger/:ID.:Token", adminapi.BadgerGenerate},
		// get pg_control output on a container database
		&rest.Route{"GET", "/monitor/container/controldata/:ID.:Token", adminapi.MonitorContainerControldata},
		// get a load test output on a container database
		&rest.Route{"GET", "/monitor/container/loadtest/:ID.:Writes.:Token", adminapi.ContainerLoadTest},

		// perform a docker start on all containers on a given server
		&rest.Route{"GET", "/admin/startall/:ID.:Token", adminapi.AdminStartServerContainers},
		// perform a docker stop on all containers on a given server
		&rest.Route{"GET", "/admin/stopall/:ID.:Token", adminapi.AdminStopServerContainers},
		// start postgres on a given container
		&rest.Route{"GET", "/admin/start-pg/:ID.:Token", adminapi.AdminStartpg},
		// perform a docker start on a given container
		&rest.Route{"GET", "/admin/start/:ID.:Token", adminapi.AdminStartNode},
		// perform a docker stop on a given container
		&rest.Route{"GET", "/admin/stop/:ID.:Token", adminapi.AdminStopNode},
		// perform a failover on a cluster standby container
		&rest.Route{"GET", "/admin/failover/:ID.:Token", adminapi.AdminFailover},
		// stop postgres on a given container
		&rest.Route{"GET", "/admin/stop-pg/:ID.:Token", adminapi.AdminStoppg},
		// join a standby container to a cluster
		&rest.Route{"GET", "/event/join-cluster/:IDList.:MasterID.:ClusterID.:Token", adminapi.EventJoinCluster},
		// perform a CPM user login
		&rest.Route{"GET", "/sec/login/:ID.:PSW", adminapi.Login},
		// perform a CPM user logout
		&rest.Route{"GET", "/sec/logout/:Token", adminapi.Logout},
		// update a CPM user
		&rest.Route{"POST", "/sec/updateuser", adminapi.UpdateUser},
		// change a CPM user password
		&rest.Route{"POST", "/sec/cp", adminapi.ChangePassword},
		// create a new CPM user
		&rest.Route{"POST", "/sec/adduser", adminapi.AddUser},
		// get a CPM user
		&rest.Route{"GET", "/sec/getuser/:ID.:Token", adminapi.GetUser},
		// get all CPM users
		&rest.Route{"GET", "/sec/getusers/:Token", adminapi.GetAllUsers},
		// delete a CPM user
		&rest.Route{"GET", "/sec/deleteuser/:ID.:Token", adminapi.DeleteUser},
		// update a CPM role
		&rest.Route{"POST", "/sec/updaterole", adminapi.UpdateRole},
		// add a CPM role
		&rest.Route{"POST", "/sec/addrole", adminapi.AddRole},
		// delete a CPM role
		&rest.Route{"GET", "/sec/deleterole/:ID.:Token", adminapi.DeleteRole},
		// get all CPM roles
		&rest.Route{"GET", "/sec/getroles/:Token", adminapi.GetAllRoles},
		// get a CPM role
		&rest.Route{"GET", "/sec/getrole/:Name.:Token", adminapi.GetRole},
		// execute a task schedule immediately
		&rest.Route{"POST", "/task/executenow", adminapi.ExecuteNow},
		// add a new task schedule for a given container
		&rest.Route{"POST", "/task/addschedule", adminapi.AddSchedule},
		// delete a task schedule for a given container
		&rest.Route{"GET", "/task/deleteschedule/:ID.:Token", adminapi.DeleteSchedule},
		// update a task schedule for a given container
		&rest.Route{"POST", "/task/updateschedule", adminapi.UpdateSchedule},
		// get all task schedules for a given container
		&rest.Route{"GET", "/task/getschedules/:ContainerName.:Token", adminapi.GetAllSchedules},
		// get a task schedule
		&rest.Route{"GET", "/task/getschedule/:ID.:Token", adminapi.GetSchedule},
		// get the status of a given task schedule
		&rest.Route{"GET", "/task/getstatus/:ID.:Token", adminapi.GetStatus},
		// get all the status of a given task schedule
		&rest.Route{"GET", "/task/getallstatus/:ID.:Token", adminapi.GetAllStatus},
		// TODO
		&rest.Route{"GET", "/task/nodes/:Token", adminapi.GetBackupNodes},
		//&rest.Route{"GET", "/mon/server/:Metric.:ServerID.:Interval.:Token", adminapi.GetServerMetrics},
		//&rest.Route{"GET", "/mon/container/pg2/:Name.:Interval.:Token", adminapi.GetPG2},
		// get the current CPM system health check output
		&rest.Route{"GET", "/mon/healthcheck/:Token", adminapi.GetHealthCheck},
		// get the current CPM version
		&rest.Route{"GET", "/version", GetVersion},
		// add a database user to a given container
		&rest.Route{"POST", "/dbuser/add", adminapi.AddContainerUser},
		// update a database user for a given container
		&rest.Route{"POST", "/dbuser/update", adminapi.UpdateContainerUser},
		// delete a database user for a given container
		&rest.Route{"GET", "/dbuser/delete/:ContainerID.:Rolname.:Token", adminapi.DeleteContainerUser},
		// get a container user for a given container
		&rest.Route{"GET", "/dbuser/get/:ContainerID.:Rolname.:Token", adminapi.GetContainerUser},
		// get all database users for a given container
		&rest.Route{"GET", "/dbuser/getall/:ID.:Token", adminapi.GetAllUsersForContainer},
		// add a CPM project
		&rest.Route{"POST", "/project/add", adminapi.AddProject},
		// update a CPM project
		&rest.Route{"POST", "/project/update", adminapi.UpdateProject},
		// get a CPM project
		&rest.Route{"GET", "/project/get/:ID.:Token", adminapi.GetProject},
		// get all CPM projects
		&rest.Route{"GET", "/project/getall/:Token", adminapi.GetAllProjects},
		// delete a CPM project
		&rest.Route{"GET", "/project/delete/:ID.:Token", adminapi.DeleteProject},
		// get an access rule
		&rest.Route{"GET", "/rules/get/:ID.:Token", adminapi.RulesGet},
		// get all access rules
		&rest.Route{"GET", "/rules/getall/:Token", adminapi.RulesGetAll},
		// delete an access rule
		&rest.Route{"GET", "/rules/delete/:ID.:Token", adminapi.RulesDelete},
		// update an access rule
		&rest.Route{"POST", "/rules/update", adminapi.RulesUpdate},
		// create a new access rule
		&rest.Route{"POST", "/rules/insert", adminapi.RulesInsert},
		// get all accessrules for a container
		&rest.Route{"GET", "/containerrules/getall/:ID.:Token", adminapi.ContainerAccessRuleGetAll},
		// update accessrules for a container
		&rest.Route{"POST", "/containerrules/update", adminapi.ContainerAccessRuleUpdate},
		// create a proxy container
		&rest.Route{"POST", "/provisionproxy", adminapi.ProvisionProxy},
		// get a proxy
		&rest.Route{"GET", "/proxy/getbycontainerid/:ContainerID.:Token", adminapi.GetProxyByContainerID},
		// update a proxy
		&rest.Route{"POST", "/proxy/update", adminapi.ProxyUpdate},
	)
	if err != nil {
		log.Fatal(err)
	}
	//	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":13001", &handler))
	log.Fatal(http.ListenAndServeTLS(":13000", "/cpmkeys/cert.pem", "/cpmkeys/key.pem", &handler))
}

func GetVersion(w rest.ResponseWriter, r *rest.Request) {

	w.(http.ResponseWriter).Write([]byte("0.9.6"))
}
