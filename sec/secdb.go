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
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	_ "github.com/lib/pq"
	"log"
)

var pguser = "postgres"

var pghost = "127.0.0.1"

//var pghost = "172.17.0.3"
var pgport = "5432"
var db = "clusteradmin"
var dbConn *sql.DB

func init() {
	//logit.Info.Println("secdb:init: called to open dbConn")
	var err error
	dbConn, err = sql.Open("postgres", "sslmode=disable user="+pguser+" host="+pghost+" port="+pgport+" dbname="+db)
	if err != nil {
		log.Fatal(err)
	}
}

func DBGetUser(Name string) (User, error) {
	user := User{}
	//logit.Info.Println("secdb:GetUser: called name=" + Name)
	var rows *sql.Rows

	roles, err := DBGetRoles()
	if err != nil {
		return user, err
	}
	user.Roles = make(map[string]Role)
	for i := range roles {
		user.Roles[roles[i].Name] = roles[i]
	}

	user.Name = Name
	err = dbConn.QueryRow(fmt.Sprintf("select password, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from secuser where name='%s'", Name)).Scan(&user.Password, &user.UpdateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Error.Println("DBGetUser:no user with that name")
		return user, err
	case err != nil:
		logit.Error.Println("DBGetuser:Get User:" + err.Error())
		return user, err
	default:
	}

	queryStr := fmt.Sprintf("select secuserrole.role  from secuser, secuserrole where secuser.name = secuserrole.username and secuser.name = '%s'", Name)

	rows, err = dbConn.Query(queryStr)
	defer rows.Close()

	if err != nil {
		return user, err
	}

	var roleName string
	var role Role
	rolesfound := false

	for rows.Next() {
		if err = rows.Scan(&roleName); err != nil {
			return user, err
		}
		rolesfound = true
		role = user.Roles[roleName]
		role.Selected = true
		user.Roles[roleName] = role
		//logit.Info.Println("setting role " + roleName + " to true for User " + Name)
	}
	if err = rows.Err(); err != nil {
		return user, err
	}

	if rolesfound == false {
		logit.Error.Println("no roles found for user " + Name)
	}

	return user, nil
}

func DBGetRole(Name string) (Role, error) {
	role := Role{}
	role.Selected = false
	//logit.Info.Println("secdb:GetRole: called")
	var rows *sql.Rows
	var err error

	//set list of permissions for this role
	//to default set and set selected to false
	var perms []Permission
	perms, err = DBGetPermissions()
	if err != nil {
		logit.Error.Println("error in DBGetRole:GetPermissions")
		return role, err
	}
	role.Permissions = make(map[string]Permission)

	for i := range perms {
		perms[i].Selected = false
		role.Permissions[perms[i].Name] = perms[i]
	}

	//query the selected permissions for this role
	queryStr := fmt.Sprintf("select secroleperm.role, secperm.name, secperm.description from secroleperm, secperm where secroleperm.perm = secperm.name and secroleperm.role = '%s'", Name)

	rows, err = dbConn.Query(queryStr)

	if err != nil {
		return role, err
	}
	defer rows.Close()

	var roleName string
	var permName string
	var permDescription string
	var perm Permission

	for rows.Next() {
		if err = rows.Scan(
			&roleName,
			&permName,
			&permDescription); err != nil {
			return role, err
		}
		perm = Permission{}
		perm.Name = permName
		perm.Selected = true
		perm.Description = permDescription
		role.Permissions[permName] = perm
		//logit.Info.Println("setting perm " + permName + " to true for role " + roleName)
	}
	if err = rows.Err(); err != nil {
		return role, err
	}
	return role, nil
}

func DBGetPermissions() ([]Permission, error) {
	slice := []Permission{}
	//logit.Info.Println("secdb:GetPermissions: called")
	var rows *sql.Rows
	var err error

	queryStr := fmt.Sprintf("select name, description from secperm order by name")

	rows, err = dbConn.Query(queryStr)

	if err != nil {
		return slice, err
	}
	defer rows.Close()

	for rows.Next() {
		perm := Permission{}
		perm.Selected = false
		if err = rows.Scan(
			&perm.Name,
			&perm.Description); err != nil {
			return slice, err
		}
		slice = append(slice, perm)
	}
	if err = rows.Err(); err != nil {
		return slice, err
	}
	return slice, nil
}

func DBDeleteRole(name string) error {
	queryStr := fmt.Sprintf("delete from secrole where name='%s' returning name", name)
	//logit.Info.Println("secdb:DeleteRole:" + queryStr)

	var theName string
	err := dbConn.QueryRow(queryStr).Scan(&theName)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return err
	default:
		logit.Info.Println("secdb:DeleteRole:role " + name + " deleted ")
	}

	return nil
}

