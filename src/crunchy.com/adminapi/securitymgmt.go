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
	"crunchy.com/sec"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/glog"
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
	ID := r.PathParam("ID")
	PSW := r.PathParam("PSW")
	if ID == "" || PSW == "" {
		glog.Errorln("Login: ID or PSW not supplied")
		rest.Error(w, "ID or PSW not supplied", http.StatusBadRequest)
	}

	glog.Infoln("Login: called")
	tokenContents, err := secimpl.Login(ID, PSW)
	if err != nil {
		glog.Errorln("Login: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
	token := LoginToken{tokenContents}
	glog.Infoln("sending back token " + token.Contents)
	w.WriteJson(&token)
}

func Logout(w rest.ResponseWriter, r *rest.Request) {
	token := r.PathParam("Token")
	if token == "" {
		glog.Errorln("Logout: Token not supplied")
		rest.Error(w, "Token not supplied", http.StatusBadRequest)
	}

	err := secimpl.Logout(token)
	if err != nil {
		glog.Errorln("Logout: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func UpdateUser(w rest.ResponseWriter, r *rest.Request) {
	glog.Infoln("UpdateUser: in UpdateUser")
	user := sec.User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		glog.Errorln("UpdateUser: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(user.Token, "perm-user")
	if err != nil {
		glog.Errorln("UpdateUser: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = secimpl.UpdateUser(user)
	if err != nil {
		glog.Errorln("UpdateUser: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AddUser(w rest.ResponseWriter, r *rest.Request) {
	glog.Infoln("AddUser: in AddUser")
	user := sec.User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		glog.Errorln("AddUser: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(user.Token, "perm-user")
	if err != nil {
		glog.Errorln("UpdateUser: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = secimpl.AddUser(user)
	if err != nil {
		glog.Errorln("AddUser: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetUser(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetUser: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		glog.Errorln("GetUser: error User ID required")
		rest.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&status)
}

func GetAllUsers(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetAllUsers: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	usersList, err := secimpl.GetAllUsers()
	if err != nil {
		glog.Errorln("GetAllUsers: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range usersList {
		glog.Infoln("GetAllUsers: secimpl.GetAllUsers userName=" + usersList[i].Name)
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&usersList)
}

func DeleteUser(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-user")
	if err != nil {
		glog.Errorln("DeleteUser: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		glog.Errorln("DeleteUser: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	err = secimpl.DeleteUser(ID)
	if err != nil {
		glog.Errorln("DeleteUser: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func UpdateRole(w rest.ResponseWriter, r *rest.Request) {
	glog.Infoln("UpdateRole: in UpdateRole")
	role := sec.Role{}
	err := r.DecodeJsonPayload(&role)
	if err != nil {
		glog.Errorln("UpdateRole: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(role.Token, "perm-user")
	if err != nil {
		glog.Errorln("GetAllRoles: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = secimpl.UpdateRole(role)
	if err != nil {
		glog.Errorln("UpdateRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AddRole(w rest.ResponseWriter, r *rest.Request) {
	glog.Infoln("AddRole: in AddRole")
	role := sec.Role{}
	err := r.DecodeJsonPayload(&role)
	if err != nil {
		glog.Errorln("AddRole: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(role.Token, "perm-user")
	if err != nil {
		glog.Errorln("GetAllRoles: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = secimpl.AddRole(role)
	if err != nil {
		glog.Errorln("AddRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func DeleteRole(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-user")
	if err != nil {
		glog.Errorln("GetAllRoles: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if r.PathParam("ID") == "" {
		glog.Errorln("DeleteRole: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	err = secimpl.DeleteRole(r.PathParam("ID"))
	if err != nil {
		glog.Errorln("DeleteRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetAllRoles(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetAllRoles: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var roles []sec.Role
	roles, err = secimpl.GetAllRoles()
	if err != nil {
		glog.Errorln("GetAllRoles: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&roles)
}

func GetRole(w rest.ResponseWriter, r *rest.Request) {
	var err error
	err = secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		glog.Errorln("GetRole: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	Name := r.PathParam("Name")
	if Name == "" {
		glog.Errorln("GetRole: error Name required")
		rest.Error(w, "Name required", http.StatusBadRequest)
		return
	}

	var role sec.Role
	role, err = secimpl.GetRole(Name)
	if err != nil {
		glog.Errorln("GetRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&role)
}

func ChangePassword(w rest.ResponseWriter, r *rest.Request) {
	glog.Infoln("ChangePassword: in ChangePassword")
	changePass := ChgPassword{}
	err := r.DecodeJsonPayload(&changePass)
	if err != nil {
		glog.Errorln("ChangePassword: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(changePass.Token, "perm-read")
	if err != nil {
		glog.Errorln("ChangePassword: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var sameUser bool
	sameUser, err = secimpl.CompareUserToToken(changePass.Username, changePass.Token)
	if err != nil {
		glog.Errorln("ChangePassword: compare UserToToken error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//enforce perm-user if the username is not the same as the token's
	//username (e.g. bob tries to change larry's password)
	if !sameUser {
		err = secimpl.Authorize(changePass.Token, "perm-user")
		if err != nil {
			glog.Errorln("ChangePassword: authorize error " + err.Error())
			rest.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}

	err = secimpl.ChangePassword(changePass.Username, changePass.Password)
	if err != nil {
		glog.Errorln("ChangePassword: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}
