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
	"github.com/crunchydata/crunchy-postgresql-manager/task"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
)

type BackupNowPost struct {
	Token       string
	ServerID    string
	ProfileName string
	ScheduleID  string
}

type AddSchedulePost struct {
	ID                string
	Token             string
	ServerID          string
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
	RestoreDbUser     string
	RestoreDbPass     string
}

const CLUSTERADMIN_DB = "clusteradmin"

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

	if postMsg.ServerID == "" {
		logit.Error.Println("node ServerID required")
		rest.Error(w, "server ID required", 400)
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

	//get the server details for where the backup should be made
	server := types.Server{}
	server, err = admindb.GetServer(dbConn, postMsg.ServerID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	request := task.TaskRequest{}
	request.ScheduleID = postMsg.ScheduleID
	request.ServerID = server.ID
	request.ContainerName = schedule.ContainerName
	request.ServerName = server.Name
	request.ServerIP = server.IPAddress
	request.ProfileName = postMsg.ProfileName
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

	if postMsg.ServerID == "" {
		logit.Error.Println("node ServerID required")
		rest.Error(w, "server ID required", 400)
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

	err = secimpl.Authorize(dbConn, postMsg.Token, "perm-backup")
	if err != nil {
		logit.Error.Println("validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	s := task.TaskSchedule{}

	s.ServerID = postMsg.ServerID
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
	s.RestoreDbUser = postMsg.RestoreDbUser
	s.RestoreDbPass = postMsg.RestoreDbPass

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

	w.WriteJson(result)

}

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
	if postMsg.ServerID == "" {
		logit.Error.Println("ServerID required")
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
	s.ServerID = postMsg.ServerID
	s.Minutes = postMsg.Minutes
	s.Hours = postMsg.Hours
	s.Enabled = postMsg.Enabled
	s.DayOfMonth = postMsg.DayOfMonth
	s.Month = postMsg.Month
	s.DayOfWeek = postMsg.DayOfWeek
	s.Name = postMsg.Name

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
			n.ServerID = results[i].ServerID
			n.Role = results[i].Role
			n.Image = results[i].Image
			n.CreateDate = results[i].CreateDate
			n.ProjectName = results[i].ProjectName
			n.ServerName = results[i].ServerName
			n.Status = "UNKNOWN"
			nodes = append(nodes, n)
		}
		i++
	}

	w.WriteJson(&nodes)

}
