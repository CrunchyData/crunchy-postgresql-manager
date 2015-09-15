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

package admindb

import (
	"database/sql"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
)

const CPMTEST_DB = "cpmtest"
const CPMTEST_USER = "cpmtest"

func GetUserCredentials(dbConn *sql.DB, node *Container) (Credential, error) {
	var err error
	cred := Credential{}

	if node.Image != "cpm-node-proxy" {
		//get port
		var pgport Setting
		pgport, err = GetSetting(dbConn, "PG-PORT")
		nodeuser, err := GetContainerUser(dbConn, node.Name, CPMTEST_USER)
		if err != nil {
			logit.Error.Println(err.Error())
			return cred, err
		}
		cred.Host = node.Name
		cred.Database = CPMTEST_DB
		cred.Username = nodeuser.Rolname
		cred.Password = nodeuser.Passwd
		cred.Port = pgport.Value
		return cred, err
	}

	//return proxy credentials
	var proxy Proxy
	proxy, err = GetProxy(dbConn, node.Name)
	if err != nil {
		logit.Error.Println(err.Error())
		return cred, err
	}

	cred.Database = proxy.Database
	cred.Host = proxy.Host
	cred.Username = proxy.Usename
	cred.Password = proxy.Passwd
	cred.Port = "5432"
	return cred, err

}
