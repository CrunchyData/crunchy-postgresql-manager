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
	"errors"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/task"
	"os"
	"time"
)

var startTime time.Time
var startTimeString string
var restoreRemotePath, restoreRemoteHost, restoreRemoteUser string
var restoreDbUser, restoreDbPass, restoreSet string
var restoreScheduleID, restoreStatusID string
var restoreContainerName, repoRemotePath string
var restoreServerName, restoreServerIP string

var StatusID = ""

func main() {

	startTime = time.Now()
	startTimeString = startTime.String()
	logit.Info.Println("backrestrestore running....")

	err := getEnvVars()
	if err != nil {
		logit.Error.Println(err.Error())
		return
	}

	s := task.TaskStatus{}
	eDuration := time.Since(startTime)
	s.StartTime = startTimeString
	s.ElapsedTime = eDuration.String()
	s.Status = "initializing"
	s.TaskSize = "n/a"
	sendStats(&s)

	logit.Info.Println("giving DNS time to register the backup job....sleeping for 7 secs")
	sleepTime, _ := time.ParseDuration("7s")
	time.Sleep(sleepTime)

	stats("restore job starting")

	stats("stopping postgres...")
	var stopResponse cpmcontainerapi.StopPGResponse
	stopResponse, err = cpmcontainerapi.StopPGClient(restoreContainerName)
	if err != nil {
		logit.Error.Println(err.Error())
		s.Status = "error in stopPG"
		sendStats(&s)
		os.Exit(1)
	}
	logit.Info.Println("StopPG....")
	logit.Info.Println(stopResponse.Output)
	logit.Info.Println(stopResponse.Status)
	logit.Info.Println("End of StopPG....")

	//wait for postgres to quit
	time.Sleep(sleepTime)

	stats("performing the restore...")
	//perform the restore
	restoreRequest := cpmcontainerapi.RestoreRequest{}
	restoreRequest.RestoreRemotePath = restoreRemotePath
	restoreRequest.RestoreRemoteHost = restoreRemoteHost
	restoreRequest.RestoreRemoteUser = restoreRemoteUser
	restoreRequest.RestoreDbUser = restoreDbUser
	restoreRequest.RestoreDbPass = restoreDbPass
	restoreRequest.RestoreSet = restoreSet

	var restoreResponse cpmcontainerapi.RestoreResponse
	restoreResponse, err = cpmcontainerapi.RestoreClient(restoreContainerName, &restoreRequest)
	if err != nil {
		logit.Error.Println(err.Error())
		s.Status = "error in restore"
		sendStats(&s)
		os.Exit(1)
	}
	logit.Info.Println("Restore....")
	logit.Info.Println(restoreResponse.Output)
	logit.Info.Println(restoreResponse.Status)
	logit.Info.Println("End of Restore....")

	stats("starting postgres after the restore...")

	var startResponse cpmcontainerapi.StartPGResponse
	startResponse, err = cpmcontainerapi.StartPGClient(restoreContainerName)
	if err != nil {
		logit.Error.Println(err.Error())
		s.Status = "error in startPG"
		sendStats(&s)
		os.Exit(1)
	}
	logit.Info.Println("StartPG....")
	logit.Info.Println(startResponse.Output)
	logit.Info.Println(startResponse.Status)
	logit.Info.Println("End of StartPG....")

	stats("seeding the database...")
	var seedResponse cpmcontainerapi.SeedResponse
	seedResponse, err = cpmcontainerapi.SeedClient(restoreContainerName)
	if err != nil {
		logit.Error.Println(err.Error())
		s.Status = "error in Seed"
		sendStats(&s)
		os.Exit(1)
	}
	logit.Info.Println("Seed....")
	logit.Info.Println(seedResponse.Output)
	logit.Info.Println(seedResponse.Status)
	logit.Info.Println("End of Seed....")

	//send final stats to backup
	finalstats("restore completed")

}

//report stats back to the cpm-admin for this job
func stats(str string) {

	sleepTime, _ := time.ParseDuration("7s")

	//logit.Info.Println("sending stats...")
	//logit.Info.Println("sleeping for 7 secs")
	time.Sleep(sleepTime)
	stats := task.TaskStatus{}
	eDuration := time.Since(startTime)
	stats.StartTime = startTimeString
	stats.ElapsedTime = eDuration.String()
	stats.Status = str
	stats.TaskSize = "n/a"
	sendStats(&stats)
}

//report stats back to the cpm-admin for this backup job
func finalstats(str string) {

	//connect to backupserver on cpm-admin
	//send stats to backup
	stats := task.TaskStatus{}
	eDuration := time.Since(startTime)
	stats.StartTime = startTimeString
	stats.ElapsedTime = eDuration.String()
	stats.Status = "completed"
	stats.TaskSize = "n/a"

	sendStats(&stats)
	logit.Info.Println("final stats here")
}

func sendStats(stats *task.TaskStatus) error {
	logit.Info.Println(stats.Status)

	stats.ContainerName = restoreContainerName
	stats.ScheduleID = restoreScheduleID
	stats.ProfileName = "pg_backrest_restore"
	stats.Path = restoreContainerName
	stats.TaskName = restoreContainerName
	stats.ID = StatusID

	var addResponse task.StatusAddResponse
	var err error

	if StatusID != "" {
		_, err = task.StatusUpdateClient(stats)
	} else {
		addResponse, err = task.StatusAddClient(stats)
		StatusID = addResponse.ID
	}
	if err != nil {
		logit.Error.Println("error in adding status:" + err.Error())
		return err
	}

	return nil
}

func getEnvVars() error {
	var err error

	logit.Info.Println("getEnvVars called\n")
	var found = true

	restoreRemotePath = os.Getenv("RestoreRemotePath")
	if restoreRemotePath == "" {
		logit.Error.Println("RestoreRemotePath env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreRemotePath=[" + restoreRemotePath + "]")

	restoreRemoteHost = os.Getenv("RestoreRemoteHost")
	if restoreRemoteHost == "" {
		logit.Error.Println("RestoreRemoteHost env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreRemoteHost=[" + restoreRemoteHost + "]")

	restoreRemoteUser = os.Getenv("RestoreRemoteUser")
	if restoreRemoteUser == "" {
		logit.Error.Println("RestoreRemoteUser env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreRemoteUser=[" + restoreRemoteUser + "]")

	restoreDbUser = os.Getenv("RestoreDbUser")
	if restoreDbUser == "" {
		logit.Error.Println("RestoreDbUser env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreDbUser=[" + restoreDbUser + "]")

	restoreDbPass = os.Getenv("RestoreDbPass")
	if restoreDbPass == "" {
		logit.Error.Println("RestoreDbPass env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreDbPass=[" + restoreDbPass + "]")

	restoreContainerName = os.Getenv("RestoreContainerName")
	if restoreContainerName == "" {
		logit.Error.Println("RestoreContainerName env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreContainerName=[" + restoreContainerName + "]")
	restoreScheduleID = os.Getenv("RestoreScheduleID")
	if restoreScheduleID == "" {
		logit.Error.Println("RestoreScheduleID env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreScheduleID=[" + restoreScheduleID + "]")
	restoreStatusID = os.Getenv("RestoreStatusID")
	if restoreStatusID == "" {
		logit.Error.Println("RestoreStatusID env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreStatusID=[" + restoreStatusID + "]")

	if !found {
		logit.Error.Println("restorecommand job missing required env vars")
		return errors.New("required env vars missing")
	}
	return err

}
