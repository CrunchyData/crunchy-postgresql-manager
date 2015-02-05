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
	"crunchy.com/logutil"
	"database/sql"
	"fmt"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

type DBSetting struct {
	Name       string
	Value      string
	UpdateDate string
}
type DBServer struct {
	ID             string
	Name           string
	IPAddress      string
	DockerBridgeIP string
	PGDataPath     string
	ServerClass    string
	CreateDate     string
	NodeCount      string
}

type DBCluster struct {
	ID          string
	Name        string
	ClusterType string
	Status      string
	CreateDate  string
}

type DBClusterNode struct {
	ID         string
	ClusterID  string
	ServerID   string
	Name       string
	Role       string
	Image      string
	CreateDate string
}

type DBLinuxStats struct {
	ID        string
	ClusterID string
	Stats     string
}

type DBPGStats struct {
	ID        string
	ClusterID string
	Stats     string
}

var dbConn *sql.DB

func SetConnection(conn *sql.DB) {
	logutil.Log("admindb:SetConnection: called to open dbConn")
	dbConn = conn

}

func GetDBServer(id string) (DBServer, error) {
	logutil.Log("GetDBServer called with id=" + id)
	server := DBServer{}

	err := dbConn.QueryRow(fmt.Sprintf("select id, name, ipaddress, dockerbip, pgdatapath, serverclass, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from server where id=%s", id)).Scan(&server.ID, &server.Name, &server.IPAddress, &server.DockerBridgeIP, &server.PGDataPath, &server.ServerClass, &server.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logutil.Log("admindb:GetDBServer:no server with that id")
		return server, err
	case err != nil:
		logutil.Log("admindb:GetDBServer:" + err.Error())
		return server, err
	default:
		logutil.Log("admindb:GetDBServer: server name returned is " + server.Name)
	}

	return server, nil
}

func GetDBCluster(id string) (DBCluster, error) {
	logutil.Log("admindb:GetDBCluster: called")
	cluster := DBCluster{}

	err := dbConn.QueryRow(fmt.Sprintf("select id, name, clustertype, status, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from cluster where id=%s", id)).Scan(&cluster.ID, &cluster.Name, &cluster.ClusterType, &cluster.Status, &cluster.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logutil.Log("admindb:GetDBCluster:no cluster with that id")
		return cluster, err
	case err != nil:
		logutil.Log("admindb:GetDBCluster:" + err.Error())
		return cluster, err
	default:
		logutil.Log("admindb:GetDBCluster: cluster name returned is " + cluster.Name)
	}

	return cluster, nil
}

