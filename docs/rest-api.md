Authentication
---------------

###GET /sec/login/:ID.:PSW

login and return an auth token

+ ID : the user id to authenticate with
+ PSW : the user password to authenticate with

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /clusters/:Token

+ Token : the generated auth token for this session

returns all clusters


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/clusters/789c31ff-b18f-47b3-bb63-1fd603895aa5
[
  {
    "ID": "2",
    "Name": "aa",
    "ClusterType": "asynchronous",
    "Status": "initialized",
    "CreateDate": "02-11-2015 09:18:13",
    "Token": ""
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /servers/:Token

+ Token : the generated auth token for this session

returns all servers


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/servers/789c31ff-b18f-47b3-bb63-1fd603895aa5
[
  {
    "ID": "1",
    "Name": "espresso",
    "IPAddress": "192.168.0.106",
    "DockerBridgeIP": "",
    "PGDataPath": "/var/lib/pgsql",
    "ServerClass": "low",
    "CreateDate": "02-10-2015 14:35:34"
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /cluster

+ ID : the generated auth token for this session
+ Name : the name to use for this cluster
+ ClusterType : the type of cluster (synchronous|asynchronous)
+ Token : the generated auth token for this session

updates or adds a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

curl --insecure -i -d '{"ID":"1","ClusterID":"1"}' https://cpm-admin.crunchy.lab:13000/event/join-cluster

###POST /autocluster

+ Name : the name to use for this cluster
+ ClusterType : the type of cluster (synchronous|asynchronous)
+ ClusterProfile : the cluster profile to use for cluster creation (SM|LG|MED)
+ Token : the generated auth token for this session

performs an auto-cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /savesettings

+ AdminURL : the URL to cpm-admin to use
+ DockerRegistry : not used yet
+ PGPort :  Postgres port to use within CPM
+ DomainName : domain name to use within CPM
+ Token : the generated auth token for this session

saves settings


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /saveprofiles

+ LargeMEM : large memory value
+ LargeCPU : large cpu value
+ MediumCPU : medium cpu value
+ MediumMEM : medium memory value
+ SmallMEM : small memory value
+ SmallCPU : small cpu value
+ Token : the generated auth token for this session

saves profiles


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /saveclusterprofiles

+ Size : 
+ Count : 
+ Algo : 
+ MasterProfile : 
+ StandbyProfile : 
+ MasterServer : 
+ StandbyServer : 
+ Token : the generated auth token for this session

saves cluster profiles


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /settings/:Token

+ Token : the generated auth token for this session

returns all settings


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/settings/789c31ff-b18f-7b3-bb63-1fd603895aa5
[
  {
    "Name": "CP-LG-ALGO",
    "Value": "round-robin",
    "UpdateDate": "02-10-2015 13:50:57"
  },
  {
    "Name": "CP-LG-COUNT",
    "Value": "1",
    "UpdateDate": "02-10-2015 13:50:57"
  },
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /addserver/:ID.:Name.:IPAddress.:DockerBridgeIP.:PGDataPath.:ServerClass.:Token

+ ID : 0 for adding a new server...non-zero is to update a server
+ Name : the server name
+ IPAddress : the server IP address
+ DockerBridgeIP : the Docker Bridge IP to use for this server
+ PGDataPath : the root file path to where PG data files will be stored
+ ServerClass : the server class we are assiging to this server (low|medium|high)
+ Token : the generated auth token for this session

add a server


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/addserver
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /server/:ID.:Token

+ ID : the unique assigned ID of a server
+ Token : the generated auth token for this session

return a server


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/1.8dc0caed-39e7-47b4-878c-de1c8b0b595d
{
  "ID": "1",
  "Name": "espresso",
  "IPAddress": "192.168.0.106",
  "DockerBridgeIP": "172.17.42.1",
  "PGDataPath": "/var/lib/pgsql",
  "ServerClass": "low",
  "CreateDate": "02-10-2015 14:35:34"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/:ID.:Token

+ ID : the unique assigned ID of a cluster
+ Token : the generated auth token for this session

return a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/cluster/2.1efbfd50-9bb243a0-8c91-b6f4a837a4f2
{
  "ID": "2",
  "Name": "aa",
  "ClusterType": "asynchronous",
  "Status": "initialized",
  "CreateDate": "02-11-2015 09:18:13",
  "Token": ""
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/configure/:ID.:Token

+ ID : the unique assigned ID of a cluster
+ Token : the generated auth token for this session

configure a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/delete/:ID.:Token

+ ID : the unique assigned ID of a cluster
+ Token : the generated auth token for this session

delete a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /deleteserver/:ID.:Token

+ ID : the unique assigned ID of a server
+ Token : the generated auth token for this session

delete a server


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /nodes/:Token

+ Token : the generated auth token for this session

return all containers


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/nodes/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
[
  {
    "ID": "9",
    "ClusterID": "2",
    "ServerID": "1",
    "Name": "aa-master",
    "Role": "master",
    "Image": "cpm-node",
    "CreateDate": "02-11-2015 09:18:14",
    "Status": "UNKNOWN"
  },
  {
    "ID": "11",
    "ClusterID": "2",
    "ServerID": "1",
    "Name": "aa-pgpool",
    "Role": "pgpool",
    "Image": "cpm-pgpool",
    "CreateDate": "02-11-2015 09:18:20",
    "Status": "UNKNOWN"
  },
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /nodes/nocluster/:Token

+ Token : the generated auth token for this session

return all containers not in a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/nodes/nocluster/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
[
  {
    "ID": "8",
    "ClusterID": "-1",
    "ServerID": "1",
    "Name": "tc",
    "Role": "unassigned",
    "Image": "cpm-node",
    "CreateDate": "02-11-2015 09:17:43",
    "Status": "UNKNOWN"
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /clusternodes/:ClusterID.:Token

+ ClusterID : the unique ID of a cluster
+ Token : the generated auth token for this session

return all containers for a given cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/clusternodes/2.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
[
  {
    "ID": "9",
    "ClusterID": "2",
    "ServerID": "1",
    "Name": "aa-master",
    "Role": "master",
    "Image": "cpm-node",
    "CreateDate": "02-11-2015 09:18:14",
    "Status": "UNKNOWN"
  },
  {
    "ID": "11",
    "ClusterID": "2",
    "ServerID": "1",
    "Name": "aa-pgpool",
    "Role": "pgpool",
    "Image": "cpm-pgpool",
    "CreateDate": "02-11-2015 09:18:20",
    "Status": "UNKNOWN"
  },
  {
    "ID": "10",
    "ClusterID": "2",
    "ServerID": "1",
    "Name": "aa-standby-0",
    "Role": "standby",
    "Image": "cpm-node",
    "CreateDate": "02-11-2015 09:18:19",
    "Status": "UNKNOWN"
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /nodes/forserver/:ServerID.:Token

+ ServerID : the unique ID for a server
+ Token : the generated auth token for this session

return all containers for a server


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/nodes/forserver/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
[
  {
    "ID": "9",
    "ClusterID": "2",
    "ServerID": "1",
    "Name": "aa-master",
    "Role": "master",
    "Image": "cpm-node",
    "CreateDate": "02-11-2015 09:18:14",
    "Status": "down"
  },
  {
    "ID": "11",
    "ClusterID": "2",
    "ServerID": "1",
    "Name": "aa-pgpool",
    "Role": "pgpool",
    "Image": "cpm-pgpool",
    "CreateDate": "02-11-2015 09:18:20",
    "Status": "down"
  },
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /provision/:Profile.:Image.:ServerID.:ContainerName.:Standalone.:Token

+ Profile : the Docker profile to use for this node
+ Image : the Docker image name to base this node on
+ ServerID : the unique ID of the server to host this container
+ ContainerName : the user picked name for this container
+ Standalone : flag for making this node available to be part of a cluster
+ Token : the generated auth token for this session

provision a new container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /node/:ID.:Token

+ Token : the generated auth token for this session

return a container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/node/8.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
{
  "ID": "8",
  "ClusterID": "-1",
  "ServerID": "1",
  "Name": "tc",
  "Role": "unassigned",
  "Image": "cpm-node",
  "CreateDate": "02-11-2015 09:17:43",
  "Status": "OFFLINE"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /kube/:Token

+ Token : the generated auth token for this session

return boolean of kube configuration


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/kube/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
{
  "URL": "KUBE_URL is not set"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /deletenode/:ID.:Token

+ Token : the generated auth token for this session

delete a container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure -X DELETE https://cpm-admin.crunchy.lab:13000/node/17
~~~~~~~~~~~~~~~~~~~~~~~~


###GET /monitor/server-getinfo/:ServerID.:Metric.:Token

+ Token : the generated auth token for this session

return server monitoring data


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/monitor/server-getinfo/.cpmdf.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
{"df":[
{"filesystem":"/dev/mapper/centos-root","total":"32G","used":"7G","available":"26G","pctused":"20%","mountpt":"/"}
,
{"filesystem":"devtmpfs","total":"8G","used":"0G","available":"8G","pctused":"0%","mountpt":"/dev"}
,
{"filesystem":"tmpfs","total":"8G","used":"1G","available":"8G","pctused":"1%","mountpt":"/dev/shm"}
,
{"filesystem":"tmpfs","total":"8G","used":"1G","available":"8G","pctused":"1%","mountpt":"/run"}
,
{"filesystem":"tmpfs","total":"8G","used":"0G","available":"8G","pctused":"0%","mountpt":"/sys/fs/cgroup"}
,
{"filesystem":"/dev/sda1","total":"1G","used":"1G","available":"1G","pctused":"29%","mountpt":"/boot"}
,
{"filesystem":"/dev/mapper/centos-home","total":"16G","used":"2G","available":"14G","pctused":"12%","mountpt":"/home"}
]}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/loadtest/:ID.:Writes.:Token

+ Token : the generated auth token for this session

perform a load test and return the results


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin:13000/monior/container/loadtest/1.1000.9a8f9a1e-9c81-4e4f-9f52-01d2ea6cd741
 

[{"operation":"inserts","count":1000,"results":44.979},{"operation":"selects","count":1000,"results":31.952},{"operation":"updates","count":1000,"results":40.442},{"operation":"deletes","count":1000,"results":30.374}]


~~~~~~~~~~~~~~~~~~~~~~~~


###GET /admin/start-pg/:ID.:Token

+ Token : the generated auth token for this session

start a container's postgres


###Example


###GET /admin/start/:ID.:Token

+ Token : the generated auth token for this session

start a container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/admin/start/8.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/stop/:ID.:Token

+ Token : the generated auth token for this session

stop a container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/admin/stop/8.1efbfd50-9b2-43a0-8c91-b6f4a837a4f2
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/failover/:ID.:Token

+ Token : the generated auth token for this session

cause a postgres fail over


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/stop-pg/:ID.:Token

+ Token : the generated auth token for this session

stop a container's postgres


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/admin/stop-pg/8.1efbfd5-9bb2-43a0-8c91-b6f4a837a4f2
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /event/join-cluster/:IDList.:MasterID.:ClusterID.:Token

+ Token : the generated auth token for this session

add a node to a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~


###GET /sec/logout/:Token

+ Token : the generated auth token for this session

logout Logout


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/logout/ab3ba695-533-43a4-a3c0-9a4c55288c14
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~


###POST /sec/updateuser

+ Token : the generated auth token for this session

update a user's account into UpdateUser


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/cp

+ Token : the generated auth token for this session

change a user's password ChangePassword


###Example


###POST /sec/adduser

+ Token : the generated auth token for this session

add a user AddUser


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getuser/:ID.:Token

+ Token : the generated auth token for this session

return a user's information GetUser


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/getuser/cpm.1efbfd5-9bb2-43a0-8c91-b6f4a837a4f2
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getusers/:Token

+ Token : the generated auth token for this session

return all users information GetAllUsers


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/getusers/1efbfd50-9b2-43a0-8c91-b6f4a837a4f2
[
  {
    "Name": "cpm",
    "Password": "dd6ced",
    "Roles": {
      "superuser": {
        "Name": "superuser",
        "Selected": true,
        "Permissions": {
          "perm-backup": {
            "Name": "perm-backup",
            "Description": "perform backups",
            "Selected": true
          },
          "perm-cluster": {
            "Name": "perm-cluster",
            "Description": "maintain clusters",
            "Selected": true
          },
          "perm-container": {
            "Name": "perm-container",
            "Description": "maintain containers",
            "Selected": true
          },
          "perm-server": {
            "Name": "perm-server",
            "Description": "maintain servers",
            "Selected": true
          },
          "perm-setting": {
            "Name": "perm-setting",
            "Description": "maintain settings",
            "Selected": true
          },
          "perm-user": {
            "Name": "perm-user",
            "Description": "maintain users",
            "Selected": true
          }
        },
        "UpdateDate": "",
        "Token": ""
      }
    },
    "UpdateDate": "",
    "Token": ""
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/deleteuser/:ID.:Token

+ Token : the generated auth token for this session

delete a user DeleteUser


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/updaterole

+ Token : the generated auth token for this session

update a CPM role UpdateRole


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/addrole

+ Token : the generated auth token for this session

add a CPM role AddRole


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/deleterole/:ID.:Token

+ Token : the generated auth token for this session

delete a CPM role DeleteRole


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getroles/:Token

+ Token : the generated auth token for this session

get all CPM roles GetAllRoles


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/getroles/1efbfd50-9b2-43a0-8c91-b6f4a837a4f2
[
  {
    "Name": "superuser",
    "Selected": false,
    "Permissions": {
      "perm-backup": {
        "Name": "perm-backup",
        "Description": "perform backups",
        "Selected": true
      },
      "perm-cluster": {
        "Name": "perm-cluster",
        "Description": "maintain clusters",
        "Selected": true
      },
      "perm-container": {
        "Name": "perm-container",
        "Description": "maintain containers",
        "Selected": true
      },
      "perm-server": {
        "Name": "perm-server",
        "Description": "maintain servers",
        "Selected": true
      },
      "perm-setting": {
        "Name": "perm-setting",
        "Description": "maintain settings",
        "Selected": true
      },
      "perm-user": {
        "Name": "perm-user",
        "Description": "maintain users",
        "Selected": true
      }
    },
    "UpdateDate": "",
    "Token": ""
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getrole/:Name.:Token

+ Token : the generated auth token for this session

get a CPM role GetRole (CURRENTLY NOT USED)


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000sec/getrole/superuser.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
{
  "Name": "",
  "Selected": false,
  "Permissions": null,
  "UpdateDate": "",
  "Token": ""
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /backup/now

+ Token : the generated auth token for this session

perform a postgres backup BackupNow


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /backup/addschedule

+ Token : the generated auth token for this session

add a new container admin schedule AddSchedule


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/deleteschedule/:ID.:Token

+ Token : the generated auth token for this session

remove a container schedule DeleteSchedule


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /backup/updateschedule

+ Token : the generated auth token for this session

update a container schedule UpdateSchedule


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/getschedules/:ContainerName.:Token

+ Token : the generated auth token for this session

get all schedules for a container GetAllSchedules


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/getschedule/:ID.:Token

+ Token : the generated auth token for this session

get a container schedule GetSchedule


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/getstatus/:ID.:Token

+ Token : the generated auth token for this session

get a schedule job's status GetStatus


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/getallstatus/:ID.:Token

+ Token : the generated auth token for this session

get all scheduled job's status for a container GetAllStatus


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/nodes/:Token

+ Token : the generated auth token for this session

GetBackupNodes


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /mon/server/:Metric.:ServerID.:Interval.:Token

+ Token : the generated auth token for this session

GetServerMetrics

###Example
~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/mon/server/cpu.1.1w.18c0bb2a-2fa2-422e-a583-9df68948802f
[
  {
    "name": "cpu",
    "columns": [
      "time",
      "sequence_number",
      "server",
      "value"
    ],
    "points": [
      [
        1.424206259245e+12,
        150001,
        "espresso",
        0.21
      ],
      [
        1.424206559231e+12,
        340001,
        "espresso",
        1.21
      ],

     ]
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /mon/container/pg2/:Name.:Interval.:Token

+ Token : the generated auth token for this session

GetPG2 - container database sizes

###Example
~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/mon/container/pg2/pga.1w.a1459da5-8db1-48ac-b1f0-3bb3427f96e2
[
 {
    "Color": "#c05020",
    "Data": [
      {
        "name": "pg2",
        "columns": [
          "time",
          "sequence_number",
          "database",
          "value"
        ],
        "points": [

          [
            1.424221017122e+12,
            9.080001e+06,
            "cpmtest",
            9
          ]
        ]
      }
    ],
    "Name": "cpmtest"
  },
 {
    "Color": "#c05020",
    "Data": [
      {
        "name": "pg2",
        "columns": [
          "time",
          "sequence_number",
          "database",
          "value"
        ],
        "points": [

          [
            1.424221017122e+12,
            9.080001e+06,
            "postgres",
            9
          ]
        ]
      }
    ],
    "Name": "postgres"
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /mon/hc1/:Token

+ Token : the generated auth token for this session

GetHC1 - health check 1 - databases down

###Example
~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/mon/h1/24c715ca-2468-4450-8fee-6e2a9f7714dc
[
  {
    "name": "hc1",
    "columns": [
      "time",
      "sequence_number",
      "seconds",
      "service",
      "servicetype",
      "status"
    ],
    "points": [
      [
        1.424353939047e+12,
        1.2470001e+07,
        1.42435393e+09,
        "pga",
        "db",
        "down"
      ],
      [
        1.424353936039e+12,
        1.2440001e+07,
        1.42435393e+09,
        "ac-standby-0",
        "db",
        "down"
      ],
      [
        1.424353933032e+12,
        1.2430001e+07,
        1.42435393e+09,
        "ac-pgpool",
        "db",
        "down"
      ],
      [
        1.424353930021e+12,
        1.2380001e+07,
        1.42435393e+09,
        "ac-master",
        "db",
        "down"
      ]
    ]
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /version

returns the CPM version number


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/version
1.0.0
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/settings/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container pg_settings data


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/monitor/container/settings/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/repl/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container pg_replication data


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/monitor/container/repl/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
[
  {
    "Pid": "39",
    "Usesysid": "10",
    "Usename": "postgres",
    "AppName": "walreceiver",
    "ClientAddr": "172.17.0.8",
    "ClientHostname": "ac-standby-0.crunchy.lab",
    "ClientPort": "58959",
    "BackendStart": "2015-03-03 11:04-02",
    "State": "streaming",
    "SentLocation": "0/307DF40",
    "WriteLocation": "0/307DF40",
    "FlushLocation": "0/307DF40",
    "ReplayLocation": "0/307DF40",
    "SyncPriority": "0",
    "SyncState": "async"
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/database/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container pg_databases data


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/monitor/container/database/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
[
  {
    "Datname": "template1",
    "BlksRead": "0",
    "TupReturned": "0",
    "TupFetched": "0",
    "TupInserted": "0",
    "TupUpdated": "0",
    "TupDeleted": "0",
    "StatsReset": " "
  },
  {
    "Datname": "template0",
    "BlksRead": "0",
    "TupReturned": "0",
    "TupFetched": "0",
    "TupInserted": "0",
    "TupUpdated": "0",
    "TupDeleted": "0",
    "StatsReset": " "
  },
  {
    "Datname": "postgres",
    "BlksRead": "151",
    "TupReturned": "19513",
    "TupFetched": "1347",
    "TupInserted": "0",
    "TupUpdated": "0",
    "TupDeleted": "0",
    "StatsReset": "2015-03-03 15:39:15"
  },
  {
    "Datname": "cpmtest",
    "BlksRead": "76",
    "TupReturned": "18125",
    "TupFetched": "694",
    "TupInserted": "0",
    "TupUpdated": "0",
    "TupDeleted": "0",
    "StatsReset": "2015-03-03 15:41:24"
  }
]
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/bgwriter/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container bgwriter data


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/monitor/container/bgwriter/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
{
  "Now": "03/03/15 04:00:38",
  "AllocMbps": ".0015",
  "CheckpointMbps": "0.",
  "CleanMbps": "0.",
  "BackendMbps": "0.",
  "WriteMbps": "0."
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/controldata/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container controldata data


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl --insecure https://cpm-admin.crunchy.lab:13000/monitor/container/controldata/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /project/add

+ ID : can be empty
+ Name : the name to use for this project
+ Desc : the project description
+ CreateDt : can be empty
+ Token : the generated auth token for this session

adds a project


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl -X POST -d @addproject.json http://cpm-admin:13001/project/add

addproject.json:

{
"ID": "1",
"Name": "donut",
"Desc": "some desc",
"UpdateDate": "02-10-2015 14:35:34",
"Token": "572107c9-294e-4668-8993-5c089febd1da"
}

~~~~~~~~~~~~~~~~~~~~~~~~

###GET /project/getall/:Token

+ Token : the generated auth token for this session

return all projects


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:13001/project/getall/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /project/get/:ID.Token

+ ID : id of a project
+ Token : the generated auth token for this session

return a single project


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:13001/project/get/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /project/delete/:ID.Token

+ ID : id of a project
+ Token : the generated auth token for this session

delete a single project


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:13001/project/delete/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /project/update

+ ID : the generated id of a project
+ Name : the name to use for this project
+ Desc : the description of the project
+ UpdateDate : can be empty
+ Token : the generated auth token for this session

updates a project


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl -X POST -d @updateproject.json http://cpm-admin.crunchy.lab:13001/project/update

updateproject.json:
{
"ID": "1",
"Name": "donut",
"Desc": "some desc two",
"UpdateDate": "",
"Token": "572107c9-294e-4668-8993-5c089febd1da"
}
~~~~~~~~~~~~~~~~~~~~~~~~
