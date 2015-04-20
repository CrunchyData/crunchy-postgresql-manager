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

type DBNodeUser struct {
	ID            string
	Containername string
	Usename       string
	Passwd        string
	UpdateDate    string
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
	dbConn = conn
}

func GetDBServer(id string) (DBServer, error) {
	//logit.Info.Println("GetDBServer called with id=" + id)
	server := DBServer{}

	err := dbConn.QueryRow(fmt.Sprintf("select id, name, ipaddress, dockerbip, pgdatapath, serverclass, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from server where id=%s", id)).Scan(&server.ID, &server.Name, &server.IPAddress, &server.DockerBridgeIP, &server.PGDataPath, &server.ServerClass, &server.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetDBServer:no server with that id")
		return server, err
	case err != nil:
		logit.Info.Println("admindb:GetDBServer:" + err.Error())
		return server, err
	default:
		logit.Info.Println("admindb:GetDBServer: server name returned is " + server.Name)
	}

	return server, nil
}

func GetDBCluster(id string) (DBCluster, error) {
	//logit.Info.Println("admindb:GetDBCluster: called")
	cluster := DBCluster{}

	err := dbConn.QueryRow(fmt.Sprintf("select id, name, clustertype, status, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from cluster where id=%s", id)).Scan(&cluster.ID, &cluster.Name, &cluster.ClusterType, &cluster.Status, &cluster.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetDBCluster:no cluster with that id")
		return cluster, err
	case err != nil:
		logit.Info.Println("admindb:GetDBCluster:" + err.Error())
		return cluster, err
	default:
		logit.Info.Println("admindb:GetDBCluster: cluster name returned is " + cluster.Name)
	}

	return cluster, nil
}

func GetAllDBClusters() ([]DBCluster, error) {
	//logit.Info.Println("admindb:GetAllDBClusters: called")
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
	//logit.Info.Println("admindb:UpdateCluster:called")
	queryStr := fmt.Sprintf("update cluster set ( name, clustertype, status) = ('%s', '%s', '%s') where id = %s returning id", cluster.Name, cluster.ClusterType, cluster.Status, cluster.ID)

	logit.Info.Println("admindb:UpdateDBCluster:update str=[" + queryStr + "]")
	var clusterid int
	err := dbConn.QueryRow(queryStr).Scan(&clusterid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:UpdateDBCluster:cluster updated " + cluster.ID)
	}
	return nil

}
func InsertDBCluster(cluster DBCluster) (int, error) {
	//logit.Info.Println("admindb:InsertCluster:called")
	queryStr := fmt.Sprintf("insert into cluster ( name, clustertype, status, createdt) values ( '%s', '%s', '%s', now()) returning id", cluster.Name, cluster.ClusterType, cluster.Status)

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

func DeleteDBCluster(id string) error {
	queryStr := fmt.Sprintf("delete from cluster where  id=%s returning id", id)
	//logit.Info.Println("admindb:DeleteDBCluster:" + queryStr)

	var clusterid int
	err := dbConn.QueryRow(queryStr).Scan(&clusterid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:DeleteDBCluster:cluster deleted " + id)
	}
	return nil
}

func GetDBNode(id string) (DBClusterNode, error) {
	//logit.Info.Println("admindb:GetDBNode:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where id=%s", id)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetDBNode:no node with that id " + id)
		return node, err
	case err != nil:
		return node, err
	}

	return node, nil
}

func GetDBNodeByName(name string) (DBClusterNode, error) {
	//logit.Info.Println("admindb:GetNodeByName:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where name='%s'", name)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetDBNodeByName:no node with that name " + name)
		return node, err
	case err != nil:
		return node, err
	}

	return node, nil
}

//find the oldest node in a cluster, used for serf join-cluster event
func GetDBNodeOldestInCluster(clusterid string) (DBClusterNode, error) {
	//logit.Info.Println("admindb:GetNodeOldestInCluster:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where createdt = (select max(createdt) from node where clusterid = %s)", clusterid)
	logit.Info.Println("admindb:GetNodeOldestInCluster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetDBNodeOldestInCluster: no node with that clusterid " + clusterid)
		return node, err
	case err != nil:
		return node, err
	}

	return node, nil
}

