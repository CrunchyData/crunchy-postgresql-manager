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
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var startTime time.Time
var startTimeString string
var backupContainerName string
var backupServerName string
var backupServerIP string
var backupProfileName string
var backupPath string
var backupProxyIP string
var backupHost string
var backupUsername string
var backupPassword string
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

	_, err := io.WriteString(file, "backupcommand running....\n")
	if err != nil {
		log.Println(err.Error())
	}

	defer closeLog()

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

	io.WriteString(file, "giving DNS time to register the backup job....sleeping for 7 secs")
	sleepTime, _ := time.ParseDuration("7s")
	time.Sleep(sleepTime)

	//perform the backup
	backupfunc("end")

	//send final stats to backup
	finalstats("end")

}

func closeLog() {
	//copy the output log into the backup directory
	file.Close()
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("mv", filename, "/pgdata/")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
		io.WriteString(file, err.Error())
	}
}

//report stats back to the cpm-admin for this backup job
func stats(str string) {

	sleepTime, _ := time.ParseDuration("7s")

	for true {
		io.WriteString(file, "sending stats...\n")
		io.WriteString(file, "sleeping for 7 secs\n")
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
	io.WriteString(file, "final stats here\n")
}

//perform the backup
func backupfunc(str string) {
	//do a pg_basebackup here

	io.WriteString(file, "doing backup on "+backupHost+"\n")

	//create base backup from master
	if backupProxyIP != "" {
		backupHost = backupProxyIP
		io.WriteString(file, "doing proxy backup to "+backupHost+"\n")
	}

	cmd := exec.Command(CPMBIN+"basebackup.sh", backupHost, backupUsername, backupPassword)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		io.WriteString(file, "backupfunc:"+err.Error()+"\n")
		io.WriteString(file, "backupfunc cmd stdout :"+out.String()+"\n")
		io.WriteString(file, "backupfunc cmd stderr :"+stderr.String()+"\n")
		return
	}

	io.WriteString(file, "basebackup output was"+out.String()+"\n")
	io.WriteString(file, " backups is completed\n")
}

func sendStats(stats *task.TaskStatus) error {
	stats.ContainerName = backupContainerName
	stats.ServerName = backupServerName
	stats.ScheduleID = scheduleID
	stats.ServerIP = backupServerIP
	stats.ProfileName = backupProfileName
	stats.Path = backupPath
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
		io.WriteString(file, "error in adding status:"+err.Error()+"\n")
		return err
	}

	//send to backup
	io.WriteString(file, "elapsed time:"+stats.ElapsedTime+"\n")
	io.WriteString(file, "tasksize :"+stats.TaskSize+"\n")
	return nil
}

func getEnvVars() {
	io.WriteString(file, "getEnvVars called\n")
	var found = true
	backupContainerName = os.Getenv("BACKUP_CONTAINERNAME")
	if backupContainerName == "" {
		io.WriteString(file, "BACKUP_CONTAINERNAME env var not set\n")
		found = false
	}
	backupUsername = os.Getenv("BACKUP_USERNAME")
	if backupUsername == "" {
		io.WriteString(file, "BACKUP_USERNAME env var not set\n")
	}
	backupPassword = os.Getenv("BACKUP_PASSWORD")
	if backupPassword == "" {
		io.WriteString(file, "BACKUP_PASSWORD env var not set\n")
	}
	backupPath = os.Getenv("BACKUP_PATH")
	if backupPath == "" {
		io.WriteString(file, "BACKUP_PATH env var not set\n")
		found = false
	}
	backupServerName = os.Getenv("BACKUP_SERVERNAME")
	if backupServerName == "" {
		io.WriteString(file, "BACKUP_SERVERNAME env var not set\n")
		found = false
	}
	backupProxyIP = os.Getenv("BACKUP_PROXY_IP")
	backupServerIP = os.Getenv("BACKUP_SERVERIP")
	if backupServerIP == "" {
		io.WriteString(file, "BACKUP_SERVERIP env var not set\n")
		found = false
	}
	backupProfileName = os.Getenv("BACKUP_PROFILENAME")
	if backupProfileName == "" {
		io.WriteString(file, "BACKUP_PROFILENAME env var not set\n")
		found = false
	}
	backupHost = os.Getenv("BACKUP_HOST")
	if backupHost == "" {
		io.WriteString(file, "BACKUP_HOST env var not set\n")
		found = false
	}

	var proxyHost = os.Getenv("BACKUP_PROXY_HOST")
	if proxyHost != "" {
		io.WriteString(file, "BACKUP_PROXY_HOST was set\n")
		backupHost = proxyHost
	}

	scheduleID = os.Getenv("BACKUP_SCHEDULEID")
	if scheduleID == "" {
		io.WriteString(file, "BACKUP_SCHEDULEID env var not set\n")
		found = false
	}
	backupPort = os.Getenv("BACKUP_PORT")
	if backupPort == "" {
		io.WriteString(file, "BACKUP_PORT env var not set\n")
		found = false
	}
	backupUser = os.Getenv("BACKUP_USER")
	if backupUser == "" {
		io.WriteString(file, "BACKUP_USER env var not set\n")
		found = false
	}
	backupAgentURL = os.Getenv("BACKUP_SERVER_URL")
	if backupAgentURL == "" {
		io.WriteString(file, "BACKUP_SERVER_URL env var not set\n")
		found = false
	}
	io.WriteString(file, "BACKUP_SERVER_URL ["+backupAgentURL+"]")

	if !found {
		panic("backup job missing required env vars")
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
		io.WriteString(file, err.Error())
		panic(err)
	}

	//only return the size from the du output
	var parsed = strings.Split(out.String(), "\t")
	return parsed[0]

}
