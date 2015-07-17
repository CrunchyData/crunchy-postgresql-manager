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
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

func AddContainerUser(w rest.ResponseWriter, r *rest.Request) {
	postMsg := NodeUser{}
	err := r.DecodeJsonPayload(&postMsg)
	if err != nil {
		logit.Error.Println("AddContainerUser: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var dbConn *sql.DB
	dbConn, err = util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, postMsg.Token, "perm-user")
	if err != nil {
		logit.Error.Println("AddSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if postMsg.ID == "" {
		logit.Error.Println("AddContainerUser: error node ID required")
		rest.Error(w, "ID required", 400)
		return
	}

	if postMsg.Rolname == "" {
		logit.Error.Println("AddContainerUser: error node Rolname required")
		rest.Error(w, "Rolname required", 400)
		return
	}
	if postMsg.Passwd == "" {
		logit.Error.Println("AddContainerUser: error node Passwd required")
		rest.Error(w, "Passwd required", 400)
		return
	}

	//create user on the container
	//get container info
	node, err := admindb.GetContainer(dbConn, postMsg.ID)
	if err != nil {
		logit.Error.Println("AddContainUser: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get connection to container's database
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

	var SUPERUSER = ""
	var INHERIT = ""
	var CREATEROLE = ""
	var CREATEDB = ""
	var LOGIN = ""
	var REPLICATION = ""

	logit.Info.Println("Rolsuper is " + strconv.FormatBool(postMsg.Rolsuper))
	if postMsg.Rolsuper {
		SUPERUSER = "SUPERUSER"
	}
	if postMsg.Rolinherit {
		INHERIT = "INHERIT"
	}
	if postMsg.Rolcreaterole {
		CREATEROLE = "CREATEROLE"
	}
	if postMsg.Rolcreatedb {
		CREATEDB = "CREATEDB"
	}
	if postMsg.Rollogin {
		LOGIN = "LOGIN"
	}
	if postMsg.Rolreplication {
		REPLICATION = "REPLICATION"
	}
	query := "create user " + postMsg.Rolname + " " +
		SUPERUSER + " " +
		INHERIT + " " +
		CREATEROLE + " " +
		CREATEDB + " " +
		LOGIN + " " +
		REPLICATION + " " +
		"PASSWORD '" + postMsg.Passwd + "'"

	logit.Info.Println(query)

	_, err = dbConn2.Query(query)
	if err != nil {
		logit.Error.Println("AddContainerUser:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//create user in the admin db
	dbuser := admindb.ContainerUser{}
	dbuser.Containername = node.Name
	dbuser.Passwd = postMsg.Passwd
	dbuser.Rolname = postMsg.Rolname

	result, err := admindb.AddContainerUser(dbConn, dbuser)
	if err != nil {
		logit.Error.Println("AddContainerUser: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	logit.Info.Printf("AddContainerUser: new ID %d\n", result)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func DeleteContainerUser(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-backup")
	if err != nil {
		logit.Error.Println("DeleteContainerUser: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ContainerID := r.PathParam("ContainerID")

	if ContainerID == "" {
		rest.Error(w, "ContainerID required", 400)
		return
	}
	rolname := r.PathParam("Rolname")

	if rolname == "" {
		rest.Error(w, "Rolname required", 400)
		return
	}

	//get node info
	node, err := admindb.GetContainer(dbConn, ContainerID)
	if err != nil {
		logit.Error.Println("AddContainUser: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = admindb.DeleteContainerUser(dbConn, node.Name, rolname)
	if err != nil {
		logit.Error.Println("DeleteContainerUser: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//get connection to container's database
	var host = node.Name
	if KubeEnv {
		host = node.Name + "-db"
	}

	//fetch cpmtest user credentials
	var cpmuser admindb.ContainerUser
	cpmuser, err = admindb.GetContainerUser(dbConn, node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")

	dbConn2, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, cpmuser.Passwd)
	defer dbConn2.Close()

	query := "drop role " + rolname

	logit.Info.Println(query)

	_, err = dbConn2.Query(query)
	if err != nil {
		logit.Error.Println("DeleteContainerUser:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func GetContainerUser(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetContainerUser: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ContainerID := r.PathParam("ContainerID")
	if ContainerID == "" {
		rest.Error(w, "ContainerID required", 400)
		return
	}

	Rolname := r.PathParam("Rolname")
	if Rolname == "" {
		rest.Error(w, "Rolname required", 400)
		return
	}

	//get container info
	node, err := admindb.GetContainer(dbConn, ContainerID)
	if err != nil {
		logit.Error.Println("AddContainUser: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get connection to container's database
	var host = node.Name
	if KubeEnv {
		host = node.Name + "-db"
	}

	//fetch  user credentials
	var nodeuser admindb.ContainerUser
	nodeuser, err = admindb.GetContainerUser(dbConn, node.Name, Rolname)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fetch cpmtest user credentials
	var cpmuser admindb.ContainerUser
	cpmuser, err = admindb.GetContainerUser(dbConn, node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")

	dbConn2, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, cpmuser.Passwd)
	defer dbConn2.Close()

	query := "select rolname::text, rolsuper::text, rolinherit::text, rolcreaterole::text, rolcreatedb::text, rolcatupdate::text, rolcanlogin::text, rolreplication::text from pg_roles where rolname = '" + Rolname + "' order by rolname"

	logit.Info.Println(query)

	err = dbConn2.QueryRow(query).Scan(
		&nodeuser.Rolname,
		&nodeuser.Rolsuper,
		&nodeuser.Rolinherit,
		&nodeuser.Rolcreaterole,
		&nodeuser.Rolcreatedb,
		&nodeuser.Rolcatupdate,
		&nodeuser.Rolcanlogin,
		&nodeuser.Rolreplication)
	if err != nil {
		logit.Error.Println("GetContainerUser:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(&nodeuser)

}

func GetAllUsersForContainer(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllUsersForContainer: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ID := r.PathParam("ID")

	if ID == "" {
		rest.Error(w, "ID required", 400)
		return
	}

	//get container info
	var node admindb.Container
	node, err = admindb.GetContainer(dbConn, ID)
	if err != nil {
		logit.Error.Println("GetAllUsersForContainer: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get connection to container's database
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
	var dbConn2 *sql.DB
	dbConn2, err = util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, nodeuser.Passwd)
	defer dbConn2.Close()

	users := make([]admindb.ContainerUser, 0)

	//query results
	var rows *sql.Rows

	rows, err = dbConn2.Query("select rolname::text, rolsuper::text, rolinherit::text, rolcreaterole::text, rolcreatedb::text, rolcatupdate::text, rolcanlogin::text, rolreplication::text from pg_roles order by rolname")
	if err != nil {
		logit.Error.Println("GetAllUsersForContainer:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	for rows.Next() {
		user := admindb.ContainerUser{}
		if err = rows.Scan(
			&user.Rolname,
			&user.Rolsuper,
			&user.Rolinherit,
			&user.Rolcreaterole,
			&user.Rolcreatedb,
			&user.Rolcatupdate,
			&user.Rolcanlogin,
			&user.Rolreplication,
		); err != nil {
			logit.Error.Println("GetAllUsersForContainer:" + err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user.Containername = node.Name
		user.ContainerID = node.ID
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		logit.Error.Println("GetAllUsersForContainer:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&users)

}

func UpdateContainerUser(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	postMsg := NodeUser{}
	err = r.DecodeJsonPayload(&postMsg)
	if err != nil {
		logit.Error.Println("UpdateContainerUser: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, postMsg.Token, "perm-user")
	if err != nil {
		logit.Error.Println("UpdateContainerUser: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if postMsg.ID == "" {
		logit.Error.Println("UpdateContainerUser: error node ID required")
		rest.Error(w, "ID required", 400)
		return
	}

	if postMsg.Rolname == "" {
		logit.Error.Println("UpdateContainerUser: error node Rolname required")
		rest.Error(w, "Rolname required", 400)
		return
	}

	//create user on the container
	//get container info
	var node admindb.Container
	node, err = admindb.GetContainer(dbConn, postMsg.ID)
	if err != nil {
		logit.Error.Println("AddContainUser: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if postMsg.Passwd == "" {
	} else {
		//update the password
	}

	//get connection to container's database
	var host = node.Name
	if KubeEnv {
		host = node.Name + "-db"
	}

	//fetch cpmtest user credentials
	var cpmuser admindb.ContainerUser
	cpmuser, err = admindb.GetContainerUser(dbConn, node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting(dbConn, "PG-PORT")
	var dbConn2 *sql.DB

	dbConn2, err = util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, cpmuser.Passwd)
	defer dbConn2.Close()

	var SUPERUSER = "SUPERUSER"
	var INHERIT = "INHERIT"
	var CREATEROLE = "CREATEROLE"
	var CREATEDB = "CREATEDB"
	var LOGIN = "LOGIN"
	var REPLICATION = "REPLICATION"

	logit.Info.Println("Rolsuper is " + strconv.FormatBool(postMsg.Rolsuper))
	if !postMsg.Rolsuper {
		SUPERUSER = "NOSUPERUSER"
	}

	if !postMsg.Rolinherit {
		INHERIT = "NOINHERIT"
	}

	if !postMsg.Rolcreaterole {
		CREATEROLE = "NOCREATEROLE"
	}
	if !postMsg.Rolcreatedb {
		CREATEDB = "NOCREATEDB"
	}

	if !postMsg.Rollogin {
		LOGIN = "NOLOGIN"
	}
	if !postMsg.Rolreplication {
		REPLICATION = "NOREPLICATION"
	}

	query := "alter user " + postMsg.Rolname + " " +
		SUPERUSER + " " +
		INHERIT + " " +
		CREATEROLE + " " +
		CREATEDB + " " +
		LOGIN + " " +
		REPLICATION + " "

	if postMsg.Passwd != "" {
		query = query + " PASSWORD '" + postMsg.Passwd + "'"
	}

	logit.Info.Println(query)

	_, err = dbConn2.Query(query)
	if err != nil {
		logit.Error.Println("UpdateContainerUser:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if postMsg.Passwd != "" {
		//update user's password
		dbuser := admindb.ContainerUser{}
		dbuser.Containername = node.Name
		dbuser.Passwd = postMsg.Passwd
		dbuser.Rolname = postMsg.Rolname

		err = admindb.UpdateContainerUser(dbConn, dbuser)
		if err != nil {
			logit.Error.Println("UpdateContainerUser: " + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}