func GetAllDBClusters() ([]DBCluster, error) {
	logutil.Log("admindb:GetAllDBClusters: called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select id, name, clustertype, status, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from cluster order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	clusters := make([]DBCluster, 0)
	for rows.Next() {
		cluster := DBCluster{}
		if err = rows.Scan(
			&cluster.ID,
			&cluster.Name,
			&cluster.ClusterType,
			&cluster.Status, &cluster.CreateDate); err != nil {
			return nil, err
		}
		clusters = append(clusters, cluster)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return clusters, nil
}

func UpdateDBCluster(cluster DBCluster) error {
	logutil.Log("admindb:UpdateCluster:called")
	queryStr := fmt.Sprintf("update cluster set ( name, clustertype, status) = ('%s', '%s', '%s') where id = %s returning id", cluster.Name, cluster.ClusterType, cluster.Status, cluster.ID)

	logutil.Log("admindb:UpdateDBCluster:update str=[" + queryStr + "]")
	var clusterid int
	err := dbConn.QueryRow(queryStr).Scan(&clusterid)
	switch {
	case err != nil:
		return err
	default:
		logutil.Log("admindb:UpdateDBCluster:cluster updated " + cluster.ID)
	}
	return nil

}
func InsertDBCluster(cluster DBCluster) (int, error) {
	logutil.Log("admindb:InsertCluster:called")
	queryStr := fmt.Sprintf("insert into cluster ( name, clustertype, status, createdt) values ( '%s', '%s', '%s', now()) returning id", cluster.Name, cluster.ClusterType, cluster.Status)

	logutil.Log("admindb:InsertCluster:" + queryStr)
	var clusterid int
	err := dbConn.QueryRow(queryStr).Scan(&clusterid)
	switch {
	case err != nil:
		logutil.Log("admindb:InsertCluster:" + err.Error())
		return -1, err
	default:
		logutil.Log("admindb:InsertCluster: cluster inserted returned is " + strconv.Itoa(clusterid))
	}

	return clusterid, nil
}

func DeleteDBCluster(id string) error {
	queryStr := fmt.Sprintf("delete from cluster where  id=%s returning id", id)
	logutil.Log("admindb:DeleteDBCluster:" + queryStr)

	var clusterid int
	err := dbConn.QueryRow(queryStr).Scan(&clusterid)
	switch {
	case err != nil:
		return err
	default:
		logutil.Log("admindb:DeleteDBCluster:cluster deleted " + id)
	}
	return nil
}

func GetDBNode(id string) (DBClusterNode, error) {
	logutil.Log("admindb:GetDBNode:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where id=%s", id)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logutil.Log("admindb:GetDBNode:no node with that id " + id)
		return node, err
	case err != nil:
		return node, err
	}

	return node, nil
}

func GetDBNodeByName(name string) (DBClusterNode, error) {
	logutil.Log("admindb:GetNodeByName:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where name='%s'", name)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logutil.Log("admindb:GetDBNodeByName:no node with that name " + name)
		return node, err
	case err != nil:
		return node, err
	}

	return node, nil
}

//find the oldest node in a cluster, used for serf join-cluster event
func GetDBNodeOldestInCluster(clusterid string) (DBClusterNode, error) {
	logutil.Log("admindb:GetNodeOldestInCluster:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where createdt = (select max(createdt) from node where clusterid = %s)", clusterid)
	logutil.Log("admindb:GetNodeOldestInCluster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logutil.Log("admindb:GetDBNodeOldestInCluster: no node with that clusterid " + clusterid)
		return node, err
	case err != nil:
		return node, err
	}

	return node, nil
}

//find the master node in a cluster, used for serf fail-over event
func GetDBNodeMaster(clusterid string) (DBClusterNode, error) {
	logutil.Log("admindb:GetDBNodeMaster:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where noderole = 'master' and clusterid = %s", clusterid)
	logutil.Log("admindb:GetDBNodeMaster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logutil.Log("admindb:GetDBNodeMaster: no master node with that clusterid " + clusterid)
		return node, err
	case err != nil:
		return node, err
	}

	return node, nil
}

//find the pgpool node in a cluster
func GetDBNodePgpool(clusterid string) (DBClusterNode, error) {
	logutil.Log("admindb:GetDBNodeMaster:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where noderole = 'pgpool' and clusterid = %s", clusterid)
	logutil.Log("admindb:GetDBNodeMaster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logutil.Log("admindb:GetDBNodeMaster: no pgpool node with that clusterid " + clusterid)
		return node, err
	case err != nil:
		return node, err
	}

	return node, nil
}

