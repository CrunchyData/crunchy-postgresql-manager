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
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/swarmapi"
	"github.com/crunchydata/crunchy-postgresql-manager/task"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
)

type BackupNowPost struct {
	Token         string
	ServerID      string
	ProjectID     string
	ContainerName string
	ProfileName   string
	ScheduleID    string
	StatusID      string
}

type AddSchedulePost struct {
	ID                string
	Token             string
	Serverip          string
	ProfileName       string
	Name              string
	Enabled           string
	ContainerName     string
	Minutes           string
	Hours             string
	DayOfMonth        string
	Month             string
	DayOfWeek         string
	RestoreSet        string
	RestoreRemotePath string
	RestoreRemoteHost string
	RestoreRemoteUser string
	RestoreDbUser     string
	RestoreDbPass     string
}

const CLUSTERADMIN_DB = "clusteradmin"

// ExecuteNow executes a task schedule on demand allowing an immediate task execution
func ExecuteNow(w rest.ResponseWriter, r *rest.Request) {

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	postMsg := BackupNowPost{}
	err = r.DecodeJsonPayload(&postMsg)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, postMsg.Token, "perm-backup")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if postMsg.ProfileName == "" {
		logit.Error.Println("node ProfileName required")
		rest.Error(w, "ProfileName required", 400)
		return
	}

	if postMsg.ScheduleID == "" {
		logit.Error.Println("schedule ID required")
		rest.Error(w, "schedule ID required", 400)
		return
	}

	schedule, err2 := task.GetSchedule(dbConn, postMsg.ScheduleID)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		rest.Error(w, err2.Error(), 400)
		return
	}

	request := task.TaskRequest{}
	request.ScheduleID = postMsg.ScheduleID

	//in the case of a restore job, the user can supply a new containername
	if postMsg.ContainerName == "" {
		request.ContainerName = schedule.ContainerName
	} else {
		request.ContainerName = postMsg.ContainerName
	}

	//the restore command requires the task schedule statusID
	request.StatusID = postMsg.StatusID

	request.ProfileName = postMsg.ProfileName

	//for restore jobs, we go ahead and create the new
	//database container here, could possible move to the
	//restore job task later on
	if postMsg.ProfileName == "restore" {
		var newid string
		provisionParams := swarmapi.DockerRunRequest{}
		provisionParams.Profile = "SM"
		provisionParams.ProjectID = postMsg.ProjectID
		provisionParams.ContainerName = postMsg.ContainerName
		provisionParams.Image = "cpm-node"
		provisionParams.IPAddress = schedule.Serverip
		logit.Info.Println("before restore provision with...")
		logit.Info.Println("profile=" + provisionParams.Profile)
		logit.Info.Println("projectid=" + provisionParams.ProjectID)
		logit.Info.Println("containername=" + provisionParams.ContainerName)
		logit.Info.Println("image=" + provisionParams.Image)
		logit.Info.Println("ipaddress=" + provisionParams.IPAddress)
		newid, err = provisionImpl(dbConn, &provisionParams, false)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logit.Info.Printf("created node for restore job id = " + newid)
	}

	output, err := task.ExecuteNowClient(&request)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logit.Info.Println("output=" + output.Output)

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

