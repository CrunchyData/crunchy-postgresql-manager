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

package template

import (
	"bytes"
	"database/sql"
	"errors"
	"flag"
	"github.com/crunchydata/crunchy-postgresql-manager/admindb"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
	"github.com/crunchydata/crunchy-postgresql-manager/util"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
)

var templateTypeFlag = flag.String("t", "", "the template type, can be 'hba','postgresql', or 'recovery'")
var pgpooluseriplistFlag = flag.String("pgpooluseriplist", "", "list of pgpool Users and IP addresses")
var useriplistFlag = flag.String("useriplist", "", "list of pg Users and IP addresses")
var bridgeiplistFlag = flag.String("bridgeiplist", "", "list of bridge IP addresses")
var pgUserFlag = flag.String("pguserid", "", "pg user ID")

var outputFile *os.File

type Rule struct {
	Type     string
	Database string
	User     string
	Address  string
	Method   string
}

type HBAParameters struct {
	ADMIN_HOST     string
	MONITOR_HOST   string
	BACKUP_HOST    string
	PG_HOST_IP     string
	USER           string
	USER_IP_LIST   []string
	PGPOOL_HOST    string
	STANDBY_LIST   []string
	BRIDGE_IP_LIST []string
	SERVER_IP_LIST []string
	RULES_LIST     []Rule
}

type PostgresqlParameters struct {
	PG_HOST_IP   string
	PG_PORT      string
	CLUSTER_TYPE string
}
type RecoveryParameters struct {
	USER       string
	PG_HOST_IP string
	PG_PORT    string
}

type PGPoolParameters struct {
	HOST_LIST []string
}

