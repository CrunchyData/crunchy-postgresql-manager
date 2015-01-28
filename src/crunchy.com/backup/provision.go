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
	"crunchy.com/logutil"
	"fmt"
	"time"
)

func ProvisionBackupJob(args *BackupRequest) error {

	logutil.Log("backup.Provision called")
	logutil.Log("with scheduleid=" + args.ScheduleID)
	logutil.Log("with serverid=" + args.ServerID)
	logutil.Log("with servername=" + args.ServerName)
	logutil.Log("with serverip=" + args.ServerIP)
	logutil.Log("with containername=" + args.ContainerName)
	logutil.Log("with profilename=" + args.ProfileName)

	var params cpmagent.DockerRunArgs
	params = cpmagent.DockerRunArgs{}
	params.Image = "crunchy-backup-job"
	params.ServerID = args.ServerID
	backupcontainername := args.ContainerName + "-backup"
	params.ContainerName = backupcontainername
	params.Standalone = "false"

	//get server info
	server, err := admindb.GetDBServer(params.ServerID)
	if err != nil {
		logutil.Log("Provision:" + err.Error())
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
		logutil.Log("Provision:" + err.Error())
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
		logutil.Log("Provision:" + err.Error())
		return err
	}
	params.EnvVars["BACKUP_PORT"] = setting.Value
	params.EnvVars["BACKUP_USER"] = "postgres"
	params.EnvVars["BACKUP_SERVER_URL"] = "cluster-backup" + "." + domain + ":" + "13010"

	//provision the volume
	var responseStr string
	responseStr, err = cpmagent.AgentCommand("/cluster/bin/provisionvolume.sh",
		params.PGDataPath,
		server.IPAddress)
	logutil.Log(responseStr)
	//run the container
	kubeEnv := false

	if !kubeEnv {
		var output string
		params.CommandPath = "/cluster/bin/docker-run-backup.sh"

		output, err = cpmagent.AgentDockerRun(params, server.IPAddress)

		if err != nil {
			logutil.Log("Provision: " + output)
			return err
		}
		logutil.Log("docker-run-backup.sh output=" + output)
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
