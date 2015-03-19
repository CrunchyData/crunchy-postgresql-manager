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
	"database/sql"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

//docker run
func Provision(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-container")
	if err != nil {
		glog.Errorln("Provision: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	params := new(cpmagent.DockerRunArgs)
	PROFILE := r.PathParam("Profile")
	params.Image = r.PathParam("Image")
	params.ServerID = r.PathParam("ServerID")
	params.ContainerName = r.PathParam("ContainerName")
	params.Standalone = r.PathParam("Standalone")

	errorStr := ""

	if PROFILE == "" {
		glog.Errorln("Provision error profile required")
		errorStr = "Profile required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if params.ServerID == "" {
		glog.Errorln("Provision error serverid required")
		errorStr = "ServerID required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if params.ContainerName == "" {
		glog.Errorln("Provision error containername required")
		errorStr = "ContainerName required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if params.Image == "" {
		glog.Errorln("Provision error image required")
		errorStr = "Image required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if params.Standalone == "" {
		glog.Errorln("Provision error standalone flag required")
		errorStr = "Standalone required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	glog.Infoln("params.Image=" + params.Image)
	glog.Infoln("params.Profile=" + PROFILE)
	glog.Infoln("params.ServerID=" + params.ServerID)
	glog.Infoln("params.ContainerName=" + params.ContainerName)
	glog.Infoln("params.Standalone=" + params.Standalone)

	err = provisionImpl(params, PROFILE, false)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func provisionImpl(params *cpmagent.DockerRunArgs, PROFILE string, standby bool) error {
	glog.Infoln("PROFILE: provisionImpl starts 1")

	var errorStr string
	//make sure the container name is not already taken
	_, err2 := admindb.GetDBNodeByName(params.ContainerName)
	if err2 != nil {
		if err2 != sql.ErrNoRows {
			return err2
		}
	} else {
		errorStr = "container name" + params.ContainerName + " already used can't provision"
		glog.Errorln("Provision error" + errorStr)
		return errors.New(errorStr)
	}

	//go get the IPAddress
	server, err := admindb.GetDBServer(params.ServerID)
	if err != nil {
		glog.Errorln("Provision:" + err.Error())
		return err
	}

	glog.Infoln("provisioning on server " + server.IPAddress)

	//for database nodes, on the target server, we need to allocate
	//a disk volume for the /pgdata container volume to work with
	//this causes a volume to be created with the directory
	//named the same as the container name

	var responseStr string

	params.PGDataPath = server.PGDataPath + "/" + params.ContainerName

	glog.Infoln("PROFILE provisionImpl 2 about to provision volume")
	if params.Image != "cpm-pgpool" {
		responseStr, err = cpmagent.AgentCommand(CPMBIN+"provisionvolume.sh",
			params.PGDataPath,
			server.IPAddress)
		if err != nil {
			glog.Errorln("Provision: problem in provisionvolume call" + err.Error())
			return err
		}
		glog.Infoln("Provision: provisionvolume call response=" + responseStr)
	}
	glog.Infoln("PROFILE provisionImpl 3 provision volume completed")

	//run docker run to create the container

	params.CPU, params.MEM, err = getDockerResourceSettings(PROFILE)
	if err != nil {
		glog.Errorln("Provision: problem in getting profiles call" + err.Error())
		return err
	}

	var output string

	if !kubeEnv {
		//remove any existing docker containers with this name
		glog.Infoln("PROFILE provisionImpl remove old container start")
		responseStr, err = cpmagent.DockerRemoveContainer(params.ContainerName,
			server.IPAddress)
		if err != nil {
			glog.Errorln("Provision:" + err.Error())
			return err
		}
		glog.Infoln("PROFILE provisionImpl remove old container end")
		params.CommandPath = CPMBIN + "docker-run.sh"
		output, err = cpmagent.AgentDockerRun(*params, server.IPAddress)

		if err != nil {
			glog.Errorln("Provision: " + output)
			return err
		}
		glog.Infoln("docker-run.sh output=" + output)
		glog.Infoln("PROFILE provisionImpl end of docker-run")
	} else {
		//delete the kube pod with this name
		//we only log an error, this is ok because
		//we can get a 'not found' as an error
		err = DeletePod(kubeURL, params.ContainerName)
		glog.Infoln("after delete pod")
		if err != nil {
			glog.Infoln("Provision:" + err.Error())
		}

		podInfo := template.KubePodParams{
			params.ContainerName,
			params.ContainerName,
			params.CPU, params.MEM,
			params.Image,
			params.PGDataPath, "13000"}
		glog.Infoln("before create pod")
		err = CreatePod(kubeURL, podInfo)
		glog.Infoln("after create pod")
		if err != nil {
			glog.Errorln("Provision:" + err.Error())
			return err
		}
		//we have to wait here since the Kube sometimes
		//is not that fast in setting up the service
		//for a pod..choosing 15 seconds to wait
		time.Sleep(15000 * time.Millisecond)
	}

	dbnode := admindb.DBClusterNode{}
	dbnode.ID = ""
	dbnode.Name = params.ContainerName
	dbnode.Image = params.Image
	dbnode.ClusterID = "-1"
	dbnode.ServerID = params.ServerID

	if params.Standalone == "true" {
		dbnode.Role = "standalone"
	} else {
		dbnode.Role = "unassigned"
	}

	var strid int
	strid, err = admindb.InsertDBNode(dbnode)
	newid := strconv.Itoa(strid)
	if err != nil {
		glog.Errorln("Provision:" + err.Error())
		return err
	}
	dbnode.ID = newid

	if params.Image == "cpm-pgpool" {
		return nil
	}

	//go get the domain name from the settings
	var domainname admindb.DBSetting
	domainname, err = admindb.GetDBSetting("DOMAIN-NAME")
	if err != nil {
		glog.Errorln("Provision:DOMAIN-NAME setting error " + err.Error())
		return err
	}
	fqdn := params.ContainerName + "." + domainname.Value

	//we are depending on a DNS entry being created shortly after
	//creating the node either in Docker or Kube!
	//you might need to wait here until you can reach the new node's agent
	glog.Infoln("PROFILE waiting till DNS ready")
	err = waitTillReady(fqdn)
	if err != nil {
		glog.Errorln("Provision:" + err.Error())
		return err
	}

	if standby {
		glog.Infoln("standby node being created, will not initdb")
	} else {
		//initdb on the new node

		glog.Infoln("PROFILE running initdb on the node")
		output, err = PGCommand(CPMBIN+"initdb.sh", fqdn)
		if err != nil {
			glog.Errorln("Provision:" + err.Error())
			return err
		}
		glog.Infoln("initdb output was" + output)
		glog.Infoln("PROFILE initdb completed")
		//create postgresql.conf
		var data string
		var mode = "standalone"
		var port = "5432"

		data, err = template.Postgresql(mode, port, "")

		//place postgresql.conf on new node
		err = RemoteWritefile("/pgdata/postgresql.conf", data, fqdn)
		if err != nil {
			glog.Errorln("Provision:" + err.Error())
			return err
		}
		//create pg_hba.conf
		data, err = template.Hba(kubeEnv, mode, params.ContainerName, port, "", domainname.Value)
		if err != nil {
			glog.Errorln("Provision:" + err.Error())
			return err
		}
		//place pg_hba.conf on new node
		err = RemoteWritefile("/pgdata/pg_hba.conf", data, fqdn)
		if err != nil {
			glog.Errorln("Provision:" + err.Error())
			return err
		}
		glog.Infoln("PROFILE templates all built and copied to node")
		//start pg on new node
		output, err = PGCommand(CPMBIN+"startpg.sh", fqdn)
		if err != nil {
			glog.Errorln("Provision:" + err.Error())
			return err
		}
		glog.Infoln("startpg output was" + output)

		//seed database with initial objects
		output, err = PGCommand(CPMBIN+"seed.sh", fqdn)
		if err != nil {
			glog.Errorln("Provision:" + err.Error())
			return err
		}
		glog.Infoln("seed output was" + output)
	}
	glog.Infoln("PROFILE node provisioning completed")

	return nil
}

func RemoteWritefile(path string, filecontents string, ipaddress string) error {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("RemoteWritefile: dialing:" + err.Error())
		return err
	}
	if client == nil {
		glog.Errorln("RemoteWritefile: dialing:" + err.Error())
		return errors.New("client was null on rpc call to " + ipaddress)
	}

	var command cpmagent.Command

	args := &cpmagent.Args{}
	args.A = filecontents
	args.B = path
	err = client.Call("Command.Writefile", args, &command)
	if err != nil {
		glog.Errorln("RemoteWritefile:  Command Writefile " + args.B + " error:" + err.Error())
		return err
	}
	glog.Infoln("RemoteWritefile: Writefile output " + command.Output)
	return nil
}

func PGCommand(pgcommand string, ipaddress string) (string, error) {
	client, err := rpc.DialHTTP("tcp", ipaddress+":13000")
	if err != nil {
		glog.Errorln("PGCommand: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		glog.Errorln("PGCommand: dialing:" + err.Error())
		return "", errors.New("client was null on pgcommand rpc to " + ipaddress)
	}

	var command cpmagent.Command

	args := &cpmagent.Args{}
	args.A = pgcommand
	err = client.Call("Command.PGCommand", args, &command)
	if err != nil {
		glog.Errorln("PGCommand:  Command PGCommand " + args.A + " error:" + err.Error())
		return "", err
	}
	return command.Output, nil
}

func waitTillReady(container string) error {

	var err error
	for i := 0; i < 40; i++ {
		err = RemoteWritefile("/tmp/waitTest", "waitTillReady was here", container)
		if err != nil {
			glog.Errorln("waitTillReady:waited for cpmagent on " + container)
			time.Sleep(2000 * time.Millisecond)
		} else {
			glog.Infoln("waitTillReady:connected to cpmagent on " + container)
			return nil
		}
	}
	glog.Infoln("waitTillReady: could not connect to cpmagent on " + container)
	return errors.New("could not connect to cpmagent on " + container)

}

//return the CPU MEM settings
func getDockerResourceSettings(size string) (string, string, error) {
	var CPU, MEM string
	var setting admindb.DBSetting
	var err error

	switch size {
	case "SM":
		setting, err = admindb.GetDBSetting("S-DOCKER-PROFILE-CPU")
		CPU = setting.Value
		setting, err = admindb.GetDBSetting("S-DOCKER-PROFILE-MEM")
		MEM = setting.Value
	case "MED":
		setting, err = admindb.GetDBSetting("M-DOCKER-PROFILE-CPU")
		CPU = setting.Value
		setting, err = admindb.GetDBSetting("M-DOCKER-PROFILE-MEM")
		MEM = setting.Value
	default:
		setting, err = admindb.GetDBSetting("L-DOCKER-PROFILE-CPU")
		CPU = setting.Value
		setting, err = admindb.GetDBSetting("L-DOCKER-PROFILE-MEM")
		MEM = setting.Value
	}

	return CPU, MEM, err

}