//
// TODO combine with GetMaster into a GetNodesByRole func
//
func GetAllDBStandbyNodes(clusterid string) ([]DBClusterNode, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where noderole = 'standby' and clusterid = %s", clusterid)
	logutil.Log("admindb:GetAllDBStandbyNodes:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	nodes := make([]DBClusterNode, 0)
	for rows.Next() {
		node := DBClusterNode{}
		if err = rows.Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

func GetAllDBNodesForServer(serverID string) ([]DBClusterNode, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where serverid = %s order by name", serverID)
	logutil.Log("admindb:GetAllDBNodesForServer:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	nodes := make([]DBClusterNode, 0)
	for rows.Next() {
		node := DBClusterNode{}
		if err = rows.Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

func GetAllDBNodesForCluster(clusterID string) ([]DBClusterNode, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where clusterid = %s order by name", clusterID)
	logutil.Log("admindb:GetAllDBNodesForCluster:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	nodes := make([]DBClusterNode, 0)
	for rows.Next() {
		node := DBClusterNode{}
		if err = rows.Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

//
// GetAllDBNodesNotInCluster is used to fetch all nodes that
// are eligible to be added into a cluster
//
func GetAllDBNodesNotInCluster() ([]DBClusterNode, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where noderole != 'standalone' and clusterid = -1 order by name")
	logutil.Log("admindb:GetAllDBNodesNotInCluster:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	nodes := make([]DBClusterNode, 0)
	for rows.Next() {
		node := DBClusterNode{}
		if err = rows.Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

func GetAllDBNodes() ([]DBClusterNode, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node order by name")
	logutil.Log("admindb:GetAllDBNodes:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	nodes := make([]DBClusterNode, 0)
	for rows.Next() {
		node := DBClusterNode{}
		if err = rows.Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

func InsertDBNode(node DBClusterNode) (int, error) {
	queryStr := fmt.Sprintf("insert into node ( name, clusterid, serverid, noderole, image, createdt) values ( '%s', %s, %s, '%s','%s', now()) returning id", node.Name, node.ClusterID, node.ServerID, node.Role, node.Image)

	logutil.Log("admindb:InsertDBNode:" + queryStr)
	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		logutil.Log("admindb:InsertDBNode:" + err.Error())
		return -1, err
	default:
		logutil.Log("admindb:InsertDBNode:node inserted returned is " + strconv.Itoa(nodeid))
	}

	return nodeid, nil
}

func DeleteDBNode(id string) error {
	queryStr := fmt.Sprintf("delete from node where  id=%s returning id", id)
	logutil.Log("admindb:DeleteDBNode:" + queryStr)

	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		log.Println(err)
		return err
	default:
		logutil.Log("admindb:DeleteDBNode:cluster deleted " + id)
	}
	return nil
}

func UpdateDBNode(node DBClusterNode) error {
	queryStr := fmt.Sprintf("update node set ( name, clusterid, serverid, noderole, image) = ('%s', %s, %s, '%s', '%s') where id = %s returning id", node.Name, node.ClusterID, node.ServerID, node.Role, node.Image, node.ID)
	logutil.Log("admindb:UpdateDBNode:" + queryStr)

	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		return err
	default:
		logutil.Log("admindb:UpdateDBNode: node updated " + node.ID)
	}
	return nil

}

func GetAllDBServers() ([]DBServer, error) {
	logutil.Log("admindb:GetAllDBServer:called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select id, name, ipaddress, dockerbip, pgdatapath, serverclass, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from server order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	servers := make([]DBServer, 0)
	for rows.Next() {
		server := DBServer{}
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

func GetAllDBServersByClassByCount() ([]DBServer, error) {
	//select s.id, s.name, s.serverclass, count(n) from server s left join node n on  s.id = n.serverid  group by s.id order by s.serverclass, count(n);

	logutil.Log("admindb:GetAllDBServerByClassByCount:called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select s.id, s.name, s.ipaddress, s.dockerbip, s.pgdatapath, s.serverclass, to_char(s.createdt, 'MM-DD-YYYY HH24:MI:SS'), count(n) from server s left join node n on s.id = n.serverid group by s.id  order by s.serverclass, count(n)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	servers := make([]DBServer, 0)
	for rows.Next() {
		server := DBServer{}
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

func UpdateDBServer(server DBServer) error {
	logutil.Log("admindb:UpdateServer:called")
	queryStr := fmt.Sprintf("update server set ( name, ipaddress, pgdatapath, serverclass, dockerbip) = ('%s', '%s', '%s', '%s', '%s') where id = %s returning id", server.Name, server.IPAddress, server.PGDataPath, server.ServerClass, server.DockerBridgeIP, server.ID)

	logutil.Log("admindb:UpdateDBServer:update str=" + queryStr)
	var serverid int
	err := dbConn.QueryRow(queryStr).Scan(&serverid)
	switch {
	case err != nil:
		return err
	default:
		logutil.Log("admindb:UpdateDBServer:server updated " + server.ID)
	}
	return nil

}
func InsertDBServer(server DBServer) (int, error) {
	logutil.Log("admindb:InsertServer:called")
	queryStr := fmt.Sprintf("insert into server ( name, ipaddress, pgdatapath, serverclass, dockerbip, createdt) values ( '%s', '%s', '%s', '%s', '%s', now()) returning id", server.Name, server.IPAddress, server.PGDataPath, server.ServerClass, server.DockerBridgeIP)

	logutil.Log("admindb:InsertServer:" + queryStr)
	var serverid int
	err := dbConn.QueryRow(queryStr).Scan(&serverid)
	switch {
	case err != nil:
		logutil.Log("admindb:InsertServer:" + err.Error())
		return -1, err
	default:
		logutil.Log("admindb:InsertServer: server inserted returned is " + strconv.Itoa(serverid))
	}

	return serverid, nil
}

func DeleteDBServer(id string) error {
	queryStr := fmt.Sprintf("delete from server where  id=%s returning id", id)
	logutil.Log("admindb:DeleteDBServer:" + queryStr)

	var serverid int
	err := dbConn.QueryRow(queryStr).Scan(&serverid)
	switch {
	case err != nil:
		return err
	default:
		logutil.Log("admindb:DeleteDBServer:server deleted " + id)
	}
	return nil
}

func GetAllDBSettings() ([]DBSetting, error) {
	logutil.Log("admindb:GetAllDBSettings: called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select name, value, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	settings := make([]DBSetting, 0)
	for rows.Next() {
		setting := DBSetting{}
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

func GetDBSetting(key string) (DBSetting, error) {
	logutil.Log("admindb:GetDBSetting:called")
	setting := DBSetting{}

	queryStr := fmt.Sprintf("select value, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings where name = '%s'", key)
	logutil.Log("admindb:GetDBSetting:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&setting.Value, &setting.UpdateDate)
	switch {
	case err == sql.ErrNoRows:
		logutil.Log("admindb:GetDBSetting: no Setting with that key " + key)
		return setting, err
	case err != nil:
		return setting, err
	}

	return setting, nil
}

func InsertDBSetting(setting DBSetting) error {
	logutil.Log("admindb:InsertSetting:called")
	queryStr := fmt.Sprintf("insert into setting ( name, value, createdt) values ( '%s', '%s', now()) returning name", setting.Name, setting.Value)

	logutil.Log("admindb:InsertSetting:" + queryStr)
	var name string
	err := dbConn.QueryRow(queryStr).Scan(&name)
	switch {
	case err != nil:
		logutil.Log("admindb:InsertSetting:" + err.Error())
		return err
	default:
	}

	return nil
}

func UpdateDBSetting(setting DBSetting) error {
	logutil.Log("admindb:UpdateSetting:called")
	queryStr := fmt.Sprintf("update settings set ( value, updatedt) = ('%s', now()) where name = '%s'  returning name", setting.Value, setting.Name)

	logutil.Log("admindb:UpdateDBSetting:update str=[" + queryStr + "]")
	var name string
	err := dbConn.QueryRow(queryStr).Scan(&name)
	switch {
	case err != nil:
		logutil.Log("admindb:UpdateDBSetting:" + err.Error())
		return err
	default:
	}
	return nil

}

func GetAllDBSettingsMap() (map[string]string, error) {
	logutil.Log("admindb:GetAllDBSettingsMap: called")
	m := make(map[string]string)

	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select name, value, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings order by name")
	if err != nil {
		return m, err
	}
	defer rows.Close()
	//settings := make([]DBSetting, 0)
	for rows.Next() {
		setting := DBSetting{}
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
	glog.Info("hi from Test")
}
