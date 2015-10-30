## Security

###GET /sec/login/:ID.:PSW

login and return an auth token

+ ID : the user id to authenticate with
+ PSW : the user password to authenticate with

~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/logout/:Token

+ Token : the generated auth token for this session

logout Logout
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/sec/logout/ab3ba695-533-43a4-a3c0-9a4c55288c14
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/updateuser

+ Token : the generated auth token for this session

update a user account into UpdateUser
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @provision-proxy.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/provisionproxy
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/cp

+ Token : the generated auth token for this session

change a user password ChangePassword
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @changepassword.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/changepassword
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/adduser

+ Token : the generated auth token for this session

add a user AddUser
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @adduser.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/adduser
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getuser/:ID.:Token

+ Token : the generated auth token for this session

return a user information GetUser
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/sec/getuser/cpm.1efbfd5-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getusers/:Token

+ Token : the generated auth token for this session

return all users information GetAllUsers
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/sec/getusers/1efbfd50-9b2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/deleteuser/:ID.:Token

+ Token : the generated auth token for this session

delete a user DeleteUser
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/sec/login/cpm.cpm
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/updaterole

+ Token : the generated auth token for this session

update a CPM role UpdateRole

~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @updaterole.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/sec/updaterole
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /sec/addrole

+ Token : the generated auth token for this session

add a CPM role AddRole
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @addrole.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/sec/addrole
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/deleterole/:ID.:Token

+ Token : the generated auth token for this session

delete a CPM role DeleteRole
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/sec/deleterole/1.ajkajdkajdfkjadf
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getroles/:Token

+ Token : the generated auth token for this session

get all CPM roles GetAllRoles
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/sec/getroles/1efbfd50-9b2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /sec/getrole/:Name.:Token

+ Token : the generated auth token for this session

get a CPM role GetRole (CURRENTLY NOT USED)
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001sec/getrole/superuser.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

## Project Information

###POST /project/add

+ ID : can be empty
+ Name : the name to use for this project
+ Desc : the project description
+ CreateDt : can be empty
+ Token : the generated auth token for this session

adds a project
~~~~~~~~~~~~~~~~~~~~~~~~
curl -X POST -d @addproject.json http://cpm-admin:13001/project/add
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /project/getall/:Token

+ Token : the generated auth token for this session

return all projects
~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:13001/project/getall/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /project/get/:ID.Token

+ ID : id of a project
+ Token : the generated auth token for this session

return a single project
~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:13001/project/get/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /project/delete/:ID.Token

+ ID : id of a project
+ Token : the generated auth token for this session

delete a single project
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
~~~~~~~~~~~~~~~~~~~~~~~~
curl -X POST -d @updateproject.json http://cpm-admin.crunchy.lab:13001/project/update
~~~~~~~~~~~~~~~~~~~~~~~~
###GET /projectnodes/:ID.:Token

+ ID : the unique assigned ID of a project
+ Token : the generated auth token for this session

return a list of containers in a project
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/projectnodes/1.8dc0caed-39e7-47b4-878c-de1c8b0b595d
~~~~~~~~~~~~~~~~~~~~~~~~

## Container Information

###GET /admin/stop-pg/:ID.:Token

+ ID : the container ID
+ Token : the generated auth token for this session

stop a container postgres
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/admin/stop-pg/8.1efbfd5-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/stop/:ID.:Token 
+ Token : the generated auth token for this session

stop a container
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/admin/stop/8.1efbfd50-9b2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/start/:ID.:Token

+ ID : the container ID
+ Token : the generated auth token for this session

start a container
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/admin/start/8.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/start-pg/:ID.:Token

+ ID : the container ID
+ Token : the generated auth token for this session

start a containers postgres database
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/admin/start-pg/1.8dc0caed-39e7-47b4-878c-de1c8b0b595d
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /node/:ID.:Token

+ Token : the generated auth token for this session

return a container
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/node/8.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /deletenode/:ID.:Token

+ Token : the generated auth token for this session

delete a container
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/deletenode/17.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /provision

+ Profile : the Docker profile to use for this node
+ Image : the Docker image name to base this node on
+ ServerID : the unique ID of the server to host this container
+ ContainerName : the user picked name for this container
+ Standalone : flag for making this node available to be part of a cluster
+ Token : the generated auth token for this session

provision a new container
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @provision.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/provision
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /nodes/nocluster/:Token

+ Token : the generated auth token for this session

return all containers not in a cluster
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/nodes/nocluster/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /nodes/:Token

+ Token : the generated auth token for this session

return all containers
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/nodes/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

## Access Rule Information

###GET /rules/get/:ID.:Token

+ ID : the access rule ID
+ Token : the generated auth token for this session

 get an access rule
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/rules/get/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /rules/getall/:Token

