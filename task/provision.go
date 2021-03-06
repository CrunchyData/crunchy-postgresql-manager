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

package task

import (
	"database/sql"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/swarmapi"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"time"
)

// ProvisionBackupJob creates a Docker container that performs the backup on a db
func ProvisionBackupJob(dbConn *sql.DB, args *TaskRequest) error {

	//get node
	node, err := admindb.GetContainerByName(dbConn, args.ContainerName)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	var credential types.Credential
	credential, err = admindb.GetUserCredentials(dbConn, &node)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("backup.Provision called" +
		"with scheduleid=" + args.ScheduleID +
		"with containername=" + args.ContainerName +
		"with profilename=" + args.ProfileName)

	params := &swarmapi.DockerRunRequest{}
	params.Image = "cpm-backup-job"
	backupcontainername := args.ContainerName + "-backup"
	params.ContainerName = backupcontainername
	params.Standalone = "false"
	params.Profile = "SM"

	//remove the prior backup job if it exists
	removeParams := &swarmapi.DockerRemoveRequest{}
	removeParams.ContainerName = params.ContainerName
	var removeResponse swarmapi.DockerRemoveResponse
	removeResponse, err = swarmapi.DockerRemove(removeParams)
	logit.Info.Println("docker remove response " + removeResponse.Output)

	var pgdatapath types.Setting
	pgdatapath, err = admindb.GetSetting(dbConn, "PG-DATA-PATH")
	params.PGDataPath = pgdatapath.Value + "/" + backupcontainername + "/" + getFormattedDate()

	//get the docker profile settings
	var setting types.Setting
	setting, err = admindb.GetSetting(dbConn, "S-DOCKER-PROFILE-CPU")
	params.CPU = setting.Value
	setting, err = admindb.GetSetting(dbConn, "S-DOCKER-PROFILE-MEM")
	params.MEM = setting.Value

	params.EnvVars = make(map[string]string)

	params.EnvVars["BACKUP_NAME"] = backupcontainername
	params.EnvVars["BACKUP_SCHEDULEID"] = args.ScheduleID
	params.EnvVars["BACKUP_PROFILENAME"] = args.ProfileName
	params.EnvVars["BACKUP_CONTAINERNAME"] = args.ContainerName
	//params.EnvVars["BACKUP_PATH"] = params.PGDataPath
	params.EnvVars["BACKUP_PATH"] = "/" + backupcontainername + "/" + getFormattedDate()
	params.EnvVars["BACKUP_USERNAME"] = credential.Username
	params.EnvVars["BACKUP_PASSWORD"] = credential.Password
	params.EnvVars["BACKUP_HOST"] = args.ContainerName

	if node.Image == "cpm-node-proxy" {
		var proxy types.Proxy
		proxy, err = admindb.GetProxy(dbConn, args.ContainerName)
		if err != nil {
			logit.Error.Println(err.Error())
			return err
		}
		params.EnvVars["BACKUP_PROXY_HOST"] = proxy.Host
	}

	setting, err = admindb.GetSetting(dbConn, "PG-PORT")
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	params.EnvVars["BACKUP_PORT"] = setting.Value
	params.EnvVars["BACKUP_USER"] = "postgres"
	params.EnvVars["BACKUP_SERVER_URL"] = "cpm-task:13001"

	//provision the volume on all CPM servers
	var infoResponse swarmapi.DockerInfoResponse

	infoResponse, err = swarmapi.DockerInfo()
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	servers := make([]types.Server, len(infoResponse.Output))

	i := 0
	for i = range infoResponse.Output {
		servers[i].ID = infoResponse.Output[i]
		servers[i].Name = infoResponse.Output[i]
		servers[i].IPAddress = infoResponse.Output[i]
		i++
	}

	request := &cpmserverapi.DiskProvisionRequest{"/tmp/foo"}
	request.Path = params.PGDataPath
	for _, each := range servers {
		_, err = cpmserverapi.DiskProvisionClient(each.Name, request)
		if err != nil {
			logit.Error.Println(err.Error())
			return err
		}
	}

	//run the container
	//params.CommandPath = "docker-run-backup.sh"
	var response swarmapi.DockerRunResponse
	response, err = swarmapi.DockerRun(params)
	if err != nil {
		logit.Error.Println(response.ID)
		return err
	}
	logit.Info.Println("docker-run-backup.sh output=" + response.ID)

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