func Postgresql(mode string, port string, clusterType string) (string, error) {
	var info PostgresqlParameters
	info.PG_PORT = port
	info.CLUSTER_TYPE = clusterType

	var path string
	switch mode {
	case "standalone", "master", "standby":
		path = util.GetBase() + "/conf/" + mode + "/postgresql.conf.template"
	default:
		return "", errors.New("invalid mode in processPostgresql of " + mode)
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("postgresql").Parse(string(contents))
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	err = tmpl.Execute(buff, info)

	return buff.String(), nil
}

func Hba(dbConn *sql.DB, mode string, hostname string, port string, clusterid string, domainname string, cars []Rule) (string, error) {

	var hbaInfo HBAParameters

	hbaInfo.PG_HOST_IP = hostname + "." + domainname
	hbaInfo.BACKUP_HOST = hostname + "-backup." + domainname
	hbaInfo.MONITOR_HOST = "cpm-mon." + domainname
	hbaInfo.ADMIN_HOST = "cpm-admin." + domainname
	hbaInfo.RULES_LIST = cars

	servers, err := admindb.GetAllServers(dbConn)
	if err != nil {
		logit.Error.Println("Hba:" + err.Error())
		return "", err
	}

	i := 0
	var allservers = ""
	var allbridges = ""
	for i = range servers {
		logit.Info.Println("Hba:" + servers[i].IPAddress)
		if allservers == "" {
			allservers = servers[i].IPAddress
			allbridges = servers[i].DockerBridgeIP
		} else {
			allservers = allservers + ":" + servers[i].IPAddress
			allbridges = allbridges + ":" + servers[i].DockerBridgeIP
		}
	}
	logit.Info.Println("Hba:processing serverlist=" + allservers)
	hbaInfo.SERVER_IP_LIST = strings.Split(allservers, ":")
	hbaInfo.BRIDGE_IP_LIST = strings.Split(allbridges, ":")

	var path string
	switch mode {
	case "unassigned":
		path = util.GetBase() + "/conf/standalone/pg_hba.conf.template"
	case "standalone", "master", "standby":
		path = util.GetBase() + "/conf/" + mode + "/pg_hba.conf.template"
	default:
		return "", errors.New("invalid mode in processHba of " + mode)
	}

	if mode == "standby" || mode == "master" {
		_, pgpoolNode, standbyList, err := getMasterValues(dbConn, clusterid, domainname)
		if err != nil {
			return "", err
		}

		hbaInfo.PGPOOL_HOST = pgpoolNode.Name + "." + domainname
		hbaInfo.STANDBY_LIST = standbyList
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("hba").Parse(string(contents))
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")

	logInfo(hbaInfo)

	err = tmpl.Execute(buff, hbaInfo)
	logit.Info.Println("Hba:" + buff.String())

	return buff.String(), nil
}

//
// getMasterValues returns a master node, pgpool node, and list of standby nodes
//
func getMasterValues(dbConn *sql.DB, clusterID string, domainname string) (admindb.Container, admindb.Container, []string, error) {
	master := admindb.Container{}
	pgpool := admindb.Container{}
	//we pass in a list of containers in this cluster
	//that will be added to the pg_hba.conf of the master
	//for allowing replication
	nodes, err1 := admindb.GetAllContainersForCluster(dbConn, clusterID)
	if err1 != nil {
		return master, pgpool, make([]string, 1), err1
	}

	masterFound := false
	pgpoolFound := false
	//nodelist := ""
	i := 0
	nodeslice := make([]string, 10, 10)
	var nodecount = 0

	for i = range nodes {
		if nodes[i].Role == "master" {
			master = nodes[i]
			masterFound = true
			nodeslice[nodecount] = nodes[i].Name + "." + domainname
			nodecount++
		} else if nodes[i].Role == "pgpool" {
			pgpool = nodes[i]
			pgpoolFound = true
		} else if nodes[i].Role == "standby" {
			nodeslice[nodecount] = nodes[i].Name + "." + domainname
			nodecount++
		}
		i++
	}

	if masterFound == false {
		return master, pgpool, make([]string, 1), errors.New("no master found in this cluster")
	}
	if pgpoolFound == false {
		return master, pgpool, make([]string, 1), errors.New("no pgpool found in this cluster")
	}

	nodelist := make([]string, nodecount, nodecount)
	copy(nodelist, nodeslice)
	return master, pgpool, nodelist, nil
}
func logInfo(info HBAParameters) {
	logit.Info.Println("HBA Parameters are:")
	logit.Info.Println("PG_HOST_IP=" + info.PG_HOST_IP)
	logit.Info.Println("USER=" + info.USER)
	i := 0
	for i = range info.USER_IP_LIST {
		logit.Info.Println("USER_IP_LIST[" + strconv.Itoa(i) + "]=" + info.USER_IP_LIST[i])
		i++
	}
	logit.Info.Println("PGPOOL_HOST=" + info.PGPOOL_HOST)
	i = 0
	for i = range info.STANDBY_LIST {
		logit.Info.Println("STANDBY_LIST[" + strconv.Itoa(i) + "]=" + info.STANDBY_LIST[i])
		i++
	}
	i = 0
	for i = range info.BRIDGE_IP_LIST {
		logit.Info.Println("BRIDGE_IP_LIST[" + strconv.Itoa(i) + "]=" + info.BRIDGE_IP_LIST[i])
		i++
	}
}

func Recovery(masterhost string, port string, masteruser string) (string, error) {
	var info RecoveryParameters
	info.PG_PORT = port
	info.USER = masteruser
	info.PG_HOST_IP = masterhost

	var path string
	path = util.GetBase() + "/conf/" + "standby/recovery.conf.template"

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("recovery").Parse(string(contents))
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	err = tmpl.Execute(buff, info)

	return buff.String(), nil
}

//
// Poolhba right now this is simple, just read the template and spit it back
// out, no substitutions are done right now, they will be in the future no doubt
//
func Poolhba() (string, error) {

	var info RecoveryParameters

	var path = util.GetBase() + "/conf/" + "pgpool/pool_hba.conf.template"

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("pgpoolhba").Parse(string(contents))
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")

	err = tmpl.Execute(buff, info)
	logit.Info.Println("Poolhba:" + buff.String())

	return buff.String(), nil
}

//
// Poolpasswd right now this is simple, just read the template and spit it back
// out, no substitutions are done right now, they will be in the future no doubt
//
func Poolpasswd() (string, error) {

	var info RecoveryParameters

	var path = util.GetBase() + "/conf/" + "pgpool/pool_passwd.template"

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("pgpoolpasswd").Parse(string(contents))
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")

	err = tmpl.Execute(buff, info)
	logit.Info.Println("Poolpasswd:" + buff.String())

	return buff.String(), nil
}

//
// Poolconf generates a pgpool.conf file from a template and
// values passed in
//
func Poolconf(poolnames []string) (string, error) {

	var poolParams PGPoolParameters
	poolParams.HOST_LIST = poolnames
	var path = util.GetBase() + "/conf/" + "pgpool/pgpool.conf.template"

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("pgpoolconf").Parse(string(contents))
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")

	err = tmpl.Execute(buff, poolParams)
	logit.Info.Println("Poolconf:" + buff.String())

	return buff.String(), nil
}
