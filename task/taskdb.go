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

package task

import (
	"database/sql"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	_ "github.com/lib/pq"
	"strconv"
)

func AddStatus(dbConn *sql.DB, status *TaskStatus) (string, error) {

	logit.Info.Println("AddStatus called")
	//logit.Info.Println("AddStatus called")

	queryStr := fmt.Sprintf("insert into taskstatus ( containername, starttime, taskname, servername, serverip, path, elapsedtime, tasksize, status, profilename, scheduleid, updatedt) values ( '%s', now(), '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %s, now()) returning id",
		status.ContainerName,
		status.TaskName,
		status.ServerName,
		status.ServerIP,
		status.Path,
		status.ElapsedTime,
		status.TaskSize,
		status.Status, status.ProfileName, status.ScheduleID)

	logit.Info.Println("AddStatus:" + queryStr)
	var theID int
	err := dbConn.QueryRow(queryStr).Scan(
		&theID)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return "", err
	default:
	}

	var strvalue string
	strvalue = strconv.Itoa(theID)
	logit.Info.Println("AddStatus returning ID=" + strvalue)
	return strvalue, nil
}

func UpdateStatus(dbConn *sql.DB, status *TaskStatus) error {

	logit.Info.Println("backup.UpdateStatus called")

	queryStr := fmt.Sprintf("update taskstatus set ( status, tasksize, elapsedtime, updatedt) = ('%s', '%s', '%s', now()) where id = %s returning containername",
		status.Status,
		status.TaskSize,
		status.ElapsedTime,
		status.ID)

	logit.Info.Println("backup:UpdateStatus:[" + queryStr + "]")
	var name string
	err := dbConn.QueryRow(queryStr).Scan(&name)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return err
	default:
	}

	return nil
}

func AddSchedule(dbConn *sql.DB, s TaskSchedule) (string, error) {

	logit.Info.Println("AddSchedule called")

	queryStr := fmt.Sprintf("insert into taskschedule ( serverid, containername, profilename, name, enabled, minutes, hours, dayofmonth, month, dayofweek, updatedt) values ( '%s','%s','%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s',  now()) returning id",
		s.ServerID,
		s.ContainerName,
		s.ProfileName,
		s.Name,
		s.Enabled,
		s.Minutes,
		s.Hours,
		s.DayOfMonth,
		s.Month,
		s.DayOfWeek)

	logit.Info.Println("AddSchedule:" + queryStr)
	var theID string
	err := dbConn.QueryRow(queryStr).Scan(
		&theID)
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}

	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return "", err
	default:
	}

	return theID, nil
}

func UpdateSchedule(dbConn *sql.DB, s TaskSchedule) error {

	logit.Info.Println("backup.UpdateSchedule called")

	queryStr := fmt.Sprintf("update taskschedule set ( enabled, serverid, name, minutes, hours, dayofmonth, month, dayofweek, updatedt) = ('%s', %s, '%s', '%s', '%s', '%s', '%s', '%s', now()) where id = %s  returning containername",
		s.Enabled,
		s.ServerID,
		s.Name,
		s.Minutes,
		s.Hours,
		s.DayOfMonth,
		s.Month,
		s.DayOfWeek,
		s.ID)

	logit.Info.Println("backup:UpdateSchedule:[" + queryStr + "]")
	var name string
	err := dbConn.QueryRow(queryStr).Scan(&name)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return err
	default:
	}

	return nil
}

func DeleteSchedule(dbConn *sql.DB, id string) error {
	queryStr := fmt.Sprintf("delete from taskschedule where id=%s returning id", id)
	logit.Info.Println("backup:DeleteSchedule:" + queryStr)

	var theID int
	err := dbConn.QueryRow(queryStr).Scan(&theID)
	switch {
	case err != nil:
		return err
	default:
	}

	return nil
}

func GetSchedule(dbConn *sql.DB, id string) (TaskSchedule, error) {
	logit.Info.Println("GetSchedule called with id=" + id)
	s := TaskSchedule{}

	err := dbConn.QueryRow(fmt.Sprintf("select a.id, a.serverid, b.name, b.ipaddress, a.containername, a.profilename, a.name, a.enabled, a.minutes, a.hours, a.dayofmonth, a.month, a.dayofweek, date_trunc('second', a.updatedt)::text from taskschedule a, server b where a.id=%s and b.id = a.serverid", id)).Scan(&s.ID, &s.ServerID, &s.ServerName, &s.ServerIP, &s.ContainerName, &s.ProfileName, &s.Name, &s.Enabled, &s.Minutes, &s.Hours, &s.DayOfMonth, &s.Month, &s.DayOfWeek, &s.UpdateDt)
	switch {
	case err == sql.ErrNoRows:
		logit.Error.Println("taskdb:GetSchedule:no schedule with that id")
		return s, err
	case err != nil:
		logit.Error.Println(err.Error())
		return s, err
	default:
	}

	return s, nil
}

