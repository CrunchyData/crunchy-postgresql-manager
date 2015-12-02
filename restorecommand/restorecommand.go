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
	"bytes"
	"errors"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/task"
	"os"
	"os/exec"
	"time"
)

var startTime time.Time
var startTimeString string

var restoreServerip string
var restoreBackupPath string
var restorePath string
var restoreContainerName string
var restoreScheduleID string
var restoreProfileName string
var restoreStatusID string

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
	//s.ElapsedTime = eDuration.String()
	s.ElapsedTime = fmt.Sprintf("%.3fs", eDuration.Seconds())
	s.Status = "initializing"
	s.TaskSize = "n/a"
	sendStats(&s)

	//	logit.Info.Println("giving DNS time to register the backup job....sleeping for 7 secs")
	sleepTime, _ := time.ParseDuration("7s")
	//	time.Sleep(sleepTime)

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
	//remove anything left in the /pgdata on the receiving container
	logit.Info.Println("removing any existing pgdata files")
	var frompath string
	frompath = "/pgdata/" + restorePath + "/*"
	logit.Info.Println("/bin/rm -rf " + frompath)
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/rm", "-rf", frompath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
		logit.Error.Println("rm stdout=" + out.String())
		logit.Error.Println("rm stderr=" + stderr.String())
		s.Status = "error in removing old files"
		sendStats(&s)
		os.Exit(1)
	}
	logit.Info.Println("remove was successful")

	//I'm choosing to do the remove here this way since this restore
	//might be a HUGE amount of data and the copy command could run a LONG
	//time, longer than an http timeout might allow for, since restorecommand
	//is running inside a container itself, this copy can run any amount of time

	//copy from the backup path all files to the /pgdata on the receiving container
	frompath = " /pgdata" + restoreBackupPath
	topath := "  /pgdata/" + restorePath
	logit.Info.Println("/var/cpm/bin/copyfiles.sh" + frompath + topath)
	cmd = exec.Command("/var/cpm/bin/copyfiles.sh", frompath, topath)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
		s.Status = "error in copying backup files"
		logit.Error.Println("cp stdout=" + out.String())
		logit.Error.Println("cp stderr=" + stderr.String())
		sendStats(&s)
		os.Exit(1)
	}

	logit.Info.Println("restore - copy of files was successful...")

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
	//stats.ElapsedTime = eDuration.String()
	stats.ElapsedTime = fmt.Sprintf("%.3fs", eDuration.Seconds())
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
	//stats.ElapsedTime = eDuration.String()
	stats.ElapsedTime = fmt.Sprintf("%.3fs", eDuration.Seconds())
	stats.Status = "completed"
	stats.TaskSize = "n/a"

	sendStats(&stats)
	logit.Info.Println("final stats here")
}

func sendStats(stats *task.TaskStatus) error {
	logit.Info.Println("restore - " + restoreBackupPath + " to " + restoreContainerName + " - " + stats.Status)

	stats.Status = "restore - " + restoreBackupPath + " to " + restoreContainerName + " - " + stats.Status
	stats.ContainerName = restoreContainerName
	stats.ScheduleID = restoreScheduleID
	stats.ProfileName = "restore"
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

	restoreServerip = os.Getenv("RestoreServerip")
	if restoreServerip == "" {
		logit.Error.Println("RestoreServerip env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreServerip=[" + restoreServerip + "]")

	restoreBackupPath = os.Getenv("RestoreBackupPath")
	if restoreBackupPath == "" {
		logit.Error.Println("RestoreBackupPath env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreBackupPath=[" + restoreBackupPath + "]")

	restorePath = os.Getenv("RestorePath")
	if restorePath == "" {
		logit.Error.Println("RestorePath env var not set\n")
		found = false
	}
	logit.Info.Println("RestorePath=[" + restorePath + "]")

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

	restoreProfileName = os.Getenv("RestoreProfileName")
	if restoreProfileName == "" {
		logit.Error.Println("RestoreProfileName env var not set\n")
		found = false
	}
	logit.Info.Println("RestoreProfileName=[" + restoreProfileName + "]")

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
