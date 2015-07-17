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
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/sec"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
)

var secimpl sec.SecInterface

func init() {
	secimpl = sec.DefaultSec{}
}

type LoginToken struct {
	Contents string
}

type ChgPassword struct {
	Username string
	Password string
	Token    string
}

func Login(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	ID := r.PathParam("ID")
	PSW := r.PathParam("PSW")
	if ID == "" || PSW == "" {
		logit.Error.Println("Login: ID or PSW not supplied")
		rest.Error(w, "ID or PSW not supplied", http.StatusBadRequest)
	}

	logit.Info.Println("Login: called")

	tokenContents, err := secimpl.Login(dbConn, ID, PSW)
	if err != nil {
		logit.Error.Println("Login: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
	token := LoginToken{tokenContents}
	logit.Info.Println("sending back token " + token.Contents)
	w.WriteJson(&token)
}

func Logout(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	token := r.PathParam("Token")
	if token == "" {
		logit.Error.Println("Logout: Token not supplied")
		rest.Error(w, "Token not supplied", http.StatusBadRequest)
	}

	err = secimpl.Logout(dbConn, token)
	if err != nil {
		logit.Error.Println("Logout: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func UpdateUser(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("UpdateUser: in UpdateUser")
	user := sec.User{}
	err = r.DecodeJsonPayload(&user)
	if err != nil {
		logit.Error.Println("UpdateUser: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logit.Info.Println("UpdateUser: Name=" + user.Name)
	logit.Info.Println("UpdateUser: token=" + user.Token)
	err = secimpl.Authorize(dbConn, user.Token, "perm-user")
	if err != nil {
		logit.Error.Println("UpdateUser: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = secimpl.UpdateUser(dbConn, user)
	if err != nil {
		logit.Error.Println("UpdateUser: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AddUser(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("AddUser: in AddUser")
	user := sec.User{}
	err = r.DecodeJsonPayload(&user)
	if err != nil {
		logit.Error.Println("AddUser: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, user.Token, "perm-user")
	if err != nil {
		logit.Error.Println("UpdateUser: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = secimpl.AddUser(dbConn, user)
	if err != nil {
		logit.Error.Println("AddUser: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetUser(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetUser: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("GetUser: error User ID required")
		rest.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

func GetAllUsers(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllUsers: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	usersList, err := secimpl.GetAllUsers(dbConn)
	if err != nil {
		logit.Error.Println("GetAllUsers: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range usersList {
		logit.Info.Println("GetAllUsers: secimpl.GetAllUsers userName=" + usersList[i].Name)
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&usersList)
}

func DeleteUser(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-user")
	if err != nil {
		logit.Error.Println("DeleteUser: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("DeleteUser: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	err = secimpl.DeleteUser(dbConn, ID)
	if err != nil {
		logit.Error.Println("DeleteUser: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func UpdateRole(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("UpdateRole: in UpdateRole")
	role := sec.Role{}
	err = r.DecodeJsonPayload(&role)
	if err != nil {
		logit.Error.Println("UpdateRole: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, role.Token, "perm-user")
	if err != nil {
		logit.Error.Println("GetAllRoles: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = secimpl.UpdateRole(dbConn, role)
	if err != nil {
		logit.Error.Println("UpdateRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AddRole(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("AddRole: in AddRole")
	role := sec.Role{}
	err = r.DecodeJsonPayload(&role)
	if err != nil {
		logit.Error.Println("AddRole: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, role.Token, "perm-user")
	if err != nil {
		logit.Error.Println("GetAllRoles: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = secimpl.AddRole(dbConn, role)
	if err != nil {
		logit.Error.Println("AddRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func DeleteRole(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-user")
	if err != nil {
		logit.Error.Println("GetAllRoles: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if r.PathParam("ID") == "" {
		logit.Error.Println("DeleteRole: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	err = secimpl.DeleteRole(dbConn, r.PathParam("ID"))
	if err != nil {
		logit.Error.Println("DeleteRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetAllRoles(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllRoles: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var roles []sec.Role
	roles, err = secimpl.GetAllRoles(dbConn)
	if err != nil {
		logit.Error.Println("GetAllRoles: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&roles)
}

func GetRole(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetRole: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	Name := r.PathParam("Name")
	if Name == "" {
		logit.Error.Println("GetRole: error Name required")
		rest.Error(w, "Name required", http.StatusBadRequest)
		return
	}

	var role sec.Role
	role, err = secimpl.GetRole(dbConn, Name)
	if err != nil {
		logit.Error.Println("GetRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&role)
}

func ChangePassword(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("ChangePassword: in ChangePassword")
	changePass := ChgPassword{}
	err = r.DecodeJsonPayload(&changePass)
	if err != nil {
		logit.Error.Println("ChangePassword: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, changePass.Token, "perm-read")
	if err != nil {
		logit.Error.Println("ChangePassword: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var sameUser bool
	sameUser, err = secimpl.CompareUserToToken(dbConn, changePass.Username, changePass.Token)
	if err != nil {
		logit.Error.Println("ChangePassword: compare UserToToken error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//enforce perm-user if the username is not the same as the token's
	//username (e.g. bob tries to change larry's password)
	if !sameUser {
		err = secimpl.Authorize(dbConn, changePass.Token, "perm-user")
		if err != nil {
			logit.Error.Println("ChangePassword: authorize error " + err.Error())
			rest.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}

	err = secimpl.ChangePassword(dbConn, changePass.Username, changePass.Password)
	if err != nil {
		logit.Error.Println("ChangePassword: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}
