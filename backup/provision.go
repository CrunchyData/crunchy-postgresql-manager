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
	"database/sql"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"time"
)

func ProvisionBackupJob(dbConn *sql.DB, args *BackupRequest) error {

	logit.Info.Println("backup.Provision called")
	logit.Info.Println("with scheduleid=" + args.ScheduleID)
	logit.Info.Println("with serverid=" + args.ServerID)
	logit.Info.Println("with servername=" + args.ServerName)
	logit.Info.Println("with serverip=" + args.ServerIP)
	logit.Info.Println("with containername=" + args.ContainerName)
	logit.Info.Println("with profilename=" + args.ProfileName)

	params := &cpmserverapi.DockerRunRequest{}
	params.Image = "crunchydata/cpm-backup-job"
	params.ServerID = args.ServerID
	backupcontainername := args.ContainerName + "-backup"
	params.ContainerName = backupcontainername
	params.Standalone = "false"

	//get server info
	server, err := admindb.GetServer(dbConn, params.ServerID)
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}

	params.PGDataPath = server.PGDataPath + "/" + backupcontainername + "/" + getFormattedDate()

	//get the docker profile settings
	var setting admindb.Setting
	setting, err = admindb.GetSetting(dbConn, "S-DOCKER-PROFILE-CPU")
	params.CPU = setting.Value
	setting, err = admindb.GetSetting(dbConn, "S-DOCKER-PROFILE-MEM")
	params.MEM = setting.Value

	params.EnvVars = make(map[string]string)

	setting, err = admindb.GetSetting(dbConn, "DOMAIN-NAME")
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

	setting, err = admindb.GetSetting(dbConn, "PG-PORT")
	if err != nil {
		logit.Error.Println("Provision:" + err.Error())
		return err
	}
	params.EnvVars["BACKUP_PORT"] = setting.Value
	params.EnvVars["BACKUP_USER"] = "postgres"
	params.EnvVars["BACKUP_SERVER_URL"] = "cpm-backup" + "." + domain + ":" + "13000"

	//provision the volume
	request := &cpmserverapi.DiskProvisionRequest{"/tmp/foo"}
	request.Path = params.PGDataPath
	var url = "http://" + server.IPAddress + ":10001"
	_, err = cpmserverapi.DiskProvisionClient(url, request)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//run the container
	params.CommandPath = "docker-run-backup.sh"
	var response cpmserverapi.DockerRunResponse
	url = "http://" + server.IPAddress + ":10001"
	response, err = cpmserverapi.DockerRunClient(url, params)
	if err != nil {
		logit.Error.Println("Provision: " + response.Output)
		return err
	}
	logit.Info.Println("docker-run-backup.sh output=" + response.Output)

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
