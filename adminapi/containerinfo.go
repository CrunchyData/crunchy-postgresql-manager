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

package adminapi

import (
	"database/sql"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var CPMTEST_DB = "cpmtest"
var CPMTEST_USER = "cpmtest"

func MonitorContainerSettings(w rest.ResponseWriter, r *rest.Request) {

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("MonitorContainerSettings: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("MonitorContainerGetInfo:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var host = node.Name
	if KubeEnv {
		host = node.Name + "-db"
	}

	//fetch cpmtest user credentials
	var nodeuser admindb.ContainerUser
	nodeuser, err = admindb.GetContainerUser(dbConn, node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("cpmtest password is " + nodeuser.Passwd)

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")

	dbConn2, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, nodeuser.Passwd)
	defer dbConn2.Close()

	settings := make([]PostgresSetting, 0)
	var rows *sql.Rows

	rows, err = dbConn2.Query("select name, current_setting(name), source from pg_settings where source not in ('default','override')")
	if err != nil {
		logit.Error.Println("MonitorContainerSettings:" + err.Error())
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
			logit.Error.Println("MonitorContainerSettings:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		settings = append(settings, setting)
	}
	if err = rows.Err(); err != nil {
		logit.Error.Println("MonitorContainerSettings:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&settings)
}

func MonitorContainerControldata(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("MonitorContainerControldata: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("MonitorContainerControldata:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings := make([]PostgresControldata, 0)

	//send the container a pg_controldata command
	var cdout cpmcontainerapi.ControldataResponse
	cdout, err = cpmcontainerapi.ControldataClient(node.Name)
	if err != nil {
		logit.Error.Println("MonitorContainerControldata:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logit.Info.Println(cdout.Output)

	lines := strings.Split(cdout.Output, "\n")
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
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("ContainerBgwriter: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("ContainerBgwriter:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var host = node.Name
	if KubeEnv {
		host = node.Name + "-db"
	}

	//get password
	var nodeuser admindb.ContainerUser
	nodeuser, err = admindb.GetContainerUser(dbConn, node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println("ContainerBgwriter:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")

	var dbConn2 *sql.DB
	dbConn2, err = util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, nodeuser.Passwd)
	defer dbConn2.Close()

	info := Bgwriter{}
	err = dbConn2.QueryRow("SELECT to_char(now(), 'mm/dd/yy HH12:MI:SS') now, to_char(block_size::numeric * buffers_alloc / (1024 * 1024 * seconds), 'FM999999999999D9999') AS alloc_mbps, to_char(block_size::numeric * buffers_checkpoint / (1024 * 1024 * seconds), 'FM999999999999D9999') AS checkpoint_mbps, to_char(block_size::numeric * buffers_clean / (1024 * 1024 * seconds), 'FM999999999999D9999') AS clean_mbps, to_char(block_size::numeric * buffers_backend/ (1024 * 1024 * seconds), 'FM999999999999D9999') AS backend_mbps, to_char(block_size::numeric * (buffers_checkpoint + buffers_clean + buffers_backend) / (1024 * 1024 * seconds), 'FM999999999999D9999') AS write_mbps FROM ( SELECT now() AS sample,now() - stats_reset AS uptime,EXTRACT(EPOCH FROM now()) - extract(EPOCH FROM stats_reset) AS seconds, b.*,p.setting::integer AS block_size FROM pg_stat_bgwriter b,pg_settings p WHERE p.name='block_size') bgw").Scan(&info.Now, &info.AllocMbps, &info.CheckpointMbps, &info.CleanMbps, &info.BackendMbps, &info.WriteMbps)
	switch {
	case err == sql.ErrNoRows:
		logit.Error.Println("ContainerBgwriter:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	case err != nil:
		logit.Error.Println("ContainerBgwriter:" + err.Error())
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
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("ContainerStatdatabase: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("ContainerStatdatabase:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var host = node.Name
	if KubeEnv {
		host = node.Name + "-db"
	}

	//get password
	var nodeuser admindb.ContainerUser
	nodeuser, err = admindb.GetContainerUser(dbConn, node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")

	dbConn2, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, nodeuser.Passwd)
	defer dbConn2.Close()

	stats := make([]Statdatabase, 0)
	var rows *sql.Rows

	rows, err = dbConn2.Query("SELECT datname, blks_read::text, tup_returned::text, tup_fetched::text, tup_inserted::text, tup_updated::text, tup_deleted::text, coalesce(to_char(stats_reset, 'YYYY-MM-DD HH24:MI:SS'), ' ') as stats_reset from pg_stat_database")
	if err != nil {
		logit.Error.Println("ContainerStatdatabase:" + err.Error())
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
			logit.Error.Println("ContainerStatdatabase:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stats = append(stats, stat)
	}
	if err = rows.Err(); err != nil {
		logit.Error.Println("ContainerStatdatabase:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&stats)
}

type Loadtestresults struct {
	Operation string
	Count     int
	Results   string
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
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("ContainerStatrepl: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("ContainerStatrepl:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var host = node.Name
	if KubeEnv {
		host = node.Name + "-db"
	}

	//fetch cpmtest user credentials
	var nodeuser admindb.ContainerUser
	nodeuser, err = admindb.GetContainerUser(dbConn, node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")

	dbConn2, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, nodeuser.Passwd)
	defer dbConn2.Close()

	stats := make([]Statrepl, 0)
	var rows *sql.Rows

	rows, err = dbConn2.Query("SELECT pid , usesysid , usename , application_name , client_addr , coalesce(client_hostname, ' ') , client_port , to_char(backend_start, 'YYYY-MM-DD HH24:MI-SS') as backend_start , state , sent_location , write_location , flush_location , replay_location , sync_priority , sync_state from pg_stat_replication")
	if err != nil {
		logit.Error.Println("ContainerStatrepl:" + err.Error())
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
			logit.Error.Println("ContainerStatrepl:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stats = append(stats, stat)
	}
	if err = rows.Err(); err != nil {
		logit.Error.Println("ContainerStatrepl:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&stats)
}

func ContainerLoadTest(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("ContainerLoadTest: authorize error " + err.Error())
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
	var writes int
	writes, err = strconv.Atoi(Writes)
	if err != nil {
		logit.Error.Println("ContainerLoadTest:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("ContainerLoadTest:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var host = node.Name
	if KubeEnv {
		host = node.Name + "-db"
	}

	results, err2 := loadtest(dbConn, node.Name, host, writes)
	if err2 != nil {
		logit.Error.Println("ContainerLoadTest:" + err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&results)
}

func loadtest(dbConn *sql.DB, nodename string, host string, writes int) ([]Loadtestresults, error) {
	var err error
	var name string
	var query string
	var id int
	var i int
	var dbConn2 *sql.DB
	var results = make([]Loadtestresults, 4)

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")

	//fetch cpmtest user credentials
	var nodeuser admindb.ContainerUser
	nodeuser, err = admindb.GetContainerUser(dbConn, nodename, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		return results, err
	}

	//get db connection
	dbConn2, err = util.GetMonitoringConnection(host, CPMTEST_USER, pgport.Value, CPMTEST_DB, nodeuser.Passwd)
	if err != nil {
		logit.Error.Println("loadtest connection error:" + err.Error())
		return results, err
	}
	defer dbConn2.Close()

	start := time.Now()

	//inserts
	results[0].Operation = "inserts"
	results[0].Count = writes
	for i = 0; i < writes; i++ {
		query = fmt.Sprintf("insert into loadtest ( id, name ) values ( %d, 'this is a row for load test') returning id ", i)
		err = dbConn2.QueryRow(query).Scan(&id)
		switch {
		case err != nil:
			logit.Error.Println("loadtest insert error:" + err.Error())
			return results, err
		}
	}

	results[0].Results = time.Since(start).String()

	start = time.Now()

	//selects
	results[1].Operation = "selects"
	results[1].Count = writes
	for i = 0; i < writes; i++ {
		err = dbConn2.QueryRow(fmt.Sprintf("select name from loadtest where id=%d", i)).Scan(&name)
		switch {
		case err == sql.ErrNoRows:
			logit.Error.Println("no row with that id")
			return results, err
		case err != nil:
			logit.Error.Println(err.Error())
			return results, err
		}
	}

	results[1].Results = time.Since(start).String()

	start = time.Now()

	//updates
	results[2].Operation = "updates"
	results[2].Count = writes
	for i = 0; i < writes; i++ {
		query = fmt.Sprintf("update loadtest set ( name ) = ('howdy' ) where id = %d returning id", i)
		err = dbConn2.QueryRow(query).Scan(&id)
		switch {
		case err != nil:
			return results, err
		}
	}

	results[2].Results = time.Since(start).String()

	start = time.Now()

	//deletes
	results[3].Operation = "deletes"
	results[3].Count = writes
	for i = 0; i < writes; i++ {
		query = fmt.Sprintf("delete from loadtest where id=%d returning id", i)
		err := dbConn2.QueryRow(query).Scan(&id)
		switch {
		case err != nil:
			return results, err
		}

	}

	results[3].Results = time.Since(start).String()

	return results, nil
}

func MonitorStatements(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("MonitorStatements: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	container, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("MonitorStatements:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var host = container.Name
	if KubeEnv {
		host = container.Name + "-db"
	}

	//fetch cpmtest user credentials
	var nodeuser admindb.ContainerUser
	nodeuser, err = admindb.GetContainerUser(dbConn, container.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")

	dbConn2, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, nodeuser.Passwd)
	defer dbConn2.Close()

	//get the list of databases
	databases := make([]string, 0)
	var rows *sql.Rows

	rows, err = dbConn2.Query(
		"SELECT datname from pg_database where datname not in ('template1', 'template0') order by datname")
	if err != nil {
		logit.Error.Println("MonitorStatements:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var datname string
	for rows.Next() {
		if err = rows.Scan(
			&datname,
		); err != nil {
			logit.Error.Println("MonitorContainerSettings:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		databases = append(databases, datname)
	}
	if err = rows.Err(); err != nil {
		logit.Error.Println("MonitorStatements:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows.Close()

	statements := make([]PostgresStatement, 0)

	for i := range databases {
		//get a database connection to a specific database
		dbConn3, err := util.GetMonitoringConnection(host, databases[i], pgport.Value, CPMTEST_USER, nodeuser.Passwd)
		defer dbConn3.Close()

		rows, err = dbConn3.Query(
			"SELECT query, calls, total_time, rows, to_char(coalesce(100.0 * shared_blks_hit / nullif(shared_blks_hit + shared_blks_read, 0), -1), '999D99') AS hit_percent FROM pg_stat_statements ORDER BY total_time DESC LIMIT 5")
		if err != nil {
			logit.Error.Println("MonitorStatements:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		for rows.Next() {
			stat := PostgresStatement{}
			stat.Database = databases[i]
			if err = rows.Scan(
				&stat.Query,
				&stat.Calls,
				&stat.TotalTime,
				&stat.Rows,
				&stat.HitPercent,
			); err != nil {
				logit.Error.Println("MonitorContainerSettings:" + err.Error())
				rest.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			statements = append(statements, stat)
		}
		if err = rows.Err(); err != nil {
			logit.Error.Println("MonitorStatements:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&statements)
}

func BadgerGenerate(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("BadgerGenerate: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	container, err := admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("BadgerGenerate:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var host = container.Name
	if KubeEnv {
		host = container.Name + "-db"
	}

	//send the container a pg_controldata command
	var cdout cpmcontainerapi.BadgerGenerateResponse
	cdout, err = cpmcontainerapi.BadgerGenerateClient(host)
	if err != nil {
		logit.Error.Println("BadgerGenerate:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logit.Info.Println(cdout.Output)

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&cdout)
}