func DBDeleteUser(name string) error {
	queryStr := fmt.Sprintf("delete from secuser where name='%s' returning name", name)
	logit.Info.Println("secdb:DeleteUser:" + queryStr)

	var theName string
	err := dbConn.QueryRow(queryStr).Scan(&theName)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return err
	default:
		logit.Info.Println("secdb:DeleteUser:User " + name + " deleted ")
	}

	return nil
}

func DBAddRole(role Role) error {
	logit.Info.Println("secdb:AddRole:called")
	queryStr := fmt.Sprintf("insert into secrole( name, updatedt) values ( '%s', now()) returning name", role.Name)

	logit.Info.Println("secdb:AddRole:" + queryStr)
	var theName string
	err := dbConn.QueryRow(queryStr).Scan(&theName)
	switch {
	case err != nil:
		logit.Error.Println("secdb:AddRole:" + err.Error())
		return err
	default:
		logit.Info.Println("secdb:AddRole: role inserted " + role.Name)
	}

	for k, v := range role.Permissions {
		if v.Selected {
			err = DBAddRolePerm(role.Name, k)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DBAddUserRole(user string, role string) error {
	logit.Info.Println("secdb:AddUserRole:called")
	queryStr := fmt.Sprintf("insert into secuserrole ( username, role) values ( '%s', '%s') returning username", user, role)

	logit.Info.Println("secdb:AddUserRole:" + queryStr)
	var theUser string
	err := dbConn.QueryRow(queryStr).Scan(&theUser)
	switch {
	case err != nil:
		logit.Error.Println("secdb:AddUserRole:" + err.Error())
		return err
	default:
		logit.Info.Println("secdb:AddUserRole: inserted user=" + user + " role=" + role)
	}

	return nil
}

func DBAddRolePerm(role string, perm string) error {
	logit.Info.Println("secdb:AddRolePerm:called")
	queryStr := fmt.Sprintf("insert into secroleperm ( role, perm) values ( '%s', '%s') returning role", role, perm)

	logit.Info.Println("secdb:AddRolePerm:" + queryStr)
	var theRole string
	err := dbConn.QueryRow(queryStr).Scan(&theRole)
	switch {
	case err != nil:
		logit.Error.Println("secdb:AddRolePerm:" + err.Error())
		return err
	default:
		logit.Info.Println("secdb:AddRolePerm: inserted role=" + role + " perm=" + perm)
	}

	return nil
}

func DBAddUser(user User) error {
	logit.Info.Println("secdb:AddUser:called")
	queryStr := fmt.Sprintf("insert into secuser ( name, password, updatedt) values ( '%s', '%s', now()) returning name", user.Name, user.Password)

	logit.Info.Println("secdb:AddUser:" + queryStr)
	var theName string
	err := dbConn.QueryRow(queryStr).Scan(&theName)
	switch {
	case err != nil:
		logit.Error.Println("secdb:AddUser:" + err.Error())
		return err
	default:
		logit.Info.Println("secdb:AddUser: inserted " + user.Name)
	}
	for k, v := range user.Roles {
		if v.Selected {
			err = DBAddUserRole(user.Name, k)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DBUpdateRole(role Role) error {
	logit.Info.Println("secdb:UpdateRole:called")
	err := DBDeleteRole(role.Name)
	if err != nil {
		return err
	}

	err = DBAddRole(role)
	if err != nil {
		return err
	}

	return nil
}

func DBUpdateUser(user User) error {
	logit.Info.Println("secdb:DBUpdateUser:called")

	err := DBDeleteUser(user.Name)
	if err != nil {
		return err
	}

	err = DBAddUser(user)
	if err != nil {
		return err
	}
	return nil
}

func DBGetAllUsers() ([]User, error) {
	userList := []User{}
	//logit.Info.Println("secdb:DBGetAllUser: called")
	var rows *sql.Rows
	var err error

	queryStr := fmt.Sprintf("select name, password from secuser order by name")

	rows, err = dbConn.Query(queryStr)

	if err != nil {
		return userList, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}
		if err = rows.Scan(
			&user.Name,
			&user.Password); err != nil {
			return userList, err
		}
		if err != nil {
			logit.Error.Println("error in GetUser" + err.Error())
			return userList, err
		}
		userList = append(userList, user)
	}
	if err = rows.Err(); err != nil {
		return userList, err
	}

	var user User
	for i := range userList {
		logit.Info.Println("fetching user info for " + userList[i].Name)
		user, err = DBGetUser(userList[i].Name)
		if err != nil {
			logit.Error.Println("error" + err.Error())
			return userList, err
		}

		LogUser(user)
		userList[i].Roles = user.Roles
	}

	return userList, nil
}

func DBGetRoles() ([]Role, error) {
	slice := []Role{}
	//logit.Info.Println("secdb:GetRoles: called")
	var rows *sql.Rows
	var err error

	queryStr := fmt.Sprintf("select name from secrole order by name")

	rows, err = dbConn.Query(queryStr)

	if err != nil {
		logit.Error.Println("error in GetRoles:" + err.Error())
		return slice, err
	}
	defer rows.Close()

	for rows.Next() {
		role := Role{}
		role.Selected = false
		if err = rows.Scan(
			&role.Name); err != nil {
			return slice, err
		}
		slice = append(slice, role)
	}
	if err = rows.Err(); err != nil {
		return slice, err
	}

	var role Role
	for i := range slice {
		//logit.Info.Println("fetching role info for " + slice[i].Name)
		role, err = DBGetRole(slice[i].Name)
		if err != nil {
			logit.Error.Println("error" + err.Error())
			return slice, err
		}

		//LogPermissions(role.Permissions)
		slice[i].Permissions = role.Permissions
	}
	return slice, nil
}

func LogUser(user User) {
	logit.Info.Println("***user***")
	logit.Info.Println("user.Name=" + user.Name + " user.Password=" + user.Password)
	for k, v := range user.Roles {
		logit.Info.Println("***role***")
		logit.Info.Println("role=" + k + " Selected=" + fmt.Sprintf("%t", v.Selected))
		for i, j := range v.Permissions {
			logit.Info.Println("perm=" + i + " desc=" + j.Description + " selected=" + fmt.Sprintf("%t", j.Selected))
		}
		logit.Info.Println("******")
	}
	logit.Info.Println("******")

}
func LogPermissions(perms map[string]Permission) {
	logit.Info.Println("***log of permissions***")
	for i, j := range perms {
		logit.Info.Println("perm=" + i + " desc=" + j.Description + " selected=" + fmt.Sprintf("%t", j.Selected))
	}
	logit.Info.Println("******")

}

func DBGetSession(token string) (Session, error) {
	session := Session{}
	//logit.Info.Println("secdb:GetSession: called token=" + token)

	queryStr := fmt.Sprintf("select token, name, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS')  from secsession where token = '%s'", token)

	err := dbConn.QueryRow(queryStr).Scan(
		&session.Token,
		&session.Name,
		&session.UpdateDate)

	switch {
	case err == sql.ErrNoRows:
		logit.Error.Println("secdb:DBGetSession:no token matched")
		return session, err
	case err != nil:
		logit.Error.Println("secdb:DBGetSession:" + err.Error())
		return session, err
	default:
		//logit.Info.Println("secdb:DBGetSession: token returned is " + session.Token)
	}

	return session, nil
}

func DBAddSession(uuid string, id string) error {
	//logit.Info.Println("secdb:DBAddSession:called")
	queryStr := fmt.Sprintf("insert into secsession ( token, name, updatedt) values ( '%s', '%s', now()) returning token", uuid, id)

	//logit.Info.Println("secdb:DBAddSession:" + queryStr)
	var theToken string
	err := dbConn.QueryRow(queryStr).Scan(&theToken)
	switch {
	case err != nil:
		logit.Error.Println("secdb:DBAddSession:" + err.Error())
		return err
	default:
		logit.Info.Println("secdb:AddSession: Session inserted " + theToken)
	}

	return nil
}

func DBDeleteSession(uuid string) error {
	logit.Info.Println("secdb:DBDeleteSession:called")

	//if the uuid is not there, return
	_, err := DBGetSession(uuid)
	if err == sql.ErrNoRows {
		return nil
	}

	queryStr := fmt.Sprintf("delete from secsession where token='%s' returning token", uuid)
	logit.Info.Println("secdb:DeleteSession:" + queryStr)

	var theToken string
	err = dbConn.QueryRow(queryStr).Scan(&theToken)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return err
	default:
		logit.Info.Println("secdb:DeleteSession " + uuid + " deleted ")
	}

	return nil
}

func DBUpdatePassword(username string, password string) error {
	logit.Info.Println("UpdatePassword:called")
	queryStr := fmt.Sprintf("update secuser set ( password, updatedt) = ('%s', now()) where name = '%s' returning name", password, username)

	logit.Info.Println("UpdatePassword: str=[" + queryStr + "]")
	var theName string
	err := dbConn.QueryRow(queryStr).Scan(&theName)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("UpdatePassword:updated " + username)
	}
	return nil
}
