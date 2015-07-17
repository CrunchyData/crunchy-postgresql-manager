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
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/sec"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

type Setting struct {
	Name       string
	Value      string
	UpdateDate string
}
type Server struct {
	ID             string
	Name           string
	IPAddress      string
	DockerBridgeIP string
	PGDataPath     string
	ServerClass    string
	CreateDate     string
	NodeCount      string
}

type Project struct {
	ID         string
	Name       string
	Desc       string
	Containers map[string]string
	Clusters   map[string]string
	UpdateDate string
}

type Container struct {
	ID          string
	ClusterID   string
	ServerID    string
	ServerName  string
	Name        string
	Role        string
	Image       string
	CreateDate  string
	ProjectID   string
	ProjectName string
	ClusterName string
}

type Cluster struct {
	ID          string
	ProjectID   string
	Name        string
	ClusterType string
	Status      string
	CreateDate  string
	Containers  map[string]string
}

type ContainerUser struct {
	ID             string
	Containername  string
	ContainerID    string
	Passwd         string
	Rolname        string
	Rolsuper       string
	Rolinherit     string
	Rolcreaterole  string
	Rolcreatedb    string
	Rolcatupdate   string
	Rolcanlogin    string
	Rolreplication string
	UpdateDate     string
}

type LinuxStats struct {
	ID        string
	ClusterID string
	Stats     string
}

type PGStats struct {
	ID        string
	ClusterID string
	Stats     string
}

func GetServer(dbConn *sql.DB, id string) (Server, error) {
	//logit.Info.Println("GetServer called with id=" + id)
	server := Server{}

	err := dbConn.QueryRow(fmt.Sprintf("select id, name, ipaddress, dockerbip, pgdatapath, serverclass, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from server where id=%s", id)).Scan(&server.ID, &server.Name, &server.IPAddress, &server.DockerBridgeIP, &server.PGDataPath, &server.ServerClass, &server.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetServer:no server with that id")
		return server, err
	case err != nil:
		logit.Info.Println("admindb:GetServer:" + err.Error())
		return server, err
	default:
		logit.Info.Println("admindb:GetServer: server name returned is " + server.Name)
	}

	return server, nil
}

func GetClusterName(dbConn *sql.DB, id string) (string, error) {
	//logit.Info.Println("admindb:GetCluster: called")
	var clustername string
	err := dbConn.QueryRow(fmt.Sprintf("select name from cluster where id=%s", id)).Scan(&clustername)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetClusterName:no cluster with that id")
		return clustername, err
	case err != nil:
		logit.Info.Println("admindb:GetClusterName:" + err.Error())
		return clustername, err
	default:
		logit.Info.Println("admindb:GetClusterName: cluster name returned is " + clustername)
	}

	return clustername, err
}

func GetCluster(dbConn *sql.DB, id string) (Cluster, error) {
	//logit.Info.Println("admindb:GetCluster: called")
	cluster := Cluster{}

	err := dbConn.QueryRow(fmt.Sprintf("select projectid, id, name, clustertype, status, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from cluster where id=%s", id)).Scan(&cluster.ProjectID, &cluster.ID, &cluster.Name, &cluster.ClusterType, &cluster.Status, &cluster.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetCluster:no cluster with that id")
		return cluster, err
	case err != nil:
		logit.Info.Println("admindb:GetCluster:" + err.Error())
		return cluster, err
	default:
		logit.Info.Println("admindb:GetCluster: cluster name returned is " + cluster.Name)
	}

	var containers []Container
	cluster.Containers = make(map[string]string)

	containers, err = GetAllContainersForCluster(dbConn, cluster.ID)
	if err != nil {
		logit.Info.Println("admindb:GetCluster:" + err.Error())
		return cluster, err
	}

	for i := range containers {
		cluster.Containers[containers[i].ID] = containers[i].Name
	}

	return cluster, nil
}

