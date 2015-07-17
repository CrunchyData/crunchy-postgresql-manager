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
)

type Permission struct {
	Name        string
	Description string
	Selected    bool
}

type Role struct {
	Name        string
	Selected    bool
	Permissions map[string]Permission
	UpdateDate  string
	Token       string
}

type User struct {
	Name       string
	Password   string
	Roles      map[string]Role
	UpdateDate string
	Token      string
}

type Session struct {
	Name       string
	Token      string
	UpdateDate string
}

type SecInterface interface {
	Authorize(*sql.DB, string, string) error
	Login(*sql.DB, string, string) (string, error)
	Logout(*sql.DB, string) error
	ChangePassword(*sql.DB, string, string) error
	CompareUserToToken(*sql.DB, string, string) (bool, error)
	UpdateUser(*sql.DB, User) error
	AddUser(*sql.DB, User) error
	GetUser(*sql.DB, string) (User, error)
	GetAllUsers(*sql.DB) ([]User, error)
	DeleteUser(*sql.DB, string) error
	UpdateRole(*sql.DB, Role) error
	AddRole(*sql.DB, Role) error
	DeleteRole(*sql.DB, string) error
	GetAllRoles(*sql.DB) ([]Role, error)
	GetRole(*sql.DB, string) (Role, error)
	LogRole(Role)
	LogUser(User)
}
