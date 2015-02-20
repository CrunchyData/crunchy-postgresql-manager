Authentication
---------------

###GET /sec/login/:ID.:PSW

login and return an auth token


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /clusters/:Token


returns all clusters


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/clusters/789c31ff-b18f-47b3-bb63-1fd603895aa5
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

returns all servers


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/servers/789c31ff-b18f-47b3-bb63-1fd603895aa5
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

updates or adds a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

curl -i -d '{"ID":"1","ClusterID":"1"}' http://cpm-admin.crunchy.lab:8080/event/join-cluster

###POST /autocluster

performs an auto-cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /savesettings

saves settings


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /saveprofiles

saves profiles


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /saveclusterprofiles

saves cluster profiles


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /settings/:Token

returns all settings


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/settings/789c31ff-b18f-7b3-bb63-1fd603895aa5
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

###GET /addserver

add a server


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/addserver
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /server/:ID.:Token

return a server


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/1.8dc0caed-39e7-47b4-878c-de1c8b0b595d
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

return a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/cluster/2.1efbfd50-9bb243a0-8c91-b6f4a837a4f2
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

configure a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/delete/:ID.:Token

delete a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /deleteserver/:ID.:Token

delete a server


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /nodes/:Token

return all containers


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/nodes/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
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

return all containers not in a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/nodes/nocluster/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
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

return all containers for a given cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/clusternodes/2.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
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

return all containers for a server


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/nodes/forserver/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
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

###POST /node
-----------------

PostNode


###Example


###GET /provision/:Profile.:Image.:ServerID.:ContainerName.:Standalone.:Token

provision a new container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /node/:ID.:Token

return a container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/node/8.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
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

return boolean of kube configuration


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/kube/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
{
  "URL": "KUBE_URL is not set"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /deletenode/:ID.:Token

delete a container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl -X DELETE http://cpm-admin.crunchy.lab:8080/node/17
~~~~~~~~~~~~~~~~~~~~~~~~


###GET /monitor/server-getinfo/:ServerID.:Metric.:Token

return server monitoring data


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/monitor/server-getinfo/.cpmdf.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
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

###GET /monitor/container-getinfo/:ID.:Metric.:Token

return container monitoring data



###Example

~~~~~~~~~~~~~~~~~~~~~~~~

~~~~~~~~~~~~~~~~~~~~~~~~


###GET /monitor/container-loadtest/:ID.:Metric.:Writes.:Token

perform a load test and return the results


###Example


###GET /admin/start-pg/:ID.:Token

start a container's postgres


###Example


###GET /admin/start/:ID.:Token

start a container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/admin/start/8.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/stop/:ID.:Token

stop a container


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/admin/stop/8.1efbfd50-9b2-43a0-8c91-b6f4a837a4f2
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/failover/:ID.:Token

cause a postgres fail over


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/stop-pg/:ID.:Token

stop a container's postgres


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/admin/stop-pg/8.1efbfd5-9bb2-43a0-8c91-b6f4a837a4f2
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /event/join-cluster/:IDList.:MasterID.:ClusterID.:Token

add a node to a cluster


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~


###GET /sec/logout/:Token

logout Logout


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/logout/ab3ba695-533-43a4-a3c0-9a4c55288c14
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~


###POST /sec/updateuser

update a user's account into UpdateUser


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/cp

change a user's password ChangePassword


###Example


###POST /sec/adduser

add a user AddUser


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getuser/:ID.:Token

return a user's information GetUser


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/getuser/cpm.1efbfd5-9bb2-43a0-8c91-b6f4a837a4f2
{
  "Status": "OK"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getusers/:Token

return all users information GetAllUsers


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/getusers/1efbfd50-9b2-43a0-8c91-b6f4a837a4f2
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

delete a user DeleteUser


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/updaterole

update a CPM role UpdateRole


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/addrole

add a CPM role AddRole


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/deleterole/:ID.:Token

delete a CPM role DeleteRole


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getroles/:Token

get all CPM roles GetAllRoles


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/getroles/1efbfd50-9b2-43a0-8c91-b6f4a837a4f2
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

get a CPM role GetRole (CURRENTLY NOT USED)


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080sec/getrole/superuser.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
{
  "Name": "",
  "Selected": false,
  "Permissions": null,
  "UpdateDate": "",
  "Token": ""
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /backup/now

perform a postgres backup BackupNow


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /backup/addschedule

add a new container admin schedule AddSchedule


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/deleteschedule/:ID.:Token

remove a container schedule DeleteSchedule


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /backup/updateschedule

update a container schedule UpdateSchedule


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/getschedules/:ContainerName.:Token

get all schedules for a container GetAllSchedules


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/getschedule/:ID.:Token

get a container schedule GetSchedule


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/getstatus/:ID.:Token

get a schedule job's status GetStatus


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/getallstatus/:ID.:Token

get all scheduled job's status for a container GetAllStatus


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /backup/nodes/:Token

GetBackupNodes


###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /mon/server/:Metric.:ServerID.:Interval.:Token

GetServerMetrics

###Example
~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/mon/server/cpu.1.1w.18c0bb2a-2fa2-422e-a583-9df68948802f
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

GetPG2 - container database sizes

###Example
~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/mon/container/pg2/pga.1w.a1459da5-8db1-48ac-b1f0-3bb3427f96e2
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

GetHC1 - health check 1 - databases down

###Example
~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/mon/h1/24c715ca-2468-4450-8fee-6e2a9f7714dc
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
