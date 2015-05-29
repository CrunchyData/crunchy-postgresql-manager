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
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserveragent"
	"github.com/crunchydata/crunchy-postgresql-manager/kubeclient"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/template"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

//docker run
func Provision(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-container")
	if err != nil {
		logit.Error.Println("Provision: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	params := new(cpmserveragent.DockerRunArgs)
	PROFILE := r.PathParam("Profile")
	params.Image = r.PathParam("Image")
	params.ServerID = r.PathParam("ServerID")
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
	logit.Info.Println("params.ContainerName=" + params.ContainerName)
	logit.Info.Println("params.Standalone=" + params.Standalone)

	err = provisionImpl(params, PROFILE, false)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = provisionImplInit(params, PROFILE, false)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func provisionImpl(params *cpmserveragent.DockerRunArgs, PROFILE string, standby bool) error {
	logit.Info.Println("PROFILE: provisionImpl starts 1")

	var errorStr string
	//make sure the container name is not already taken
	_, err2 := admindb.GetContainerByName(params.ContainerName)
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
	server, err := admindb.GetServer(params.ServerID)
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
		responseStr, err = cpmserveragent.AgentCommand("provisionvolume.sh",
			params.PGDataPath,
			server.IPAddress)
		if err != nil {
			logit.Error.Println("Provision: problem in provisionvolume call" + err.Error())
			return err
		}
		logit.Info.Println("Provision: provisionvolume call response=" + responseStr)
	}
	logit.Info.Println("PROFILE provisionImpl 3 provision volume completed")

	//run docker run to create the container

	params.CPU, params.MEM, err = getDockerResourceSettings(PROFILE)
	if err != nil {
		logit.Error.Println("Provision: problem in getting profiles call" + err.Error())
		return err
	}
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting("PG-PORT")
	if err != nil {
		logit.Error.Println("Provision:PG-PORT setting error " + err.Error())
		return err
	}

	var output string

	if !KubeEnv {
		//remove any existing docker containers with this name
		logit.Info.Println("PROFILE provisionImpl remove old container start")
		responseStr, err = cpmserveragent.DockerRemoveContainer(params.ContainerName,
			server.IPAddress)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		logit.Info.Println("PROFILE provisionImpl remove old container end")
		params.CommandPath = "docker-run.sh"
		output, err = cpmserveragent.AgentDockerRun(*params, server.IPAddress)

		if err != nil {
			logit.Error.Println("Provision: " + output)
			return err
		}
		logit.Info.Println("docker-run.sh output=" + output)
		logit.Info.Println("PROFILE provisionImpl end of docker-run")
	} else {
		//delete the kube pod with this name
		//we only log an error, this is ok because
		//we can get a 'not found' as an error
		err = kubeclient.DeletePod(KubeURL, params.ContainerName)
		logit.Info.Println("after delete pod")
		if err != nil {
			logit.Info.Println("Provision:" + err.Error())
		}

		err = kubeclient.DeleteService(KubeURL, params.ContainerName)
		if err != nil {
			logit.Info.Println("Provision:" + err.Error())
		}

		err = kubeclient.DeleteService(KubeURL, params.ContainerName+"-db")
		if err != nil {
			logit.Info.Println("Provision:" + err.Error())
		}

		podInfo := template.KubePodParams{
			ID:                   params.ContainerName,
			PODID:                params.ContainerName,
			CPU:                  params.CPU,
			MEM:                  params.MEM,
			IMAGE:                params.Image,
			VOLUME:               params.PGDataPath,
			PORT:                 "13000",
			BACKUP_NAME:          "",
			BACKUP_SERVERNAME:    "",
			BACKUP_SERVERIP:      "",
			BACKUP_SCHEDULEID:    "",
			BACKUP_PROFILENAME:   "",
			BACKUP_CONTAINERNAME: "",
			BACKUP_PATH:          "",
			BACKUP_HOST:          "",
			BACKUP_PORT:          "",
			BACKUP_USER:          "",
			BACKUP_SERVER_URL:    "",
		}
		err = kubeclient.CreatePod(KubeURL, podInfo)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}

		//create the service to the admin port 13000
		err = kubeclient.CreateService(KubeURL, podInfo)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}

		//create the service to the PG port
		podInfo.PORT = pgport.Value
		podInfo.ID = podInfo.ID + "-db"
		err = kubeclient.CreateService(KubeURL, podInfo)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		//we have to wait here since the Kube sometimes
		//is not that fast in setting up the service
		//for a pod..choosing 15 seconds to wait
		time.Sleep(15000 * time.Millisecond)
	}

	dbnode := admindb.Container{}
	dbnode.ID = ""
	dbnode.Name = params.ContainerName
	dbnode.Image = params.Image
	dbnode.ClusterID = "-1"
	dbnode.ProjectID = "1"
	dbnode.ServerID = params.ServerID

	if params.Standalone == "true" {
		dbnode.Role = "standalone"
	} else {
		dbnode.Role = "unassigned"
	}

	var strid int
	strid, err = admindb.InsertContainer(dbnode)
	newid := strconv.Itoa(strid)
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}
	dbnode.ID = newid

	//register default db users on the new node
	err = createDBUsers(dbnode)

	return err

}

