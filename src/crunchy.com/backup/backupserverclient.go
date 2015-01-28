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
	"errors"
	"net/rpc"
)

//called by backup jobs as they execute
func AddStatusClient(ipaddress string, status BackupStatus) (string, error) {

	logutil.Log("AddStatus called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logutil.Log("AddStatus: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logutil.Log("AddStatus: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.AddStatus", &status, &command)
	if err != nil {
		logutil.Log("AddStatus: error " + err.Error())
		return "", err
	}
	logutil.Log("status.ID=" + status.ID)

	return command.Output, nil
}

//called by backup jobs as they execute
func UpdateStatusClient(ipaddress string, status BackupStatus) (string, error) {

	logutil.Log("UpdateStatus called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logutil.Log("UpdateStatus: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logutil.Log("UpdateStatus: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.UpdateStatus", &status, &command)
	if err != nil {
		logutil.Log("UpdateStatus: error " + err.Error())
		return "", err
	}

	return command.Output, nil
}

//called by admin do perform an adhoc backup job
func BackupNowClient(ipaddress string, request BackupRequest) (string, error) {

	logutil.Log("BackupNow called ip=" + ipaddress)
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logutil.Log("BackupNow: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logutil.Log("BackupNow: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.BackupNow", &request, &command)
	if err != nil {
		logutil.Log("BackupNow: error " + err.Error())
		return "", err
	}

	return command.Output, nil
}

//called by admin to add to reload schedules in the backup server
func ReloadClient(ipaddress string, sched BackupSchedule) (string, error) {

	logutil.Log("ReloadClient called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logutil.Log("ReloadClient: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logutil.Log("ReloadClient: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.Reload", &sched, &command)
	if err != nil {
		logutil.Log("ReloadError: error " + err.Error())
		return "", err
	}

	return command.Output, nil
}

/*
//called by admin to add to the schedule
func AddScheduleClient(ipaddress string, sched BackupSchedule) (string, error) {

	logutil.Log("AddSchedule called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logutil.Log("AddSchedule: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logutil.Log("AddSchedule: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.AddSchedule", &sched, &command)
	if err != nil {
		logutil.Log("AddSchedule: error " + err.Error())
		return "", err
	}

	return command.Output, nil
}

//called by admin to update the schedule
func UpdateScheduleClient(ipaddress string, sched BackupSchedule) (string, error) {

	logutil.Log("UpdateSchedule called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logutil.Log("UpdateSchedule: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logutil.Log("UpdateSchedule: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.UpdateSchedule", &sched, &command)
	if err != nil {
		logutil.Log("UpdateSchedule: error " + err.Error())
		return "", err
	}

	return command.Output, nil
}

//called by admin to delete the schedule
func DeleteScheduleClient(ipaddress string, sched BackupSchedule) (string, error) {

	logutil.Log("DeleteSchedule called")
	client, err := rpc.DialHTTP("tcp", ipaddress)
	if err != nil {
		logutil.Log("DeleteSchedule: dialing:" + err.Error())
		return "", err
	}
	if client == nil {
		logutil.Log("DeleteSchedule: client was nil")
		return "", errors.New("client was nil from rpc dial")
	}

	var command Command

	err = client.Call("Command.DeleteSchedule", &sched, &command)
	if err != nil {
		logutil.Log("DeleteSchedule: error " + err.Error())
		return "", err
	}

	return command.Output, nil
}
*/
