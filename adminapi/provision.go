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
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/template"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"strconv"
	"time"
)

//docker run
//TODO:  convert this to POST
func Provision(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-container")
	if err != nil {
		logit.Error.Println("Provision: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	params := &cpmserverapi.DockerRunRequest{}
	PROFILE := r.PathParam("Profile")
	params.Image = r.PathParam("Image")
	params.ServerID = r.PathParam("ServerID")
	params.ProjectID = r.PathParam("ProjectID")
	params.ContainerName = r.PathParam("ContainerName")
	params.Standalone = r.PathParam("Standalone")

	errorStr := ""

	if PROFILE == "" {
		logit.Error.Println("Provision error profile required")
		errorStr = "Profile required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if params.ServerID == "" {
		logit.Error.Println("Provision error serverid required")
		errorStr = "ServerID required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if params.ProjectID == "" {
		logit.Error.Println("Provision error ProjectID required")
		errorStr = "ProjectID required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if params.ContainerName == "" {
		logit.Error.Println("Provision error containername required")
		errorStr = "ContainerName required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if params.Image == "" {
		logit.Error.Println("Provision error image required")
		errorStr = "Image required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if params.Standalone == "" {
		logit.Error.Println("Provision error standalone flag required")
		errorStr = "Standalone required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	logit.Info.Println("params.Image=" + params.Image)
	logit.Info.Println("params.Profile=" + PROFILE)
	logit.Info.Println("params.ServerID=" + params.ServerID)
	logit.Info.Println("params.ProjectID=" + params.ProjectID)
	logit.Info.Println("params.ContainerName=" + params.ContainerName)
	logit.Info.Println("params.Standalone=" + params.Standalone)

	err = provisionImpl(dbConn, params, PROFILE, false)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = provisionImplInit(dbConn, params, PROFILE, false)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func provisionImpl(dbConn *sql.DB, params *cpmserverapi.DockerRunRequest, PROFILE string, standby bool) error {
	logit.Info.Println("PROFILE: provisionImpl starts 1")

	var errorStr string
	//make sure the container name is not already taken
	_, err2 := admindb.GetContainerByName(dbConn, params.ContainerName)
	if err2 != nil {
		if err2 != sql.ErrNoRows {
			return err2
		}
	} else {
		errorStr = "container name" + params.ContainerName + " already used can't provision"
		logit.Error.Println("Provision error" + errorStr)
		return errors.New(errorStr)
	}

	//go get the IPAddress
	server, err := admindb.GetServer(dbConn, params.ServerID)
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}

	logit.Info.Println("provisioning on server " + server.IPAddress)

	//for database nodes, on the target server, we need to allocate
	//a disk volume for the /pgdata container volume to work with
	//this causes a volume to be created with the directory
	//named the same as the container name

	var responseStr string

	params.PGDataPath = server.PGDataPath + "/" + params.ContainerName

	logit.Info.Println("PROFILE provisionImpl 2 about to provision volume")
	if params.Image != "cpm-pgpool" {
		preq := &cpmserverapi.DiskProvisionRequest{}
		preq.Path = params.PGDataPath
		var url = "http://" + server.IPAddress + ":10001"
		_, err = cpmserverapi.DiskProvisionClient(url, preq)
		if err != nil {
			logit.Error.Println("Provision: problem in provisionvolume call" + err.Error())
			return err
		}
		logit.Info.Println("Provision: provisionvolume call response=" + responseStr)
	}
	logit.Info.Println("PROFILE provisionImpl 3 provision volume completed")

	//run docker run to create the container

	params.CPU, params.MEM, err = getDockerResourceSettings(dbConn, PROFILE)
	if err != nil {
		logit.Error.Println("Provision: problem in getting profiles call" + err.Error())
		return err
	}
	/*
		var pgport admindb.Setting
		pgport, err = admindb.GetSetting(dbConn, "PG-PORT")
		if err != nil {
			logit.Error.Println("Provision:PG-PORT setting error " + err.Error())
			return err
		}
	*/

	var output string

	//remove any existing docker containers with this name
	logit.Info.Println("PROFILE provisionImpl remove old container start")
	rreq := &cpmserverapi.DockerRemoveRequest{}
	rreq.ContainerName = params.ContainerName
	var url = "http://" + server.IPAddress + ":10001"
	_, err = cpmserverapi.DockerRemoveClient(url, rreq)
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}
	logit.Info.Println("PROFILE provisionImpl remove old container end")
	params.CommandPath = "docker-run.sh"
	_, err = cpmserverapi.DockerRunClient(url, params)
	if err != nil {
		logit.Error.Println("Provision: " + output)
		return err
	}
	logit.Info.Println("docker-run.sh output=" + output)
	logit.Info.Println("PROFILE provisionImpl end of docker-run")

	dbnode := admindb.Container{}
	dbnode.ID = ""
	dbnode.Name = params.ContainerName
	dbnode.Image = params.Image
	dbnode.ClusterID = "-1"
	dbnode.ProjectID = params.ProjectID
	dbnode.ServerID = params.ServerID

	if params.Standalone == "true" {
		dbnode.Role = "standalone"
	} else {
		dbnode.Role = "unassigned"
	}

	var strid int
	strid, err = admindb.InsertContainer(dbConn, dbnode)
	newid := strconv.Itoa(strid)
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}
	dbnode.ID = newid

	//register default db users on the new node
	err = createDBUsers(dbConn, dbnode)

	return err

}

//currently we define default DB users (postgres, cpmtest, pgpool)
//for all database containers
func createDBUsers(dbConn *sql.DB, dbnode admindb.Container) error {
	var err error
	var password admindb.Setting

	//get the postgres password
	password, err = admindb.GetSetting(dbConn, "POSTGRESPSW")
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	//register postgres user
	var user = admindb.ContainerUser{}
	user.Containername = dbnode.Name
	user.Rolname = "postgres"
	user.Passwd = password.Value
	_, err = admindb.AddContainerUser(dbConn, user)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//cpmtest and pgpool users are created by the node-setup.sql script
	//here, we just register them when we create a new node

	//get the cpmtest password
	password, err = admindb.GetSetting(dbConn, "CPMTESTPSW")
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	//register cpmtest user
	user.Containername = dbnode.Name
	user.Rolname = "cpmtest"
	user.Passwd = password.Value
	_, err = admindb.AddContainerUser(dbConn, user)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//get the pgpool password
	password, err = admindb.GetSetting(dbConn, "PGPOOLPSW")
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	user.Containername = dbnode.Name
	user.Rolname = "pgpool"
	user.Passwd = password.Value
	//register pgpool user
	_, err = admindb.AddContainerUser(dbConn, user)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	return err
}

func provisionImplInit(dbConn *sql.DB, params *cpmserverapi.DockerRunRequest, PROFILE string, standby bool) error {
	//go get the domain name from the settings
	var domainname admindb.Setting
	var pgport admindb.Setting
	var err error

	domainname, err = admindb.GetSetting(dbConn, "DOMAIN-NAME")
	if err != nil {
		logit.Error.Println("Provision:DOMAIN-NAME setting error " + err.Error())
		return err
	}
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")
	if err != nil {
		logit.Error.Println("Provision:PG-PORT setting error " + err.Error())
		return err
	}

	fqdn := params.ContainerName + "." + domainname.Value

	//we are depending on a DNS entry being created shortly after
	//creating the node in Docker
	//you might need to wait here until you can reach the new node's agent
	logit.Info.Println("PROFILE waiting till DNS ready")
	err = waitTillReady(fqdn)
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}
	logit.Info.Println("checkpt 1")

	if standby {
		logit.Info.Println("standby node being created, will not initdb")
	} else {
		//initdb on the new node

		logit.Info.Println("PROFILE running initdb on the node")
		var resp cpmcontainerapi.InitdbResponse

		logit.Info.Println("checkpt 2")
		resp, err = cpmcontainerapi.InitdbClient(fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		logit.Info.Println("checkpt 3")
		logit.Info.Println("initdb output was" + resp.Output)
		logit.Info.Println("PROFILE initdb completed")
		//create postgresql.conf
		var data string
		var mode = "standalone"

		data, err = template.Postgresql(mode, pgport.Value, "")

		//place postgresql.conf on new node
		_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/postgresql.conf", data, fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		//create pg_hba.conf
		rules := make([]template.Rule, 0)
		data, err = template.Hba(dbConn, mode, params.ContainerName, pgport.Value, "", domainname.Value, rules)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		//place pg_hba.conf on new node
		_, err = cpmcontainerapi.RemoteWritefileClient("/pgdata/pg_hba.conf", data, fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		logit.Info.Println("PROFILE templates all built and copied to node")
		//start pg on new node
		var startResp cpmcontainerapi.StartPGResponse
		startResp, err = cpmcontainerapi.StartPGClient(fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		logit.Info.Println("startpg output was" + startResp.Output)

		//seed database with initial objects
		var seedResp cpmcontainerapi.SeedResponse
		seedResp, err = cpmcontainerapi.SeedClient(fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		logit.Info.Println("seed output was" + seedResp.Output)
	}
	logit.Info.Println("PROFILE node provisioning completed")

	return nil
}

func waitTillReady(container string) error {

	var err error
	for i := 0; i < 40; i++ {
		_, err = cpmcontainerapi.RemoteWritefileClient("/tmp/waitTest", "waitTillReady was here", container)
		if err != nil {
			logit.Error.Println("waitTillReady:waited for cpmcontainerapi on " + container)
			time.Sleep(2000 * time.Millisecond)
		} else {
			logit.Info.Println("waitTillReady:connected to cpmcontainerapi on " + container)
			return nil
		}
	}
	logit.Info.Println("waitTillReady: could not connect to cpmcontainerapi on " + container)
	return errors.New("could not connect to cpmcontainerapi on " + container)

}

//return the CPU MEM settings
func getDockerResourceSettings(dbConn *sql.DB, size string) (string, string, error) {
	var CPU, MEM string
	var setting admindb.Setting
	var err error

	switch size {
	case "SM":
		setting, err = admindb.GetSetting(dbConn, "S-DOCKER-PROFILE-CPU")
		CPU = setting.Value
		setting, err = admindb.GetSetting(dbConn, "S-DOCKER-PROFILE-MEM")
		MEM = setting.Value
	case "MED":
		setting, err = admindb.GetSetting(dbConn, "M-DOCKER-PROFILE-CPU")
		CPU = setting.Value
		setting, err = admindb.GetSetting(dbConn, "M-DOCKER-PROFILE-MEM")
		MEM = setting.Value
	default:
		setting, err = admindb.GetSetting(dbConn, "L-DOCKER-PROFILE-CPU")
		CPU = setting.Value
		setting, err = admindb.GetSetting(dbConn, "L-DOCKER-PROFILE-MEM")
		MEM = setting.Value
	}

	return CPU, MEM, err

}
