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
	"github.com/crunchydata/crunchy-postgresql-manager/backup"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var startTime time.Time
var backupContainerName string
var backupServerName string
var backupServerIP string
var backupProfileName string
var backupPath string
var backupHost string
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
	s := backup.BackupStatus{}
	eDuration := time.Since(startTime)
	s.StartTime = startTime.String()
	s.StartTime = startTime.String()
	s.ElapsedTime = eDuration.String()
	s.Status = "initializing"
	s.BackupSize = du()
	sendStats(s)

	//kick off stats reporting in a separate thread
	go stats("hi")

	io.WriteString(file, "giving DNS time to register the backup job....sleeping for 7 secs")
	time.Sleep(7000 * time.Millisecond)

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

	for true {
		io.WriteString(file, "sending stats...\n")
		io.WriteString(file, "sleeping for 7 secs\n")
		time.Sleep(7000 * time.Millisecond)
		stats := backup.BackupStatus{}
		eDuration := time.Since(startTime)
		stats.StartTime = startTime.String()
		stats.StartTime = startTime.String()
		stats.ElapsedTime = eDuration.String()
		stats.Status = "running"
		stats.BackupSize = du()
		sendStats(stats)
	}
}

//report stats back to the cpm-admin for this backup job
func finalstats(str string) {

	//connect to backupserver on cpm-admin
	//send stats to backup
	stats := backup.BackupStatus{}
	eDuration := time.Since(startTime)
	stats.StartTime = startTime.String()
	stats.ElapsedTime = eDuration.String()
	stats.Status = "completed"
	stats.BackupSize = du()

	sendStats(stats)
	io.WriteString(file, "final stats here\n")
}

//perform the backup
func backupfunc(str string) {
	//do a pg_basebackup here

	io.WriteString(file, "doing backup on "+backupHost+"\n")

	//create base backup from master
	cmd := exec.Command(CPMBIN+"basebackup.sh",
		backupHost)
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

func sendStats(stats backup.BackupStatus) error {
	stats.ContainerName = backupContainerName
	stats.ServerName = backupServerName
	stats.ScheduleID = scheduleID
	stats.ServerIP = backupServerIP
	stats.ProfileName = backupProfileName
	stats.Path = backupPath
	stats.BackupName = backupHost
	stats.ID = StatusID

	var err error
	if StatusID != "" {
		_, err = backup.UpdateStatusClient(backupAgentURL, stats)
	} else {
		StatusID, err = backup.AddStatusClient(backupAgentURL, stats)
	}
	if err != nil {
		io.WriteString(file, "error in adding status:"+err.Error()+"\n")
		return err
	}

	//send to backup
	io.WriteString(file, "elapsed time:"+stats.ElapsedTime+"\n")
	io.WriteString(file, "backupsize :"+stats.BackupSize+"\n")
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
