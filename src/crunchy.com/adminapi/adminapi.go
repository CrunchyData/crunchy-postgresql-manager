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
	"crunchy.com/backup"
	"crunchy.com/cpmagent"
	"crunchy.com/util"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func init() {
	fmt.Println("before parsing in init")
	flag.Parse()
}

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
		//&rest.Route{"POST", "/node", PostNode},
		&rest.Route{"GET", "/provision/:Profile.:Image.:ServerID.:ContainerName.:Standalone.:Token", Provision},
		&rest.Route{"GET", "/node/:ID.:Token", GetNode},
		&rest.Route{"GET", "/kube/:Token", Kube},
		&rest.Route{"GET", "/deletenode/:ID.:Token", DeleteNode},
		&rest.Route{"GET", "/monitor/server-getinfo/:ServerID.:Metric.:Token", MonitorServerGetInfo},
		&rest.Route{"GET", "/monitor/container-getinfo/:ID.:Metric.:Token", MonitorContainerGetInfo},
		&rest.Route{"GET", "/monitor/container-loadtest/:ID.:Metric.:Writes.:Token", MonitorContainerGetInfo},
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
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8080", &handler))
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

func Kube(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("Kube: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
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

func MonitorContainerGetInfo(w rest.ResponseWriter, r *rest.Request) {
	var err error

	err = secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("MonitorContainerGetInfo: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	Metric := r.PathParam("Metric")

	if ID == "" {
		rest.Error(w, "ID required", 400)
		return
	}
	if Metric == "" {
		rest.Error(w, "Metric required", 400)
		return
	}

	if Metric == "statreplication" ||
		Metric == "loadtest" || Metric == "bgwriter" ||
		Metric == "statdatabase" {
	} else {
		glog.Errorln("MonitorContainerGetInfo: error-invalid metric type")
		err = errors.New("invalid metric type")
		rest.Error(w, err.Error(), 400)
		return
	}

	var InsertCount = ""
	if Metric == "loadtest" {
		InsertCount = r.PathParam("Writes")
		if InsertCount == "" {
			rest.Error(w, "Writes param required", 400)
			return
		}
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("MonitorContainerGetInfo:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//$1 - host node.Name
	//$2 - port "5432"
	//$3 - user "cpmtest"
	//$4 - password "cpmtest"
	//$5 - database "cpmtest"
	//$6 - insert count - only valid for loadtest Metric

	//hardcoded for now, TODO pull from metadata
	var port = "5432"
	var user = "cpmtest"
	var password = "cpmtest"
	var database = "cpmtest"

	cmd := exec.Command(CPMBIN+Metric,
		node.Name,
		port,
		user,
		password,
		database,
		InsertCount)

	for i := 0; i < len(cmd.Args); i++ {
		glog.Infoln("MonitorContainerGetInfo:" + cmd.Args[i])
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		glog.Errorln("MonitorContainerGetInfo:" + err.Error())
		glog.Flush()
		errorString := fmt.Sprintf("%s\n%s\n%s\n", err.Error(), out.String(), stderr.String())
		rest.Error(w, errorString, 400)
		return
	}
	glog.Infoln("MonitorContainerGetInfo: command output was " + out.String())

	//w.(http.ResponseWriter).Write([]byte(output))
	w.(http.ResponseWriter).Write([]byte(out.String()))
	w.WriteHeader(http.StatusOK)
}

func MonitorContainerLoadtest(w rest.ResponseWriter, r *rest.Request) {
	ID := r.PathParam("ID")
	Writes := r.PathParam("Writes")

	if ID == "" {
		rest.Error(w, "ID required", 400)
		return
	}
	if Writes == "" {
		rest.Error(w, "Writes required", 400)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("MonitorContainerGetInfo:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	server, err2 := admindb.GetDBServer(node.ServerID)
	if err2 != nil {
		glog.Errorln("MonitorContainerGetInfo:" + err2.Error())
		rest.Error(w, err2.Error(), 400)
		return
	}

	var output string
	var port = "5432"

	output, err = cpmagent.AgentCommandConfigureNode(CPMBIN+"loadtest", node.Name,
		port, Writes, "", "", "", "", server.IPAddress)
	if err != nil {
		glog.Errorln("MonitorContainerGetInfo:" + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	w.(http.ResponseWriter).Write([]byte(output))
	w.WriteHeader(http.StatusOK)
}
