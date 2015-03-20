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

package backup

import (
	"crunchy.com/admindb"
	"crunchy.com/cpmagent"
	"crunchy.com/kubeclient"
	"crunchy.com/template"
	"fmt"
	"github.com/golang/glog"
	"time"
)

func ProvisionBackupJob(args *BackupRequest) error {

	glog.Infoln("backup.Provision called")
	glog.Infoln("with scheduleid=" + args.ScheduleID)
	glog.Infoln("with serverid=" + args.ServerID)
	glog.Infoln("with servername=" + args.ServerName)
	glog.Infoln("with serverip=" + args.ServerIP)
	glog.Infoln("with containername=" + args.ContainerName)
	glog.Infoln("with profilename=" + args.ProfileName)
	glog.Flush()

	var params cpmagent.DockerRunArgs
	params = cpmagent.DockerRunArgs{}
	params.Image = "crunchydata/cpm-backup-job"
	params.ServerID = args.ServerID
	backupcontainername := args.ContainerName + "-backup"
	params.ContainerName = backupcontainername
	params.Standalone = "false"

	//get server info
	server, err := admindb.GetDBServer(params.ServerID)
	if err != nil {
		glog.Errorln("Provision:" + err.Error())
		glog.Flush()
		return err
	}

	params.PGDataPath = server.PGDataPath + "/" + backupcontainername + "/" + getFormattedDate()

	//get the docker profile settings
	var setting admindb.DBSetting
	setting, err = admindb.GetDBSetting("S-DOCKER-PROFILE-CPU")
	params.CPU = setting.Value
	setting, err = admindb.GetDBSetting("S-DOCKER-PROFILE-MEM")
	params.MEM = setting.Value

	params.EnvVars = make(map[string]string)

	setting, err = admindb.GetDBSetting("DOMAIN-NAME")
	if err != nil {
		glog.Errorln("Provision:" + err.Error())
		glog.Flush()
		return err
	}
	var domain = setting.Value

	params.EnvVars["BACKUP_NAME"] = backupcontainername
	params.EnvVars["BACKUP_SERVERNAME"] = server.Name
	params.EnvVars["BACKUP_SERVERIP"] = server.IPAddress
	params.EnvVars["BACKUP_SCHEDULEID"] = args.ScheduleID
	params.EnvVars["BACKUP_PROFILENAME"] = args.ProfileName
	params.EnvVars["BACKUP_CONTAINERNAME"] = args.ContainerName
	params.EnvVars["BACKUP_PATH"] = params.PGDataPath
	params.EnvVars["BACKUP_HOST"] = args.ContainerName + "." + setting.Value

	setting, err = admindb.GetDBSetting("PG-PORT")
	if err != nil {
		glog.Errorln("Provision:" + err.Error())
		glog.Flush()
		return err
	}
	params.EnvVars["BACKUP_PORT"] = setting.Value
	params.EnvVars["BACKUP_USER"] = "postgres"
	params.EnvVars["BACKUP_SERVER_URL"] = "cpm-backup" + "." + domain + ":" + "13000"

	//provision the volume
	var responseStr string
	responseStr, err = cpmagent.AgentCommand(CPMBIN+"provisionvolume.sh",
		params.PGDataPath,
		server.IPAddress)
	glog.Infoln(responseStr)
	glog.Flush()

	//run the container
	if kubeEnv {
		//create a pod template to run the cpm-backup-job
		//create the pod
		err = kubeclient.DeletePod(kubeURL, backupcontainername)
		glog.Flush()
		var podInfo = template.KubePodParams{
			ID:                   backupcontainername,
			PODID:                "",
			CPU:                  "0",
			MEM:                  "0",
			IMAGE:                params.Image,
			VOLUME:               params.EnvVars["BACKUP_PATH"],
			PORT:                 "5432",
			BACKUP_NAME:          params.EnvVars["BACKUP_NAME"],
			BACKUP_SERVERNAME:    params.EnvVars["BACKUP_SERVERNAME"],
			BACKUP_SERVERIP:      params.EnvVars["BACKUP_SERVERIP"],
			BACKUP_SCHEDULEID:    params.EnvVars["BACKUP_SCHEDULEID"],
			BACKUP_PROFILENAME:   params.EnvVars["BACKUP_PROFILENAME"],
			BACKUP_CONTAINERNAME: params.EnvVars["BACKUP_CONTAINERNAME"],
			BACKUP_PATH:          params.EnvVars["BACKUP_PATH"],
			BACKUP_HOST:          params.EnvVars["BACKUP_HOST"],
			BACKUP_PORT:          params.EnvVars["BACKUP_PORT"],
			BACKUP_USER:          params.EnvVars["BACKUP_USER"],
			BACKUP_SERVER_URL:    params.EnvVars["BACKUP_SERVER_URL"],
		}

		err = kubeclient.CreatePod(kubeURL, podInfo)
		if err != nil {
			glog.Errorln(err.Error())
			glog.Flush()
			return err
		}
	} else {
		var output string
		params.CommandPath = CPMBIN + "docker-run-backup.sh"

		output, err = cpmagent.AgentDockerRun(params, server.IPAddress)

		if err != nil {
			glog.Errorln("Provision: " + output)
			glog.Flush()
			return err
		}
		glog.Infoln("docker-run-backup.sh output=" + output)
	}

	return nil
}

func getFormattedDate() string {
	t := time.Now()
	yyyy := t.Year()
	mm := t.Month()
	dd := t.Day()
	hh := t.Hour()
	min := t.Minute()
	output := fmt.Sprintf("%d%d%d%d%d", yyyy, mm, dd, hh, min)

	return output
}
