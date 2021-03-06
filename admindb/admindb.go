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
	"github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/sec"
	"github.com/crunchydata/crunchy-postgresql-manager/types"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	_ "github.com/lib/pq"
	"strconv"
)

// GetClusterName returns the name of a cluster based for a given cluster ID
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

// GetCluster returns a cluster object from the database for a given cluster ID
func GetCluster(dbConn *sql.DB, id string) (types.Cluster, error) {
	//logit.Info.Println("admindb:GetCluster: called")
	cluster := types.Cluster{}

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

	var containers []types.Container
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

// GetAllClustersForProject returns a list of cluster objects from the database for a given project
func GetAllClustersForProject(dbConn *sql.DB, projectId string) ([]types.Cluster, error) {
	//logit.Info.Println("admindb:GetAllClusters: called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select id, projectid, name, clustertype, status, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from cluster where projectid = " + projectId + " order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var containers []types.Container
	clusters := make([]types.Cluster, 0)
	for rows.Next() {
		cluster := types.Cluster{}
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

// UpdateCluster updates a given cluster object in the database
func UpdateCluster(dbConn *sql.DB, cluster types.Cluster) error {
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

// InsertCluster inserts a new cluster object into the database
func InsertCluster(dbConn *sql.DB, cluster types.Cluster) (int, error) {
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

// DeleteCluster deletes a given cluster from the database
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

// GetContainer returns a container object based on a container ID
func GetContainer(dbConn *sql.DB, id string) (types.Container, error) {
	//logit.Info.Println("admindb:GetContainer:called")
	container := types.Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name from container c, project p where c.id=%s and p.id = c.projectid ", id)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName)
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

// GetContainerByName returns a container object based on a container name
func GetContainerByName(dbConn *sql.DB, name string) (types.Container, error) {
	//logit.Info.Println("admindb:GetNodeByName:called")
	container := types.Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name from project p, container c where c.name='%s' and c.projectid = p.id ", name)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName)
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

// GetContainerOldestInCluster find the oldest container in a cluster
func GetContainerOldestInCluster(dbConn *sql.DB, clusterid string) (types.Container, error) {
	//logit.Info.Println("admindb:GetNodeOldestInCluster:called")
	container := types.Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid,  c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, s.name from project p, container c where  c.projectid = p.id and c.createdt = (select max(createdt) from container where clusterid = %s)", clusterid)
	logit.Info.Println("admindb:GetNodeOldestInCluster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName)
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

// GetContainerMaster find the master container in a cluster
func GetContainerMaster(dbConn *sql.DB, clusterid string) (types.Container, error) {
	//logit.Info.Println("admindb:GetContainerMaster:called")
	container := types.Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name from project p, container c  where c.role = 'master' and c.clusterid = %s and c.projectid = p.id ", clusterid)
	logit.Info.Println("admindb:GetContainerMaster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName)
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

// GetContainerPgpool find the pgpool container in a cluster
func GetContainerPgpool(dbConn *sql.DB, clusterid string) (types.Container, error) {
	//logit.Info.Println("admindb:GetContainerMaster:called")
	container := types.Container{}

	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name from project p, container c  where c.role = 'pgpool' and c.clusterid = %s and c.projectid = p.id", clusterid)
	logit.Info.Println("admindb:GetContainerPgpool:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&container.ID, &container.Name, &container.ClusterID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName)
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

// GetAllStandbyContainers returns a list of container standby objects for a given cluster
// TODO combine with GetMaster into a GetContainersByRole func
//
func GetAllStandbyContainers(dbConn *sql.DB, clusterid string) ([]types.Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name from project p, container c  where c.role = 'standby' and c.clusterid = %s and c.projectid = p.id ", clusterid)
	logit.Info.Println("admindb:GetAllStandbyContainers:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clustername string
	containers := make([]types.Container, 0)
	for rows.Next() {
		container := types.Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.ClusterID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName); err != nil {
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

// GetAllContainersForServer returns a list of container objects for a given server
func GetAllContainersForServer(dbConn *sql.DB, serverID string) ([]types.Container, error) {
	containers := make([]types.Container, 0)
	return containers, nil
}

// GetAllContainersForCluster returns a list of container objects for a given cluster
func GetAllContainersForCluster(dbConn *sql.DB, clusterID string) ([]types.Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name from project p, container c  where c.clusterid = %s  and c.projectid = p.id order by c.name", clusterID)
	logit.Info.Println("admindb:GetAllContainersForCluster:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clustername string
	containers := make([]types.Container, 0)
	for rows.Next() {
		container := types.Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.ClusterID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName); err != nil {
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

// GetAllContainersForProject returns a list of container objects for a given project
func GetAllContainersForProject(dbConn *sql.DB, projectID string) ([]types.Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name from container c where c.id not in (select p.containerid from proxy p) and c.projectid = %s order by c.name", projectID)
	logit.Info.Println("admindb:GetAllContainersForProject:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	containers := make([]types.Container, 0)
	for rows.Next() {
		container := types.Container{}
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

// GetAllContainersNotInCluster is used to fetch all nodes that are eligible to be added into a cluster
func GetAllContainersNotInCluster(dbConn *sql.DB) ([]types.Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name, l.name from project p, container c left join cluster l on c.clusterid = l.id where c.role != 'standalone' and c.clusterid = -1 and c.projectid = p.id  order by c.name")
	logit.Info.Println("admindb:GetAllContainersNotInCluster:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	containers := make([]types.Container, 0)
	for rows.Next() {
		container := types.Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.ClusterID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName, &container.ClusterName); err != nil {
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

// GetAllContainers returns a list of all containers
func GetAllContainers(dbConn *sql.DB) ([]types.Container, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select c.id, c.name, c.clusterid, c.role, c.image, to_char(c.createdt, 'MM-DD-YYYY HH24:MI:SS'), p.id, p.name from project p, container c where c.projectid = p.id   order by c.name")
	logit.Info.Println("admindb:GetAllContainers:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clustername string
	containers := make([]types.Container, 0)
	for rows.Next() {
		container := types.Container{}
		if err = rows.Scan(&container.ID, &container.Name, &container.ClusterID, &container.Role, &container.Image, &container.CreateDate, &container.ProjectID, &container.ProjectName); err != nil {
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

// InsertContainer inserts a new container object and returns the container id
func InsertContainer(dbConn *sql.DB, container types.Container) (int, error) {
	queryStr := fmt.Sprintf("insert into container ( name, clusterid, role, image, createdt, projectid) values ( '%s', %s,  '%s','%s', now(), %s) returning id", container.Name, container.ClusterID, container.Role, container.Image, container.ProjectID)

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

// DeleteContainer deleles a given container
func DeleteContainer(dbConn *sql.DB, id string) error {
	queryStr := fmt.Sprintf("delete from container where  id=%s returning id", id)
	logit.Info.Println("admindb:DeleteContainer:" + queryStr)

	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return err
	default:
		logit.Info.Println("admindb:DeleteContainer:cluster deleted " + id)
	}
	return nil
}

// UpdateContainer updates a given container
func UpdateContainer(dbConn *sql.DB, container types.Container) error {
	queryStr := fmt.Sprintf("update container set ( name, clusterid, role, image) = ('%s', %s, '%s', '%s') where id = %s returning id", container.Name, container.ClusterID, container.Role, container.Image, container.ID)
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

// GetAllGeneralSettings returns a list of all settings of 'general' types
func GetAllGeneralSettings(dbConn *sql.DB) ([]types.Setting, error) {
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select name, value, description, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings where name in ('PG-PORT', 'DOMAIN-NAME', 'DOCKER-REGISTRY', 'ADMIN-URL') order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	settings := make([]types.Setting, 0)
	for rows.Next() {
		setting := types.Setting{}
		if err = rows.Scan(
			&setting.Name,
			&setting.Value,
			&setting.Description,
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

// GetAllSettings returns a list of all settings
func GetAllSettings(dbConn *sql.DB) ([]types.Setting, error) {
	//logit.Info.Println("admindb:GetAllSettings: called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select name, value, description, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	settings := make([]types.Setting, 0)
	for rows.Next() {
		setting := types.Setting{}
		if err = rows.Scan(
			&setting.Name,
			&setting.Value,
			&setting.Description,
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

// GetSetting returns a particular setting based on it's key
func GetSetting(dbConn *sql.DB, key string) (types.Setting, error) {
	//logit.Info.Println("admindb:GetSetting:called")
	setting := types.Setting{}

	queryStr := fmt.Sprintf("select value, description, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings where name = '%s'", key)
	//logit.Info.Println("admindb:GetSetting:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&setting.Value, &setting.Description, &setting.UpdateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetSetting: no Setting with that key " + key)
		return setting, err
	case err != nil:
		return setting, err
	}

	return setting, nil
}

func InsertSetting(dbConn *sql.DB, setting types.Setting) error {
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

// UpdateSetting updates a given setting value
func UpdateSetting(dbConn *sql.DB, setting types.Setting) error {
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

// GetAllSettingsMap returns a map of all settings
func GetAllSettingsMap(dbConn *sql.DB) (map[string]string, error) {
	logit.Info.Println("admindb:GetAllSettingsMap: called")
	m := make(map[string]string)

	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select name, value, description, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings order by name")
	if err != nil {
		return m, err
	}
	defer rows.Close()
	//settings := make([]Setting, 0)
	for rows.Next() {
		setting := types.Setting{}
		if err = rows.Scan(
			&setting.Name,
			&setting.Value,
			&setting.Description,
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

// AddContainerUser inserts a new database user for a given container and returns the new ID
func AddContainerUser(dbConn *sql.DB, s types.ContainerUser) (int, error) {

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
		logit.Error.Println(err.Error())
		return theID, err
	}

	switch {
	case err != nil:
		logit.Error.Println(err.Error())
		return theID, err
	default:
	}

	return theID, nil
}

// DeleteContainerUser deletes a container database user
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

// GetContainerUser returns a container user for a given container and user name
func GetContainerUser(dbConn *sql.DB, containername string, usename string) (types.ContainerUser, error) {
	var user types.ContainerUser
	var err error
	queryStr := fmt.Sprintf("select id, passwd, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from containeruser where usename = '%s' and containername = '%s'", usename, containername)
	logit.Info.Println("admindb:GetContainerUser:" + queryStr)
	err = dbConn.QueryRow(queryStr).Scan(&user.ID, &user.Passwd, &user.UpdateDate)
	switch {
	case err != nil:
		return user, err
	default:
		logit.Info.Println("GetContainerUser got a row back")
	}

	user.Rolname = usename
	user.Containername = containername

	var unencrypted string
	unencrypted, err = sec.DecryptPassword(user.Passwd)
	if err != nil {
		return user, err
	}
	user.Passwd = unencrypted
	return user, nil
}

// UpdateContainerUser updates a given container database user
func UpdateContainerUser(dbConn *sql.DB, user types.ContainerUser) error {
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

// GetProject returns a given project object
func GetProject(dbConn *sql.DB, id string) (types.Project, error) {
	//logit.Info.Println("GetProject called with id=" + id)
	project := types.Project{}

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
	var containers []types.Container
	containers, err = GetAllContainersForProject(dbConn, project.ID)
	if err != nil {
		logit.Info.Println("admindb:GetProject:" + err.Error())
		return project, err
	}

	for i := range containers {
		project.Containers[containers[i].ID] = containers[i].Name
	}
	project.Clusters = make(map[string]string)
	var clusters []types.Cluster
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

// GetAllProjects returns a list of all project objects
func GetAllProjects(dbConn *sql.DB) ([]types.Project, error) {
	//logit.Info.Println("admindb:GetAllProjects: called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select id, name, description, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from project order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	projects := make([]types.Project, 0)
	var containers []types.Container
	var clusters []types.Cluster
	var proxies []types.Proxy

	for rows.Next() {
		project := types.Project{}
		if err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Desc,
			&project.UpdateDate); err != nil {
			return nil, err
		}

		project.Proxies = make(map[string]string)

		proxies, err = GetAllProxiesForProject(dbConn, project.ID)
		if err != nil {
			logit.Info.Println("admindb:GetAllProjects:" + err.Error())
			return projects, err
		}

		for i := range proxies {
			project.Proxies[proxies[i].ContainerID] = proxies[i].ContainerName
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

// UpdateProject updates a given project object
func UpdateProject(dbConn *sql.DB, project types.Project) error {
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

// DeleteProject deletes a given project
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

// InsertProject inserts a new project object
func InsertProject(dbConn *sql.DB, project types.Project) (int, error) {
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

// GetAllProxiesForProject returns a list of proxy objects for a given project
func GetAllProxiesForProject(dbConn *sql.DB, projectID string) ([]types.Proxy, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select p.containerid, c.name from proxy p, container c where p.containerid = c.id and p.projectid = %s order by c.name", projectID)
	logit.Info.Println("admindb:GetAllproxiesForProject:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	proxies := make([]types.Proxy, 0)
	for rows.Next() {
		proxy := types.Proxy{}
		if err = rows.Scan(&proxy.ContainerID, &proxy.ContainerName); err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return proxies, nil
}

// GetProxy returns a given proxy object by container name
func GetProxy(dbConn *sql.DB, containername string) (types.Proxy, error) {
	var rows *sql.Rows
	proxy := types.Proxy{}
	var err error

	queryStr := fmt.Sprintf("select p.usename , p.passwd, p.port, p.host, p.databasename from proxy p, container c where p.containerid = c.id and c.name = '%s'", containername)

	logit.Info.Println("GetProxy:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return proxy, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&proxy.Usename, &proxy.Passwd,
			&proxy.Port, &proxy.Host, &proxy.Database); err != nil {
			return proxy, err
		}
	}

	proxy.ContainerName = containername

	if err = rows.Err(); err != nil {
		return proxy, err
	}
	var unencrypted string
	unencrypted, err = sec.DecryptPassword(proxy.Passwd)
	if err != nil {
		return proxy, err
	}
	proxy.Passwd = unencrypted
	return proxy, nil
}

// GetProxyByContainerID returns a proxy object by container ID
func GetProxyByContainerID(dbConn *sql.DB, containerID string) (types.Proxy, error) {
	var rows *sql.Rows
	proxy := types.Proxy{}
	var err error

	queryStr := fmt.Sprintf("select p.usename, p.passwd, p.projectid, p.id,  p.containerid, c.name , p.port, p.host, p.databasename from proxy p, container c where p.containerid = c.id and c.id = %s", containerID)

	logit.Info.Println("GetProxyByContainerID:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return proxy, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(
			&proxy.Usename,
			&proxy.Passwd,
			&proxy.ProjectID,
			&proxy.ID, &proxy.ContainerID,
			&proxy.ContainerName, &proxy.Port, &proxy.Host, &proxy.Database); err != nil {
			return proxy, err
		}
	}
	if err = rows.Err(); err != nil {
		return proxy, err
	}
	var unencrypted string
	unencrypted, err = sec.DecryptPassword(proxy.Passwd)
	if err != nil {
		return proxy, err
	}
	proxy.Passwd = unencrypted

	//see if the container is up
	var resp cpmcontainerapi.StatusResponse
	resp, err = cpmcontainerapi.StatusClient(proxy.ContainerName)
	proxy.ContainerStatus = resp.Status
	if err != nil {
		return proxy, err
	}

	//test the remote host connectivity
	var hoststatus string
	hoststatus, err = util.FastPing(proxy.Port, proxy.Host)
	if err != nil {
		return proxy, err
	}
	if hoststatus == "OFFLINE" {
		proxy.Status = hoststatus
		return proxy, err
	}

	//test the database port on the remote host
	proxy.Status, err = GetDatabaseStatus(proxy, containerID)
	if err != nil {
		return proxy, err
	}

	return proxy, nil
}

// GetDatabaseStatus returns a simple status of a container database
func GetDatabaseStatus(proxy types.Proxy, containerid string) (string, error) {
	/**
	node, err := GetContainer(dbConn, containerid)
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}

	var credential types.Credential
	credential, err = GetUserCredentials(dbConn, &node)
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}
	*/

	dbConn2, err := util.GetMonitoringConnection(proxy.Host, proxy.Usename, proxy.Port, proxy.Database, proxy.Passwd)
	defer dbConn2.Close()

	var value string

	err = dbConn2.QueryRow(fmt.Sprintf("select now()::text")).Scan(&value)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("getProxyStatus  no rows returned")
		return "OFFLINE", nil
	case err != nil:
		logit.Info.Println("getProxyStatus error " + err.Error())
		return "OFFLINE", nil
	default:
		logit.Info.Println("getProxyStatus returned " + value)
	}

	return "RUNNING", nil

}
