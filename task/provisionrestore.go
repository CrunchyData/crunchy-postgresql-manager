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
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/swarmapi"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
)

// ProvisionRestoreJob creates a docker container to orchestrate a restore job
func ProvisionRestoreJob(dbConn *sql.DB, args *TaskRequest) error {

	logit.Info.Println("task.ProvisionRestoreJob called")
	logit.Info.Println("with scheduleid=" + args.ScheduleID)
	logit.Info.Println("with containername=" + args.ContainerName)
	logit.Info.Println("with profilename=" + args.ProfileName)
	logit.Info.Println("with statusid=" + args.StatusID)

	restorecontainername := args.ContainerName + "-restore-job"

	//remove any existing container with the same name
	inspectReq := &swarmapi.DockerInspectRequest{}
	inspectReq.ContainerName = restorecontainername
	inspectResponse, err := swarmapi.DockerInspect(inspectReq)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	if inspectResponse.RunningState != "not-found" {
		rreq := &swarmapi.DockerRemoveRequest{}
		rreq.ContainerName = restorecontainername
		_, err = swarmapi.DockerRemove(rreq)
		if err != nil {
			logit.Error.Println(err.Error())
			return err
		}
	}

	//create the new container
	params := &swarmapi.DockerRunRequest{}
	params.Image = "cpm-restore-job"
	params.ContainerName = restorecontainername
	params.Standalone = "false"
	params.Profile = "SM"

	schedule, err := GetSchedule(dbConn, args.ScheduleID)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("schedule serverip is " + schedule.Serverip)

	var taskstatus TaskStatus
	taskstatus, err = GetStatus(dbConn, args.StatusID)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	//params.PGDataPath = server.PGDataPath + "/" + restorecontainername + "/" + getFormattedDate()

	//get the docker profile settings
	var setting types.Setting
	setting, err = admindb.GetSetting(dbConn, "S-DOCKER-PROFILE-CPU")
	params.CPU = setting.Value
	setting, err = admindb.GetSetting(dbConn, "S-DOCKER-PROFILE-MEM")
	params.MEM = setting.Value
	setting, err = admindb.GetSetting(dbConn, "PG-DATA-PATH")
	var datapath string
	datapath = setting.Value

	//this gets mounted under /pgdata and allows us to access
	//both the backup files and the restored containers files
	params.PGDataPath = datapath

	params.EnvVars = make(map[string]string)

	params.EnvVars["RestoreServerip"] = schedule.Serverip
	params.EnvVars["RestoreBackupPath"] = taskstatus.Path
	params.EnvVars["RestorePath"] = args.ContainerName
	params.EnvVars["RestoreContainerName"] = args.ContainerName
	params.EnvVars["RestoreScheduleID"] = args.ScheduleID
	params.EnvVars["RestoreProfileName"] = args.ProfileName
	params.EnvVars["RestoreStatusID"] = args.StatusID

	//run the container
	//params.CommandPath = "docker-run-restore.sh"
	var response swarmapi.DockerRunResponse
	//var url = "http://" + server.IPAddress + ":10001"
	response, err = swarmapi.DockerRun(params)
	if err != nil {
		logit.Error.Println(response.ID)
		return err
	}
	logit.Info.Println("docker run output=" + response.ID)

	return nil
}
