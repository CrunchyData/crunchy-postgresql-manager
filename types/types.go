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

package types

import ()

type Credential struct {
	Host     string //container name
	Database string //database name
	Username string //database user name
	Password string //database user password
	Port     string //database port number
}

type Setting struct {
	Name        string //setting name (key)
	Value       string //setting value
	Description string
	UpdateDate  string
	Token       string
}

type Settings struct {
	AdminURL       string
	DockerRegistry string
	PGPort         string
	DomainName     string
	Token          string
}
type Server struct {
	ID             string //unique key
	Name           string //server name
	IPAddress      string //server ip address
	DockerBridgeIP string //the docker bridge IP range that is assigned to this server
	ServerClass    string //the class of the server, small, medium, large
	CreateDate     string
	NodeCount      string //calculated value of the number of containers assigned to this server
}

type Project struct {
	ID         string            //unique key
	Name       string            //project name
	Desc       string            //description
	Containers map[string]string //map of containers in this project
	Clusters   map[string]string //map of clusters in this project
	Proxies    map[string]string //map of proxies in this project
	UpdateDate string
	Token      string
}

type Container struct {
	ID          string //unique key
	ClusterID   string //foreign key to cluster
	Name        string //container name, also used as the host name
	Role        string //role of this container (master, standby, standalone)
	Image       string //docker image this container is based on
	CreateDate  string
	ProjectID   string //foreign key to project
	ProjectName string //project name this container is in
	ClusterName string //cluster name this contaienr is part of
}

type Cluster struct {
	ID          string //unique key
	ProjectID   string //foreign key to project
	Name        string //cluster name
	ClusterType string //async or sync type of replication
	Status      string //either up or down or initialized
	CreateDate  string
	Containers  map[string]string //map of containers within this cluster
	Token       string
}

type ContainerUser struct {
	ID             string //unique key
	Containername  string //container name
	ContainerID    string //foreign key to the container
	Passwd         string //database password of this user
	Rolname        string //postgres role name
	Rolsuper       string //postgres superuser permission flag
	Rolinherit     string //postgres inherit permission flag
	Rolcreaterole  string //postgres createrole permission flag
	Rolcreatedb    string //postgres createdb permission flag
	Rolcanlogin    string //postgres login permission flag
	Rolreplication string //postgres replication permission flag
	UpdateDate     string
}

type Proxy struct {
	ID              string //unique key
	ContainerUserID string //foreign key to the container user
	Database        string //database name
	Host            string //database host name
	Usename         string //database user name
	Passwd          string //database password
	ContainerID     string //foreign key to container
	ContainerName   string //container name
	ServerName      string //server name
	Status          string //database status either down or running
	ContainerStatus string //container status either down or running
	ProjectID       string //foreign key to project
	Port            string //database port number
	UpdateDate      string
	Token           string
}

type HealthCheck struct {
	ID             string //unique key
	ProjectName    string //project name
	ProjectID      string //foreign key to project
	ContainerName  string //container name
	ContainerID    string //foreign key to container
	ContainerRole  string //container role
	ContainerImage string //docker image of container
	Status         string //either up or down
	UpdateDate     string
}

type ClusterProfiles struct {
	Size           string //either small, medium, or large
	Count          string //number of standby nodes
	Algo           string //round-robin
	MasterProfile  string //the docker profile to use for the master node
	StandbyProfile string //the docker profile to use for the standby nodes
	MasterServer   string //the server class to place the master on
	StandbyServer  string //the server class to place the standby nodes on
	Token          string
}
type Profiles struct {
	SmallCPU  string //small profile cpu shares setting
	SmallMEM  string //small profile memory setting
	MediumCPU string //medium profile cpu shares setting
	MediumMEM string //medium profile docker memory setting
	LargeCPU  string //large profile docker cpu shares setting
	LargeMEM  string //large profile docker memory setting
	Token     string
}

type ClusterNode struct {
	ID          string //unique key
	ClusterID   string //foreign key to cluster
	Name        string //cluster name
	Role        string //role of node in cluster
	Image       string //docker image of node
	CreateDate  string
	Status      string
	ProjectID   string //foreign key to project
	ProjectName string //project name
	ClusterName string //cluster name
}

type ProvisionStatus struct {
	Status string
	ID     string
}

type SimpleStatus struct {
	Status string
}

type PostgresStatement struct {
	Database   string //pg_stat_statements database value
	Query      string //pg_stat_statements query value
	Calls      string //pg_stat_statements calls value
	TotalTime  string //pg_stat_statements totaltime value
	Rows       string //pg_stat_statements rows value
	HitPercent string //pg_stat_statements hitpercent value
}

type PostgresSetting struct {
	Name           string
	CurrentSetting string
	Source         string
}

type PostgresControldata struct {
	Name  string
	Value string
}

type NodeUser struct {
	ID             string //unique key
	Containername  string //container name
	Rolname        string //role name of this suser
	Passwd         string //password of this user
	Updatedt       string
	Token          string
	Rolsuper       bool //superuser permission flag
	Rolinherit     bool //inherit permission flag
	Rolcreaterole  bool //createrole permission flag
	Rolcreatedb    bool //createdb permission flag
	Rollogin       bool //login permission flag
	Rolreplication bool //replication permission flag
}

type MonitorServerParam struct {
	ServerID string //foreign key to server
	Metric   string
}
type MonitorContainerParam struct {
	ID           string //unique key
	Metric       string //metric name
	DatabaseName string //database name
}

type MonitorOutput struct {
	Metric   string
	Response string
}