//currently we define default DB users (postgres, cpmtest, pgpool)
//for all database containers
func createDBUsers(dbnode admindb.Container) error {
	var err error
	var password admindb.Setting

	//get the postgres password
	password, err = admindb.GetSetting("POSTGRESPSW")
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	//register postgres user
	var user = admindb.ContainerUser{}
	user.Containername = dbnode.Name
	user.Usename = "postgres"
	user.Passwd = password.Value
	_, err = admindb.AddContainerUser(user)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//cpmtest and pgpool users are created by the node-setup.sql script
	//here, we just register them when we create a new node

	//get the cpmtest password
	password, err = admindb.GetSetting("CPMTESTPSW")
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	//register cpmtest user
	user.Containername = dbnode.Name
	user.Usename = "cpmtest"
	user.Passwd = password.Value
	_, err = admindb.AddContainerUser(user)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//get the pgpool password
	password, err = admindb.GetSetting("PGPOOLPSW")
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	user.Containername = dbnode.Name
	user.Usename = "pgpool"
	user.Passwd = password.Value
	//register pgpool user
	_, err = admindb.AddContainerUser(user)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	return err
}

func provisionImplInit(params *cpmserveragent.DockerRunArgs, PROFILE string, standby bool) error {
	//go get the domain name from the settings
	var domainname admindb.Setting
	var pgport admindb.Setting
	var err error

	domainname, err = admindb.GetSetting("DOMAIN-NAME")
	if err != nil {
		logit.Error.Println("Provision:DOMAIN-NAME setting error " + err.Error())
		return err
	}
	pgport, err = admindb.GetSetting("PG-PORT")
	if err != nil {
		logit.Error.Println("Provision:PG-PORT setting error " + err.Error())
		return err
	}

	fqdn := params.ContainerName + "." + domainname.Value

	//we are depending on a DNS entry being created shortly after
	//creating the node either in Docker or Kube!
	//you might need to wait here until you can reach the new node's agent
	logit.Info.Println("PROFILE waiting till DNS ready")
	err = waitTillReady(fqdn)
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}

	if standby {
		logit.Info.Println("standby node being created, will not initdb")
	} else {
		//initdb on the new node

		logit.Info.Println("PROFILE running initdb on the node")
		var output string

		output, err = PGCommand("initdb.sh", fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		logit.Info.Println("initdb output was" + output)
		logit.Info.Println("PROFILE initdb completed")
		//create postgresql.conf
		var data string
		var mode = "standalone"

		data, err = template.Postgresql(mode, pgport.Value, "")

		//place postgresql.conf on new node
		err = RemoteWritefile("/pgdata/postgresql.conf", data, fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		//create pg_hba.conf
		data, err = template.Hba(KubeEnv, mode, params.ContainerName, pgport.Value, "", domainname.Value)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		//place pg_hba.conf on new node
		err = RemoteWritefile("/pgdata/pg_hba.conf", data, fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		logit.Info.Println("PROFILE templates all built and copied to node")
		//start pg on new node
		output, err = PGCommand("startpg.sh", fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		logit.Info.Println("startpg output was" + output)

		//seed database with initial objects
		output, err = PGCommand("seed.sh", fqdn)
		if err != nil {
			logit.Error.Println("Provision:" + err.Error())
			return err
		}
		logit.Info.Println("seed output was" + output)
	}
	logit.Info.Println("PROFILE node provisioning completed")

	return nil
}

func RemoteWritefile(path string, filecontents string, ipaddress string) error {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		logit.Error.Println("RemoteWritefile: dialing:" + err.Error())
		return err
	}
	if client == nil {
		logit.Error.Println("RemoteWritefile: dialing:" + err.Error())
		return errors.New("client was null on rpc call to " + ipaddress)
	}

	var command cpmserveragent.Command

	args := &cpmserveragent.Args{}
	args.A = filecontents
	args.B = path
	err = client.Call("Command.Writefile", args, &command)
	if err != nil {
		logit.Error.Println("RemoteWritefile:  Command Writefile " + args.B + " error:" + err.Error())
		return err
	}
	logit.Info.Println("RemoteWritefile: Writefile output " + command.Output)
	return nil
}

func PGCommand(pgcommand string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		logit.Error.Println("PGCommand: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logit.Error.Println("PGCommand: dialing:" + err.Error())
		return "", errors.New("client was null on pgcommand rpc to " + ipaddress)
	}

	var command cpmserveragent.Command

	args := &cpmserveragent.Args{}
	args.A = util.GetBase() + "/bin/" + pgcommand
	err = client.Call("Command.PGCommand", args, &command)
	if err != nil {
		logit.Error.Println("PGCommand:  Command PGCommand " + args.A + " error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func waitTillReady(container string) error {

	var err error
	for i := 0; i < 40; i++ {
		err = RemoteWritefile("/tmp/waitTest", "waitTillReady was here", container)
		if err != nil {
			logit.Error.Println("waitTillReady:waited for cpmnodeagent on " + container)
			time.Sleep(2000 * time.Millisecond)
		} else {
			logit.Info.Println("waitTillReady:connected to cpmnodeagent on " + container)
			return nil
		}
	}
	logit.Info.Println("waitTillReady: could not connect to cpmnodeagent on " + container)
	return errors.New("could not connect to cpmnodeagent on " + container)

}

//return the CPU MEM settings
func getDockerResourceSettings(size string) (string, string, error) {
	var CPU, MEM string
	var setting admindb.Setting
	var err error

	switch size {
	case "SM":
		setting, err = admindb.GetSetting("S-DOCKER-PROFILE-CPU")
		CPU = setting.Value
		setting, err = admindb.GetSetting("S-DOCKER-PROFILE-MEM")
		MEM = setting.Value
	case "MED":
		setting, err = admindb.GetSetting("M-DOCKER-PROFILE-CPU")
		CPU = setting.Value
		setting, err = admindb.GetSetting("M-DOCKER-PROFILE-MEM")
		MEM = setting.Value
	default:
		setting, err = admindb.GetSetting("L-DOCKER-PROFILE-CPU")
		CPU = setting.Value
		setting, err = admindb.GetSetting("L-DOCKER-PROFILE-MEM")
		MEM = setting.Value
	}

	return CPU, MEM, err

}
