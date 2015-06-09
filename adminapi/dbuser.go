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
)

func AddContainerUser(w rest.ResponseWriter, r *rest.Request) {
	postMsg := NodeUser{}
	err := r.DecodeJsonPayload(&postMsg)
	if err != nil {
		logit.Error.Println("AddContainerUser: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(postMsg.Token, "perm-user")
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

	if postMsg.Usename == "" {
		logit.Error.Println("AddContainerUser: error node Usename required")
		rest.Error(w, "Usename required", 400)
		return
	}
	if postMsg.Passwd == "" {
		logit.Error.Println("AddContainerUser: error node Passwd required")
		rest.Error(w, "Passwd required", 400)
		return
	}

	//create user on the container
	//get container info
	node, err := admindb.GetContainer(postMsg.ID)
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
	nodeuser, err = admindb.GetContainerUser(node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("cpmtest password is " + nodeuser.Passwd)

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting("PG-PORT")

	dbConn, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, nodeuser.Passwd)
	defer dbConn.Close()

	query := "create user " + postMsg.Usename + " " +
		postMsg.Rolsuper + " " +
		postMsg.Rolinherit + " " +
		postMsg.Rolcreaterole + " " +
		postMsg.Rolcreatedb + " " +
		postMsg.Rollogin + " " +
		postMsg.Rolreplication + " " +
		"PASSWORD '" + postMsg.Passwd + "'"

	logit.Info.Println(query)

	_, err = dbConn.Query(query)
	if err != nil {
		logit.Error.Println("AddContainerUser:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//create user in the admin db
	dbuser := admindb.ContainerUser{}
	dbuser.Containername = node.Name
	dbuser.Passwd = postMsg.Passwd
	dbuser.Rolname = postMsg.Usename

	result, err := admindb.AddContainerUser(dbuser)
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
	err := secimpl.Authorize(r.PathParam("Token"), "perm-backup")
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
	node, err := admindb.GetContainer(ContainerID)
	if err != nil {
		logit.Error.Println("AddContainUser: " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = admindb.DeleteContainerUser(node.Name, rolname)
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
	cpmuser, err = admindb.GetContainerUser(node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting("PG-PORT")

	dbConn, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, cpmuser.Passwd)
	defer dbConn.Close()

	query := "drop role " + rolname

	logit.Info.Println(query)

	_, err = dbConn.Query(query)
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
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
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

	usename := r.PathParam("Usename")
	if usename == "" {
		rest.Error(w, "Usename required", 400)
		return
	}

	//get container info
	node, err := admindb.GetContainer(ContainerID)
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
	nodeuser, err = admindb.GetContainerUser(node.Name, usename)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fetch cpmtest user credentials
	var cpmuser admindb.ContainerUser
	cpmuser, err = admindb.GetContainerUser(node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting("PG-PORT")

	dbConn, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, cpmuser.Passwd)
	defer dbConn.Close()

	query := "select rolname::text, rolsuper::text, rolinherit::text, rolcreaterole::text, rolcreatedb::text, rolcatupdate::text, rolcanlogin::text, rolreplication::text from pg_roles where rolname = '" + usename + "' order by rolname"

	logit.Info.Println(query)

	err = dbConn.QueryRow(query).Scan(
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
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
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
	node, err := admindb.GetContainer(ID)
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
	nodeuser, err = admindb.GetContainerUser(node.Name, CPMTEST_USER)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logit.Info.Println("cpmtest password is " + nodeuser.Passwd)

	//get port
	var pgport admindb.Setting
	pgport, err = admindb.GetSetting("PG-PORT")

	dbConn, err := util.GetMonitoringConnection(host, CPMTEST_DB, pgport.Value, CPMTEST_USER, nodeuser.Passwd)
	defer dbConn.Close()

	users := make([]admindb.ContainerUser, 0)

	//query results
	var rows *sql.Rows

	rows, err = dbConn.Query("select rolname::text, rolsuper::text, rolinherit::text, rolcreaterole::text, rolcreatedb::text, rolcatupdate::text, rolcanlogin::text, rolreplication::text from pg_roles order by rolname")
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
