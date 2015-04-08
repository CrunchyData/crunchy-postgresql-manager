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
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"net/http"
)

func GetAllSettings(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllSettings: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	results, err := admindb.GetAllDBSettings()
	if err != nil {
		logit.Error.Println("GetAllSettings: error-" + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	settings := make([]Setting, len(results))
	i := 0
	for i = range results {
		settings[i].Name = results[i].Name
		settings[i].Value = results[i].Value
		settings[i].UpdateDate = results[i].UpdateDate
		i++
	}

	w.WriteJson(&settings)
}

func SaveSettings(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("SaveSettings:")
	settings := Settings{}
	err := r.DecodeJsonPayload(&settings)
	if err != nil {
		logit.Error.Println("SaveSettings: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = secimpl.Authorize(settings.Token, "perm-setting")
	if err != nil {
		logit.Error.Println("SaveSettings: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logit.Info.Println("SaveSettings: DockerRegistry=" + settings.DockerRegistry)
	logit.Info.Println("SaveSettings: PGPort=" + settings.PGPort)
	logit.Info.Println("SaveSettings: DomainName=" + settings.DomainName)

	dbsetting := admindb.DBSetting{"DOCKER-REGISTRY", settings.DockerRegistry, ""}
	err2 := admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"PG-PORT", settings.PGPort, ""}
	err2 = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"DOMAIN-NAME", settings.DomainName, ""}
	err2 = admindb.UpdateDBSetting(dbsetting)
	if err2 != nil {
		logit.Error.Println("SaveSettings: error in UpdateDBSetting " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func SaveProfiles(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("SaveProfiles:")
	profiles := Profiles{}
	err := r.DecodeJsonPayload(&profiles)
	if err != nil {
		logit.Error.Println("SaveProfiles: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(profiles.Token, "perm-setting")
	if err != nil {
		logit.Error.Println("SaveProfiles: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logit.Info.Println("SaveProfiles: smallCPU=" + profiles.SmallCPU + " smallMEM=" + profiles.SmallMEM)
	logit.Info.Println("SaveProfiles: mediumCPU=" + profiles.MediumCPU + " mediumMEM=" + profiles.MediumMEM)
	logit.Info.Println("SaveProfiles: largeCPU=" + profiles.LargeCPU + " largeMEM=" + profiles.LargeMEM)

	dbsetting := admindb.DBSetting{"S-DOCKER-PROFILE-CPU", profiles.SmallCPU, ""}
	err2 := admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"S-DOCKER-PROFILE-MEM", profiles.SmallMEM, ""}
	err2 = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"M-DOCKER-PROFILE-CPU", profiles.MediumCPU, ""}
	err2 = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"M-DOCKER-PROFILE-MEM", profiles.MediumMEM, ""}
	err2 = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"L-DOCKER-PROFILE-CPU", profiles.LargeCPU, ""}
	err2 = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"L-DOCKER-PROFILE-MEM", profiles.LargeMEM, ""}
	err2 = admindb.UpdateDBSetting(dbsetting)
	if err2 != nil {
		logit.Error.Println("SaveProfiles: sql error " + err2.Error())
		rest.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func SaveClusterProfiles(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("SaveProfiles:")
	profiles := ClusterProfiles{}
	var err error
	err = r.DecodeJsonPayload(&profiles)
	if err != nil {
		logit.Error.Println("SaveProfiles: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(profiles.Token, "perm-setting")
	if err != nil {
		logit.Error.Println("SaveProfiles: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logit.Info.Println("SaveClusterProfiles: size=" + profiles.Size)

	dbsetting := admindb.DBSetting{"CP-" + profiles.Size + "-COUNT", profiles.Count, ""}
	err = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"CP-" + profiles.Size + "-ALGO", profiles.Algo, ""}
	err = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"CP-" + profiles.Size + "-M-PROFILE", profiles.MasterProfile, ""}
	err = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"CP-" + profiles.Size + "-S-PROFILE", profiles.StandbyProfile, ""}
	err = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"CP-" + profiles.Size + "-M-SERVER", profiles.MasterServer, ""}
	err = admindb.UpdateDBSetting(dbsetting)
	dbsetting = admindb.DBSetting{"CP-" + profiles.Size + "-S-SERVER", profiles.StandbyServer, ""}
	err = admindb.UpdateDBSetting(dbsetting)

	if err != nil {
		logit.Error.Println("SaveClusterProfiles: sql error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func getClusterProfileInfo(sz string) (ClusterProfiles, error) {
	prof := ClusterProfiles{}

	results, err := admindb.GetAllDBSettingsMap()
	if err != nil {
		logit.Error.Println("GetAllSettings: error-" + err.Error())
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
