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

package adminapi

import (
	"database/sql"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
)

// GetAllGeneralSettings return a list of general settings
func GetAllGeneralSettings(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err := admindb.GetAllGeneralSettings(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	settings := make([]types.Setting, len(results))
	i := 0
	for i = range results {
		settings[i].Name = results[i].Name
		settings[i].Value = results[i].Value
		settings[i].UpdateDate = results[i].UpdateDate
		i++
	}

	w.WriteJson(&settings)
}

// GetAllSettings return all settings
func GetAllSettings(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err := admindb.GetAllSettings(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	settings := make([]types.Setting, len(results))
	i := 0
	for i = range results {
		settings[i].Name = results[i].Name
		settings[i].Value = results[i].Value
		settings[i].Description = results[i].Description
		settings[i].UpdateDate = results[i].UpdateDate
		i++
	}

	w.WriteJson(&settings)
}

// SaveSetting update an existing setting value
func SaveSetting(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	//logit.Info.Println("SaveSetting:")
	setting := types.Setting{}
	err = r.DecodeJsonPayload(&setting)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = secimpl.Authorize(dbConn, setting.Token, "perm-setting")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dbSetting := types.Setting{}
	dbSetting.Name = setting.Name
	dbSetting.Value = setting.Value

	err2 := admindb.UpdateSetting(dbConn, dbSetting)
	if err2 != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func getClusterProfileInfo(dbConn *sql.DB, sz string) (types.ClusterProfiles, error) {

	prof := types.ClusterProfiles{}

	results, err := admindb.GetAllSettingsMap(dbConn)
	if err != nil {
		logit.Error.Println(err.Error())
		return prof, err
	}

	prof.Size = sz
	prof.Count = results["CP-"+prof.Size+"-COUNT"]
	prof.Algo = results["CP-"+prof.Size+"-ALGO"]
	prof.MasterProfile = results["CP-"+prof.Size+"-M-PROFILE"]
	prof.StandbyProfile = results["CP-"+prof.Size+"-S-PROFILE"]
	prof.MasterServer = results["CP-"+prof.Size+"-M-SERVER"]
	prof.StandbyServer = results["CP-"+prof.Size+"-S-SERVER"]

	return prof, nil
}
