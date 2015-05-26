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

type Project struct {
	ID         string
	Name       string
	Desc       string
	UpdateDate string
	Token      string
}

func UpdateProject(w rest.ResponseWriter, r *rest.Request) {
	logit.Info.Println("UpdateProject: in UpdateProject")
	project := Project{}
	err := r.DecodeJsonPayload(&project)
	if err != nil {
		logit.Error.Println("UpdateProject: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logit.Info.Println("UpdateProject: ID=" + project.ID)
	logit.Info.Println("UpdateProject: Name=" + project.Name)
	logit.Info.Println("UpdateProject: Desc=" + project.Desc)
	logit.Info.Println("UpdateProject: token=" + project.Token)

	err = secimpl.Authorize(project.Token, "perm-container")
	if err != nil {
		logit.Error.Println("UpdateProject: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dbproject := admindb.Project{project.ID, project.Name, project.Desc, project.UpdateDate}
	err = admindb.UpdateProject(dbproject)
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
	logit.Info.Println("AddProject: in AddProject")
	project := Project{}
	err := r.DecodeJsonPayload(&project)
	if err != nil {
		logit.Error.Println("AddProject: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = secimpl.Authorize(project.Token, "perm-container")
	if err != nil {
		logit.Error.Println("UpdateProject: authorize error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dbproject := admindb.Project{project.ID, project.Name, project.Desc, project.UpdateDate}
	var newid int
	newid, err = admindb.InsertProject(dbproject)
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
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
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

	results, err := admindb.GetProject(ID)
	if err != nil {
		logit.Error.Println("GetProject:" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	project := Project{results.ID, results.Name, results.Desc, results.UpdateDate, ""}

	w.WriteJson(&project)
}

func GetAllProjects(w rest.ResponseWriter, r *rest.Request) {
	err := secimpl.Authorize(r.PathParam("Token"), "perm-read")
	if err != nil {
		logit.Error.Println("GetAllProjects: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	projectsList, err := admindb.GetAllProjects()
	if err != nil {
		logit.Error.Println("GetAllProjects: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	projects := make([]Project, len(projectsList))
	i := 0
	for i = range projectsList {
		projects[i].ID = projectsList[i].ID
		projects[i].Name = projectsList[i].Name
		projects[i].Desc = projectsList[i].Desc
		projects[i].UpdateDate = projectsList[i].UpdateDate
		i++
	}

	w.WriteJson(&projects)
}

func DeleteProject(w rest.ResponseWriter, r *rest.Request) {

	err := secimpl.Authorize(r.PathParam("Token"), "perm-container")
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
	err = admindb.DeleteProject(ID)
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
