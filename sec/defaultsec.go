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
	"errors"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/rs/xid"
)

//this is the default security implementation for CPM, others
//can be swapped in as required if they implement the
//security interface as defined in secinterface.go

type DefaultSec struct {
}

// Login perform a login using a password and user id returning the security token if successful
func (d DefaultSec) Login(dbConn *sql.DB, id string, psw string) (string, error) {
	logit.Info.Println("DefaultSec.Login")
	var uuid string
	var err error
	var user User
	var unencryptedPsw string
	user, err = DBGetUser(dbConn, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logit.Error.Println(err.Error())
			return "", errors.New("user not found")
		} else {
			logit.Error.Println(err.Error())
			return "", err
		}
	}
	logit.Info.Println("Login checkpoint 1")

	unencryptedPsw, err = DecryptPassword(user.Password)
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}

	if unencryptedPsw != psw {
		return "", errors.New("incorrect password")
	}
	logit.Info.Println("Login checkpoint 2")

	/**
	uuid, err = newUUID()
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}
	*/
	guid := xid.New()
	uuid = guid.String()

	logit.Info.Println("Login checkpoint 3")
	//register the session
	err = DBAddSession(dbConn, uuid, id)
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}

	logit.Info.Println("secimpl Login returning uuid " + uuid)
	return uuid, nil
}

// Logout logout the user using the security token
func (d DefaultSec) Logout(dbConn *sql.DB, uuid string) error {
	logit.Info.Println("DefaultSec.Logout")
	err := DBDeleteSession(dbConn, uuid)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	logit.Info.Println("DefaultSec.Logout ok for " + uuid)
	return nil
}

// UpdateUser update the user object
func (d DefaultSec) UpdateUser(dbConn *sql.DB, user User) error {
	logit.Info.Println("DefaultSec.UpdateUser")
	err := DBUpdateUser(dbConn, user)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	return nil
}

// AddUser create a new user object
func (d DefaultSec) AddUser(dbConn *sql.DB, user User) error {
	logit.Info.Println("DefaultSec.AddUser")
	encryptedPsw, err := EncryptPassword(user.Password)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	user.Password = encryptedPsw

	err = DBAddUser(dbConn, user)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	return nil
}

// GetUser return a given user by ID
func (d DefaultSec) GetUser(dbConn *sql.DB, id string) (User, error) {
	logit.Info.Println("DefaultSec.GetUser id=" + id)
	user, err := DBGetUser(dbConn, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logit.Error.Println("no user found " + id)
			return user, err
		} else {
			logit.Error.Println(err.Error())
			return user, err
		}
	}
	return user, nil
}

// GetAllUsers return a list of all users
func (d DefaultSec) GetAllUsers(dbConn *sql.DB) ([]User, error) {
	logit.Info.Println("DefaultSec.GetAllUsers")
	var users []User
	var err error
	users, err = DBGetAllUsers(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
		return users, err
	}
	return users, err
}

// DeleteUser delete a user
func (d DefaultSec) DeleteUser(dbConn *sql.DB, id string) error {
	logit.Info.Println("DefaultSec.DeleteUser id=" + id)
	err := DBDeleteUser(dbConn, id)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	return nil
}

// UpdateRole update a role
func (d DefaultSec) UpdateRole(dbConn *sql.DB, role Role) error {
	logit.Info.Println("DefaultSec.UpdateRole")
	err := DBUpdateRole(dbConn, role)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	return nil
}

// AddRole add a role
func (d DefaultSec) AddRole(dbConn *sql.DB, role Role) error {
	logit.Info.Println("DefaultSec.AddRole")
	err := DBAddRole(dbConn, role)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	return nil
}

// DeleteRole delete a role by name
func (d DefaultSec) DeleteRole(dbConn *sql.DB, name string) error {
	logit.Info.Println("DefaultSec.DeleteRole name=" + name)
	err := DBDeleteRole(dbConn, name)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}
	return nil
}

// GetAllRoles return a list of all roles
func (d DefaultSec) GetAllRoles(dbConn *sql.DB) ([]Role, error) {
	logit.Info.Println("DefaultSec.GetAllRoles")
	roles := []Role{}
	var err error
	roles, err = DBGetRoles(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
		return roles, err
	}

	return roles, nil
}

// GetRole return a role by name
func (d DefaultSec) GetRole(dbConn *sql.DB, name string) (Role, error) {
	logit.Info.Println("DefaultSec.GetRole Name=" + name)
	permissions := make(map[string]string)
	permissions["perm1"] = "perm1 desc"
	role := Role{}
	return role, nil
}

// LogRole print to stdout a role
func (d DefaultSec) LogRole(role Role) {
	logit.Info.Println("***role***")
	logit.Info.Println("role=" + role.Name + " Selected=" + fmt.Sprintf("%t", role.Selected))
	for k, v := range role.Permissions {
		logit.Info.Println("perm=" + k + " desc=" + v.Description + " selected=" + fmt.Sprintf("%t", v.Selected))
	}
	logit.Info.Println("******")

}

// LogUser print to stdout a user
func (d DefaultSec) LogUser(user User) {
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

// Authorize perform an authorization based on a security token and requested action
func (d DefaultSec) Authorize(dbConn *sql.DB, token string, action string) error {
	var err error

	if token == "" {
		return errors.New("user login required")
	}

	var session Session
	session, err = DBGetSession(dbConn, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("expired user session, new user login required")
		} else {
			logit.Info.Println("error in DefaultSec.Authorize: " + err.Error())
			return errors.New("error authorizing user session")
		}
	}

	//var user User
	var user User
	user, err = DBGetUser(dbConn, session.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("security error, contact CPM admin")
		} else {
			logit.Error.Println(err.Error())
			return errors.New("error authorizing user session - u")
		}
	}

	//authorize all read-only actions
	//logit.Info.Println("Authorize:  action=[" + action + "]")
	if action == "perm-read" {
		return nil
	}

	//look at selected roles and permissions and see if we have a match
	found := false
	for _, r := range user.Roles {
		if r.Selected {
			for i, p := range r.Permissions {
				if p.Selected {
					if i == action {
						found = true
					}
				}
			}
		}
	}

	if !found {
		return errors.New("unauthorized action")
	}

	return err
}

// ChangePassword change a users password
func (d DefaultSec) ChangePassword(dbConn *sql.DB, username string, newpass string) error {
	encryptedPsw, err := EncryptPassword(newpass)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	err = DBUpdatePassword(dbConn, username, encryptedPsw)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	return nil
}

// CompareUserToToken test to see if a token matches a user id
func (d DefaultSec) CompareUserToToken(dbConn *sql.DB, username string, token string) (bool, error) {
	var err error
	var session Session
	session, err = DBGetSession(dbConn, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("expired user session, new user login required")
		} else {
			logit.Error.Println(err.Error())
			return false, err
		}
	}

	//var user User
	var user User
	user, err = DBGetUser(dbConn, session.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("security error, contact CPM admin 2")
		} else {
			logit.Error.Println(err.Error())
			return false, err
		}
	}

	logit.Info.Println("comparing [" + username + "] to [" + user.Name + "]")
	if username == user.Name {
		return true, nil
	}

	return false, nil
}
