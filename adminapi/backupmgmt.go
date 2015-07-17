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
	"github.com/crunchydata/crunchy-postgresql-manager/backup"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
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
	ID            string
	Token         string
	ServerID      string
	ProfileName   string
	Name          string
	Enabled       string
	ContainerName string
	Minutes       string
	Hours         string
	DayOfMonth    string
	Month         string
	DayOfWeek     string
}

const CLUSTERADMIN_DB = "clusteradmin"

func BackupNow(w rest.ResponseWriter, r *rest.Request) {

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	postMsg := BackupNowPost{}
	err = r.DecodeJsonPayload(&postMsg)
	if err != nil {
		logit.Error.Println("BackupNow: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, postMsg.Token, "perm-backup")
	if err != nil {
		logit.Error.Println("BackupNow: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if postMsg.ServerID == "" {
		logit.Error.Println("BackupNow: error node ServerID required")
		rest.Error(w, "server ID required", 400)
		return
	}
	if postMsg.ProfileName == "" {
		logit.Error.Println("BackupNow: error node ProfileName required")
		rest.Error(w, "ProfileName required", 400)
		return
	}

	if postMsg.ScheduleID == "" {
		logit.Error.Println("BackupNow: error schedule ID required")
		rest.Error(w, "schedule ID required", 400)
		return
	}

	schedule, err2 := backup.GetSchedule(dbConn, postMsg.ScheduleID)
	if err2 != nil {
		logit.Error.Println("BackupNow: " + err2.Error())
		rest.Error(w, err2.Error(), 400)
		return
	}

	//get the server details for where the backup should be made
	server := admindb.Server{}
	server, err = admindb.GetServer(dbConn, postMsg.ServerID)
	if err != nil {
		logit.Error.Println("BackupNow: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//get the domain name
	//get domain name
	var domainname admindb.Setting
	domainname, err = admindb.GetSetting(dbConn, "DOMAIN-NAME")
	if err != nil {
		logit.Error.Println("BackupNow: DOMAIN-NAME err " + err.Error())
	}

	request := backup.BackupRequest{}
	request.ScheduleID = postMsg.ScheduleID
	request.ServerID = server.ID
	if KubeEnv {
		request.ContainerName = schedule.ContainerName + "-db"
	} else {
		request.ContainerName = schedule.ContainerName
	}
	request.ServerName = server.Name
	request.ServerIP = server.IPAddress
	request.ProfileName = postMsg.ProfileName
	backupServerURL := "cpm-backup." + domainname.Value + ":13000"
	output, err := backup.BackupNowClient(backupServerURL, request)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logit.Info.Println("output=" + output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AddSchedule(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	postMsg := AddSchedulePost{}
	err = r.DecodeJsonPayload(&postMsg)
	if err != nil {
		logit.Error.Println("AddSchedule: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if postMsg.ServerID == "" {
		logit.Error.Println("AddSchedule: error node ServerID required")
		rest.Error(w, "server ID required", 400)
		return
	}

	if postMsg.ContainerName == "" {
		logit.Error.Println("AddSchedule: error node ContainerName required")
		rest.Error(w, "ContainerName required", 400)
		return
	}
	if postMsg.ProfileName == "" {
		logit.Error.Println("AddSchedule: error node ProfileName required")
		rest.Error(w, "ProfileName required", 400)
		return
	}
	if postMsg.Name == "" {
		logit.Error.Println("AddSchedule: error node Name required")
		rest.Error(w, "Name required", 400)
		return
	}

	err = secimpl.Authorize(dbConn, postMsg.Token, "perm-backup")
	if err != nil {
		logit.Error.Println("AddSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	s := backup.BackupSchedule{}

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

	result, err := backup.AddSchedule(dbConn, s)
	if err != nil {
		logit.Error.Println("GetNode: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	logit.Info.Println("AddSchedule: new ID " + result)

	//we choose by design to not notify the backup server
	//on schedule adds, instead we mark any new schedule
	//as DISABLED, forcing the user to change the defaults
	//and use the UpdateSchedule which does force a notify
	//to the backup server

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func DeleteSchedule(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-backup")
	if err != nil {
		logit.Error.Println("DeleteSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")

	if ID == "" {
		rest.Error(w, "schedule ID required", 400)
		return
	}

	err = backup.DeleteSchedule(dbConn, ID)
	if err != nil {
		logit.Error.Println("DeleteSchedule: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//notify backup server to reload schedules

	//get the domain name
	//get domain name
	var domainname admindb.Setting
	domainname, err = admindb.GetSetting(dbConn, "DOMAIN-NAME")
	if err != nil {
		logit.Error.Println("DeleteSchedule: DOMAIN-NAME err " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}

	s := backup.BackupSchedule{}

	backupServerURL := "cpm-backup." + domainname.Value + ":13000"
	var output string
	output, err = backup.ReloadClient(backupServerURL, s)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	logit.Info.Println("reload output=" + output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func GetSchedule(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")

	if ID == "" {
		rest.Error(w, "ID required", 400)
		return
	}

	result, err := backup.GetSchedule(dbConn, ID)
	if err != nil {
		logit.Error.Println("GetNode: " + err.Error())
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
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, Token, "perm-read")
	if err != nil {
		logit.Error.Println("GetAllSchedules: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ContainerName := r.PathParam("ContainerName")

	if ContainerName == "" {
		rest.Error(w, "ContainerName required", 400)
		return
	}

	schedules, err := backup.GetAllSchedules(dbConn, ContainerName)
	if err != nil {
		logit.Error.Println("GetAllSchedules: " + err.Error())
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
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	err = secimpl.Authorize(dbConn, Token, "perm-read")
	if err != nil {
		logit.Error.Println("GetStatus: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", 400)
		return
	}
	stat, err := backup.GetStatus(dbConn, ID)
	if err != nil {
		logit.Error.Println("GetStatus: " + err.Error())
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
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}

	err = secimpl.Authorize(dbConn, Token, "perm-read")
	if err != nil {
		logit.Error.Println("GetAllStatus: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "Schedule ID required", 400)
		return
	}

	stats, err := backup.GetAllStatus(dbConn, ID)
	if err != nil {
		logit.Error.Println("GetAllStatus: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(stats)

}

func UpdateSchedule(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	postMsg := AddSchedulePost{}
	err = r.DecodeJsonPayload(&postMsg)
	if err != nil {
		logit.Error.Println("UpdateSchedule: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, postMsg.Token, "perm-backup")
	if err != nil {
		logit.Error.Println("UpdateSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if postMsg.ID == "" {
		logit.Error.Println("UpdateSchedule: error schedule ID required")
		rest.Error(w, "schedule ID required", 400)
		return
	}
	if postMsg.ServerID == "" {
		logit.Error.Println("UpdateSchedule: error ServerID required")
		rest.Error(w, "schedule ID required", 400)
		return
	}
	if postMsg.Enabled == "" {
		logit.Error.Println("UpdateSchedule: error Enabled required")
		rest.Error(w, "Enabled required", 400)
		return
	}
	if postMsg.Minutes == "" {
		logit.Error.Println("UpdateSchedule: error Minutes required")
		rest.Error(w, "schedule Minutes required", 400)
		return
	}
	if postMsg.Hours == "" {
		logit.Error.Println("UpdateSchedule: error Hours required")
		rest.Error(w, "schedule Hours required", 400)
		return
	}
	if postMsg.DayOfMonth == "" {
		logit.Error.Println("UpdateSchedule: error DayOfMonth required")
		rest.Error(w, "schedule DayOfMonth required", 400)
		return
	}
	if postMsg.Month == "" {
		logit.Error.Println("UpdateSchedule: error Month required")
		rest.Error(w, "schedule Month required", 400)
		return
	}
	if postMsg.DayOfWeek == "" {
		logit.Error.Println("UpdateSchedule: error DayOfWeek required")
		rest.Error(w, "schedule DayOfWeek required", 400)
		return
	}
	if postMsg.Name == "" {
		logit.Error.Println("UpdateSchedule: error Name required")
		rest.Error(w, "schedule Name required", 400)
		return
	}

	s := backup.BackupSchedule{}
	s.ID = postMsg.ID
	s.ServerID = postMsg.ServerID
	s.Minutes = postMsg.Minutes
	s.Hours = postMsg.Hours
	s.Enabled = postMsg.Enabled
	s.DayOfMonth = postMsg.DayOfMonth
	s.Month = postMsg.Month
	s.DayOfWeek = postMsg.DayOfWeek
	s.Name = postMsg.Name

	err = backup.UpdateSchedule(dbConn, s)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//notify backup server to reload it's schedules

	//get the domain name
	//get domain name
	var domainname admindb.Setting
	domainname, err = admindb.GetSetting(dbConn, "DOMAIN-NAME")
	if err != nil {
		logit.Error.Println("BackupNow: DOMAIN-NAME err " + err.Error())
	}
	backupServerURL := "cpm-backup." + domainname.Value + ":13000"
	output, err := backup.ReloadClient(backupServerURL, s)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	logit.Info.Println("reload output=" + output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetBackupNodes(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllNodes: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err2 := admindb.GetAllContainers(dbConn)
	if err2 != nil {
		logit.Error.Println("GetAllNodes: " + err2.Error())
		rest.Error(w, err2.Error(), 400)
	}
	i := 0
	//nodes := make([]ClusterNode, found)
	nodes := []ClusterNode{}
	for i = range results {
		if results[i].Role == "unassigned" ||
			results[i].Role == "master" {
			n := ClusterNode{}
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
