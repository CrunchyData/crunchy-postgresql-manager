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
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/adminapi"
	"log"
	"net/http"
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
		&rest.Route{"GET", "/projectclusters/:ID.:Token", adminapi.GetAllClustersForProject},
		&rest.Route{"GET", "/clusters/:Token", adminapi.GetAllClusters},
		&rest.Route{"GET", "/servers/:Token", adminapi.GetAllServers},
		&rest.Route{"POST", "/cluster", adminapi.PostCluster},
		&rest.Route{"POST", "/autocluster", adminapi.AutoCluster},
		&rest.Route{"POST", "/savesetting", adminapi.SaveSetting},
		&rest.Route{"POST", "/savesettings", adminapi.SaveSettings},
		&rest.Route{"POST", "/saveprofiles", adminapi.SaveProfiles},
		&rest.Route{"POST", "/saveclusterprofiles", adminapi.SaveClusterProfiles},
		&rest.Route{"GET", "/settings/:Token", adminapi.GetAllSettings},
		&rest.Route{"GET", "/generalsettings/:Token", adminapi.GetAllGeneralSettings},
		&rest.Route{"GET", "/addserver/:ID.:Name.:IPAddress.:DockerBridgeIP.:PGDataPath.:ServerClass.:Token", adminapi.AddServer},
		&rest.Route{"GET", "/server/:ID.:Token", adminapi.GetServer},
		&rest.Route{"GET", "/cluster/:ID.:Token", adminapi.GetCluster},
		&rest.Route{"GET", "/cluster/configure/:ID.:Token", adminapi.ConfigureCluster},
		&rest.Route{"GET", "/cluster/scale/:ID.:Token", adminapi.ScaleUpCluster},
		&rest.Route{"GET", "/cluster/delete/:ID.:Token", adminapi.DeleteCluster},
		&rest.Route{"GET", "/deleteserver/:ID.:Token", adminapi.DeleteServer},
		&rest.Route{"GET", "/projectnodes/:ID.:Token", adminapi.GetAllNodesForProject},
		&rest.Route{"GET", "/nodes/:Token", adminapi.GetAllNodes},
		&rest.Route{"GET", "/nodes/nocluster/:Token", adminapi.GetAllNodesNotInCluster},
		&rest.Route{"GET", "/clusternodes/:ClusterID.:Token", adminapi.GetAllNodesForCluster},
		&rest.Route{"GET", "/nodes/forserver/:ServerID.:Token", adminapi.GetAllNodesForServer},
		&rest.Route{"GET", "/provision/:Profile.:Image.:ServerID.:ProjectID.:ContainerName.:Standalone.:Token", adminapi.Provision},
		&rest.Route{"GET", "/node/:ID.:Token", adminapi.GetNode},
		&rest.Route{"GET", "/deletenode/:ID.:Token", adminapi.DeleteNode},
		&rest.Route{"GET", "/monitor/server-getinfo/:ServerID.:Metric.:Token", adminapi.MonitorServerGetInfo},
		&rest.Route{"GET", "/monitor/container/settings/:ID.:Token", adminapi.MonitorContainerSettings},
		&rest.Route{"GET", "/monitor/container/statements/:ID.:Token", adminapi.MonitorStatements},
		&rest.Route{"GET", "/monitor/container/repl/:ID.:Token", adminapi.ContainerInfoStatrepl},
		&rest.Route{"GET", "/monitor/container/database/:ID.:Token", adminapi.ContainerInfoStatdatabase},
		&rest.Route{"GET", "/monitor/container/bgwriter/:ID.:Token", adminapi.ContainerInfoBgwriter},
		&rest.Route{"GET", "/monitor/container/badger/:ID.:Token", adminapi.BadgerGenerate},
		&rest.Route{"GET", "/monitor/container/controldata/:ID.:Token", adminapi.MonitorContainerControldata},
		&rest.Route{"GET", "/monitor/container/loadtest/:ID.:Writes.:Token", adminapi.ContainerLoadTest},
		&rest.Route{"GET", "/admin/startall/:ID.:Token", adminapi.AdminStartServerContainers},
		&rest.Route{"GET", "/admin/stopall/:ID.:Token", adminapi.AdminStopServerContainers},
		&rest.Route{"GET", "/admin/start-pg/:ID.:Token", adminapi.AdminStartpg},
		&rest.Route{"GET", "/admin/start/:ID.:Token", adminapi.AdminStartNode},
		&rest.Route{"GET", "/admin/stop/:ID.:Token", adminapi.AdminStopNode},
		&rest.Route{"GET", "/admin/failover/:ID.:Token", adminapi.AdminFailover},
		&rest.Route{"GET", "/admin/stop-pg/:ID.:Token", adminapi.AdminStoppg},

		&rest.Route{"GET", "/event/join-cluster/:IDList.:MasterID.:ClusterID.:Token", adminapi.EventJoinCluster},
		&rest.Route{"GET", "/sec/login/:ID.:PSW", adminapi.Login},
		&rest.Route{"GET", "/sec/logout/:Token", adminapi.Logout},
		&rest.Route{"POST", "/sec/updateuser", adminapi.UpdateUser},
		&rest.Route{"POST", "/sec/cp", adminapi.ChangePassword},
		&rest.Route{"POST", "/sec/adduser", adminapi.AddUser},
		&rest.Route{"GET", "/sec/getuser/:ID.:Token", adminapi.GetUser},
		&rest.Route{"GET", "/sec/getusers/:Token", adminapi.GetAllUsers},
		&rest.Route{"GET", "/sec/deleteuser/:ID.:Token", adminapi.DeleteUser},
		&rest.Route{"POST", "/sec/updaterole", adminapi.UpdateRole},
		&rest.Route{"POST", "/sec/addrole", adminapi.AddRole},
		&rest.Route{"GET", "/sec/deleterole/:ID.:Token", adminapi.DeleteRole},
		&rest.Route{"GET", "/sec/getroles/:Token", adminapi.GetAllRoles},
		&rest.Route{"GET", "/sec/getrole/:Name.:Token", adminapi.GetRole},
		&rest.Route{"POST", "/backup/now", adminapi.BackupNow},
		&rest.Route{"POST", "/backup/addschedule", adminapi.AddSchedule},
		&rest.Route{"GET", "/backup/deleteschedule/:ID.:Token", adminapi.DeleteSchedule},
		&rest.Route{"POST", "/backup/updateschedule", adminapi.UpdateSchedule},
		&rest.Route{"GET", "/backup/getschedules/:ContainerName.:Token", adminapi.GetAllSchedules},
		&rest.Route{"GET", "/backup/getschedule/:ID.:Token", adminapi.GetSchedule},
		&rest.Route{"GET", "/backup/getstatus/:ID.:Token", adminapi.GetStatus},
		&rest.Route{"GET", "/backup/getallstatus/:ID.:Token", adminapi.GetAllStatus},
		&rest.Route{"GET", "/backup/nodes/:Token", adminapi.GetBackupNodes},
		//&rest.Route{"GET", "/mon/server/:Metric.:ServerID.:Interval.:Token", adminapi.GetServerMetrics},
		//&rest.Route{"GET", "/mon/container/pg2/:Name.:Interval.:Token", adminapi.GetPG2},
		&rest.Route{"GET", "/mon/healthcheck/:Token", adminapi.GetHealthCheck},
		&rest.Route{"GET", "/version", adminapi.GetVersion},
		&rest.Route{"POST", "/dbuser/add", adminapi.AddContainerUser},
		&rest.Route{"POST", "/dbuser/update", adminapi.UpdateContainerUser},
		&rest.Route{"GET", "/dbuser/delete/:ContainerID.:Rolname.:Token", adminapi.DeleteContainerUser},
		&rest.Route{"GET", "/dbuser/get/:ContainerID.:Rolname.:Token", adminapi.GetContainerUser},
		&rest.Route{"GET", "/dbuser/getall/:ID.:Token", adminapi.GetAllUsersForContainer},
		&rest.Route{"POST", "/project/add", adminapi.AddProject},
		&rest.Route{"POST", "/project/update", adminapi.UpdateProject},
		&rest.Route{"GET", "/project/get/:ID.:Token", adminapi.GetProject},
		&rest.Route{"GET", "/project/getall/:Token", adminapi.GetAllProjects},
		&rest.Route{"GET", "/project/delete/:ID.:Token", adminapi.DeleteProject},
		&rest.Route{"GET", "/rules/get/:ID.:Token", adminapi.RulesGet},
		&rest.Route{"GET", "/rules/getall/:Token", adminapi.RulesGetAll},
		&rest.Route{"GET", "/rules/delete/:ID.:Token", adminapi.RulesDelete},
		&rest.Route{"POST", "/rules/update", adminapi.RulesUpdate},
		&rest.Route{"POST", "/rules/insert", adminapi.RulesInsert},
		&rest.Route{"GET", "/containerrules/getall/:ID.:Token", adminapi.ContainerAccessRuleGetAll},
		&rest.Route{"POST", "/containerrules/update", adminapi.ContainerAccessRuleUpdate},
	)
	if err != nil {
		log.Fatal(err)
	}
	//	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":13001", &handler))
	log.Fatal(http.ListenAndServeTLS(":13000", "/cpmkeys/cert.pem", "/cpmkeys/key.pem", &handler))
}
