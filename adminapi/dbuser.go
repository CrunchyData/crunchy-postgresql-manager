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

	err = secimpl.Authorize(postMsg.Token, "perm-backup")
	if err != nil {
		logit.Error.Println("AddSchedule: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if postMsg.Containername == "" {
		logit.Error.Println("AddContainerUser: error node Containername required")
		rest.Error(w, "Containername required", 400)
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

	dbuser := admindb.ContainerUser{}
	dbuser.Containername = postMsg.Containername
	dbuser.Passwd = postMsg.Passwd
	dbuser.Usename = postMsg.Usename

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
	ID := r.PathParam("ID")

	if ID == "" {
		rest.Error(w, "NodeUser ID required", 400)
		return
	}

	err = admindb.DeleteContainerUser(ID)
	if err != nil {
		logit.Error.Println("DeleteContainerUser: " + err.Error())
		rest.Error(w, err.Error(), 400)
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
	containername := r.PathParam("Containername")
	if containername == "" {
		rest.Error(w, "Containername required", 400)
		return
	}

	usename := r.PathParam("Usename")
	if usename == "" {
		rest.Error(w, "Usename required", 400)
		return
	}

	result, err := admindb.GetContainerUser(containername, usename)
	if err != nil {
		logit.Error.Println("GetContainerUser: " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	w.WriteJson(result)

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

	rows, err = dbConn.Query("select usename, usesysid, usecreatedb, usesuper, usecatupd, userepl, coalesce(valuntil::text, '') from pg_user order by usename")
	if err != nil {
		logit.Error.Println("GetAllUsersForContainer:" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	for rows.Next() {
		user := admindb.ContainerUser{}
		if err = rows.Scan(
			&user.Usename,
			&user.Usesysid,
			&user.Usecreatedb,
			&user.Usesuper,
			&user.Usecatupd,
			&user.Userepl,
			&user.Valuntil,
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
