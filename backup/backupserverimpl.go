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
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"

	"database/sql"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/robfig/cron"
	"net/http"
)

type BackupProfile struct {
	ID   string
	Name string
}

type BackupStatus struct {
	Token         string
	ID            string
	ContainerName string
	StartTime     string
	BackupName    string
	ProfileName   string
	ServerName    string
	ServerIP      string
	ScheduleID    string
	Path          string
	ElapsedTime   string
	BackupSize    string
	Status        string
	UpdateDt      string
}

type BackupRequest struct {
	ScheduleID    string
	ServerID      string
	ServerName    string
	ProfileName   string
	ServerIP      string
	ContainerName string
}

type BackupSchedule struct {
	ID            string
	ServerID      string
	ServerName    string
	ServerIP      string
	ContainerName string
	ProfileName   string
	Name          string
	Enabled       string
	Minutes       string
	Hours         string
	DayOfMonth    string
	Month         string
	DayOfWeek     string
	UpdateDt      string
}

type StatusAddResponse struct {
	ID string
}
type StatusUpdateResponse struct {
	Output string
}

type ReloadRequest struct {
	Name string
}

type BackupNowResponse struct {
	Output string
}

type ReloadResponse struct {
	Output string
}

//global cron instance that gets started, stopped, restarted
var CRONInstance *cron.Cron

const CLUSTERADMIN_DB = "clusteradmin"

//called by backup jobs as they execute
func StatusAdd(w rest.ResponseWriter, r *rest.Request) {

	logit.Info.Println("StatusAdd called")

	request := BackupStatus{}
	err := r.DecodeJsonPayload(&request)
	if err != nil {
		logit.Error.Println("StatusAdd: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("StatusAdd: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	defer dbConn.Close()

	var id string
	id, err = AddStatus(dbConn, &request)
	if err != nil {
		logit.Error.Println("AddStatus error " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := StatusAddResponse{}
	response.ID = id
	w.WriteJson(&response)
}

//called by backup jobs as they execute
func StatusUpdate(w rest.ResponseWriter, r *rest.Request) {

	request := BackupStatus{}
	err := r.DecodeJsonPayload(&request)
	if err != nil {
		logit.Error.Println("StatusUpdate: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dbConn *sql.DB
	dbConn, err = util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("StatusUpdate: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("StatusUpdate called")

	err = UpdateStatus(dbConn, &request)
	if err != nil {
		logit.Error.Println("UpdateStatus error " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := StatusUpdateResponse{}
	response.Output = "ok"
	w.WriteJson(&response)
}

//called by admin do perform an adhoc backup job
func BackupNow(w rest.ResponseWriter, r *rest.Request) {
	request := BackupRequest{}
	err := r.DecodeJsonPayload(&request)
	if err != nil {
		logit.Error.Println("BackupNow: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dbConn *sql.DB
	dbConn, err = util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer dbConn.Close()
	logit.Info.Println("BackupNow.impl called")

	err = ProvisionBackupJob(dbConn, &request)
	if err != nil {
		logit.Error.Println("BackupNow.impl error:" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logit.Info.Println("BackupNow.impl completed")

	response := BackupNowResponse{}
	response.Output = "ok"
	w.WriteJson(&response)
}

//called by admin to cause a reload of the cron jobs
func Reload(w rest.ResponseWriter, r *rest.Request) {

	logit.Info.Println("Reload called")

	err := LoadSchedules()
	if err != nil {
		logit.Error.Println("Reload error " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := ReloadResponse{}
	response.Output = "ok"
	w.WriteJson(&response)

}

func LoadSchedules() error {

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		return err

	}
	defer dbConn.Close()

	logit.Info.Println("LoadSchedules called")

	var schedules []BackupSchedule
	schedules, err = GetSchedules(dbConn)
	if err != nil {
		logit.Error.Println("LoadSchedules error " + err.Error())
		return err
	}

	if CRONInstance != nil {
		logit.Info.Println("stopping current cron instance...")
		CRONInstance.Stop()
	}

	//kill off the old cron, garbage collect it
	CRONInstance = nil

	//create a new cron
	logit.Info.Println("creating cron instance...")
	CRONInstance = cron.New()

	var cronexp string
	for i := 0; i < len(schedules); i++ {
		cronexp = getCron(schedules[i])
		logit.Info.Println("would have loaded schedule..." + cronexp)
		if schedules[i].Enabled == "YES" {
			logit.Info.Println("schedule " + schedules[i].ID + " was enabled so adding it")
			x := DefaultJob{}
			x.request = BackupRequest{}
			x.request.ScheduleID = schedules[i].ID
			x.request.ServerID = schedules[i].ServerID
			x.request.ServerName = schedules[i].ServerName
			x.request.ServerIP = schedules[i].ServerIP
			x.request.ContainerName = schedules[i].ContainerName
			x.request.ProfileName = schedules[i].ProfileName

			CRONInstance.AddJob(cronexp, x)
		} else {
			logit.Info.Println("schedule " + schedules[i].ID + " NOT enabled so dropping it")
		}

	}

	logit.Info.Println("starting new CRONInstance")
	CRONInstance.Start()

	return err
}

func getCron(s BackupSchedule) string {
	//leave seconds field with 0 as a default
	var cronexp = "0"
	cronexp = cronexp + " "

	cronexp = cronexp + s.Minutes
	cronexp = cronexp + " "

	cronexp = cronexp + s.Hours
	cronexp = cronexp + " "

	cronexp = cronexp + s.DayOfMonth
	cronexp = cronexp + " "

	cronexp = cronexp + s.Month
	cronexp = cronexp + " "

	cronexp = cronexp + s.DayOfWeek

	return cronexp
}
