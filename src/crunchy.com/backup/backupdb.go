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
	"crunchy.com/mylog"
	"database/sql"
	"fmt"
	//"github.com/golang/logger.Info"
	_ "github.com/lib/pq"
	"strconv"
)

var logger mylog.MyLogger
var dbConn *sql.DB

func SetLogger(log mylog.MyLogger) {
	logger = log
}

func SetConnection(conn *sql.DB) {
	//logger.Info.Println("backupdb:SetConnection: called to open dbConn")
	dbConn = conn
}

func DBAddStatus(status BackupStatus) (string, error) {

	logger.Info.Println("DBAddStatus called")
	//logger.Info.Println("DBAddStatus called")

	queryStr := fmt.Sprintf("insert into backupstatus ( containername, starttime, backupname, servername, serverip, path, elapsedtime, backupsize, status, profilename, scheduleid, updatedt) values ( '%s', now(), '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %s, now()) returning id",
		status.ContainerName,
		status.BackupName,
		status.ServerName,
		status.ServerIP,
		status.Path,
		status.ElapsedTime,
		status.BackupSize,
		status.Status, status.ProfileName, status.ScheduleID)

	logger.Info.Println("DBAddStatus:" + queryStr)
	var theID int
	err := dbConn.QueryRow(queryStr).Scan(
		&theID)
	switch {
	case err != nil:
		logger.Error.Println("DBAddStatus: error " + err.Error())
		return "", err
	default:
	}

	var strvalue string
	strvalue = strconv.Itoa(theID)
	logger.Info.Println("DBAddStatus returning ID=" + strvalue)
	return strvalue, nil
}

func DBUpdateStatus(status BackupStatus) error {

	logger.Info.Println("backup.DBUpdateStatus called")

	queryStr := fmt.Sprintf("update backupstatus set ( status, backupsize, elapsedtime, updatedt) = ('%s', '%s', '%s', now()) where id = %s returning containername",
		status.Status,
		status.BackupSize,
		status.ElapsedTime,
		status.ID)

	logger.Info.Println("backup:DBUpdateStatus:[" + queryStr + "]")
	var name string
	err := dbConn.QueryRow(queryStr).Scan(&name)
	switch {
	case err != nil:
		logger.Error.Println("backup:DBUpdateStatus:" + err.Error())
		return err
	default:
	}

	return nil
}

