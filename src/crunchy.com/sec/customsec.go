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

package sec

import (
	"github.com/golang/glog"
)

//this is a place holder for any future security implementations
//such as one to LDAP or a custom security solution

type CustomSec struct {
}

func (d CustomSec) Login(id string, psw string) (string, error) {
	glog.Infoln("CustomSec.Login")
	return "", nil
}

func (d CustomSec) Logout(id string) error {
	glog.Infoln("CustomSec.Logout")
	return nil
}

func (d CustomSec) UpdateUser(user User) error {
	glog.Infoln("CustomSec.UpdateUser")
	return nil
}

func (d CustomSec) AddUser(user User) error {
	glog.Infoln("CustomSec.AddUser")
	return nil
}

func (d CustomSec) GetUser(id string) (User, error) {
	user := User{Name: "myname", Password: "mypass"}
	glog.Infoln("CustomSec.GetUser id=" + id)
	return user, nil
}

func (d CustomSec) GetAllUsers() ([]User, error) {
	user := User{Name: "myname", Password: "mypass"}
	users := []User{user}
	glog.Infoln("CustomSec.GetAllUsers")
	return users, nil
}

func (d CustomSec) DeleteUser(id string) error {
	glog.Infoln("CustomSec.DeleteUser id=" + id)
	return nil
}

func (d CustomSec) UpdateRole(role Role) error {
	glog.Infoln("CustomSec.UpdateRole")
	return nil
}

func (d CustomSec) AddRole(role Role) error {
	glog.Infoln("CustomSec.AddRole")
	return nil
}

func (d CustomSec) DeleteRole(name string) error {
	glog.Infoln("CustomSec.DeleteRole name=" + name)
	return nil
}

func (d CustomSec) GetAllRoles() ([]Role, error) {
	glog.Infoln("CustomSec.GetAllRoles")
	roles := []Role{}
	return roles, nil
}

func (d CustomSec) GetRole(name string) (Role, error) {
	glog.Infoln("CustomSec.GetRole Name=" + name)
	permissions := make(map[string]string)
	permissions["perm1"] = "perm1 desc"
	permissions["perm2"] = "perm2 desc"
	role := Role{}
	return role, nil
}

func (d CustomSec) LogRole(role Role) {
}

func (d CustomSec) LogUser(user User) {
}

func (d CustomSec) Authorize(token string, action string) error {
	var err error
	return err
}
func (d CustomSec) ChangePassword(username string, newpass string) error {
	var err error
	return err
}

func (d CustomSec) CompareUserToToken(string, string) (bool, error) {
	var err error
	return false, err
}
