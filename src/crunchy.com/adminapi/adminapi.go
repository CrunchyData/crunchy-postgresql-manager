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
	"crunchy.com/backup"
	"crunchy.com/util"
	"database/sql"
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	fmt.Println("before parsing in init")
	flag.Parse()
}

var PGBIN = "/usr/pgsql-9.4/bin/"
var CPMDIR = "/opt/cpm/"
var CPMBIN = CPMDIR + "bin/"

func main() {

	fmt.Println("at top of adminapi main")
	//flag.Parse()
	//glog.Flush()

	//glog.Info("called flag.Parse\n")
	//glog.Flush()

	var dbConn *sql.DB
	found := false
	var err error

	for i := 0; i < 10; i++ {
		dbConn, err = util.GetConnection("clusteradmin")
		if err != nil {
			glog.Errorln(err.Error())
			glog.Errorln("could not get initial database connection, will retry in 5 seconds")
			time.Sleep(time.Millisecond * 5000)
		} else {
			//glog.Infoln("got db connection")
			found = true
			break
		}
	}

	admindb.SetConnection(dbConn)
	backup.SetConnection(dbConn)

	if !found {
		panic("could not connect to clusteradmin db")
	}

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
		&rest.Route{"GET", "/clusters/:Token", GetAllClusters},
		&rest.Route{"GET", "/servers/:Token", GetAllServers},
		&rest.Route{"POST", "/cluster", PostCluster},
		&rest.Route{"POST", "/autocluster", AutoCluster},
		&rest.Route{"POST", "/savesettings", SaveSettings},
		&rest.Route{"POST", "/saveprofiles", SaveProfiles},
		&rest.Route{"POST", "/saveclusterprofiles", SaveClusterProfiles},
		&rest.Route{"GET", "/settings/:Token", GetAllSettings},
		&rest.Route{"GET", "/addserver/:ID.:Name.:IPAddress.:DockerBridgeIP.:PGDataPath.:ServerClass.:Token", AddServer},
		&rest.Route{"GET", "/server/:ID.:Token", GetServer},
		&rest.Route{"GET", "/cluster/:ID.:Token", GetCluster},
		&rest.Route{"GET", "/cluster/configure/:ID.:Token", ConfigureCluster},
		&rest.Route{"GET", "/cluster/delete/:ID.:Token", DeleteCluster},
		&rest.Route{"GET", "/deleteserver/:ID.:Token", DeleteServer},
		&rest.Route{"GET", "/nodes/:Token", GetAllNodes},
		&rest.Route{"GET", "/nodes/nocluster/:Token", GetAllNodesNotInCluster},
		&rest.Route{"GET", "/clusternodes/:ClusterID.:Token", GetAllNodesForCluster},
		&rest.Route{"GET", "/nodes/forserver/:ServerID.:Token", GetAllNodesForServer},
		&rest.Route{"GET", "/provision/:Profile.:Image.:ServerID.:ContainerName.:Standalone.:Token", Provision},
		&rest.Route{"GET", "/node/:ID.:Token", GetNode},
		&rest.Route{"GET", "/kube/:Token", Kube},
		&rest.Route{"GET", "/deletenode/:ID.:Token", DeleteNode},
		&rest.Route{"GET", "/monitor/server-getinfo/:ServerID.:Metric.:Token", MonitorServerGetInfo},
		&rest.Route{"GET", "/monitor/container/settings/:ID.:Token", MonitorContainerSettings},
		&rest.Route{"GET", "/monitor/container/repl/:ID.:Token", ContainerInfoStatrepl},
		&rest.Route{"GET", "/monitor/container/database/:ID.:Token", ContainerInfoStatdatabase},
		&rest.Route{"GET", "/monitor/container/bgwriter/:ID.:Token", ContainerInfoBgwriter},
		&rest.Route{"GET", "/monitor/container/controldata/:ID.:Token", MonitorContainerControldata},
		&rest.Route{"GET", "/monitor/container/loadtest/:ID.:Writes.:Token", ContainerLoadTest},
		&rest.Route{"GET", "/admin/start-pg/:ID.:Token", AdminStartpg},
		&rest.Route{"GET", "/admin/start/:ID.:Token", AdminStartNode},
		&rest.Route{"GET", "/admin/stop/:ID.:Token", AdminStopNode},
		&rest.Route{"GET", "/admin/failover/:ID.:Token", AdminFailover},
		&rest.Route{"GET", "/admin/stop-pg/:ID.:Token", AdminStoppg},

		&rest.Route{"GET", "/event/join-cluster/:IDList.:MasterID.:ClusterID.:Token", EventJoinCluster},
		&rest.Route{"GET", "/sec/login/:ID.:PSW", Login},
		&rest.Route{"GET", "/sec/logout/:Token", Logout},
		&rest.Route{"POST", "/sec/updateuser", UpdateUser},
		&rest.Route{"POST", "/sec/cp", ChangePassword},
		&rest.Route{"POST", "/sec/adduser", AddUser},
		&rest.Route{"GET", "/sec/getuser/:ID.:Token", GetUser},
		&rest.Route{"GET", "/sec/getusers/:Token", GetAllUsers},
		&rest.Route{"GET", "/sec/deleteuser/:ID.:Token", DeleteUser},
		&rest.Route{"POST", "/sec/updaterole", UpdateRole},
		&rest.Route{"POST", "/sec/addrole", AddRole},
		&rest.Route{"GET", "/sec/deleterole/:ID.:Token", DeleteRole},
		&rest.Route{"GET", "/sec/getroles/:Token", GetAllRoles},
		&rest.Route{"GET", "/sec/getrole/:Name.:Token", GetRole},
		&rest.Route{"POST", "/backup/now", BackupNow},
		&rest.Route{"POST", "/backup/addschedule", AddSchedule},
		&rest.Route{"GET", "/backup/deleteschedule/:ID.:Token", DeleteSchedule},
		&rest.Route{"POST", "/backup/updateschedule", UpdateSchedule},
		&rest.Route{"GET", "/backup/getschedules/:ContainerName.:Token", GetAllSchedules},
		&rest.Route{"GET", "/backup/getschedule/:ID.:Token", GetSchedule},
		&rest.Route{"GET", "/backup/getstatus/:ID.:Token", GetStatus},
		&rest.Route{"GET", "/backup/getallstatus/:ID.:Token", GetAllStatus},
		&rest.Route{"GET", "/backup/nodes/:Token", GetBackupNodes},
		&rest.Route{"GET", "/mon/server/:Metric.:ServerID.:Interval.:Token", GetServerMetrics},
		&rest.Route{"GET", "/mon/container/pg2/:Name.:Interval.:Token", GetPG2},
		&rest.Route{"GET", "/mon/hc1/:Token", GetHC1},
		&rest.Route{"GET", "/version", GetVersion},
		&rest.Route{"GET", "/testcreate/:Token", TestCreate},
		&rest.Route{"GET", "/testdelete/:Token", TestDelete},
	)
	if err != nil {
		log.Fatal(err)
	}
	//	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServeTLS(":13000", "/cpmkeys/cert.pem", "/cpmkeys/key.pem", &handler))
}