+ Token : the generated auth token for this session

 get all access rules
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/rules/getall/1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /rules/delete/:ID.:Token

+ ID : the access rule ID
+ Token : the generated auth token for this session

 delete an access rule
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/rules/delete/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

## Server Information

###GET /admin/startall/:ID.:Token

+ ID : the unique ID for a server
+ Token : the generated auth token for this session

perform a docker start on all containers on a given server
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/admin/startall/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/stopall/:ID.:Token

+ ID : the unique ID for a server
+ Token : the generated auth token for this session

perform a docker stop on all containers on a given server
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/admin/stopall/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /nodes/forserver/:ServerID.:Token

+ ServerID : the unique ID for a server
+ Token : the generated auth token for this session

return all containers for a server
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/nodes/forserver/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /server/:ID.:Token

+ ID : the unique assigned ID of a server
+ Token : the generated auth token for this session

return a server
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/1.8dc0caed-39e7-47b4-878c-de1c8b0b595d
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /deleteserver/:ID.:Token

+ ID : the unique assigned ID of a server
+ Token : the generated auth token for this session

delete a server
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/sec/login/cpm.cpm
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /servers/:Token

+ Token : the generated auth token for this session

returns all servers
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/servers/789c31ff-b18f-47b3-bb63-1fd603895aa5
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /servers/:Token

+ Token : the security token used for auth

Get all the servers defined in CPM
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/servers/789c31ff-b18f-47b3-bb63-1fd603895aa5
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
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/addserver/1.foo.192-168-0-104.171-10-10-17.
~~~~~~~~~~~~~~~~~~~~~~~~

## Database User Information

###POST /dbuser/add

add a database user to a given container
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @dbuseradd.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/dbuser/add
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /dbuser/update

update a database user to a given container
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @dbuserupdate.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/dbuser/update
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /dbuser/delete/:ContainerID.:Rolname.:Token

+ ContainerID : the container ID
+ Rolname : the role name we are deleting
+ Token : the generated auth token for this session

delete a database user for a given container
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/dbuser/delete/1.foo.kjakdjfkajdkfj
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /dbuser/get/:ContainerID.:Rolname.:Token

+ ContainerID : the container ID
+ Rolname : the role name we are fetching
+ Token : the generated auth token for this session

get a database user for a given container
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/dbuser/get/1.foo.kjakdjfkajdkfj
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /dbuser/getall/:ID.:Token

+ ContainerID : the container ID
+ Token : the generated auth token for this session

get all database users for a given container
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/dbuser/getall/1.kjakdjfkajdkfj
~~~~~~~~~~~~~~~~~~~~~~~~

## Cluster Information

###GET /event/join-cluster/:IDList.:MasterID.:ClusterID.:Token

+ Token : the generated auth token for this session

add a node to a cluster
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/event/join-cluster/1.1.1.789c31ff-b18f-47b3-bb63-1fd603895aa5
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /admin/failover/:ID.:Token

+ ID : the container ID
+ Token : the generated auth token for this session

cause a postgres fail over on a given container
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/admin/failover/1.789c31ff-b18f-47b3-bb63-1fd603895aa5
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /clusternodes/:ClusterID.:Token

+ ClusterID : the unique ID of a cluster
+ Token : the generated auth token for this session

return all containers for a given cluster
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/clusternodes/2.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/stop/:ID.:Token

+ ID : the unique assigned ID of a cluster
+ Token : the generated auth token for this session

perform a docker stop on a given clusters set of containers
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/cluster/stop/2.1efbfd50-9bb243a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/start/:ID.:Token

+ ID : the unique assigned ID of a cluster
+ Token : the generated auth token for this session

perform a docker start on a given clusters set of containers
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/cluster/start/2.1efbfd50-9bb243a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/:ID.:Token

+ ID : the unique assigned ID of a cluster
+ Token : the generated auth token for this session

return a cluster
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/cluster/2.1efbfd50-9bb243a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/configure/:ID.:Token

+ ID : the unique assigned ID of a cluster
+ Token : the generated auth token for this session

configure a cluster
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/cluster/configure/2.1efbfd50-9bb243a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/delete/:ID.:Token

+ ID : the unique assigned ID of a cluster
+ Token : the generated auth token for this session

delete a cluster and its containers
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/cluster/delete/1.1efbfd50-9bb243a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /projectclusters/:ID.:Token

+ ID : the user id to authenticate with
+ Token : the security token used for auth

Get all the clusters for a given project
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/projectclusters/1.789c31ff-b18f-47b3-bb63-1fd603895aa5
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /cluster

updates or adds a cluster
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @postcluster.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/cluster
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /autocluster

+ Name : the name to use for this cluster
+ ClusterType : the type of cluster (synchronous|asynchronous)
+ ClusterProfile : the cluster profile to use for cluster creation (SM|LG|MED)
+ Token : the generated auth token for this session