func GetAllClustersForProject(dbConn *sql.DB, projectId string) ([]Cluster, error) {
	//logit.Info.Println("admindb:GetAllClusters: called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select id, projectid, name, clustertype, status, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from cluster where projectid = " + projectId + " order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var containers []Container
	clusters := make([]Cluster, 0)
	for rows.Next() {
		cluster := Cluster{}
		if err = rows.Scan(
			&cluster.ID,
			&cluster.ProjectID,
			&cluster.Name,
			&cluster.ClusterType,
			&cluster.Status, &cluster.CreateDate); err != nil {
			return nil, err
		}

		cluster.Containers = make(map[string]string)
		containers, err = GetAllContainersForCluster(dbConn, cluster.ID)
		if err != nil {
			logit.Info.Println("admindb:GetCluster:" + err.Error())
		}

		for i := range containers {
			cluster.Containers[containers[i].ID] = containers[i].Name
			logit.Info.Println("admindb:GetCluster: add to map " + cluster.Containers[containers[i].ID])
		}

		clusters = append(clusters, cluster)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return clusters, nil
}
func GetAllClusters(dbConn *sql.DB) ([]Cluster, error) {
	//logit.Info.Println("admindb:GetAllClusters: called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select id, projectid, name, clustertype, status, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from cluster order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var containers []Container
	clusters := make([]Cluster, 0)
	for rows.Next() {
		cluster := Cluster{}
		if err = rows.Scan(
			&cluster.ID,
			&cluster.ProjectID,
			&cluster.Name,
			&cluster.ClusterType,
			&cluster.Status, &cluster.CreateDate); err != nil {
			return nil, err
		}

		cluster.Containers = make(map[string]string)
		containers, err = GetAllContainersForCluster(dbConn, cluster.ID)
		if err != nil {
			logit.Info.Println("admindb:GetCluster:" + err.Error())
		}

		for i := range containers {
			cluster.Containers[containers[i].ID] = containers[i].Name
			logit.Info.Println("admindb:GetCluster: add to map " + cluster.Containers[containers[i].ID])
		}

		clusters = append(clusters, cluster)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return clusters, nil
}

func UpdateCluster(dbConn *sql.DB, cluster Cluster) error {
	//logit.Info.Println("admindb:UpdateCluster:called")
	queryStr := fmt.Sprintf("update cluster set ( name, clustertype, status) = ('%s', '%s', '%s') where id = %s returning id", cluster.Name, cluster.ClusterType, cluster.Status, cluster.ID)

	logit.Info.Println("admindb:UpdateCluster:update str=[" + queryStr + "]")
	var clusterid int
	err := dbConn.QueryRow(queryStr).Scan(&clusterid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:UpdateCluster:cluster updated " + cluster.ID)
	}
	return nil

}
func InsertCluster(dbConn *sql.DB, cluster Cluster) (int, error) {
	//logit.Info.Println("admindb:InsertCluster:called")
	queryStr := fmt.Sprintf("insert into cluster ( name, projectid, clustertype, status, createdt) values ( '%s', %s, '%s', '%s', now()) returning id", cluster.Name, cluster.ProjectID, cluster.ClusterType, cluster.Status)

	logit.Info.Println("admindb:InsertCluster:" + queryStr)
	var clusterid int
	err := dbConn.QueryRow(queryStr).Scan(&clusterid)
	switch {
	case err != nil:
		logit.Info.Println("admindb:InsertCluster:" + err.Error())
		return -1, err
	default:
		logit.Info.Println("admindb:InsertCluster: cluster inserted returned is " + strconv.Itoa(clusterid))
	}

	return clusterid, nil
}

func DeleteCluster(dbConn *sql.DB, id string) error {
	queryStr := fmt.Sprintf("delete from cluster where  id=%s returning id", id)
	//logit.Info.Println("admindb:DeleteCluster:" + queryStr)

	var clusterid int
	err := dbConn.QueryRow(queryStr).Scan(&clusterid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:DeleteCluster:cluster deleted " + id)
	}
	return nil
}

func GetContainer(dbConn *sql.DB, id string) (Container, error) {
	//logit.Info.Println("admindb:GetContainer:called")
	container := Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from container c, project p, server s where c.id=%s and p.id = c.projectid and c.serverid = s.id", id)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetContainer:no container with that id " + id)
		return container, err
	case err != nil:
		return container, err
	}

	if container.ClusterID != "-1" {
		var clustername string
		clustername, err = GetClusterName(dbConn, container.ClusterID)
		if err != nil {
			logit.Info.Println("admindb:GetContainer:error " + err.Error())
			return container, err
		}
		container.ClusterName = clustername
	}
	return container, nil
}

