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
	"crunchy.com/backup"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
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

func BackupNow(w rest.ResponseWriter, r *rest.Request) {
	postMsg := BackupNowPost{}
	err := r.DecodeJsonPayload(&postMsg)
	if err != nil {
		glog.Errorln("BackupNow: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(postMsg.Token, "perm-backup")
	if err != nil {
		glog.Errorln("BackupNow: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	if postMsg.ServerID == "" {
		glog.Errorln("BackupNow: error node ServerID required")
		rest.Error(w, "server ID required", 400)
		return
	}
	if postMsg.ProfileName == "" {
		glog.Errorln("BackupNow: error node ProfileName required")
		rest.Error(w, "ProfileName required", 400)
		return
	}

	if postMsg.ScheduleID == "" {
		glog.Errorln("BackupNow: error schedule ID required")
		rest.Error(w, "schedule ID required", 400)
		return
	}

	schedule, err := backup.DBGetSchedule(postMsg.ScheduleID)
	if err != nil {
		glog.Errorln("BackupNow: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//get the server details for where the backup should be made
	server := admindb.DBServer{}
	server, err = admindb.GetDBServer(postMsg.ServerID)
	if err != nil {
		glog.Errorln("BackupNow: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//get the domain name
	//get domain name
	var domainname admindb.DBSetting
	domainname, err = admindb.GetDBSetting("DOMAIN-NAME")
	if err != nil {
		glog.Errorln("BackupNow: DOMAIN-NAME err " + err.Error())
	}

	request := backup.BackupRequest{}
	request.ScheduleID = postMsg.ScheduleID
	request.ServerID = server.ID
	request.ContainerName = schedule.ContainerName
	request.ServerName = server.Name
	request.ServerIP = server.IPAddress
	request.ProfileName = postMsg.ProfileName
	backupServerURL := "cpm-backup." + domainname.Value + ":13010"
	output, err := backup.BackupNowClient(backupServerURL, request)
	if err != nil {
		glog.Errorln(err.Error())
	}
	glog.Infoln("output=" + output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AddSchedule(w rest.ResponseWriter, r *rest.Request) {
	postMsg := AddSchedulePost{}
	err := r.DecodeJsonPayload(&postMsg)
	if err != nil {
		glog.Errorln("AddSchedule: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(postMsg.Token, "perm-backup")
	if err != nil {
		glog.Errorln("AddSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	if postMsg.ServerID == "" {
		glog.Errorln("AddSchedule: error node ServerID required")
		rest.Error(w, "server ID required", 400)
		return
	}

	if postMsg.ContainerName == "" {
		glog.Errorln("AddSchedule: error node ContainerName required")
		rest.Error(w, "ContainerName required", 400)
		return
	}
	if postMsg.ProfileName == "" {
		glog.Errorln("AddSchedule: error node ProfileName required")
		rest.Error(w, "ProfileName required", 400)
		return
	}
	if postMsg.Name == "" {
		glog.Errorln("AddSchedule: error node Name required")
		rest.Error(w, "Name required", 400)
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

	result, err := backup.DBAddSchedule(s)
	if err != nil {
		glog.Errorln("GetNode: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	glog.Infoln("AddSchedule: new ID " + result)

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
	err := secimpl.Authorize(r.PathParam("Token"), "perm-backup")
	if err != nil {
		glog.Errorln("DeleteSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	ID := r.PathParam("ID")

	if ID == "" {
		rest.Error(w, "schedule ID required", 400)
		return
	}

	err = backup.DBDeleteSchedule(ID)
	if err != nil {
		glog.Errorln("DeleteSchedule: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//notify backup server to reload schedules

	//get the domain name
	//get domain name
	var domainname admindb.DBSetting
	domainname, err = admindb.GetDBSetting("DOMAIN-NAME")
	if err != nil {
		glog.Errorln("DeleteSchedule: DOMAIN-NAME err " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}

	s := backup.BackupSchedule{}

	backupServerURL := "cpm-backup." + domainname.Value + ":13010"
	output, err := backup.ReloadClient(backupServerURL, s)
	if err != nil {
		glog.Errorln(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	glog.Infoln("reload output=" + output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func GetSchedule(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	ID := r.PathParam("ID")

	if ID == "" {
		rest.Error(w, "ID required", 400)
		return
	}

	result, err := backup.DBGetSchedule(ID)
	if err != nil {
		glog.Errorln("GetNode: " + err.Error())
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
	err := secimpl.Authorize(Token, "perm-read")
	if err != nil {
		glog.Errorln("GetAllSchedules: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	ContainerName := r.PathParam("ContainerName")

	if ContainerName == "" {
		rest.Error(w, "ContainerName required", 400)
		return
	}

	schedules, err := backup.DBGetAllSchedules(ContainerName)
	if err != nil {
		glog.Errorln("GetAllSchedules: " + err.Error())
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
	err := secimpl.Authorize(Token, "perm-read")
	if err != nil {
		glog.Errorln("GetStatus: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", 400)
		return
	}
	stat, err := backup.DBGetStatus(ID)
	if err != nil {
		glog.Errorln("GetStatus: " + err.Error())
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
	err := secimpl.Authorize(Token, "perm-read")
	if err != nil {
		glog.Errorln("GetAllStatus: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "Schedule ID required", 400)
		return
	}

	stats, err := backup.DBGetAllStatus(ID)
	if err != nil {
		glog.Errorln("GetAllStatus: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(stats)

}

func UpdateSchedule(w rest.ResponseWriter, r *rest.Request) {
	postMsg := AddSchedulePost{}
	err := r.DecodeJsonPayload(&postMsg)
	if err != nil {
		glog.Errorln("UpdateSchedule: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(postMsg.Token, "perm-backup")
	if err != nil {
		glog.Errorln("UpdateSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	if postMsg.ID == "" {
		glog.Errorln("UpdateSchedule: error schedule ID required")
		rest.Error(w, "schedule ID required", 400)
		return
	}
	if postMsg.ServerID == "" {
		glog.Errorln("UpdateSchedule: error ServerID required")
		rest.Error(w, "schedule ID required", 400)
		return
	}
	if postMsg.Enabled == "" {
		glog.Errorln("UpdateSchedule: error Enabled required")
		rest.Error(w, "Enabled required", 400)
		return
	}
	if postMsg.Minutes == "" {
		glog.Errorln("UpdateSchedule: error Minutes required")
		rest.Error(w, "schedule Minutes required", 400)
		return
	}
	if postMsg.Hours == "" {
		glog.Errorln("UpdateSchedule: error Hours required")
		rest.Error(w, "schedule Hours required", 400)
		return
	}
	if postMsg.DayOfMonth == "" {
		glog.Errorln("UpdateSchedule: error DayOfMonth required")
		rest.Error(w, "schedule DayOfMonth required", 400)
		return
	}
	if postMsg.Month == "" {
		glog.Errorln("UpdateSchedule: error Month required")
		rest.Error(w, "schedule Month required", 400)
		return
	}
	if postMsg.DayOfWeek == "" {
		glog.Errorln("UpdateSchedule: error DayOfWeek required")
		rest.Error(w, "schedule DayOfWeek required", 400)
		return
	}
	if postMsg.Name == "" {
		glog.Errorln("UpdateSchedule: error Name required")
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

	err = backup.DBUpdateSchedule(s)
	if err != nil {
		glog.Errorln(err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//notify backup server to reload it's schedules

	//get the domain name
	//get domain name
	var domainname admindb.DBSetting
	domainname, err = admindb.GetDBSetting("DOMAIN-NAME")
	if err != nil {
		glog.Errorln("BackupNow: DOMAIN-NAME err " + err.Error())
	}
	backupServerURL := "cpm-backup." + domainname.Value + ":13010"
	output, err := backup.ReloadClient(backupServerURL, s)
	if err != nil {
		glog.Errorln(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	glog.Infoln("reload output=" + output)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetBackupNodes(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetAllNodes: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	results, err := admindb.GetAllDBNodes()
	if err != nil {
		glog.Errorln("GetAllNodes: " + err.Error())
		rest.Error(w, err.Error(), 400)
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
			n.Status = "UNKNOWN"
			nodes = append(nodes, n)
		}
		i++
	}

	w.WriteJson(&nodes)

}