performs an auto-cluster
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @autocluster.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/autocluster
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /clusters/:Token

+ Token : the generated auth token for this session

returns all clusters
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/clusters/789c31ff-b18f-47b3-bb63-1fd603895aa5
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /cluster/scale/:ID.:Token

+ ID : unique id of a given cluster
+ Token : the generated auth token for this session

add a standby node to a given cluster
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/cluster/scale/1.789c31ff-b18f-47b3-bb63-1fd603895aa5
~~~~~~~~~~~~~~~~~~~~~~~~

## Task Information

###POST /task/executenow

+ Token : the generated auth token for this session

execute a task schedule immediately
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @executenow.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/task/executenow
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /task/addschedule

+ Token : the generated auth token for this session

add a new container admin schedule AddSchedule
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @addschedule.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/task/addschedule
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /task/deleteschedule/:ID.:Token

+ Token : the generated auth token for this session

remove a container schedule DeleteSchedule
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/task/deleteschedule/1.kjkjadfjkajdfkjadksf
~~~~~~~~~~~~~~~~~~~~~~~~

###POST /task/updateschedule

+ Token : the generated auth token for this session

update a container schedule UpdateSchedule
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @updateschedule.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/task/updateschedule
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /task/getschedules/:ContainerName.:Token

+ Token : the generated auth token for this session

get all schedules for a container GetAllSchedules
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/task/getschedules/foo.kjadkfjkjakdjfkadjf
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /task/getschedule/:ID.:Token

+ Token : the generated auth token for this session

get a container schedule GetSchedule
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/task/getschedule/1.fkjkjadkfjkjadsfjkdaf
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /task/getstatus/:ID.:Token

+ Token : the generated auth token for this session

get a schedule job status GetStatus
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/task/getstatus/1.kjakdjfkajkdjfkjadfasdf
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /task/getallstatus/:ID.:Token

+ Token : the generated auth token for this session

get all scheduled job status for a container GetAllStatus
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/task/getallstatus/1.kjakjadfjkjaksdjfkajdf
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /task/nodes/:Token

+ Token : the generated auth token for this session

TODO
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/task/nodes/kjakjfjkadjfkjkajdf
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


~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @saveprofiles.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/saveprofiles
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


~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @saveclusterprofiles.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/saveclusterprofiles
~~~~~~~~~~~~~~~~~~~~~~~~

## Settings

###POST /savesetting

update a setting value
~~~~~~~~~~~~~~~~~~~~~~~~
curl --data @savesetting.json -H "Content-Type: application/json" -X POST \
http://cpm-admin:13001/savesetting
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /settings/:Token

+ Token : the generated auth token for this session

returns all settings
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/settings/789c31ff-b18f-7b3-bb63-1fd603895aa5
~~~~~~~~~~~~~~~~~~~~~~~~

## Monitoring

###GET /mon/healthcheck/:Token

+ Token : the generated auth token for this session

GetHC1 - health check 1 - databases down
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/mon/healthcheck/24c715ca-2468-4450-8fee-6e2a9f7714dc
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/settings/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container pg_settings data
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/monitor/container/settings/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/repl/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container pg_replication data
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/monitor/container/repl/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/database/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container pg_databases data
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/monitor/container/database/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/bgwriter/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container bgwriter data
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/monitor/container/bgwriter/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/controldata/:ID.Token

+ ID : the container ID
+ Token : the generated auth token for this session

return container controldata data
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/monitor/container/controldata/1.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/server-getinfo/:ServerID.:Metric.:Token

+ Token : the generated auth token for this session

return server monitoring data
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/monitor/server-getinfo/1.cpmdf.1efbfd50-9bb2-43a0-8c91-b6f4a837a4f2
~~~~~~~~~~~~~~~~~~~~~~~~

###GET /monitor/container/loadtest/:ID.:Writes.:Token

+ Token : the generated auth token for this session

perform a load test and return the results
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin:13001/monior/container/loadtest/1.1000.9a8f9a1e-9c81-4e4f-9f52-01d2ea6cd741
~~~~~~~~~~~~~~~~~~~~~~~~


###GET /mon/server/:Metric.:ServerID.:Interval.:Token

+ Token : the generated auth token for this session

GetServerMetrics


~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/mon/server/cpu.1.1w.18c0bb2a-2fa2-422e-a583-9df68948802f
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


~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/mon/container/pg2/pga.1w.a1459da5-8db1-48ac-b1f0-3bb3427f96e2
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


###GET /version

returns the CPM version number
~~~~~~~~~~~~~~~~~~~~~~~~
curl  http://cpm-admin.crunchy.lab:13001/version
~~~~~~~~~~~~~~~~~~~~~~~~
