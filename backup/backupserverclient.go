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
	"errors"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"net/rpc"
)

//called by backup jobs as they execute
func AddStatusClient(ipaddress string, status BackupStatus) (string, error) {

	logit.Info.Println("AddStatus called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logit.Error.Println("AddStatus: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logit.Error.Println("AddStatus: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.AddStatus", &status, &command)
	if err != nil {
		logit.Error.Println("AddStatus: error " + err.Error())
		return "", err
	}
	logit.Info.Println("status.ID=" + status.ID)

	return command.Output, nil
}

//called by backup jobs as they execute
func UpdateStatusClient(ipaddress string, status BackupStatus) (string, error) {

	logit.Info.Println("UpdateStatus called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logit.Error.Println("UpdateStatus: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logit.Error.Println("UpdateStatus: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.UpdateStatus", &status, &command)
	if err != nil {
		logit.Error.Println("UpdateStatus: error " + err.Error())
		return "", err
	}

	return command.Output, nil
}

//called by admin do perform an adhoc backup job
func BackupNowClient(ipaddress string, request BackupRequest) (string, error) {

	logit.Info.Println("BackupNow called ip=" + ipaddress)
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logit.Error.Println("BackupNow: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logit.Error.Println("BackupNow: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.BackupNow", &request, &command)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		return "", err
	}

	return command.Output, nil
}

//called by admin to add to reload schedules in the backup server
func ReloadClient(ipaddress string, sched BackupSchedule) (string, error) {

	logit.Info.Println("ReloadClient called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logit.Error.Println("ReloadClient: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logit.Error.Println("ReloadClient: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.Reload", &sched, &command)
	if err != nil {
		logit.Error.Println("ReloadError: error " + err.Error())
		return "", err
	}

	return command.Output, nil
}
