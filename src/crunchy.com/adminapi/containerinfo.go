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
	"crunchy.com/admindb"
	"crunchy.com/cpmagent"
	"crunchy.com/util"
	"database/sql"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"net/http"
	"os/exec"
	"strings"
)

func MonitorContainerLoadtest(w rest.ResponseWriter, r *rest.Request) {
	ID := r.PathParam("ID")
	Writes := r.PathParam("Writes")

	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	if Writes == "" {
		rest.Error(w, "Writes required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("MonitorContainerGetInfo:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server, err2 := admindb.GetDBServer(node.ServerID)
	if err2 != nil {
		glog.Errorln("MonitorContainerGetInfo:" + err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	var output string
	var port = "5432"

	output, err = cpmagent.AgentCommandConfigureNode(CPMBIN+"loadtest", node.Name,
		port, Writes, "", "", "", "", server.IPAddress)
	if err != nil {
		glog.Errorln("MonitorContainerGetInfo:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.(http.ResponseWriter).Write([]byte(output))
	w.WriteHeader(http.StatusOK)
}

func MonitorContainerSettings(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("MonitorContainerSettings: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("MonitorContainerGetInfo:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConn, err := util.GetMonitoringConnection(node.Name, "postgres", "5432", "postgres", "")
	defer dbConn.Close()

	settings := make([]PostgresSetting, 0)
	var rows *sql.Rows

	rows, err = dbConn.Query("select name, current_setting(name), source from pg_settings where source not in ('default','override')")
	if err != nil {
		glog.Errorln("MonitorContainerSettings:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()
	for rows.Next() {
		setting := PostgresSetting{}
		if err = rows.Scan(
			&setting.Name,
			&setting.CurrentSetting,
			&setting.Source,
		); err != nil {
			glog.Errorln("MonitorContainerSettings:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		settings = append(settings, setting)
	}
	if err = rows.Err(); err != nil {
		glog.Errorln("MonitorContainerSettings:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&settings)
}

func MonitorContainerControldata(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var output string
	err = secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("MonitorContainerControldata: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("MonitorContainerControldata:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings := make([]PostgresControldata, 0)

	//send the container a pg_controldata command
	output, err = cpmagent.AgentCommand(PGBIN+"pg_controldata", "/pgdata", node.Name)
	if err != nil {
		glog.Errorln("MonitorContainerControldata:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	glog.Infoln(output)

	lines := strings.Split(output, "\n")
	//fmt.Println(len(lines))
	for i := range lines {
		//fmt.Println(len(lines[i]))
		if len(lines[i]) > 1 {
			setting := PostgresControldata{}
			columns := strings.Split(lines[i], ":")
			setting.Name = strings.TrimSpace(columns[0])
			setting.Value = strings.TrimSpace(columns[1])
			//fmt.Println("name=[" + strings.TrimSpace(columns[0]) + "] value=[" + strings.TrimSpace(columns[1]) + "]")
			settings = append(settings, setting)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&settings)
}

type Bgwriter struct {
	Now            string
	AllocMbps      string
	CheckpointMbps string
	CleanMbps      string
	BackendMbps    string
	WriteMbps      string
}

func ContainerInfoBgwriter(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("ContainerBgwriter: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("ContainerBgwriter:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConn, err := util.GetMonitoringConnection(node.Name, "postgres", "5432", "postgres", "")
	defer dbConn.Close()

	info := Bgwriter{}
	err = dbConn.QueryRow("SELECT to_char(now(), 'mm/dd/yy HH12:MI:SS') now, to_char(block_size::numeric * buffers_alloc / (1024 * 1024 * seconds), 'FM999999999999D9999') AS alloc_mbps, to_char(block_size::numeric * buffers_checkpoint / (1024 * 1024 * seconds), 'FM999999999999D9999') AS checkpoint_mbps, to_char(block_size::numeric * buffers_clean / (1024 * 1024 * seconds), 'FM999999999999D9999') AS clean_mbps, to_char(block_size::numeric * buffers_backend/ (1024 * 1024 * seconds), 'FM999999999999D9999') AS backend_mbps, to_char(block_size::numeric * (buffers_checkpoint + buffers_clean + buffers_backend) / (1024 * 1024 * seconds), 'FM999999999999D9999') AS write_mbps FROM ( SELECT now() AS sample,now() - stats_reset AS uptime,EXTRACT(EPOCH FROM now()) - extract(EPOCH FROM stats_reset) AS seconds, b.*,p.setting::integer AS block_size FROM pg_stat_bgwriter b,pg_settings p WHERE p.name='block_size') bgw").Scan(&info.Now, &info.AllocMbps, &info.CheckpointMbps, &info.CleanMbps, &info.BackendMbps, &info.WriteMbps)
	switch {
	case err == sql.ErrNoRows:
		glog.Errorln("ContainerBgwriter:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	case err != nil:
		glog.Errorln("ContainerBgwriter:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&info)
}

type Statdatabase struct {
	Datname     string
	BlksRead    string
	TupReturned string
	TupFetched  string
	TupInserted string
	TupUpdated  string
	TupDeleted  string
	StatsReset  string
}

func ContainerInfoStatdatabase(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("ContainerStatdatabase: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("ContainerStatdatabase:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConn, err := util.GetMonitoringConnection(node.Name, "postgres", "5432", "postgres", "")
	defer dbConn.Close()

	stats := make([]Statdatabase, 0)
	var rows *sql.Rows

	rows, err = dbConn.Query("SELECT datname, blks_read::text, tup_returned::text, tup_fetched::text, tup_inserted::text, tup_updated::text, tup_deleted::text, coalesce(to_char(stats_reset, 'YYYY-MM-DD HH24:MI:SS'), ' ') as stats_reset from pg_stat_database")
	if err != nil {
		glog.Errorln("ContainerStatdatabase:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()
	for rows.Next() {
		stat := Statdatabase{}
		if err = rows.Scan(
			&stat.Datname,
			&stat.BlksRead,
			&stat.TupReturned,
			&stat.TupFetched,
			&stat.TupInserted,
			&stat.TupUpdated,
			&stat.TupDeleted,
			&stat.StatsReset,
		); err != nil {
			glog.Errorln("ContainerStatdatabase:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stats = append(stats, stat)
	}
	if err = rows.Err(); err != nil {
		glog.Errorln("ContainerStatdatabase:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&stats)
}

type Statrepl struct {
	Pid            string
	Usesysid       string
	Usename        string
	AppName        string
	ClientAddr     string
	ClientHostname string
	ClientPort     string
	BackendStart   string
	State          string
	SentLocation   string
	WriteLocation  string
	FlushLocation  string
	ReplayLocation string
	SyncPriority   string
	SyncState      string
}

func ContainerInfoStatrepl(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("ContainerStatrepl: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("ContainerStatrepl:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConn, err := util.GetMonitoringConnection(node.Name, "postgres", "5432", "postgres", "")
	defer dbConn.Close()

	stats := make([]Statrepl, 0)
	var rows *sql.Rows

	rows, err = dbConn.Query("SELECT pid , usesysid , usename , application_name , client_addr , client_hostname , client_port , to_char(backend_start, 'YYYY-MM-DD HH24:MI-SS') as backend_start , state , sent_location , write_location , flush_location , replay_location , sync_priority , sync_state from pg_stat_replication")
	if err != nil {
		glog.Errorln("ContainerStatrepl:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()
	for rows.Next() {
		stat := Statrepl{}
		if err = rows.Scan(
			&stat.Pid,
			&stat.Usesysid,
			&stat.Usename,
			&stat.AppName,
			&stat.ClientAddr,
			&stat.ClientHostname,
			&stat.ClientPort,
			&stat.BackendStart,
			&stat.State,
			&stat.SentLocation,
			&stat.WriteLocation,
			&stat.FlushLocation,
			&stat.ReplayLocation,
			&stat.SyncPriority,
			&stat.SyncState,
		); err != nil {
			glog.Errorln("ContainerStatrepl:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stats = append(stats, stat)
	}
	if err = rows.Err(); err != nil {
		glog.Errorln("ContainerStatrepl:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&stats)
}

func ContainerLoadTest(w rest.ResponseWriter, r *rest.Request) {
	var err error

	err = secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("ContainerLoadTest: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")

	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	Writes := r.PathParam("Writes")

	if Writes == "" {
		rest.Error(w, "Writes required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetDBNode(ID)
	if err != nil {
		glog.Errorln("ContainerLoadTest:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//$1 - host node.Name
	//$2 - port "5432"
	//$3 - user "cpmtest"
	//$4 - password "cpmtest"
	//$5 - database "cpmtest"
	//$6 - insert count - only valid for loadtest Metric

	//hardcoded for now, TODO pull from metadata
	var port = "5432"
	var user = "cpmtest"
	var password = "cpmtest"
	var database = "cpmtest"

	cmd := exec.Command(CPMBIN+"loadtest",
		node.Name,
		port,
		user,
		password,
		database,
		Writes)

	for i := 0; i < len(cmd.Args); i++ {
		glog.Infoln("ContainerLoadTest:" + cmd.Args[i])
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		glog.Errorln("ContainerLoadTest:" + err.Error())
		glog.Flush()
		errorString := fmt.Sprintf("%s\n%s\n%s\n", err.Error(), out.String(), stderr.String())
		rest.Error(w, errorString, http.StatusBadRequest)
		return
	}
	glog.Infoln("ContainerLoadTest: command output was " + out.String())

	//w.(http.ResponseWriter).Write([]byte(output))
	w.(http.ResponseWriter).Write([]byte(out.String()))
	w.WriteHeader(http.StatusOK)
}