//find the master node in a cluster, used for serf fail-over event
func GetDBNodeMaster(clusterid string) (DBClusterNode, error) {
	//logit.Info.Println("admindb:GetDBNodeMaster:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where noderole = 'master' and clusterid = %s", clusterid)
	logit.Info.Println("admindb:GetDBNodeMaster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetDBNodeMaster: no master node with that clusterid " + clusterid)
		return node, err
	case err != nil:
		return node, err
	}

	return node, nil
}

//find the pgpool node in a cluster
func GetDBNodePgpool(clusterid string) (DBClusterNode, error) {
	//logit.Info.Println("admindb:GetDBNodeMaster:called")
	node := DBClusterNode{}

	queryStr := fmt.Sprintf("select id, name, clusterid, serverid, noderole, image, to_char(createdt, 'MM-DD-YYYY HH24:MI:SS') from node where noderole = 'pgpool' and clusterid = %s", clusterid)
	logit.Info.Println("admindb:GetDBNodeMaster:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&node.ID, &node.Name, &node.ClusterID, &node.ServerID, &node.Role, &node.Image, &node.CreateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetDBNodeMaster: no pgpool node with that clusterid " + clusterid)
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
	logit.Info.Println("admindb:GetAllDBStandbyNodes:" + queryStr)
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
	logit.Info.Println("admindb:GetAllDBNodesForServer:" + queryStr)
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
	logit.Info.Println("admindb:GetAllDBNodesForCluster:" + queryStr)
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
	logit.Info.Println("admindb:GetAllDBNodesNotInCluster:" + queryStr)
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
	logit.Info.Println("admindb:GetAllDBNodes:" + queryStr)
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

	logit.Info.Println("admindb:InsertDBNode:" + queryStr)
	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		logit.Info.Println("admindb:InsertDBNode:" + err.Error())
		return -1, err
	default:
		logit.Info.Println("admindb:InsertDBNode:node inserted returned is " + strconv.Itoa(nodeid))
	}

	return nodeid, nil
}

func DeleteDBNode(id string) error {
	queryStr := fmt.Sprintf("delete from node where  id=%s returning id", id)
	logit.Info.Println("admindb:DeleteDBNode:" + queryStr)

	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		logit.Error.Println(err)
		return err
	default:
		logit.Info.Println("admindb:DeleteDBNode:cluster deleted " + id)
	}
	return nil
}

func UpdateDBNode(node DBClusterNode) error {
	queryStr := fmt.Sprintf("update node set ( name, clusterid, serverid, noderole, image) = ('%s', %s, %s, '%s', '%s') where id = %s returning id", node.Name, node.ClusterID, node.ServerID, node.Role, node.Image, node.ID)
	logit.Info.Println("admindb:UpdateDBNode:" + queryStr)

	var nodeid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:UpdateDBNode: node updated " + node.ID)
	}
	return nil

}

