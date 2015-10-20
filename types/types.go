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
	Host     string
	Database string
	Username string
	Password string
	Port     string
}

type Setting struct {
	Name        string
	Value       string
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
	Proxies    map[string]string
	UpdateDate string
	Token      string
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
	Token       string
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
	Rolcanlogin    string
	Rolreplication string
	UpdateDate     string
}

type Proxy struct {
	ID              string
	ContainerUserID string
	Database        string
	Host            string
	Usename         string
	Passwd          string
	ContainerID     string
	ContainerName   string
	ServerName      string
	Status          string
	ContainerStatus string
	ProjectID       string
	Port            string
	UpdateDate      string
	Token           string
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

type HealthCheck struct {
	ID             string
	ProjectName    string
	ProjectID      string
	ContainerName  string
	ContainerID    string
	ContainerRole  string
	ContainerImage string
	Status         string
	UpdateDate     string
}

type ClusterProfiles struct {
	Size           string
	Count          string
	Algo           string
	MasterProfile  string
	StandbyProfile string
	MasterServer   string
	StandbyServer  string
	Token          string
}
type Profiles struct {
	SmallCPU  string
	SmallMEM  string
	MediumCPU string
	MediumMEM string
	LargeCPU  string
	LargeMEM  string
	Token     string
}

type ClusterNode struct {
	ID          string
	ClusterID   string
	ServerID    string
	Name        string
	Role        string
	Image       string
	CreateDate  string
	Status      string
	ProjectID   string
	ProjectName string
	ServerName  string
	ClusterName string
}

type ProvisionStatus struct {
	Status string
	ID     string
}

type SimpleStatus struct {
	Status string
}

type PostgresStatement struct {
	Database   string
	Query      string
	Calls      string
	TotalTime  string
	Rows       string
	HitPercent string
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
	ID             string
	Containername  string
	Rolname        string
	Passwd         string
	Updatedt       string
	Token          string
	Rolsuper       bool
	Rolinherit     bool
	Rolcreaterole  bool
	Rolcreatedb    bool
	Rollogin       bool
	Rolreplication bool
}

type MonitorServerParam struct {
	ServerID string
	Metric   string
}
type MonitorContainerParam struct {
	ID           string
	Metric       string
	DatabaseName string
}

type MonitorOutput struct {
	Metric   string
	Response string
}
