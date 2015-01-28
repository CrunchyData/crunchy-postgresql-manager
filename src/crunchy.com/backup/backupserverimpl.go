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
	"crunchy.com/logutil"
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

//called by backup jobs as they execute
func (t *Command) AddStatus(status *BackupStatus, reply *Command) error {

	logutil.Log("AddStatus called")

	id, err := DBAddStatus(*status)
	if err != nil {
		logutil.Log("AddStatus error " + err.Error())
	}
	reply.Output = id
	return err
}

//called by backup jobs as they execute
func (t *Command) UpdateStatus(status *BackupStatus, reply *Command) error {

	logutil.Log("UpdateStatus called")

	err := DBUpdateStatus(*status)
	if err != nil {
		logutil.Log("UpdateStatus error " + err.Error())
	}
	return err
}

//called by admin do perform an adhoc backup job
func (t *Command) BackupNow(args *BackupRequest, reply *Command) error {

	logutil.Log("BackupNow.impl called")
	err := ProvisionBackupJob(args)
	if err != nil {
		logutil.Log("BackupNow.impl error:" + err.Error())
	}
	logutil.Log("BackupNow.impl completed")
	return err
}

/*
//called by admin to add to the schedule
func (t *Command) AddSchedule(args *BackupSchedule, reply *Command) error {

	logutil.Log("AddSchedule called")
	id, err := DBAddSchedule(*args)
	if err != nil {
		logutil.Log("AddSchedule error " + err.Error())
		return err
	}
	reply.Output = id
	return nil
}
*/

//called by admin to cause a reload of the cron jobs
func (t *Command) Reload(schedule *BackupSchedule, reply *Command) error {

	logutil.Log("Reload called")

	err := LoadSchedules()
	if err != nil {
		logutil.Log("Reload error " + err.Error())
	}

	return err
}

/*
//called by admin to update a schedule
func (t *Command) UpdateSchedule(schedule *BackupSchedule, reply *Command) error {

	logutil.Log("UpdateSchedule called")

	err := DBUpdateSchedule(*schedule)
	if err != nil {
		logutil.Log("UpdateSchedule error " + err.Error())
	}
	return err
}
*/

/*
//called by admin to delete a schedule
func (t *Command) DeleteSchedule(schedule *BackupSchedule, reply *Command) error {

	logutil.Log("DeleteSchedule called")

	err := DBDeleteSchedule(schedule.ID)
	if err != nil {
		logutil.Log("DeleteSchedule error " + err.Error())
	}
	return err
}
*/

func LoadSchedules() error {

	var err error
	logutil.Log("LoadSchedules called")

	schedules, err := DBGetSchedules()
	if err != nil {
		logutil.Log("LoadSchedules error " + err.Error())
	}

	if CRONInstance != nil {
		logutil.Log("stopping current cron instance...")
		CRONInstance.Stop()
	}

	//kill off the old cron, garbage collect it
	CRONInstance = nil

	//create a new cron
	logutil.Log("creating cron instance...")
	CRONInstance = cron.New()

	var cronexp string
	for i := 0; i < len(schedules); i++ {
		cronexp = getCron(schedules[i])
		logutil.Log("would have loaded schedule..." + cronexp)
		if schedules[i].Enabled == "YES" {
			logutil.Log("schedule " + schedules[i].ID + " was enabled so adding it")
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
			logutil.Log("schedule " + schedules[i].ID + " NOT enabled so dropping it")
		}

	}

	logutil.Log("starting new CRONInstance")
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