func DBAddSchedule(s BackupSchedule) (string, error) {

	logger.Info.Println("DBAddSchedule called")

	queryStr := fmt.Sprintf("insert into backupschedule ( serverid, containername, profilename, name, enabled, minutes, hours, dayofmonth, month, dayofweek, updatedt) values ( '%s','%s','%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s',  now()) returning id",
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

	logger.Info.Println("DBAddSchedule:" + queryStr)
	var theID string
	err := dbConn.QueryRow(queryStr).Scan(
		&theID)
	if err != nil {
		logger.Error.Println("error in DBAddSchedule query " + err.Error())
		return "", err
	}

	switch {
	case err != nil:
		logger.Error.Println("DBAddSchedule: error " + err.Error())
		return "", err
	default:
	}

	return theID, nil
}

func DBUpdateSchedule(s BackupSchedule) error {

	logger.Info.Println("backup.DBUpdateSchedule called")

	queryStr := fmt.Sprintf("update backupschedule set ( enabled, serverid, name, minutes, hours, dayofmonth, month, dayofweek, updatedt) = ('%s', %s, '%s', '%s', '%s', '%s', '%s', '%s', now()) where id = %s  returning containername",
		s.Enabled,
		s.ServerID,
		s.Name,
		s.Minutes,
		s.Hours,
		s.DayOfMonth,
		s.Month,
		s.DayOfWeek,
		s.ID)

	logger.Info.Println("backup:DBUpdateSchedule:[" + queryStr + "]")
	var name string
	err := dbConn.QueryRow(queryStr).Scan(&name)
	switch {
	case err != nil:
		logger.Error.Println("backup:DBUpdateSchedule:" + err.Error())
		return err
	default:
	}

	return nil
}

func DBDeleteSchedule(id string) error {
	queryStr := fmt.Sprintf("delete from backupschedule where id=%s returning id", id)
	logger.Info.Println("backup:DBDeleteSchedule:" + queryStr)

	var theID int
	err := dbConn.QueryRow(queryStr).Scan(&theID)
	switch {
	case err != nil:
		return err
	default:
	}

	return nil
}

func DBGetSchedule(id string) (BackupSchedule, error) {
	logger.Info.Println("DBGetSchedule called with id=" + id)
	s := BackupSchedule{}

	err := dbConn.QueryRow(fmt.Sprintf("select a.id, a.serverid, b.name, b.ipaddress, a.containername, a.profilename, a.name, a.enabled, a.minutes, a.hours, a.dayofmonth, a.month, a.dayofweek, date_trunc('second', a.updatedt)::text from backupschedule a, server b where a.id=%s and b.id = a.serverid", id)).Scan(&s.ID, &s.ServerID, &s.ServerName, &s.ServerIP, &s.ContainerName, &s.ProfileName, &s.Name, &s.Enabled, &s.Minutes, &s.Hours, &s.DayOfMonth, &s.Month, &s.DayOfWeek, &s.UpdateDt)
	switch {
	case err == sql.ErrNoRows:
		logger.Error.Println("backupdb:DBGetSchedule:no schedule with that id")
		return s, err
	case err != nil:
		logger.Error.Println("backupdb:DBGetSchedule:" + err.Error())
		return s, err
	default:
	}

	return s, nil
}

func DBGetAllSchedules(containerid string) ([]BackupSchedule, error) {
	logger.Info.Println("DBGetAllSchedules called with id=" + containerid)
	var rows *sql.Rows
	var err error

	rows, err = dbConn.Query(fmt.Sprintf("select a.id, a.serverid, s.name, s.ipaddress, a.containername, a.profilename, a.name, a.enabled, a.minutes, a.hours, a.dayofmonth, a.month, a.dayofweek, date_trunc('second', a.updatedt)::text from backupschedule a, node b, server s where a.containername= b.name and b.id = %s and a.serverid = s.id", containerid))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	schedules := make([]BackupSchedule, 0)
	for rows.Next() {
		s := BackupSchedule{}
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

func DBGetAllStatus(scheduleid string) ([]BackupStatus, error) {
	logger.Info.Println("DBGetAllStatus called with scheduleid=" + scheduleid)
	var rows *sql.Rows
	var err error

	rows, err = dbConn.Query(fmt.Sprintf("select id, containername, date_trunc('second', starttime)::text, backupname, servername, serverip, path, elapsedtime, backupsize, status, date_trunc('second', updatedt)::text from backupstatus where scheduleid=%s order by starttime", scheduleid))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stats := make([]BackupStatus, 0)
	for rows.Next() {
		s := BackupStatus{}
		if err = rows.Scan(
			&s.ID,
			&s.ContainerName,
			&s.StartTime,
			&s.BackupName,
			&s.ServerName,
			&s.ServerIP,
			&s.Path,
			&s.ElapsedTime,
			&s.BackupSize,
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

func DBGetStatus(id string) (BackupStatus, error) {
	logger.Info.Println("DBGetStatus called with id=" + id)
	s := BackupStatus{}

	err := dbConn.QueryRow(fmt.Sprintf("select id, containername, date_trunc('second', starttime), backupname, servername, serverip, path, elapsedtime, backupsize, status, date_trunc('second', updatedt) from backupstatus where id=%s", id)).Scan(&s.ID, &s.ContainerName, &s.StartTime, &s.BackupName, &s.ServerName, &s.ServerIP, &s.Path, &s.ElapsedTime, &s.BackupSize, &s.Status, &s.UpdateDt)
	switch {
	case err == sql.ErrNoRows:
		logger.Error.Println("backupdb:DBGetStatus:no status with that id")
		return s, err
	case err != nil:
		logger.Error.Println("backupdb:DBGetStatus:" + err.Error())
		return s, err
	default:
	}

	return s, nil
}

func DBGetSchedules() ([]BackupSchedule, error) {
	logger.Info.Println("DBGetSchedules called")
	var rows *sql.Rows
	var err error

	rows, err = dbConn.Query(fmt.Sprintf("select a.id, a.serverid, a.containername, a.profilename, a.name, a.enabled, a.minutes, a.hours, a.dayofmonth, a.month, a.dayofweek, date_trunc('second', a.updatedt)::text from backupschedule a "))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	schedules := make([]BackupSchedule, 0)
	for rows.Next() {
		s := BackupSchedule{}
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
