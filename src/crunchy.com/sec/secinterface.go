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
	Authorize(string, string) error
	Login(string, string) (string, error)
	Logout(string) error
	ChangePassword(string, string) error
	CompareUserToToken(string, string) (bool, error)
	UpdateUser(User) error
	AddUser(User) error
	GetUser(string) (User, error)
	GetAllUsers() ([]User, error)
	DeleteUser(string) error
	UpdateRole(Role) error
	AddRole(Role) error
	DeleteRole(string) error
	GetAllRoles() ([]Role, error)
	GetRole(string) (Role, error)
	LogRole(Role)
	LogUser(User)
}