type MonitorServerParam struct {
	ServerID string
	Metric   string
}
type MonitorContainerParam struct {
	ID           string
	Metric       string
	DatabaseName string
}

type MonitorOutput struct {
	Metric   string
	Response string
}

type Server struct {
	ID             string
	Name           string
	IPAddress      string
	DockerBridgeIP string
	PGDataPath     string
	ServerClass    string
	CreateDate     string
}

type ClusterProfiles struct {
	Size           string
	Count          string
	Algo           string
	MasterProfile  string
	StandbyProfile string
	MasterServer   string
	StandbyServer  string
	Token          string
}
type Profiles struct {
	SmallCPU  string
	SmallMEM  string
	MediumCPU string
	MediumMEM string
	LargeCPU  string
	LargeMEM  string
	Token     string
}

type Setting struct {
	Name       string
	Value      string
	UpdateDate string
}

type Settings struct {
	AdminURL       string
	DockerRegistry string
	PGPort         string
	DomainName     string
	Token          string
}

type Cluster struct {
	ID          string
	Name        string
	ClusterType string
	Status      string
	CreateDate  string
	Token       string
}

type ClusterNode struct {
	ID         string
	ClusterID  string
	ServerID   string
	Name       string
	Role       string
	Image      string
	CreateDate string
	Status     string
}

type LinuxStats struct {
	ID        string
	ClusterID string
	Stats     string
}

type PGStats struct {
	ID        string
	ClusterID string
	Stats     string
}
type SimpleStatus struct {
	Status string
}

type KubeResponse struct {
	URL string
}

type PostgresSetting struct {
	Name           string
	CurrentSetting string
	Source         string
}

type PostgresControldata struct {
	Name  string
	Value string
}

func Kube(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("Kube: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := KubeResponse{}
	kubeURL := os.Getenv("KUBE_URL")
	if kubeURL == "" {
		response.URL = "KUBE_URL is not set"
	} else {
		response.URL = kubeURL
	}
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&response)
}

func GetVersion(w rest.ResponseWriter, r *rest.Request) {

	w.(http.ResponseWriter).Write([]byte("0.9.0"))
}
