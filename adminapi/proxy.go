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
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/swarmapi"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"net/http"
	"strconv"
)

type ProxyRequest struct {
	Token         string
	ID            string
	Profile       string
	Image         string
	ServerID      string
	ProjectID     string
	ContainerName string
	Standalone    string
	Host          string
	Usename       string
	Passwd        string
	Port          string
	Database      string
}

// ProvisionProxy creates a Docker image for a proxy node definition
func ProvisionProxy(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	proxyrequest := ProxyRequest{}
	err = r.DecodeJsonPayload(&proxyrequest)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logit.Info.Println("ProvisionProxy:  Token=[" + proxyrequest.Token + "]")

	err = secimpl.Authorize(dbConn, proxyrequest.Token, "perm-container")
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	errorStr := ""

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
	if proxyrequest.Host == "" {
		logit.Error.Println("ProvisionProxy error DatabaseHost required")
		errorStr = "DatabaseHost required"
		rest.Error(w, errorStr, http.StatusBadRequest)
		return
	}
	if proxyrequest.Host == "127.0.0.1" || proxyrequest.Host == "localhost" {
		logit.Error.Println("ProvisionProxy error DatabaseHost can not be localhost")
		errorStr = "DatabaseHost can not be localhost"
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
	logit.Info.Println("ProjectID=" + proxyrequest.ProjectID)
	logit.Info.Println("ContainerName=" + proxyrequest.ContainerName)
	logit.Info.Println("Standalone=" + proxyrequest.Standalone)

	params := &swarmapi.DockerRunRequest{}
	params.Image = proxyrequest.Image
	params.ProjectID = proxyrequest.ProjectID
	params.ContainerName = proxyrequest.ContainerName
	params.Standalone = proxyrequest.Standalone
	params.Profile = proxyrequest.Profile

	_, err = provisionImpl(dbConn, params, false)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = insertProxy(&proxyrequest)
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

func insertProxy(request *ProxyRequest) error {

	var containerUserID int

	//create user in the admin db
	dbuser := types.ContainerUser{}
	dbuser.Containername = request.ContainerName
	dbuser.Passwd = request.Passwd
	dbuser.Rolname = request.Usename

	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		return err

	}
	defer dbConn.Close()

	containerUserID, err = admindb.AddContainerUser(dbConn, dbuser)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	logit.Info.Printf("insertProxy: new ID %d\n ", containerUserID)

	var container types.Container
	container, err = admindb.GetContainerByName(dbConn, request.ContainerName)
	if err != nil {
		logit.Error.Println(err.Error())
		return err
	}

	proxy := types.Proxy{}
	proxy.ContainerUserID = strconv.Itoa(containerUserID)
	proxy.ContainerID = container.ID
	proxy.Host = request.Host
	proxy.Database = request.Database
	proxy.ProjectID = request.ProjectID
	proxy.Port = request.Port

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

// GetProxyByContainerID returns a proxy node defintion for a given container
func GetProxyByContainerID(w rest.ResponseWriter, r *rest.Request) {
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

	ContainerID := r.PathParam("ContainerID")
	if ContainerID == "" {
		logit.Error.Println("ContainerID is required")
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	proxy, err := admindb.GetProxyByContainerID(dbConn, ContainerID)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request := &swarmapi.DockerInspectRequest{}
	var inspectInfo swarmapi.DockerInspectResponse
	request.ContainerName = proxy.ContainerName
	inspectInfo, err = swarmapi.DockerInspect(request)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	proxy.ServerName = inspectInfo.ServerID

	w.WriteJson(&proxy)
}

// ProxyUpdate updates a proxy node definition
func ProxyUpdate(w rest.ResponseWriter, r *rest.Request) {
	dbConn, err := util.GetConnection(CLUSTERADMIN_DB)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), 400)
		return

	}
	defer dbConn.Close()

	req := ProxyRequest{}
	err = r.DecodeJsonPayload(&req)
	if err != nil {
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logit.Info.Println("ID=" + req.ID)
	logit.Info.Println("Port" + req.Port)
	logit.Info.Println("Host" + req.Host)
	logit.Info.Println("Database" + req.Database)

	queryStr := fmt.Sprintf("update proxy set ( port, host, databasename, updatedt) = ( '%s', '%s', '%s', now()) where id = %s returning id",
		req.Port, req.Host, req.Database, req.ID)

	logit.Info.Println("UpdateProxy:" + queryStr)
	var proxyid int
	err = dbConn.QueryRow(queryStr).Scan(&proxyid)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		logit.Info.Println("UpdateProxy: update " + strconv.Itoa(proxyid))
	}

	w.WriteHeader(http.StatusOK)
	status := types.SimpleStatus{}
	status.Status = "OK"
	w.WriteJson(&status)
}
