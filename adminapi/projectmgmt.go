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
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
)

type Project struct {
	ID         string
	Name       string
	Desc       string
	UpdateDate string
	Token      string
	Containers map[string]string
	Clusters   map[string]string
}

type Child2 struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	ID        string `json:"id"`
	ProjectID string `json:"projectid"`
}
type Child struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	ID        string  `json:"id"`
	ProjectID string  `json:"projectid"`
	Children  []Child `json:"children"`
}

type Project2 struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	ID       string  `json:"id"`
	Children []Child `json:"children"`
}

func UpdateProject(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("UpdateProject: in UpdateProject")
	project := Project{}
	err = r.DecodeJsonPayload(&project)
	if err != nil {
		logit.Error.Println("UpdateProject: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logit.Info.Println("UpdateProject: ID=" + project.ID)
	logit.Info.Println("UpdateProject: Name=" + project.Name)
	logit.Info.Println("UpdateProject: Desc=" + project.Desc)
	logit.Info.Println("UpdateProject: token=" + project.Token)

	err = secimpl.Authorize(dbConn, project.Token, "perm-container")
	if err != nil {
		logit.Error.Println("UpdateProject: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dbproject := admindb.Project{
		ID:         project.ID,
		Name:       project.Name,
		Desc:       project.Desc,
		UpdateDate: project.UpdateDate,
		Containers: project.Containers,
		Clusters:   project.Clusters}

	err = admindb.UpdateProject(dbConn, dbproject)
	if err != nil {
		logit.Error.Println("UpdateProject: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func AddProject(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	logit.Info.Println("AddProject: in AddProject")
	project := Project{}
	err = r.DecodeJsonPayload(&project)
	if err != nil {
		logit.Error.Println("AddProject: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(dbConn, project.Token, "perm-container")
	if err != nil {
		logit.Error.Println("UpdateProject: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dbproject := admindb.Project{
		ID:         project.ID,
		Name:       project.Name,
		Desc:       project.Desc,
		UpdateDate: project.UpdateDate,
		Containers: project.Containers}
	var newid int
	newid, err = admindb.InsertProject(dbConn, dbproject)
	if err != nil {
		logit.Error.Println("AddProject: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logit.Info.Printf("added project id= %d\n", newid)

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}

func GetProject(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetProject: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("GetProject: error Project ID required")
		rest.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}

	results, err := admindb.GetProject(dbConn, ID)
	if err != nil {
		logit.Error.Println("GetProject:" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	project := Project{results.ID, results.Name, results.Desc, results.UpdateDate, "", results.Containers, results.Clusters}

	w.WriteJson(&project)
}

func GetAllProjects(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllProjects: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	projectsList, err := admindb.GetAllProjects(dbConn)
	if err != nil {
		logit.Error.Println("GetAllProjects: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	projects := make([]Project2, len(projectsList))
	i := 0
	for i = range projectsList {
		projects[i].ID = projectsList[i].ID
		projects[i].Name = projectsList[i].Name
		projects[i].Type = "project"

		projectchildren := make([]Child, 2)
		projects[i].Children = projectchildren

		projects[i].Children[0].Name = "Clusters"
		projects[i].Children[0].Type = "label"
		projects[i].Children[0].ProjectID = projects[i].ID
		//projects[i].Children[0].Children := make([]Child, len(projectsList[i].Clusters))
		j := 0
		clusterchilds := make([]Child, len(projectsList[i].Clusters))
		for jk, jv := range projectsList[i].Clusters {
			ch := Child{
				Name: jv,
				Type: "cluster",
				ID:   jk,
			}
			ch.ProjectID = projects[i].ID
			clusterchilds[j] = ch
			j++
		}
		projects[i].Children[0].Children = clusterchilds

		projects[i].Children[1].Name = "Databases"
		projects[i].Children[1].Type = "label"
		projects[i].Children[1].ProjectID = projects[i].ID
		dbchilds := make([]Child, len(projectsList[i].Containers))
		//projects[i].Children[1].Children := make([]Child, len(projectsList[i].Containers))
		k := 0
		for kk, kv := range projectsList[i].Containers {
			ch := Child{
				Name: kv,
				Type: "database",
				ID:   kk,
			}
			ch.ProjectID = projects[i].ID
			dbchilds[k] = ch
			k++
		}
		projects[i].Children[1].Children = dbchilds

		i++
	}

	w.WriteJson(&projects)
}

func DeleteProject(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("BackupNow: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-container")
	if err != nil {
		logit.Error.Println("DeleteProject: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ID := r.PathParam("ID")
	if ID == "" {
		logit.Error.Println("DeleteProject: error ID required")
		rest.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	err = admindb.DeleteProject(dbConn, ID)
	if err != nil {
		logit.Error.Println("DeleteProject: error secimpl call" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}
