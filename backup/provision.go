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
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserveragent"
	"github.com/crunchydata/crunchy-postgresql-manager/kubeclient"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/template"
	"time"
)

func ProvisionBackupJob(args *BackupRequest) error {

	logit.Info.Println("backup.Provision called")
	logit.Info.Println("with scheduleid=" + args.ScheduleID)
	logit.Info.Println("with serverid=" + args.ServerID)
	logit.Info.Println("with servername=" + args.ServerName)
	logit.Info.Println("with serverip=" + args.ServerIP)
	logit.Info.Println("with containername=" + args.ContainerName)
	logit.Info.Println("with profilename=" + args.ProfileName)

	var params cpmserveragent.DockerRunArgs
	params = cpmserveragent.DockerRunArgs{}
	params.Image = "crunchydata/cpm-backup-job"
	params.ServerID = args.ServerID
	backupcontainername := args.ContainerName + "-backup"
	params.ContainerName = backupcontainername
	params.Standalone = "false"

	//get server info
	server, err := admindb.GetServer(params.ServerID)
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}

	params.PGDataPath = server.PGDataPath + "/" + backupcontainername + "/" + getFormattedDate()

	//get the docker profile settings
	var setting admindb.Setting
	setting, err = admindb.GetSetting("S-DOCKER-PROFILE-CPU")
	params.CPU = setting.Value
	setting, err = admindb.GetSetting("S-DOCKER-PROFILE-MEM")
	params.MEM = setting.Value

	params.EnvVars = make(map[string]string)

	setting, err = admindb.GetSetting("DOMAIN-NAME")
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
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

	setting, err = admindb.GetSetting("PG-PORT")
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}
	params.EnvVars["BACKUP_PORT"] = setting.Value
	params.EnvVars["BACKUP_USER"] = "postgres"
	params.EnvVars["BACKUP_SERVER_URL"] = "cpm-backup" + "." + domain + ":" + "13000"

	//provision the volume
	var responseStr string
	responseStr, err = cpmserveragent.AgentCommand("provisionvolume.sh",
		params.PGDataPath,
		server.IPAddress)
	logit.Info.Println(responseStr)

	//run the container
	if kubeEnv {
		//create a pod template to run the cpm-backup-job
		//create the pod
		err = kubeclient.DeletePod(kubeURL, backupcontainername)
		var podInfo = template.KubePodParams{
			ID:                   backupcontainername,
			PODID:                "",
			CPU:                  "0",
			MEM:                  "0",
			IMAGE:                params.Image,
			VOLUME:               params.EnvVars["BACKUP_PATH"],
			PORT:                 params.EnvVars["BACKUP_PORT"],
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
			logit.Error.Println(err.Error())
			return err
		}
	} else {
		var output string
		params.CommandPath = "docker-run-backup.sh"

		output, err = cpmserveragent.AgentDockerRun(params, server.IPAddress)

		if err != nil {
			logit.Error.Println("Provision: " + output)
			return err
		}
		logit.Info.Println("docker-run-backup.sh output=" + output)
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
