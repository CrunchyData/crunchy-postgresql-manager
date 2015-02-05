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
	"github.com/golang/glog"
)

//this is the default security implementation for CPM, others
//can be swapped in as required if they implement the
//security interface as defined in secinterface.go

type DefaultSec struct {
}

func (d DefaultSec) Login(id string, psw string) (string, error) {
	glog.Infoln("DefaultSec.Login")
	var uuid string
	var err error
	var user User
	var unencryptedPsw string
	user, err = DBGetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			glog.Errorln("DefaultSec.login: " + err.Error())
			return "", errors.New("user not found")
		} else {
			glog.Errorln("error in DefaultSec.login: " + err.Error())
			return "", err
		}
	}
	glog.Infoln("Login checkpoint 1")

	unencryptedPsw, err = DecryptPassword(user.Password)
	if err != nil {
		glog.Errorln(err.Error())
		return "", err
	}

	if unencryptedPsw != psw {
		return "", errors.New("incorrect password")
	}
	glog.Infoln("Login checkpoint 2")

	uuid, err = newUUID()
	if err != nil {
		glog.Errorln("error in DefaultSec.login: " + err.Error())
		return "", err
	}

	glog.Infoln("Login checkpoint 3")
	//register the session
	err = DBAddSession(uuid, id)
	if err != nil {
		glog.Errorln("error in DefaultSec.login add session: " + err.Error())
		return "", err
	}

	glog.Infoln("secimpl Login returning uuid " + uuid)
	return uuid, nil
}

func (d DefaultSec) Logout(uuid string) error {
	glog.Infoln("DefaultSec.Logout")
	err := DBDeleteSession(uuid)
	if err != nil {
		glog.Errorln("error in DefaultSec.logout session: " + err.Error())
		return err
	}
	glog.Infoln("DefaultSec.Logout ok for " + uuid)
	return nil
}

func (d DefaultSec) UpdateUser(user User) error {
	glog.Infoln("DefaultSec.UpdateUser")
	err := DBUpdateUser(user)
	if err != nil {
		glog.Errorln("error in UpdateUser: " + err.Error())
		return err
	}

	return nil
}

func (d DefaultSec) AddUser(user User) error {
	glog.Infoln("DefaultSec.AddUser")
	encryptedPsw, err := EncryptPassword(user.Password)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}
	user.Password = encryptedPsw

	err = DBAddUser(user)
	if err != nil {
		glog.Errorln("error in AddUser: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) GetUser(id string) (User, error) {
	glog.Infoln("DefaultSec.GetUser id=" + id)
	user, err := DBGetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			glog.Errorln("no user found " + id)
			return user, err
		} else {
			glog.Errorln("error in GetUser: " + err.Error())
			return user, err
		}
	}
	return user, nil
}

func (d DefaultSec) GetAllUsers() ([]User, error) {
	glog.Infoln("DefaultSec.GetAllUsers")
	var users []User
	var err error
	users, err = DBGetAllUsers()
	if err != nil {
		glog.Errorln("error in GetAllUsers: " + err.Error())
		return users, err
	}
	return users, err
}

func (d DefaultSec) DeleteUser(id string) error {
	glog.Infoln("DefaultSec.DeleteUser id=" + id)
	err := DBDeleteUser(id)
	if err != nil {
		glog.Errorln("error in DeleteUser: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) UpdateRole(role Role) error {
	glog.Infoln("DefaultSec.UpdateRole")
	err := DBUpdateRole(role)
	if err != nil {
		glog.Errorln("error in UpdateRole: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) AddRole(role Role) error {
	glog.Infoln("DefaultSec.AddRole")
	err := DBAddRole(role)
	if err != nil {
		glog.Errorln("error in AddRole: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) DeleteRole(name string) error {
	glog.Infoln("DefaultSec.DeleteRole name=" + name)
	err := DBDeleteRole(name)
	if err != nil {
		glog.Errorln("error in DeleteRole: " + err.Error())
		return err
	}
	return nil
}

func (d DefaultSec) GetAllRoles() ([]Role, error) {
	glog.Infoln("DefaultSec.GetAllRoles")
	roles := []Role{}
	var err error
	roles, err = DBGetRoles()
	if err != nil {
		glog.Errorln("error in GetAllRoles: " + err.Error())
		return roles, err
	}

	return roles, nil
}

func (d DefaultSec) GetRole(name string) (Role, error) {
	glog.Infoln("DefaultSec.GetRole Name=" + name)
	permissions := make(map[string]string)
	permissions["perm1"] = "perm1 desc"
	role := Role{}
	return role, nil
}

func (d DefaultSec) LogRole(role Role) {
	glog.Infoln("***role***")
	glog.Infoln("role=" + role.Name + " Selected=" + fmt.Sprintf("%t", role.Selected))
	for k, v := range role.Permissions {
		glog.Infoln("perm=" + k + " desc=" + v.Description + " selected=" + fmt.Sprintf("%t", v.Selected))
	}
	glog.Infoln("******")

}

func (d DefaultSec) LogUser(user User) {
	glog.Infoln("***user***")
	glog.Infoln("user.Name=" + user.Name + " user.Password=" + user.Password)
	for k, v := range user.Roles {
		glog.Infoln("***role***")
		glog.Infoln("role=" + k + " Selected=" + fmt.Sprintf("%t", v.Selected))
		for i, j := range v.Permissions {
			glog.Infoln("perm=" + i + " desc=" + j.Description + " selected=" + fmt.Sprintf("%t", j.Selected))
		}
		glog.Infoln("******")
	}
	glog.Infoln("******")

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
			glog.Infoln("error in DefaultSec.Authorize: " + err.Error())
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
			glog.Errorln("error in DefaultSec.Authorize: " + err.Error())
			return errors.New("error authorizing user session - u")
		}
	}

	//authorize all read-only actions
	glog.Infoln("Authorize:  action=[" + action + "]")
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
		glog.Errorln(err.Error())
		return err
	}

	err = DBUpdatePassword(username, encryptedPsw)
	if err != nil {
		glog.Errorln(err.Error())
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
			glog.Errorln("error in CompareUserToToken: " + err.Error())
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
			glog.Errorln("error in CompareUserToToken: " + err.Error())
			return false, err
		}
	}

	glog.Infoln("comparing [" + username + "] to [" + user.Name + "]")
	if username == user.Name {
		glog.Errorln("compare returning true")
		return true, nil
	}

	return false, nil
}