func GetAllSchedules(dbConn *sql.DB, containerid string) ([]TaskSchedule, error) {
	logit.Info.Println("GetAllSchedules called with id=" + containerid)
	var rows *sql.Rows
	var err error

	rows, err = dbConn.Query(fmt.Sprintf("select a.id, a.serverid, s.name, s.ipaddress, a.containername, a.profilename, a.name, a.enabled, a.minutes, a.hours, a.dayofmonth, a.month, a.dayofweek, date_trunc('second', a.updatedt)::text from taskschedule a, container b, server s where a.containername= b.name and b.id = %s and a.serverid = s.id", containerid))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	schedules := make([]TaskSchedule, 0)
	for rows.Next() {
		s := TaskSchedule{}
		if err = rows.Scan(
			&s.ID,
			&s.ServerID,
			&s.ServerName,
			&s.ServerIP,
			&s.ContainerName,
			&s.ProfileName,
			&s.Name,
			&s.Enabled,
			&s.Minutes,
			&s.Hours,
			&s.DayOfMonth,
			&s.Month,
			&s.DayOfWeek,
			&s.UpdateDt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func GetAllStatus(dbConn *sql.DB, scheduleid string) ([]TaskStatus, error) {
	logit.Info.Println("GetAllStatus called with scheduleid=" + scheduleid)
	var rows *sql.Rows
	var err error

	rows, err = dbConn.Query(fmt.Sprintf("select id, containername, date_trunc('second', starttime)::text, taskname, servername, serverip, path, elapsedtime, tasksize, status, date_trunc('second', updatedt)::text from taskstatus where scheduleid=%s order by starttime", scheduleid))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stats := make([]TaskStatus, 0)
	for rows.Next() {
		s := TaskStatus{}
		if err = rows.Scan(
			&s.ID,
			&s.ContainerName,
			&s.StartTime,
			&s.TaskName,
			&s.ServerName,
			&s.ServerIP,
			&s.Path,
			&s.ElapsedTime,
			&s.TaskSize,
			&s.Status,
			&s.UpdateDt); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func GetStatus(dbConn *sql.DB, id string) (TaskStatus, error) {
	logit.Info.Println("GetStatus called with id=" + id)
	s := TaskStatus{}

	err := dbConn.QueryRow(fmt.Sprintf("select id, containername, date_trunc('second', starttime), taskname, servername, serverip, path, elapsedtime, tasksize, status, date_trunc('second', updatedt) from taskstatus where id=%s", id)).Scan(&s.ID, &s.ContainerName, &s.StartTime, &s.TaskName, &s.ServerName, &s.ServerIP, &s.Path, &s.ElapsedTime, &s.TaskSize, &s.Status, &s.UpdateDt)
	switch {
	case err == sql.ErrNoRows:
		logit.Error.Println("taskdb:GetStatus:no status with that id")
		return s, err
	case err != nil:
		logit.Error.Println(err.Error())
		return s, err
	default:
	}

	return s, nil
}

func GetSchedules(dbConn *sql.DB) ([]TaskSchedule, error) {
	logit.Info.Println("GetSchedules called")
	var rows *sql.Rows
	var err error

	rows, err = dbConn.Query(fmt.Sprintf("select a.id, a.serverid, a.containername, a.profilename, a.name, a.enabled, a.minutes, a.hours, a.dayofmonth, a.month, a.dayofweek, date_trunc('second', a.updatedt)::text from taskschedule a "))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	schedules := make([]TaskSchedule, 0)
	for rows.Next() {
		s := TaskSchedule{}
		if err = rows.Scan(
			&s.ID,
			&s.ServerID,
			&s.ContainerName,
			&s.ProfileName,
			&s.Name,
			&s.Enabled,
			&s.Minutes,
			&s.Hours,
			&s.DayOfMonth,
			&s.Month,
			&s.DayOfWeek,
			&s.UpdateDt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}
