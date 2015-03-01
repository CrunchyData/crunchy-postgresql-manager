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
	"github.com/golang/glog"
	"github.com/robfig/cron"
)

type Command struct {
	Output string
}

type BackupProfile struct {
	ID   string
	Name string
}

type BackupStatus struct {
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

//global cron instance that gets started, stopped, restarted
var CRONInstance *cron.Cron

var CPMBIN = "/opt/cpm/bin/"

//called by backup jobs as they execute
func (t *Command) AddStatus(status *BackupStatus, reply *Command) error {

	glog.Infoln("AddStatus called")

	id, err := DBAddStatus(*status)
	if err != nil {
		glog.Errorln("AddStatus error " + err.Error())
	}
	reply.Output = id
	return err
}

//called by backup jobs as they execute
func (t *Command) UpdateStatus(status *BackupStatus, reply *Command) error {

	glog.Infoln("UpdateStatus called")

	err := DBUpdateStatus(*status)
	if err != nil {
		glog.Errorln("UpdateStatus error " + err.Error())
	}
	return err
}

//called by admin do perform an adhoc backup job
func (t *Command) BackupNow(args *BackupRequest, reply *Command) error {

	glog.Infoln("BackupNow.impl called")
	err := ProvisionBackupJob(args)
	if err != nil {
		glog.Errorln("BackupNow.impl error:" + err.Error())
	}
	glog.Infoln("BackupNow.impl completed")
	return err
}

//called by admin to cause a reload of the cron jobs
func (t *Command) Reload(schedule *BackupSchedule, reply *Command) error {

	glog.Infoln("Reload called")

	err := LoadSchedules()
	if err != nil {
		glog.Errorln("Reload error " + err.Error())
	}

	return err
}

func LoadSchedules() error {

	var err error
	glog.Infoln("LoadSchedules called")

	schedules, err := DBGetSchedules()
	if err != nil {
		glog.Errorln("LoadSchedules error " + err.Error())
	}

	if CRONInstance != nil {
		glog.Infoln("stopping current cron instance...")
		CRONInstance.Stop()
	}

	//kill off the old cron, garbage collect it
	CRONInstance = nil

	//create a new cron
	glog.Infoln("creating cron instance...")
	CRONInstance = cron.New()

	var cronexp string
	for i := 0; i < len(schedules); i++ {
		cronexp = getCron(schedules[i])
		glog.Infoln("would have loaded schedule..." + cronexp)
		if schedules[i].Enabled == "YES" {
			glog.Infoln("schedule " + schedules[i].ID + " was enabled so adding it")
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
			glog.Infoln("schedule " + schedules[i].ID + " NOT enabled so dropping it")
		}

	}

	glog.Infoln("starting new CRONInstance")
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