// AddSchedule creates a new task schedule
func AddSchedule(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	postMsg := AddSchedulePost{}
	err = r.DecodeJsonPayload(&postMsg)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if postMsg.ContainerName == "" {
		logit.Error.Println("node ContainerName required")
		rest.Error(w, "ContainerName required", 400)
		return
	}
	if postMsg.ProfileName == "" {
		logit.Error.Println("node ProfileName required")
		rest.Error(w, "ProfileName required", 400)
		return
	}
	if postMsg.Name == "" {
		logit.Error.Println("node Name required")
		rest.Error(w, "Name required", 400)
		return
	}

	logit.Info.Println("in adminapi.backupmgmt.AddSchedule got serverIP of " + postMsg.Serverip)

	err = secimpl.Authorize(dbConn, postMsg.Token, "perm-backup")
	if err != nil {
		logit.Error.Println("validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	s := task.TaskSchedule{}

	s.ContainerName = postMsg.ContainerName
	s.ProfileName = postMsg.ProfileName
	s.Name = postMsg.Name

	//defaults for any new schedule
	s.Enabled = "NO"
	s.Minutes = "00"
	s.Hours = "11"
	s.DayOfMonth = "1"
	s.Month = "*"
	s.DayOfWeek = "*"
	s.RestoreSet = postMsg.RestoreSet
	s.RestoreRemotePath = postMsg.RestoreRemotePath
	s.RestoreRemoteHost = postMsg.RestoreRemoteHost
	s.RestoreRemoteUser = postMsg.RestoreRemoteUser
	s.RestoreDbUser = postMsg.RestoreDbUser
	s.RestoreDbPass = postMsg.RestoreDbPass
	s.Serverip = postMsg.Serverip

	result, err := task.AddSchedule(dbConn, s)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	logit.Info.Println("AddSchedule: new ID " + result)

	//we choose by design to not notify the task server
	//on schedule adds, instead we mark any new schedule
	//as DISABLED, forcing the user to change the defaults
	//and use the UpdateSchedule which does force a notify
	//to the task server

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

// DeleteSchedule deletes an existing task schedule
func DeleteSchedule(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-backup")
	if err != nil {
		logit.Error.Println("validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")

	if ID == "" {
		rest.Error(w, "schedule ID required", 400)
		return
	}

	err = task.DeleteSchedule(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//notify task server to reload schedules

	var output task.ReloadResponse
	output, err = task.ReloadClient()
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	logit.Info.Println("reload output=" + output.Output)

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

// GetSchedule returns a task schedule
func GetSchedule(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")

	if ID == "" {
		rest.Error(w, "ID required", 400)
		return
	}

	result, err := task.GetSchedule(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	logit.Info.Println("GetSchedule api returns serverip of " + result.Serverip)

	w.WriteJson(result)

}

// GetAllSchedules returns a list of task schedules for a given container
func GetAllSchedules(w rest.ResponseWriter, r *rest.Request) {
	Token := r.PathParam("Token")
	if Token == "" {
		rest.Error(w, "Token required", 400)
		return
	}
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, Token, "perm-read")
	if err != nil {
		logit.Error.Println("validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ContainerName := r.PathParam("ContainerName")

	if ContainerName == "" {
		rest.Error(w, "ContainerName required", 400)
		return
	}

	schedules, err := task.GetAllSchedules(dbConn, ContainerName)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	w.WriteJson(schedules)

}

// GetStatus returns the status of a given task schedule
func GetStatus(w rest.ResponseWriter, r *rest.Request) {
	Token := r.PathParam("Token")
	if Token == "" {
		rest.Error(w, "Token required", 400)
		return
	}
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, Token, "perm-read")
	if err != nil {
		logit.Error.Println("validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", 400)
		return
	}
	stat, err := task.GetStatus(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(stat)
}

// GetAllStatus returns the full set of status results for a given task schedule
func GetAllStatus(w rest.ResponseWriter, r *rest.Request) {
	Token := r.PathParam("Token")
	if Token == "" {
		rest.Error(w, "Token required", 400)
		return
	}
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}

	err = secimpl.Authorize(dbConn, Token, "perm-read")
	if err != nil {
		logit.Error.Println("validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "Schedule ID required", 400)
		return
	}

	stats, err := task.GetAllStatus(dbConn, ID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(stats)

}

// UpdateSchedule updates a given task schedule
func UpdateSchedule(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	postMsg := AddSchedulePost{}
	err = r.DecodeJsonPayload(&postMsg)
	if err != nil {
		logit.Error.Println("decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, postMsg.Token, "perm-backup")
	if err != nil {
		logit.Error.Println("validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if postMsg.ID == "" {
		logit.Error.Println("schedule ID required")
		rest.Error(w, "schedule ID required", 400)
		return
	}
	if postMsg.Enabled == "" {
		logit.Error.Println("Enabled required")
		rest.Error(w, "Enabled required", 400)
		return
	}
	if postMsg.Minutes == "" {
		logit.Error.Println("Minutes required")
		rest.Error(w, "schedule Minutes required", 400)
		return
	}
	if postMsg.Hours == "" {
		logit.Error.Println("Hours required")
		rest.Error(w, "schedule Hours required", 400)
		return
	}
	if postMsg.DayOfMonth == "" {
		logit.Error.Println("DayOfMonth required")
		rest.Error(w, "schedule DayOfMonth required", 400)
		return
	}
	if postMsg.Month == "" {
		logit.Error.Println("Month required")
		rest.Error(w, "schedule Month required", 400)
		return
	}
	if postMsg.DayOfWeek == "" {
		logit.Error.Println("DayOfWeek required")
		rest.Error(w, "schedule DayOfWeek required", 400)
		return
	}
	if postMsg.Name == "" {
		logit.Error.Println("Name required")
		rest.Error(w, "schedule Name required", 400)
		return
	}

	s := task.TaskSchedule{}
	s.ID = postMsg.ID
	s.Minutes = postMsg.Minutes
	s.Hours = postMsg.Hours
	s.Enabled = postMsg.Enabled
	s.DayOfMonth = postMsg.DayOfMonth
	s.Month = postMsg.Month
	s.DayOfWeek = postMsg.DayOfWeek
	s.Name = postMsg.Name
	s.RestoreSet = postMsg.RestoreSet
	s.RestoreRemotePath = postMsg.RestoreRemotePath
	s.RestoreRemoteHost = postMsg.RestoreRemoteHost
	s.RestoreRemoteUser = postMsg.RestoreRemoteUser
	s.RestoreDbUser = postMsg.RestoreDbUser
	s.RestoreDbPass = postMsg.RestoreDbPass
	s.Serverip = postMsg.Serverip

	err = task.UpdateSchedule(dbConn, s)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//notify task server to reload it's schedules

	output, err := task.ReloadClient()
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	logit.Info.Println("reload output=" + output.Output)

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

// TODO
func GetBackupNodes(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err2 := admindb.GetAllContainers(dbConn)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		rest.Error(w, err2.Error(), 400)
	}
	i := 0
	//nodes := make([]ClusterNode, found)
	nodes := []types.ClusterNode{}
	for i = range results {
		if results[i].Role == "unassigned" ||
			results[i].Role == "master" {
			n := types.ClusterNode{}
			n.ID = results[i].ID
			n.Name = results[i].Name
			n.ClusterID = results[i].ClusterID
			n.Role = results[i].Role
			n.Image = results[i].Image
			n.CreateDate = results[i].CreateDate
			n.ProjectName = results[i].ProjectName
			n.Status = "UNKNOWN"
			nodes = append(nodes, n)
		}
		i++
	}

	w.WriteJson(&nodes)

}