func GetAllDBServers() ([]DBServer, error) {
	logit.Info.Println("admindb:GetAllDBServer:called")
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

	logit.Info.Println("admindb:GetAllDBServerByClassByCount:called")
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
	logit.Info.Println("admindb:UpdateServer:called")
	queryStr := fmt.Sprintf("update server set ( name, ipaddress, pgdatapath, serverclass, dockerbip) = ('%s', '%s', '%s', '%s', '%s') where id = %s returning id", server.Name, server.IPAddress, server.PGDataPath, server.ServerClass, server.DockerBridgeIP, server.ID)

	logit.Info.Println("admindb:UpdateDBServer:update str=" + queryStr)
	var serverid int
	err := dbConn.QueryRow(queryStr).Scan(&serverid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:UpdateDBServer:server updated " + server.ID)
	}
	return nil

}
func InsertDBServer(server DBServer) (int, error) {
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

func DeleteDBServer(id string) error {
	queryStr := fmt.Sprintf("delete from server where  id=%s returning id", id)
	logit.Info.Println("admindb:DeleteDBServer:" + queryStr)

	var serverid int
	err := dbConn.QueryRow(queryStr).Scan(&serverid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:DeleteDBServer:server deleted " + id)
	}
	return nil
}

func GetAllDBSettings() ([]DBSetting, error) {
	//logit.Info.Println("admindb:GetAllDBSettings: called")
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
	//logit.Info.Println("admindb:GetDBSetting:called")
	setting := DBSetting{}

	queryStr := fmt.Sprintf("select value, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from settings where name = '%s'", key)
	//logit.Info.Println("admindb:GetDBSetting:" + queryStr)
	err := dbConn.QueryRow(queryStr).Scan(&setting.Value, &setting.UpdateDate)
	switch {
	case err == sql.ErrNoRows:
		logit.Info.Println("admindb:GetDBSetting: no Setting with that key " + key)
		return setting, err
	case err != nil:
		return setting, err
	}

	return setting, nil
}

func InsertDBSetting(setting DBSetting) error {
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

func UpdateDBSetting(setting DBSetting) error {
	logit.Info.Println("admindb:UpdateSetting:called")
	queryStr := fmt.Sprintf("update settings set ( value, updatedt) = ('%s', now()) where name = '%s'  returning name", setting.Value, setting.Name)

	logit.Info.Println("admindb:UpdateDBSetting:update str=[" + queryStr + "]")
	var name string
	err := dbConn.QueryRow(queryStr).Scan(&name)
	switch {
	case err != nil:
		logit.Info.Println("admindb:UpdateDBSetting:" + err.Error())
		return err
	default:
	}
	return nil

}

func GetAllDBSettingsMap() (map[string]string, error) {
	logit.Info.Println("admindb:GetAllDBSettingsMap: called")
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
	logit.Info.Println("hi from Test")
}

func GetDomain() (string, error) {
	tmp, err := GetDBSetting("DOMAIN-NAME")

	if err != nil {
		return "", err
	}
	//we trim off any leading . characters
	domain := strings.Trim(tmp.Value, ".")

	return domain, nil
}

func DBAddNodeUser(s DBNodeUser) (int, error) {

	//logit.Info.Println("DBAddNodeUser called")

	//encrypt the password...passwords at rest are encrypted
	encrypted, err := sec.EncryptPassword(s.Passwd)

	queryStr := fmt.Sprintf("insert into nodeuser ( containername, usename, passwd, updatedt) values ( '%s', '%s', '%s',  now()) returning id",
		s.Containername,
		s.Usename,
		encrypted)

	logit.Info.Println("DBAddNodeUser:" + queryStr)
	var theID int
	err = dbConn.QueryRow(queryStr).Scan(
		&theID)
	if err != nil {
		logit.Error.Println("error in DBAddNodeUser query " + err.Error())
		return theID, err
	}

	switch {
	case err != nil:
		logit.Error.Println("DBAddNodeUser: error " + err.Error())
		return theID, err
	default:
	}

	return theID, nil
}

func DBDeleteNodeUser(id string) error {
	queryStr := fmt.Sprintf("delete from nodeuser where id=%s returning id", id)
	//logit.Info.Println("admindb:DeleteDBCluster:" + queryStr)

	var nodeuserid int
	err := dbConn.QueryRow(queryStr).Scan(&nodeuserid)
	switch {
	case err != nil:
		return err
	default:
		logit.Info.Println("admindb:DBDeleteNodeUser: deleted " + id)
	}
	return nil
}

func GetAllUsersForNode(containerName string) ([]DBNodeUser, error) {
	var rows *sql.Rows
	var err error
	queryStr := fmt.Sprintf("select id, usename, passwd, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from nodeuser where containername = '%s' order by usename", containerName)
	logit.Info.Println("admindb:GetAllUsersForNode:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]DBNodeUser, 0)
	for rows.Next() {
		user := DBNodeUser{}
		user.Containername = containerName
		if err = rows.Scan(&user.ID, &user.Usename, &user.Passwd, &user.UpdateDate); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetNodeUser(containername string, usename string) (DBNodeUser, error) {
	var rows *sql.Rows
	var user DBNodeUser
	var err error
	queryStr := fmt.Sprintf("select id, passwd, to_char(updatedt, 'MM-DD-YYYY HH24:MI:SS') from nodeuser where usename = '%s' and containername = '%s'", usename, containername)
	logit.Info.Println("admindb:GetNodeUser:" + queryStr)
	rows, err = dbConn.Query(queryStr)
	if err != nil {
		return user, err
	}
	defer rows.Close()
	for rows.Next() {
		user.Usename = usename
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