func GetContainerByName(dbConn *sql.DB, name string) (Container, error) {
	//logit.Info.Println("admindb:GetNodeByName:called")
	container := Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from project p, server s, container c where c.name='%s' and c.projectid = p.id and c.serverid = s.id", name)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetContainerByName:no container with that name " + name)
		return container, err
	case err != nil:
		return container, err
	}
	if container.ClusterID != "-1" {
		var clustername string
		clustername, err = GetClusterName(dbConn, container.ClusterID)
		if err != nil {
			logit.Info.Println("admindb:GetContainerByName:error " + err.Error())
			return container, err
		}
		container.ClusterName = clustername
	}

	return container, nil
}

//find the oldest container in a cluster, used for serf join-cluster event
func GetContainerOldestInCluster(dbConn *sql.DB, clusterid string) (Container, error) {
	//logit.Info.Println("admindb:GetNodeOldestInCluster:called")
	container := Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from project p, server s, container c where  c.projectid = p.id and c.serverid = s.id and c.createdt = (select max(createdt) from container where clusterid = %s)", clusterid)
	logit.Info.Println("admindb:GetNodeOldestInCluster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetContainerOldestInCluster: no container with that clusterid " + clusterid)
		return container, err
	case err != nil:
		return container, err
	}

	if container.ClusterID != "-1" {
		var clustername string
		clustername, err = GetClusterName(dbConn, container.ClusterID)
		if err != nil {
			logit.Info.Println("admindb:GetContainerOldest:error " + err.Error())
			return container, err
		}
		container.ClusterName = clustername
	}

	return container, nil
}

//find the master container in a cluster, used for serf fail-over event
func GetContainerMaster(dbConn *sql.DB, clusterid string) (Container, error) {
	//logit.Info.Println("admindb:GetContainerMaster:called")
	container := Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from project p, server s, container c   where c.role = 'master' and c.clusterid = %s and c.projectid = p.id and c.serverid = s.id", clusterid)
	logit.Info.Println("admindb:GetContainerMaster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetContainerMaster: no master container with that clusterid " + clusterid)
		return container, err
	case err != nil:
		return container, err
	}

	if container.ClusterID != "-1" {
		var clustername string
		clustername, err = GetClusterName(dbConn, container.ClusterID)
		if err != nil {
			logit.Info.Println("admindb:GetContainerGetMaster:error " + err.Error())
			return container, err
		}
		logit.Info.Println("admindb:GetContainerGetMaster:clustername  " + clustername)
		container.ClusterName = clustername
	}

	return container, nil
}

//find the pgpool container in a cluster
func GetContainerPgpool(dbConn *sql.DB, clusterid string) (Container, error) {
	//logit.Info.Println("admindb:GetContainerMaster:called")
	container := Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from project p, server s, container c  where c.serverid = s.id and c.role = 'pgpool' and c.clusterid = %s and c.projectid = p.id", clusterid)
	logit.Info.Println("admindb:GetContainerPgpool:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetContainerPgpool: no pgpool container with that clusterid " + clusterid)
		return container, err
	case err != nil:
		return container, err
	}
	if container.ClusterID != "-1" {
		var clustername string
		clustername, err = GetClusterName(dbConn, container.ClusterID)
		if err != nil {
			logit.Info.Println("admindb:GetContainerPgPool:error " + err.Error())
			return container, err
		}
		container.ClusterName = clustername
	}

	return container, nil
}

