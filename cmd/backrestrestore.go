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
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/task"
	"os"
	"os/exec"
	"strings"
	"time"
)

var startTime time.Time
var startTimeString string
var nodeName, repoRemotePath, backupHost, backupSet, backrestKeyPass string
var nodeProfileName, backupServerName, backupServerIP string

var scheduleID string
var backupPort string
var backupUser string
var backupAgentURL string
var StatusID = ""
var filename = "/tmp/backupjob.log"
var CPMBIN = "/var/cpm/bin/"
var file *os.File

func init() {
	startTime = time.Now()
	startTimeString = startTime.String()
	var err error
	file, err = os.Create(filename)
	if err != nil {
		logit.Error.Println(err.Error())
	}
}

func main() {

	logit.Info.Println("backrestrestore running....")

	getEnvVars()
	s := task.TaskStatus{}
	eDuration := time.Since(startTime)
	s.StartTime = startTimeString
	s.ElapsedTime = eDuration.String()
	s.Status = "initializing"
	s.TaskSize = du()
	sendStats(&s)

	//kick off stats reporting in a separate thread
	go stats("hi")

	logit.Info.Println("giving DNS time to register the backup job....sleeping for 7 secs")
	sleepTime, _ := time.ParseDuration("7s")
	time.Sleep(sleepTime)

	//perform the restore
	restore("end")

	//send final stats to backup
	finalstats("end")

}

//report stats back to the cpm-admin for this backup job
func stats(str string) {

	sleepTime, _ := time.ParseDuration("7s")

	for true {
		logit.Info.Println(file, "sending stats...")
		logit.Info.Println(file, "sleeping for 7 secs")
		time.Sleep(sleepTime)
		stats := task.TaskStatus{}
		eDuration := time.Since(startTime)
		stats.StartTime = startTimeString
		stats.ElapsedTime = eDuration.String()
		stats.Status = "running"
		stats.TaskSize = du()
		sendStats(&stats)
	}
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
	stats.TaskSize = du()

	sendStats(&stats)
	logit.Info.Println("final stats here")
}

//perform the restore
func restore(str string) {
	//do a backrest restore

	logit.Info.Println("doing restore")

	cmd := exec.Command("pg_backrest", "restore", "--stanza=main", backupSet)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
		logit.Error.Println("restore cmd stdout :" + out.String())
		logit.Error.Println("restores cmd stderr :" + stderr.String())
		return
	}

	logit.Info.Println("restore output was" + out.String())
	logit.Info.Println("restore is completed\n")
}

func sendStats(stats *task.TaskStatus) error {
	stats.ContainerName = nodeName
	stats.ServerName = backupServerName
	stats.ScheduleID = scheduleID
	stats.ServerIP = backupServerIP
	stats.ProfileName = nodeProfileName
	stats.Path = nodeName
	stats.TaskName = backupHost
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

	//send to backup
	logit.Info.Println("elapsed time:" + stats.ElapsedTime)
	logit.Info.Println("tasksize :" + stats.TaskSize)
	return nil
}

func getEnvVars() {
	logit.Info.Println("getEnvVars called\n")
	var found = true
	nodeName = os.Getenv("NODE_NAME")
	if nodeName == "" {
		logit.Error.Println("NODE_NAME env var not set\n")
		found = false
	}
	repoRemotePath = os.Getenv("REPO_REMOTE_PATH")
	if repoRemotePath == "" {
		logit.Error.Println("REPO_REMOTE_PATH env var not set\n")
	}
	backupHost = os.Getenv("BACKUP_HOST")
	if backupHost == "" {
		logit.Error.Println("BACKUP_HOST env var not set\n")
	}
	backupUser = os.Getenv("BACKUP_USER")
	if backupUser == "" {
		logit.Error.Println("BACKUP_USER env var not set\n")
		found = false
	}
	backupSet = os.Getenv("BACKUP_SET")
	if backupSet == "" {
		logit.Error.Println("BACKUP_SET env var not set\n")
		found = false
	}
	backupSet = "--set=" + backupSet

	backrestKeyPass = os.Getenv("BACKREST_KEY_PASS")
	if backrestKeyPass == "" {
		logit.Error.Println("BACKREST_KEY_PASS env var not set\n")
		found = false
	}

	backupServerName = os.Getenv("BACKUP_SERVERNAME")
	if backupServerName == "" {
		logit.Error.Println("BACKUP_SERVERNAME env var not set")
		found = false
	}

	backupServerIP = os.Getenv("BACKUP_SERVERIP")
	if backupServerIP == "" {
		logit.Error.Println("BACKUP_SERVERIP env var not set")
		found = false
	}
	nodeProfileName = os.Getenv("NODE_PROFILENAME")
	if nodeProfileName == "" {
		logit.Error.Println("NODE_PROFILENAME env var not set")
		found = false
	}

	if !found {
		panic("backrestrestore job missing required env vars")
	}

}

func du() string {
	cmd := exec.Command("du", "-hs", "/pgdata")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logit.Error.Println(err.Error())
		panic(err)
	}

	//only return the size from the du output
	var parsed = strings.Split(out.String(), "\t")
	return parsed[0]

}
