Authentication
---------------

/sec/login/:ID.:PSW
-----------------
login and return an auth token

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/clusters/:Token
-----------------

returns all clusters

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/servers/:Token
-----------------
returns all servers

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/cluster
-----------------
updates or adds a cluster

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

curl -i -d '{"ID":"1","ClusterID":"1"}' http://cpm-admin.crunchy.lab:8080/event/join-cluster

/autocluster
-----------------

performs an auto-cluster

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/savesettings
-----------------

saves settings

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/saveprofiles
-----------------

saves profiles

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/saveclusterprofiles
-----------------

saves cluster profiles

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/settings/:Token
-----------------

returns all settings

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/addserver
-----------------

add a server

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/server/:ID.:Token
-----------------

return a server

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/cluster/:ID.:Token
-----------------

return a cluster

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/cluster/configure/:ID.:Token
-----------------

configure a cluster

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/cluster/delete/:ID.:Token
-----------------

delete a cluster

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/deleteserver/:ID.:Token
-----------------

delete a server

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/nodes/:Token
-----------------

return all containers

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/nodes/nocluster/:Token
-----------------

return all containers not in a cluster

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/clusternodes/:ClusterID.:Token
-----------------

return all containers for a given cluster

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/nodes/forserver/:ServerID.:Token
-----------------

return all containers for a server

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/node
-----------------

PostNode

POST

###Example


/provision/:Profile.:Image.:ServerID.:ContainerName.:Standalone.:Token
-----------------

provision a new container

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/node/:ID.:Token
-----------------

return a container

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/kube/:Token
-----------------

return boolean of kube configuration

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/deletenode/:ID.:Token
-----------------

delete a container

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl -X DELETE http://cpm-admin.crunchy.lab:8080/node/17
~~~~~~~~~~~~~~~~~~~~~~~~


/monitor/server-getinfo/:ServerID.:Metric.:Token
-----------------

return server monitoring data

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/monitor/container-getinfo/:ID.:Metric.:Token
-----------------

return container monitoring data


###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~


/monitor/container-loadtest/:ID.:Metric.:Writes.:Token
-----------------

perform a load test and return the results

###GET

###Example


/admin/start-pg/:ID.:Token
-----------------

start a container's postgres

###GET

###Example


/admin/start/:ID.:Token
-----------------

start a container

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/admin/stop/:ID.:Token
-----------------

stop a container

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/admin/failover/:ID.:Token
-----------------

cause a postgres fail over

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/admin/stop-pg/:ID.:Token
-----------------

stop a container's postgres

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/event/join-cluster/:IDList.:MasterID.:ClusterID.:Token
-----------------

add a node to a cluster

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~


/sec/logout/:Token
-----------------

logout Logout

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~


/sec/updateuser
-----------------

update a user's account into UpdateUser

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/sec/cp
-----------------

change a user's password ChangePassword

POST

###Example


/sec/adduser
-----------------

add a user AddUser

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/sec/getuser/:ID.:Token
-----------------

return a user's information GetUser

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/sec/getusers/:Token
-----------------

return all users information GetAllUsers

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/sec/deleteuser/:ID.:Token
-----------------

delete a user DeleteUser

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/sec/updaterole
-----------------

update a CPM role UpdateRole

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/sec/addrole
-----------------

add a CPM role AddRole

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/sec/deleterole/:ID.:Token
-----------------

delete a CPM role DeleteRole

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/sec/getroles/:Token
-----------------

get all CPM roles GetAllRoles

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/sec/getrole/:Token
-----------------

get a CPM role GetRole

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/backup/now
-----------------

perform a postgres backup BackupNow

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/backup/addschedule
-----------------

add a new container admin schedule AddSchedule

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/backup/deleteschedule/:ID.:Token
-----------------

remove a container schedule DeleteSchedule

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/backup/updateschedule
-----------------

update a container schedule UpdateSchedule

POST

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/backup/getschedules/:ContainerName.:Token
-----------------

get all schedules for a container GetAllSchedules

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/backup/getschedule/:ID.:Token
-----------------

get a container schedule GetSchedule

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/backup/getstatus/:ID.:Token
-----------------

get a schedule job's status GetStatus

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/backup/getallstatus/:ID.:Token
-----------------

get all scheduled job's status for a container GetAllStatus

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~

/backup/nodes/:Token
-----------------

GetBackupNodes

###GET

###Example

~~~~~~~~~~~~~~~~~~~~~~~~
curl http://cpm-admin.crunchy.lab:8080/sec/login/cpm.cpm
{
  "Contents": "789c31ff-b18f-47b3-bb63-1fd603895aa5"
}
~~~~~~~~~~~~~~~~~~~~~~~~
