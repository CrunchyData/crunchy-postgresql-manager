#!/bin/bash

#
# restore job - a task that orchestrates a restore on a container node
#
# Passed in Env Vars:
#	RestoreRemotePath - the pg_backrest remote path used for the restore
#	RestoreRemoteHost - the pg_backrest host to contact for the restore
#	RestoreRemoteUser - the pg_backrest user to connect to the backup host
#	RestoreSet - the pg_backrest backup set to restore
#

source /var/cpm/bin/setenv.sh

echo ***Environment vars***
echo RestoreRemotePath=$RestoreRemotePath
echo RestoreRemoteHost=$RestoreRemoteHost
echo RestoreRemoteUser=$RestoreRemoteUser
echo RestoreDbUser=$RestoreDbUser
echo RestoreDbPass=$RestoreDbPass
echo RestoreSet=$RestoreSet

echo RestoreContainerName=$RestoreContainerName
echo RestoreScheduleID=$RestoreScheduleID
echo RestoreServerName=$RestoreServerName
echo RestoreServerIP=$RestoreServerIP
echo RestoreProfileName=$RestoreProfileName
echo ***end of Environment vars***

env

restorecommand
