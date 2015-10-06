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
	"encoding/json"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/task"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var startTime time.Time
var startTimeString string
var nodeName, repoRemotePath, backupHost, backupSet, backrestKeyPass string
var nodeProfileName, backupServerName, backupServerIP string
var dataDir, serverID, token, dockerProfileName, projectID string

var scheduleID string
var backupPort string
var backupUser string
var backupAgentURL string
var StatusID = ""
var CPMBIN = "/var/cpm/bin/"
var CPM_ADMIN = "cpm-admin:13001"

func init() {
	startTime = time.Now()
	startTimeString = startTime.String()
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
	//go stats("hi")

	logit.Info.Println("giving DNS time to register the backup job....sleeping for 7 secs")
	sleepTime, _ := time.ParseDuration("7s")
	time.Sleep(sleepTime)

	//perform the restore
	restore("end")

	//start pg locally to complete the recovery
	//exec startpg.sh

	newid, err := createContainer()
	if err != nil {
		s.Status = "error in createContainer"
		sendStats(&s)
		return
	}

	logit.Info.Println("created container ID=" + newid)

	err = stopContainer()
	if err != nil {
		s.Status = "error in stopContainer"
		sendStats(&s)
		return
	}
	err = switchPaths()
	if err != nil {
		s.Status = "error in switchPaths"
		sendStats(&s)
		return
	}

	err = startContainer()
	if err != nil {
		s.Status = "error in startContainer"
		sendStats(&s)
		return
	}
	err = seedDatabase()
	if err != nil {
		s.Status = "error in seedDatabase"
		sendStats(&s)
		return
	}

	//send final stats to backup
	finalstats("end")

}

//report stats back to the cpm-admin for this backup job
func stats(str string) {

	sleepTime, _ := time.ParseDuration("7s")

	for true {
		logit.Info.Println("sending stats...")
		logit.Info.Println("sleeping for 7 secs")
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

	cmd := exec.Command("pg_backrest", "restore", "--type=none", "--stanza=main", backupSet)
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

	scheduleID = os.Getenv("BACKUP_SCHEDULEID")
	if scheduleID == "" {
		logit.Error.Println("BACKUP_SCHEDULEID env var not set")
		found = false
	}

	projectID = os.Getenv("BACKUP_PROJECTID")
	if projectID == "" {
		logit.Error.Println("BACKUP_PROJECTID env var not set")
		found = false
	}

	serverID = os.Getenv("BACKUP_SERVERID")
	if serverID == "" {
		logit.Error.Println("BACKUP_SERVERID env var not set")
		found = false
	}

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
	token = os.Getenv("TOKEN")
	if token == "" {
		logit.Error.Println("TOKEN env var not set")
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
	dockerProfileName = os.Getenv("DOCKER_PROFILENAME")
	if dockerProfileName == "" {
		logit.Error.Println("DOCKER_PROFILENAME env var not set")
		found = false
	}
	dataDir = os.Getenv("DATA_DIR")
	if dataDir == "" {
		logit.Error.Println("DATA_DIR env var not set")
		found = false
	}

	if !found {
		panic("backrestrestore job missing required env vars")
	}

}

func du() string {
	path := os.Getenv("PGDATA")
	logit.Info.Println("PGDATA is " + path)
	cmd := exec.Command("du", "-hs", path)
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
func createContainer() (string, error) {
	var err error

	params := &cpmserverapi.DockerRunRequest{}
	params.Image = "cpm-node"
	params.ServerID = serverID
	params.ProjectID = projectID
	params.ContainerName = nodeName
	params.Standalone = "true"

	response := types.ProvisionStatus{}
	url := "http://" + CPM_ADMIN + "/provision/" +
		dockerProfileName + "." +
		params.Image + "." +
		params.ServerID + "." +
		params.ProjectID + "." +
		params.ContainerName + "." +
		params.Standalone + "." +
		token
	logit.Info.Println("url=" + url)
	r, err := http.Get(url)
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}
	rawresponse, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawresponse, &response)
	logit.Info.Println(string(rawresponse))

	return response.ID, err
}

func stopContainer() error {
	req := cpmserverapi.DockerStopRequest{}
	req.ContainerName = nodeName
	response, err := cpmserverapi.DockerStopClient("http://"+backupHost+":10001", &req)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Println("stopContainer:" + response.Output)
	return err
}

func switchPaths() error {
	var err error

	var req = cpmserverapi.SwitchPathRequest{}
	req.DataDir = dataDir
	req.ContainerName = nodeName

	var url = "http://" + backupHost + ":10001"
	response, err := cpmserverapi.SwitchPathClient(url, &req)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	logit.Info.Println(response.Output)

	return err
}

func startContainer() error {
	req := cpmserverapi.DockerStartRequest{}
	req.ContainerName = nodeName
	response, err := cpmserverapi.DockerStartClient("http://"+backupHost+":10001", &req)

	logit.Info.Println("startContainer:" + response.Output)
	return err
}

func seedDatabase() error {
	var err error
	return err
}