//
// TODO combine with GetMaster into a GetContainersByRole func
//
func GetAllStandbyContainers(dbConn *sql.DB, clusterid string) ([]Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from project p, server s, container c  where c.role = 'standby' and c.clusterid = %s and c.projectid = p.id and c.serverid = s.id", clusterid)
	logit.Info.Println("admindb:GetAllStandbyContainers:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clustername string
	containers := make([]Container, 0)
	for rows.Next() {
		container := Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName); err != nil {
			return nil, err
		}
		if container.ClusterID != "-1" {
			clustername, err = GetClusterName(dbConn, container.ClusterID)
			if err != nil {
				logit.Info.Println("admindb:GetContainerPgPool:error " + err.Error())
				return nil, err
			}
			container.ClusterName = clustername
		}
		containers = append(containers, container)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return containers, nil
}

func GetAllContainersForServer(dbConn *sql.DB, serverID string) ([]Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from project p, server s, container c where c.serverid = %s  and c.projectid = p.id  and c.serverid = s.id order by c.name", serverID)
	logit.Info.Println("admindb:GetAllContainersForServer:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clustername string
	containers := make([]Container, 0)
	for rows.Next() {
		container := Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName); err != nil {
			return nil, err
		}
		if container.ClusterID != "-1" {
			clustername, err = GetClusterName(dbConn, container.ClusterID)
			if err != nil {
				logit.Info.Println("admindb:GetContainerForServer:error " + err.Error())
				return nil, err
			}
			container.ClusterName = clustername
		}
		containers = append(containers, container)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return containers, nil
}

func GetAllContainersForCluster(dbConn *sql.DB, clusterID string) ([]Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from project p, server s, container c  where c.clusterid = %s  and c.projectid = p.id and c.serverid = s.id order by c.name", clusterID)
	logit.Info.Println("admindb:GetAllContainersForCluster:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clustername string
	containers := make([]Container, 0)
	for rows.Next() {
		container := Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName); err != nil {
			return nil, err
		}
		if container.ClusterID != "-1" {
			clustername, err = GetClusterName(dbConn, container.ClusterID)
			if err != nil {
				logit.Info.Println("admindb:GetContainerForcluster:error " + err.Error())
				return nil, err
			}
			container.ClusterName = clustername
		}
		containers = append(containers, container)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return containers, nil
}

func GetAllContainersForProject(dbConn *sql.DB, projectID string) ([]Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select id, name from container where projectid = %s order by name", projectID)
	logit.Info.Println("admindb:GetAllContainersForProject:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	containers := make([]Container, 0)
	for rows.Next() {
		container := Container{}
		if err = rows.Scan(&container.ID, &container.Name); err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return containers, nil
}

//
// GetAllContainersNotInCluster is used to fetch all nodes that
// are eligible to be added into a cluster
//
func GetAllContainersNotInCluster(dbConn *sql.DB) ([]Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name, l.name from project p, server s , container c left join cluster l on c.clusterid = l.id where c.role != 'standalone' and c.clusterid = -1 and c.projectid = p.id  and c.serverid = s.id order by c.name")
	logit.Info.Println("admindb:GetAllContainersNotInCluster:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	containers := make([]Container, 0)
	for rows.Next() {
		container := Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName, &container.ClusterName); err != nil {
			return nil, err
		}
		container.ClusterName = container.ClusterID
		containers = append(containers, container)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return containers, nil
}

func GetAllContainers(dbConn *sql.DB) ([]Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.serverid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from project p, server s , container c where c.projectid = p.id  and c.serverid = s.id order by c.name")
	logit.Info.Println("admindb:GetAllContainers:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clustername string
	containers := make([]Container, 0)
	for rows.Next() {
		container := Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.ClusterID, &container.ServerID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ServerName); err != nil {
			return nil, err
		}

		logit.Info.Println("cluster id is [" + container.ClusterID + "]")
		if container.ClusterID != "-1" {
			clustername, err = GetClusterName(dbConn, container.ClusterID)
			if err != nil {
				logit.Info.Println("admindb:GetAllContainers:error " + err.Error())
				return nil, err
			}
			container.ClusterName = clustername
		}
		containers = append(containers, container)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return containers, nil
}

func InsertContainer(dbConn *sql.DB, container Container) (int, error) {
	queryStr := fmt.Sprintf("insert into container ( name, clusterid, serverid, role, image, createdt, projectid) values ( '%s', %s, %s, '%s','%s', now(), %s) returning id", container.Name, container.ClusterID, container.ServerID, container.Role, container.Image, container.ProjectID)

	logit.Info.Println("admindb:InsertContainer:" + queryStr)
	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		logit.Info.Println("admindb:InsertContainer:" + err.Error())
		return -1, err
	default:
		logit.Info.Println("admindb:InsertContainer:container inserted returned is " + strconv.Itoa(nodeid))
	}

	return nodeid, nil
}

func DeleteContainer(dbConn *sql.DB, id string) error {
	queryStr := fmt.Sprintf("delete from container where  id=%s returning id", id)
	logit.Info.Println("admindb:DeleteContainer:" + queryStr)

	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		logit.Error.Println(err)
		return err
	default:
		logit.Info.Println("admindb:DeleteContainer:cluster deleted " + id)
	}
	return nil
}

func UpdateContainer(dbConn *sql.DB, container Container) error {
	queryStr := fmt.Sprintf("update container set ( name, clusterid, serverid, role, image) = ('%s', %s, %s, '%s', '%s') where id = %s returning id", container.Name, container.ClusterID, container.ServerID, container.Role, container.Image, container.ID)
	logit.Info.Println("admindb:UpdateContainer:" + queryStr)

	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:UpdateContainer: container updated " + container.Name)
	}
	return nil

}

func GetAllServers(dbConn *sql.DB) ([]Server, error) {
	logit.Info.Println("admindb:GetAllServer:called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select id, name, ipaddress, dockerbip, pgdatapath, serverclass, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from server order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	servers := make([]Server, 0)
	for rows.Next() {
		server := Server{}
		if err = rows.Scan(&server.ID, &server.Name,
			&server.IPAddress, &server.DockerBridgeIP, &server.PGDataPath, &server.ServerClass, &server.CreateDate); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return servers, nil
}

func GetAllServersByClassByCount(dbConn *sql.DB) ([]Server, error) {
	//select s.id, s.name, s.serverclass, count(n) from server s left join node n on  s.id = n.serverid  group by s.id order by s.serverclass, count(n);

	logit.Info.Println("admindb:GetAllServerByClassByCount:called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select s.id, s.name, s.ipaddress, s.dockerbip, s.pgdatapath, s.serverclass, to_char(s.createdt, 'MM-DD-YYYY HH24:MI:SS'), count(n) from server s left join container n on s.id = n.serverid group by s.id  order by s.serverclass, count(n)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	servers := make([]Server, 0)
	for rows.Next() {
		server := Server{}
		if err = rows.Scan(&server.ID, &server.Name,
			&server.IPAddress, &server.DockerBridgeIP, &server.PGDataPath, &server.ServerClass, &server.CreateDate, &server.NodeCount); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return servers, nil
}

func UpdateServer(dbConn *sql.DB, server Server) error {
	logit.Info.Println("admindb:UpdateServer:called")
	queryStr := fmt.Sprintf("update server set ( name, ipaddress, pgdatapath, serverclass, dockerbip) = ('%s', '%s', '%s', '%s', '%s') where id = %s returning id", server.Name, server.IPAddress, server.PGDataPath, server.ServerClass, server.DockerBridgeIP, server.ID)

	logit.Info.Println("admindb:UpdateServer:update str=" + queryStr)
	var serverid int
	err := dbConn.QueryRow(queryStr).Scan(&serverid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:UpdateServer:server updated " + server.ID)
	}
	return nil

}
func InsertServer(dbConn *sql.DB, server Server) (int, error) {
	//logit.Info.Println("admindb:InsertServer:called")
	queryStr := fmt.Sprintf("insert into server ( name, ipaddress, pgdatapath, serverclass, dockerbip, createdt) values ( '%s', '%s', '%s', '%s', '%s', now()) returning id", server.Name, server.IPAddress, server.PGDataPath, server.ServerClass, server.DockerBridgeIP)

	logit.Info.Println("admindb:InsertServer:" + queryStr)
	var serverid int
	err := dbConn.QueryRow(queryStr).Scan(&serverid)
	switch {
	case err != nil:
		logit.Info.Println("admindb:InsertServer:" + err.Error())
		return -1, err
	default:
		logit.Info.Println("admindb:InsertServer: server inserted returned is " + strconv.Itoa(serverid))
	}

	return serverid, nil
}

func DeleteServer(dbConn *sql.DB, id string) error {
	queryStr := fmt.Sprintf("delete from server where  id=%s returning id", id)
	logit.Info.Println("admindb:DeleteServer:" + queryStr)

	var serverid int
	err := dbConn.QueryRow(queryStr).Scan(&serverid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:DeleteServer:server deleted " + id)
	}
	return nil
}

func GetAllGeneralSettings(dbConn *sql.DB) ([]Setting, error) {
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select name, value, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings where name in ('PG-PORT', 'DOMAIN-NAME', 'DOCKER-REGISTRY', 'ADMIN-URL') order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	settings := make([]Setting, 0)
	for rows.Next() {
		setting := Setting{}
		if err = rows.Scan(
			&setting.Name,
			&setting.Value,
			&setting.UpdateDate); err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return settings, nil
}

func GetAllSettings(dbConn *sql.DB) ([]Setting, error) {
	//logit.Info.Println("admindb:GetAllSettings: called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select name, value, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	settings := make([]Setting, 0)
	for rows.Next() {
		setting := Setting{}
		if err = rows.Scan(
			&setting.Name,
			&setting.Value,
			&setting.UpdateDate); err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return settings, nil
}

func GetSetting(dbConn *sql.DB, key string) (Setting, error) {
	//logit.Info.Println("admindb:GetSetting:called")
	setting := Setting{}

	queryStr := fmt.Sprintf("select value, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings where name = '%s'", key)
	//logit.Info.Println("admindb:GetSetting:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&setting.Value, &setting.UpdateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetSetting: no Setting with that key " + key)
		return setting, err
	case err != nil:
		return setting, err
	}

	return setting, nil
}

func InsertSetting(dbConn *sql.DB, setting Setting) error {
	logit.Info.Println("admindb:InsertSetting:called")
	queryStr := fmt.Sprintf("insert into setting ( name, value, createdt) values ( '%s', '%s', now()) returning name", setting.Name, setting.Value)

	logit.Info.Println("admindb:InsertSetting:" + queryStr)
	var name string
	err := dbConn.QueryRow(queryStr).Scan(&name)
	switch {
	case err != nil:
		logit.Info.Println("admindb:InsertSetting:" + err.Error())
		return err
	default:
	}

	return nil
}

func UpdateSetting(dbConn *sql.DB, setting Setting) error {
	logit.Info.Println("admindb:UpdateSetting:called")
	queryStr := fmt.Sprintf("update settings set ( value, updatedt) = ('%s', now()) where name = '%s'  returning name", setting.Value, setting.Name)

	logit.Info.Println("admindb:UpdateSetting:update str=[" + queryStr + "]")
	var name string
	err := dbConn.QueryRow(queryStr).Scan(&name)
	switch {
	case err != nil:
		logit.Info.Println("admindb:UpdateSetting:" + err.Error())
		return err
	default:
	}
	return nil

}

func GetAllSettingsMap(dbConn *sql.DB) (map[string]string, error) {
	logit.Info.Println("admindb:GetAllSettingsMap: called")
	m := make(map[string]string)

	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select name, value, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings order by name")
	if err != nil {
		return m, err
	}
	defer rows.Close()
	//settings := make([]Setting, 0)
	for rows.Next() {
		setting := Setting{}
		if err = rows.Scan(
			&setting.Name,
			&setting.Value,
			&setting.UpdateDate); err != nil {
			return m, err
		}
		m[setting.Name] = setting.Value
		//settings = append(settings, setting)
	}
	if err = rows.Err(); err != nil {
		return m, err
	}
	return m, nil
}

func Test() {
	logit.Info.Println("hi from Test")
}

func GetDomain(dbConn *sql.DB) (string, error) {
	tmp, err := GetSetting(dbConn, "DOMAIN-NAME")

	if err != nil {
		return "", err
	}
	//we trim off any leading . characters
	domain := strings.Trim(tmp.Value, ".")

	return domain, nil
}

func AddContainerUser(dbConn *sql.DB, s ContainerUser) (int, error) {

	//logit.Info.Println("AddContainerUser called")

	//encrypt the password...passwords at rest are encrypted
	encrypted, err := sec.EncryptPassword(s.Passwd)

	queryStr := fmt.Sprintf("insert into containeruser ( containername, usename, passwd, updatedt) values ( '%s', '%s', '%s',  now()) returning id",
		s.Containername,
		s.Rolname,
		encrypted)

	logit.Info.Println("AddContainerUser:" + queryStr)
	var theID int
	err = dbConn.QueryRow(queryStr).Scan(
		&theID)
	if err != nil {
		logit.Error.Println("error in AddContainerUser query " + err.Error())
		return theID, err
	}

	switch {
	case err != nil:
		logit.Error.Println("AddContainerUser: error " + err.Error())
		return theID, err
	default:
	}

	return theID, nil
}

func DeleteContainerUser(dbConn *sql.DB, containername string, rolname string) error {
	queryStr := fmt.Sprintf("delete from containeruser where  containername='%s'  and usename = '%s' returning id", containername, rolname)
	logit.Info.Println("admindb:DeleteContainerUser:" + queryStr)

	var id int
	err := dbConn.QueryRow(queryStr).Scan(&id)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:DeleteContainerUser:deleted " + containername + ":" + rolname)
	}
	return nil
}

func GetContainerUser(dbConn *sql.DB, containername string, usename string) (ContainerUser, error) {
	var rows *sql.Rows
	var user ContainerUser
	var err error
	queryStr := fmt.Sprintf("select id, passwd, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from containeruser where usename = '%s' and containername = '%s'", usename, containername)
	logit.Info.Println("admindb:GetContainerUser:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return user, err
	}
	defer rows.Close()
	for rows.Next() {
		user.Rolname = usename
		user.Containername = containername
		if err = rows.Scan(&user.ID, &user.Passwd, &user.UpdateDate); err != nil {
			return user, err
		}
	}
	if err = rows.Err(); err != nil {
		return user, err
	}
	var unencrypted string
	unencrypted, err = sec.DecryptPassword(user.Passwd)
	if err != nil {
		return user, err
	}
	user.Passwd = unencrypted
	return user, nil
}

func UpdateContainerUser(dbConn *sql.DB, user ContainerUser) error {
	logit.Info.Println("admindb:UpdateContainerUser encrypting password of " + user.Passwd)
	encrypted, err := sec.EncryptPassword(user.Passwd)
	queryStr := fmt.Sprintf("update containeruser set ( passwd, updatedt) = ('%s', now()) where usename = '%s' returning id", encrypted, user.Rolname)

	logit.Info.Println("[" + queryStr + "]")
	var userid int
	err = dbConn.QueryRow(queryStr).Scan(&userid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:UpdateContainerUser:updated " + user.ID)
	}
	return nil

}

func GetProject(dbConn *sql.DB, id string) (Project, error) {
	//logit.Info.Println("GetProject called with id=" + id)
	project := Project{}

	err := dbConn.QueryRow(fmt.Sprintf("select id, name, description, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from project where id=%s", id)).Scan(&project.ID, &project.Name, &project.Desc, &project.UpdateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetProject:no server with that id")
		return project, err
	case err != nil:
		logit.Info.Println("admindb:GetProject:" + err.Error())
		return project, err
	default:
		logit.Info.Println("admindb:GetProject: Project name returned is " + project.Name)
	}

	project.Containers = make(map[string]string)
	var containers []Container
	containers, err = GetAllContainersForProject(dbConn, project.ID)
	if err != nil {
		logit.Info.Println("admindb:GetProject:" + err.Error())
		return project, err
	}

	for i := range containers {
		project.Containers[containers[i].ID] = containers[i].Name
	}
	project.Clusters = make(map[string]string)
	var clusters []Cluster
	clusters, err = GetAllClustersForProject(dbConn, project.ID)
	if err != nil {
		logit.Info.Println("admindb:GetProject:" + err.Error())
		return project, err
	}

	for i := range clusters {
		project.Clusters[clusters[i].ID] = clusters[i].Name
	}

	return project, nil
}

func GetAllProjects(dbConn *sql.DB) ([]Project, error) {
	//logit.Info.Println("admindb:GetAllProjects: called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select id, name, description, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from project order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	projects := make([]Project, 0)
	var containers []Container
	var clusters []Cluster

	for rows.Next() {
		project := Project{}
		if err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Desc,
			&project.UpdateDate); err != nil {
			return nil, err
		}
		project.Containers = make(map[string]string)

		containers, err = GetAllContainersForProject(dbConn, project.ID)
		if err != nil {
			logit.Info.Println("admindb:GetAllProjects:" + err.Error())
			return projects, err
		}

		for i := range containers {
			project.Containers[containers[i].ID] = containers[i].Name
		}
		project.Clusters = make(map[string]string)

		clusters, err = GetAllClustersForProject(dbConn, project.ID)
		if err != nil {
			logit.Info.Println("admindb:GetAllProjects:" + err.Error())
			return projects, err
		}

		for i := range clusters {
			project.Clusters[clusters[i].ID] = clusters[i].Name
		}
		projects = append(projects, project)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func UpdateProject(dbConn *sql.DB, project Project) error {
	//logit.Info.Println("admindb:UpdateProject:called")
	queryStr := fmt.Sprintf("update project set ( name, description, updatedt) = ('%s', '%s', now()) where id = %s returning id", project.Name, project.Desc, project.ID)

	logit.Info.Println("admindb:UpdateProject:update str=[" + queryStr + "]")
	var projectid int
	err := dbConn.QueryRow(queryStr).Scan(&projectid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:UpdateProject:project updated " + project.ID)
	}
	return nil

}

func DeleteProject(dbConn *sql.DB, id string) error {
	queryStr := fmt.Sprintf("delete from project where  id=%s returning id", id)
	//logit.Info.Println("admindb:DeleteProject:" + queryStr)

	var projectid int
	err := dbConn.QueryRow(queryStr).Scan(&projectid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:DeleteProject:project deleted " + id)
	}
	return nil
}

func InsertProject(dbConn *sql.DB, project Project) (int, error) {
	//logit.Info.Println("admindb:InsertProject:called")
	queryStr := fmt.Sprintf("insert into project ( name, description, updatedt) values ( '%s', '%s', now()) returning id", project.Name, project.Desc)

	logit.Info.Println("admindb:InsertProject:" + queryStr)
	var projectid int
	err := dbConn.QueryRow(queryStr).Scan(&projectid)
	switch {
	case err != nil:
		logit.Info.Println("admindb:InsertProject:" + err.Error())
		return -1, err
	default:
		logit.Info.Println("admindb:InsertProject: Project inserted returned is " + strconv.Itoa(projectid))
	}

	return projectid, nil
}
