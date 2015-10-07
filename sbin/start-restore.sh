#!/bin/bash

#
# restore a database from a pg_backrest backup
#
# Passed in Env Vars:
#	RestoreRemotePath - the pg_backrest remote path used for the restore
#	RestoreRemoteHost - the pg_backrest host to contact for the restore
#	RestoreRemoteUser - the pg_backrest user to connect to the backup host
#	RestoreSet - the pg_backrest backup set to restore
#
# Mounted paths:
# 	/pgdata - the db path to restore files to
# 	/keys - the pg_backrest ssh keys to use for connecting to the backup host
#

source /var/cpm/bin/setenv.sh

echo ***Environment vars***
echo RestoreRemotePath=$RestoreRemotePath
echo RestoreRemoteHost=$RestoreRemoteHost
echo RestoreRemoteUser=$RestoreRemoteUser
echo RestoreDbUser=$RestoreDbUser
echo RestoreDbPass=$RestoreDbPass
echo RestoreSet=$RestoreSet
echo PGDATA=$PGDATA
echo ***end of Environment vars***


#
# function to create the backrest config file
#
createBackrestConfig() {
sed "s+RestoreRemotePath+$RestoreRemotePath+g; \
	s/RestoreRemoteHost/$RestoreRemoteHost/g; \
	s+PGDATA+$PGDATA+g; \
	s/RestoreRemoteUser/$RestoreRemoteUser/g" \
	/var/cpm/conf/pg_backrest.conf.template  > /tmp/pg_backrest.conf
}

# ssh public key to use for backrest..../tmp/keys is mounted
# copy the private and public keys to the user's directory and set perms

mkdir /root/.ssh
cp /keys/id* /root/.ssh/
chown -R root:root /root/.ssh
chmod 700 /root/.ssh
chmod 400 /root/.ssh/id*

# create the local backrest repo which is specified in the 
# local backrest config file
mkdir /tmp/backrest-repo

# location to restore database to.../pgdata is mounted
DB_PATH=$PGDATA

# write /etc/pg_backrest.conf
createBackrestConfig

# test ssh
#ssh -i /keys/id_rsa -o StrictHostKeyChecking=no $RestoreRemoteUser@$RestoreRemoteHost
#ssh -i /keys/id_rsa -o StrictHostKeyChecking=no $RestoreRemoteUser@$RestoreRemoteHost 'hostname'
setupssh.sh $RestoreRemoteHost $RestoreRemoteUser

# run restore daemon as postgres user
# daemon needs to run 'pg_backrest restore --stanza=main --set=$RestoreSet'
# and send back stats to cpm-admin database
# and provision a cpm-node

pg_backrest --config=/tmp/pg_backrest.conf --stanza=main --target=$RestoreSet restore 

exit 0
