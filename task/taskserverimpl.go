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
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"

	"database/sql"
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/robfig/cron"
	"log"
	"net/http"
)

type TaskProfile struct {
	ID   string
	Name string
}

type TaskStatus struct {
	Token         string
	ID            string
	ContainerName string
	StartTime     string
	TaskName      string
	ProfileName   string
	ScheduleID    string
	Path          string
	ElapsedTime   string
	TaskSize      string
	Status        string
	UpdateDt      string
}

type TaskRequest struct {
	ScheduleID    string
	ProfileName   string
	ContainerName string
}

type TaskSchedule struct {
	ID                string
	ContainerName     string
	ProfileName       string
	Name              string
	Enabled           string
	Minutes           string
	Hours             string
	DayOfMonth        string
	Month             string
	DayOfWeek         string
	UpdateDt          string
	RestoreSet        string
	RestoreRemotePath string
	RestoreRemoteHost string
	RestoreRemoteUser string
	RestoreDbUser     string
	RestoreDbPass     string
	Serverip          string
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

type ExecuteNowResponse struct {
	Output string
}

type ReloadResponse struct {
	Output string
}

//global cron instance that gets started, stopped, restarted
var CRONInstance *cron.Cron

const CLUSTERADMIN_DB = "clusteradmin"

// StatusAdd called by backup jobs as they execute to write new status info
func StatusAdd(w rest.ResponseWriter, r *rest.Request) {

	logit.Info.Println("StatusAdd called")

	request := TaskStatus{}
	err := r.DecodeJsonPayload(&request)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	defer dbConn.Close()

	var id string
	id, err = AddStatus(dbConn, &request)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := StatusAddResponse{}
	response.ID = id
	w.WriteJson(&response)
}

// StatusUpdate called by backup jobs as they execute
func StatusUpdate(w rest.ResponseWriter, r *rest.Request) {

	request := TaskStatus{}
	err := r.DecodeJsonPayload(&request)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dbConn *sql.DB
	dbConn, err = util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("StatusUpdate called")

	err = UpdateStatus(dbConn, &request)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := StatusUpdateResponse{}
	response.Output = "ok"
	w.WriteJson(&response)
}

// ExecuteNow called by admin do perform an adhoc task
func ExecuteNow(w rest.ResponseWriter, r *rest.Request) {
	request := TaskRequest{}
	err := r.DecodeJsonPayload(&request)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dbConn *sql.DB
	dbConn, err = util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer dbConn.Close()
	log.Println("-log- ExecuteNow.impl called profile=" + request.ProfileName)
	fmt.Println("-fmt- ExecuteNow.impl called profile=" + request.ProfileName)
	logit.Info.Println("-logit- ExecuteNow.impl called profile=" + request.ProfileName)

	if request.ProfileName == "pg_basebackup" {
		err = ProvisionBackupJob(dbConn, &request)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if request.ProfileName == "pg_backrest_restore" {
		logit.Info.Println("doing pg_backrest_restore job on...")
		err = ProvisionBackrestRestoreJob(dbConn, &request)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if request.ProfileName == "restore" {
		logit.Info.Println("doing restore job on...")
		err = ProvisionRestoreJob(dbConn, &request)
		if err != nil {
			logit.Error.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = errors.New("invalid profile name found:" + request.ProfileName)
		logit.Error.Println(err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logit.Info.Println("ExecuteNow.impl completed")

	response := ExecuteNowResponse{}
	response.Output = "ok"
	w.WriteJson(&response)
}

// Reload called by admin to cause a reload of the cron jobs
func Reload(w rest.ResponseWriter, r *rest.Request) {

	logit.Info.Println("Reload called")

	err := LoadSchedules()
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := ReloadResponse{}
	response.Output = "ok"
	w.WriteJson(&response)

}

// LoadSchedules loads the initial set of task schedules
func LoadSchedules() error {

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		return err

	}
	defer dbConn.Close()

	logit.Info.Println("LoadSchedules called")

	var schedules []TaskSchedule
	schedules, err = GetSchedules(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
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
			x.request = TaskRequest{}
			x.request.ScheduleID = schedules[i].ID
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

func getCron(s TaskSchedule) string {
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
