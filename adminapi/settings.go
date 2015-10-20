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

func SaveSetting(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	logit.Info.Println("SaveSetting:")
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

/**
func SaveSettings(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	logit.Info.Println("SaveSettings:")
	settings := types.Settings{}
	err = r.DecodeJsonPayload(&settings)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = secimpl.Authorize(dbConn, settings.Token, "perm-setting")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logit.Info.Println("SaveSettings: DockerRegistry=" + settings.DockerRegistry)
	logit.Info.Println("SaveSettings: PGPort=" + settings.PGPort)
	logit.Info.Println("SaveSettings: DomainName=" + settings.DomainName)

	dbsetting := types.Setting{"DOCKER-REGISTRY", settings.DockerRegistry, "", ""}
	err2 := admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"PG-PORT", settings.PGPort, "", ""}
	err2 = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"DOMAIN-NAME", settings.DomainName, "", ""}
	err2 = admindb.UpdateSetting(dbConn, dbsetting)
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

func SaveProfiles(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("SaveProfiles:")
	profiles := types.Profiles{}
	err = r.DecodeJsonPayload(&profiles)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, profiles.Token, "perm-setting")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logit.Info.Println("SaveProfiles: smallCPU=" + profiles.SmallCPU + " smallMEM=" + profiles.SmallMEM)
	logit.Info.Println("SaveProfiles: mediumCPU=" + profiles.MediumCPU + " mediumMEM=" + profiles.MediumMEM)
	logit.Info.Println("SaveProfiles: largeCPU=" + profiles.LargeCPU + " largeMEM=" + profiles.LargeMEM)

	dbsetting := types.Setting{"S-DOCKER-PROFILE-CPU", profiles.SmallCPU, "", ""}
	err2 := admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"S-DOCKER-PROFILE-MEM", profiles.SmallMEM, "", ""}
	err2 = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"M-DOCKER-PROFILE-CPU", profiles.MediumCPU, "", ""}
	err2 = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"M-DOCKER-PROFILE-MEM", profiles.MediumMEM, "", ""}
	err2 = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"L-DOCKER-PROFILE-CPU", profiles.LargeCPU, "", ""}
	err2 = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"L-DOCKER-PROFILE-MEM", profiles.LargeMEM, "", ""}
	err2 = admindb.UpdateSetting(dbConn, dbsetting)
	if err2 != nil {
		logit.Error.Println(err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func SaveClusterProfiles(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()
	logit.Info.Println("SaveProfiles:")
	profiles := types.ClusterProfiles{}
	err = r.DecodeJsonPayload(&profiles)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, profiles.Token, "perm-setting")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logit.Info.Println("SaveClusterProfiles: size=" + profiles.Size)

	dbsetting := types.Setting{"CP-" + profiles.Size + "-COUNT", profiles.Count, "", ""}
	err = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"CP-" + profiles.Size + "-ALGO", profiles.Algo, "", ""}
	err = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"CP-" + profiles.Size + "-M-PROFILE", profiles.MasterProfile, "", ""}
	err = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"CP-" + profiles.Size + "-S-PROFILE", profiles.StandbyProfile, "", ""}
	err = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"CP-" + profiles.Size + "-M-SERVER", profiles.MasterServer, "", ""}
	err = admindb.UpdateSetting(dbConn, dbsetting)
	dbsetting = types.Setting{"CP-" + profiles.Size + "-S-SERVER", profiles.StandbyServer, "", ""}
	err = admindb.UpdateSetting(dbConn, dbsetting)

	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}
*/

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
