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
	"crunchy.com/logutil"
	"database/sql"
	"errors"
	"fmt"
)

//this is the default security implementation for CPM, others
//can be swapped in as required if they implement the
//security interface as defined in secinterface.go

type DefaultSec struct {
}

func (d DefaultSec) Login(id string, psw string) (string, error) {
	logutil.Log("DefaultSec.Login")
	var uuid string
	var err error
	var user User
	var unencryptedPsw string
	user, err = DBGetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			logutil.Log("DefaultSec.login: " + err.Error())
			return "", errors.New("user not found")
		} else {
			logutil.Log("error in DefaultSec.login: " + err.Error())
			return "", err
		}
	}
	logutil.Log("Login checkpoint 1")

	unencryptedPsw, err = DecryptPassword(user.Password)
	if err != nil {
		logutil.Log(err.Error())
		return "", err
	}

	if unencryptedPsw != psw {
		return "", errors.New("incorrect password")
	}
	logutil.Log("Login checkpoint 2")

	uuid, err = newUUID()
	if err != nil {
		logutil.Log("error in DefaultSec.login: " + err.Error())
		return "", err
	}

	logutil.Log("Login checkpoint 3")
	//register the session
	err = DBAddSession(uuid, id)
	if err != nil {
		logutil.Log("error in DefaultSec.login add session: " + err.Error())
		return "", err
	}

	logutil.Log("secimpl Login returning uuid " + uuid)
	return uuid, nil
}

func (d DefaultSec) Logout(uuid string) error {
	logutil.Log("DefaultSec.Logout")
	err := DBDeleteSession(uuid)
	if err != nil {
		logutil.Log("error in DefaultSec.logout session: " + err.Error())
		return err
	}
	logutil.Log("DefaultSec.Logout ok for " + uuid)
	return nil
}

func (d DefaultSec) UpdateUser(user User) error {
	logutil.Log("DefaultSec.UpdateUser")
	err := DBUpdateUser(user)
	if err != nil {
		logutil.Log("error in UpdateUser: " + err.Error())
		return err
	}

	return nil
}

func (d DefaultSec) AddUser(user User) error {
	logutil.Log("DefaultSec.AddUser")
	encryptedPsw, err := EncryptPassword(user.Password)
	if err != nil {
		logutil.Log(err.Error())
		return err
	}
	user.Password = encryptedPsw

	err = DBAddUser(user)
	if err != nil {
		logutil.Log("error in AddUser: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) GetUser(id string) (User, error) {
	logutil.Log("DefaultSec.GetUser id=" + id)
	user, err := DBGetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			logutil.Log("no user found " + id)
			return user, err
		} else {
			logutil.Log("error in GetUser: " + err.Error())
			return user, err
		}
	}
	return user, nil
}

func (d DefaultSec) GetAllUsers() ([]User, error) {
	logutil.Log("DefaultSec.GetAllUsers")
	var users []User
	var err error
	users, err = DBGetAllUsers()
	if err != nil {
		logutil.Log("error in GetAllUsers: " + err.Error())
		return users, err
	}
	return users, err
}

func (d DefaultSec) DeleteUser(id string) error {
	logutil.Log("DefaultSec.DeleteUser id=" + id)
	err := DBDeleteUser(id)
	if err != nil {
		logutil.Log("error in DeleteUser: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) UpdateRole(role Role) error {
	logutil.Log("DefaultSec.UpdateRole")
	err := DBUpdateRole(role)
	if err != nil {
		logutil.Log("error in UpdateRole: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) AddRole(role Role) error {
	logutil.Log("DefaultSec.AddRole")
	err := DBAddRole(role)
	if err != nil {
		logutil.Log("error in AddRole: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) DeleteRole(name string) error {
	logutil.Log("DefaultSec.DeleteRole name=" + name)
	err := DBDeleteRole(name)
	if err != nil {
		logutil.Log("error in DeleteRole: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) GetAllRoles() ([]Role, error) {
	logutil.Log("DefaultSec.GetAllRoles")
	roles := []Role{}
	var err error
	roles, err = DBGetRoles()
	if err != nil {
		logutil.Log("error in GetAllRoles: " + err.Error())
		return roles, err
	}

	return roles, nil
}

func (d DefaultSec) GetRole(name string) (Role, error) {
	logutil.Log("DefaultSec.GetRole Name=" + name)
	permissions := make(map[string]string)
	permissions["perm1"] = "perm1 desc"
	role := Role{}
	return role, nil
}

func (d DefaultSec) LogRole(role Role) {
	logutil.Log("***role***")
	logutil.Log("role=" + role.Name + " Selected=" + fmt.Sprintf("%t", role.Selected))
	for k, v := range role.Permissions {
		logutil.Log("perm=" + k + " desc=" + v.Description + " selected=" + fmt.Sprintf("%t", v.Selected))
	}
	logutil.Log("******")

}

func (d DefaultSec) LogUser(user User) {
	logutil.Log("***user***")
	logutil.Log("user.Name=" + user.Name + " user.Password=" + user.Password)
	for k, v := range user.Roles {
		logutil.Log("***role***")
		logutil.Log("role=" + k + " Selected=" + fmt.Sprintf("%t", v.Selected))
		for i, j := range v.Permissions {
			logutil.Log("perm=" + i + " desc=" + j.Description + " selected=" + fmt.Sprintf("%t", j.Selected))
		}
		logutil.Log("******")
	}
	logutil.Log("******")

}

func (d DefaultSec) Authorize(token string, action string) error {
	var err error

	if token == "" {
		return errors.New("user login required")
	}

	var session Session
	session, err = DBGetSession(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("expired user session, new user login required")
		} else {
			logutil.Log("error in DefaultSec.Authorize: " + err.Error())
			return errors.New("error authorizing user session")
		}
	}

	//var user User
	var user User
	user, err = DBGetUser(session.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("security error, contact CPM admin")
		} else {
			logutil.Log("error in DefaultSec.Authorize: " + err.Error())
			return errors.New("error authorizing user session - u")
		}
	}

	//authorize all read-only actions
	logutil.Log("Authorize:  action=[" + action + "]")
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

func (d DefaultSec) ChangePassword(username string, newpass string) error {
	encryptedPsw, err := EncryptPassword(newpass)
	if err != nil {
		logutil.Log(err.Error())
		return err
	}

	err = DBUpdatePassword(username, encryptedPsw)
	if err != nil {
		logutil.Log(err.Error())
		return err
	}

	return nil
}

func (d DefaultSec) CompareUserToToken(username string, token string) (bool, error) {
	var err error
	var session Session
	session, err = DBGetSession(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("expired user session, new user login required")
		} else {
			logutil.Log("error in CompareUserToToken: " + err.Error())
			return false, err
		}
	}

	//var user User
	var user User
	user, err = DBGetUser(session.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("security error, contact CPM admin 2")
		} else {
			logutil.Log("error in CompareUserToToken: " + err.Error())
			return false, err
		}
	}

	logutil.Log("comparing [" + username + "] to [" + user.Name + "]")
	if username == user.Name {
		logutil.Log("compare returning true")
		return true, nil
	}

	return false, nil
}
