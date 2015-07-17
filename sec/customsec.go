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
	"database/sql"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
)

//this is a place holder for any future security implementations
//such as one to LDAP or a custom security solution

type CustomSec struct {
}

func (d CustomSec) Login(dbConn *sql.DB, id string, psw string) (string, error) {
	logit.Info.Println("CustomSec.Login")
	return "", nil
}

func (d CustomSec) Logout(dbConn *sql.DB, id string) error {
	logit.Info.Println("CustomSec.Logout")
	return nil
}

func (d CustomSec) UpdateUser(dbConn *sql.DB, user User) error {
	logit.Info.Println("CustomSec.UpdateUser")
	return nil
}

func (d CustomSec) AddUser(dbConn *sql.DB, user User) error {
	logit.Info.Println("CustomSec.AddUser")
	return nil
}

func (d CustomSec) GetUser(dbConn *sql.DB, id string) (User, error) {
	user := User{Name: "myname", Password: "mypass"}
	logit.Info.Println("CustomSec.GetUser id=" + id)
	return user, nil
}

func (d CustomSec) GetAllUsers(dbConn *sql.DB) ([]User, error) {
	user := User{Name: "myname", Password: "mypass"}
	users := []User{user}
	logit.Info.Println("CustomSec.GetAllUsers")
	return users, nil
}

func (d CustomSec) DeleteUser(dbConn *sql.DB, id string) error {
	logit.Info.Println("CustomSec.DeleteUser id=" + id)
	return nil
}

func (d CustomSec) UpdateRole(dbConn *sql.DB, role Role) error {
	logit.Info.Println("CustomSec.UpdateRole")
	return nil
}

func (d CustomSec) AddRole(dbConn *sql.DB, role Role) error {
	logit.Info.Println("CustomSec.AddRole")
	return nil
}

func (d CustomSec) DeleteRole(dbConn *sql.DB, name string) error {
	logit.Info.Println("CustomSec.DeleteRole name=" + name)
	return nil
}

func (d CustomSec) GetAllRoles(dbConn *sql.DB) ([]Role, error) {
	logit.Info.Println("CustomSec.GetAllRoles")
	roles := []Role{}
	return roles, nil
}

func (d CustomSec) GetRole(dbConn *sql.DB, name string) (Role, error) {
	logit.Info.Println("CustomSec.GetRole Name=" + name)
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

func (d CustomSec) Authorize(dbConn *sql.DB, token string, action string) error {
	var err error
	return err
}
func (d CustomSec) ChangePassword(dbConn *sql.DB, username string, newpass string) error {
	var err error
	return err
}

func (d CustomSec) CompareUserToToken(string, string) (bool, error) {
	var err error
	return false, err
}
