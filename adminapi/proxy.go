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
	"strconv"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	//"github.com/crunchydata/crunchy-postgresql-manager/template"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	//"strconv"
	//"time"
	"fmt"
)

type ProxyRequest struct {
	Token string
	Profile string
	Image string
	ServerID string
	ProjectID string
	ContainerName string
	Standalone string
	DatabaseHost string
	DatabaseUserID string
	DatabaseUserPassword string
	DatabasePort string
	Database string
}

func ProvisionProxy(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println("ProvisionProxy: error " + err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()


	proxyrequest := ProxyRequest{}
  	err = r.DecodeJsonPayload(&proxyrequest)
        if err != nil {
                logit.Error.Println("ProvisionProxy: error in decode" + err.Error())
                rest.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

	logit.Info.Println("ProvisionProxy:  Token=[" + proxyrequest.Token + "]")

	err = secimpl.Authorize(dbConn, proxyrequest.Token, "perm-container")
	if err != nil {
		logit.Error.Println("ProvisionProxy: validate token error " + err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	errorStr := ""

	if proxyrequest.ServerID == "" {
		logit.Error.Println("ProvisionProxy error serverid required")
		errorStr = "ServerID required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}

	if proxyrequest.ProjectID == "" {
		logit.Error.Println("ProvisionProxy error ProjectID required")
		errorStr = "ProjectID required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}

	if proxyrequest.ContainerName == "" {
		logit.Error.Println("ProvisionProxy error containername required")
		errorStr = "ContainerName required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}

	if proxyrequest.Image == "" {
		logit.Error.Println("ProvisionProxy error image required")
		errorStr = "Image required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if proxyrequest.DatabaseHost == "" {
		logit.Error.Println("ProvisionProxy error DatabaseHost required")
		errorStr = "DatabaseHost required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if proxyrequest.Database == "" {
		logit.Error.Println("ProvisionProxy error Database required")
		errorStr = "Database required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}


	logit.Info.Println("Image=" + proxyrequest.Image)
	logit.Info.Println("Profile=" + proxyrequest.Profile)
	logit.Info.Println("ServerID=" + proxyrequest.ServerID)
	logit.Info.Println("ProjectID=" + proxyrequest.ProjectID)
	logit.Info.Println("ContainerName=" + proxyrequest.ContainerName)
	logit.Info.Println("Standalone=" + proxyrequest.Standalone)

 	params := &cpmserverapi.DockerRunRequest{}
        params.Image = proxyrequest.Image
        params.ServerID = proxyrequest.ServerID
        params.ProjectID = proxyrequest.ProjectID
        params.ContainerName = proxyrequest.ContainerName
        params.Standalone = proxyrequest.Standalone

	err = provisionImpl(dbConn, params, proxyrequest.Profile, false)
	if err != nil {
		logit.Error.Println("ProvisionProxy error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = insertProxy(&proxyrequest)
	if err != nil {
		logit.Error.Println("ProvisionProxy error " + err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	status := SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)

}

func insertProxy(request *ProxyRequest) (error) {

	var containerUserID int

        //create user in the admin db
        dbuser := admindb.ContainerUser{}
        dbuser.Containername = request.ContainerName
        dbuser.Passwd = request.DatabaseUserPassword
        dbuser.Rolname = request.DatabaseUserID

        dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
        if err != nil {
                logit.Error.Println("insertProxy: error " + err.Error())
                return err

        }
        defer dbConn.Close()

        containerUserID, err = admindb.AddContainerUser(dbConn, dbuser)
        if err != nil {
                logit.Error.Println("insertProxy: " + err.Error())
                return err
        }

        logit.Info.Printf("insertProxy: new ID %d\n " ,  containerUserID)

	var container admindb.Container
        container, err = admindb.GetContainerByName(dbConn, request.ContainerName)
        if err != nil {
                logit.Error.Println("insertProxy: " + err.Error())
                return err
        }

	proxy := Proxy{}
        proxy.ContainerUserID = strconv.Itoa(containerUserID)
        proxy.ContainerID = container.ID
        proxy.Host = request.DatabaseHost
        proxy.Database = request.Database
        proxy.ProjectID = request.ProjectID
        proxy.Port = request.DatabasePort

 	queryStr := fmt.Sprintf("insert into proxy ( containeruserid, containerid, projectid, port, host, databasename, updatedt) values ( %s, %s, %s, '%s', '%s', '%s', now()) returning id", 
	proxy.ContainerUserID, proxy.ContainerID, proxy.ProjectID, proxy.Port, proxy.Host, proxy.Database) 
 
        logit.Info.Println("insertProxy:" + queryStr) 
        var proxyid int 
        err = dbConn.QueryRow(queryStr).Scan(&proxyid) 
        switch { 
        case err != nil: 
                logit.Info.Println("insertProxy:" + err.Error()) 
                return err 
        default: 
                logit.Info.Println("insertProxy: inserted returned is " + strconv.Itoa(proxyid)) 
        }

	return err
}

func GetProxyByContainerID(w rest.ResponseWriter, r *rest.Request) {
        dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
        if err != nil {
                logit.Error.Println("GetProxyByContainerID: error " + err.Error())
                rest.Error(w, err.Error(), 400)
                return

        }
        defer dbConn.Close()
        err = secimpl.Authorize(dbConn, r.PathParam("Token"), "perm-read")
        if err != nil {
                logit.Error.Println("GetProxyByContainerID: authorize error " + err.Error())
                rest.Error(w, err.Error(), http.StatusUnauthorized)
                return
        }

        ContainerID := r.PathParam("ContainerID")
        if ContainerID == "" {
                logit.Error.Println("GetProxyByContainerID: ContainerID is required")
                rest.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

 	proxy, err := admindb.GetProxyByContainerID(dbConn, ContainerID)
        if err != nil {
                logit.Error.Println("GetProxyByContainerID:" + err.Error())
                rest.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

        w.WriteJson(&proxy)
}


