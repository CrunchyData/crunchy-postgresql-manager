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
	"crunchy.com/logutil"
	"crunchy.com/sec"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

var secimpl sec.SecInterface

func init() {
	logutil.Log("securitymgmt init called")
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
		logutil.Log("Login: ID or PSW not supplied")
		rest.Error(w, "ID or PSW not supplied", 400)
	}

	logutil.Log("Login: called")
	tokenContents, err := secimpl.Login(ID, PSW)
	if err != nil {
		logutil.Log("Login: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
	token := LoginToken{tokenContents}
	logutil.Log("sending back token " + token.Contents)
	w.WriteJson(&token)
}

func Logout(w rest.ResponseWriter, r *rest.Request) {
	token := r.PathParam("Token")
	if token == "" {
		logutil.Log("Logout: Token not supplied")
		rest.Error(w, "Token not supplied", 400)
	}

	err := secimpl.Logout(token)
	if err != nil {
		logutil.Log("Logout: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func UpdateUser(w rest.ResponseWriter, r *rest.Request) {
	logutil.Log("UpdateUser: in UpdateUser")
	user := sec.User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		logutil.Log("UpdateUser: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(user.Token, "perm-user")
	if err != nil {
		logutil.Log("UpdateUser: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	err = secimpl.UpdateUser(user)
	if err != nil {
		logutil.Log("UpdateUser: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AddUser(w rest.ResponseWriter, r *rest.Request) {
	logutil.Log("AddUser: in AddUser")
	user := sec.User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		logutil.Log("AddUser: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(user.Token, "perm-user")
	if err != nil {
		logutil.Log("UpdateUser: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	err = secimpl.AddUser(user)
	if err != nil {
		logutil.Log("AddUser: error secimpl call" + err.Error())
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
		logutil.Log("GetUser: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logutil.Log("GetUser: error User ID required")
		rest.Error(w, "User ID required", 400)
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
		logutil.Log("GetAllUsers: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	usersList, err := secimpl.GetAllUsers()
	if err != nil {
		logutil.Log("GetAllUsers: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	for i := range usersList {
		logutil.Log("GetAllUsers: secimpl.GetAllUsers userName=" + usersList[i].Name)
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(&usersList)
}

func DeleteUser(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-user")
	if err != nil {
		logutil.Log("DeleteUser: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logutil.Log("DeleteUser: error ID required")
		rest.Error(w, "ID required", 400)
		return
	}
	err = secimpl.DeleteUser(ID)
	if err != nil {
		logutil.Log("DeleteUser: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func UpdateRole(w rest.ResponseWriter, r *rest.Request) {
	logutil.Log("UpdateRole: in UpdateRole")
	role := sec.Role{}
	err := r.DecodeJsonPayload(&role)
	if err != nil {
		logutil.Log("UpdateRole: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(role.Token, "perm-user")
	if err != nil {
		logutil.Log("GetAllRoles: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	err = secimpl.UpdateRole(role)
	if err != nil {
		logutil.Log("UpdateRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AddRole(w rest.ResponseWriter, r *rest.Request) {
	logutil.Log("AddRole: in AddRole")
	role := sec.Role{}
	err := r.DecodeJsonPayload(&role)
	if err != nil {
		logutil.Log("AddRole: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(role.Token, "perm-user")
	if err != nil {
		logutil.Log("GetAllRoles: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	err = secimpl.AddRole(role)
	if err != nil {
		logutil.Log("AddRole: error secimpl call" + err.Error())
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
		logutil.Log("GetAllRoles: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	if r.PathParam("ID") == "" {
		logutil.Log("DeleteRole: error ID required")
		rest.Error(w, "ID required", 400)
		return
	}

	err = secimpl.DeleteRole(r.PathParam("ID"))
	if err != nil {
		logutil.Log("DeleteRole: error secimpl call" + err.Error())
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
		logutil.Log("GetAllRoles: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}
	var roles []sec.Role
	roles, err = secimpl.GetAllRoles()
	if err != nil {
		logutil.Log("GetAllRoles: error secimpl call" + err.Error())
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
		logutil.Log("GetRole: validate token error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	Name := r.PathParam("Name")
	if Name == "" {
		logutil.Log("GetRole: error Name required")
		rest.Error(w, "Name required", 400)
		return
	}

	var role sec.Role
	role, err = secimpl.GetRole(Name)
	if err != nil {
		logutil.Log("GetRole: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.WriteJson(&role)
}

func ChangePassword(w rest.ResponseWriter, r *rest.Request) {
	logutil.Log("ChangePassword: in ChangePassword")
	changePass := ChgPassword{}
	err := r.DecodeJsonPayload(&changePass)
	if err != nil {
		logutil.Log("ChangePassword: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(changePass.Token, "perm-read")
	if err != nil {
		logutil.Log("ChangePassword: authorize error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	var sameUser bool
	sameUser, err = secimpl.CompareUserToToken(changePass.Username, changePass.Token)
	if err != nil {
		logutil.Log("ChangePassword: compare UserToToken error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return
	}

	//enforce perm-user if the username is not the same as the token's
	//username (e.g. bob tries to change larry's password)
	if !sameUser {
		err = secimpl.Authorize(changePass.Token, "perm-user")
		if err != nil {
			logutil.Log("ChangePassword: authorize error " + err.Error())
			rest.Error(w, err.Error(), 400)
			return
		}
	}

	err = secimpl.ChangePassword(changePass.Username, changePass.Password)
	if err != nil {
		logutil.Log("ChangePassword: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}
