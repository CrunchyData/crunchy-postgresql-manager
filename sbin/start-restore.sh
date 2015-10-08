#!/bin/bash

#
# restore a database from a pg_backrest backup
#
# Passed in Env Vars:
#	$1 RestoreRemotePath - the pg_backrest remote path used for the restore
#	$2 RestoreRemoteHost - the pg_backrest host to contact for the restore
#	$3 RestoreRemoteUser - the pg_backrest user to connect to the backup host
#	$4 RestoreDbUser - a known database user in the database we are restoring
#	$5 RestoreDbPass - a known database user password
#	$6 RestoreSet - the pg_backrest backup set to restore
#
# Mounted paths:
# 	/pgdata - the db path to restore files to
# 	/keys - the pg_backrest ssh keys to use for connecting to the backup host
#

source /var/cpm/bin/setenv.sh

echo ***Environment vars***
echo RestoreRemotePath=$1
export RestoreRemotePath=$1
echo RestoreRemoteHost=$2
export RestoreRemoteHost=$2
echo RestoreRemoteUser=$3
export RestoreRemoteUser=$3
echo RestoreDbUser=$4
export RestoreDbUser=$4
echo RestoreDbPass=$5
export RestoreDbPass=$5
echo RestoreSet=$6
export RestoreSet=$6
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
	/var/cpm/conf/pg_backrest.conf.template  > /etc/pg_backrest.conf
}

# ssh public key to use for backrest..../tmp/keys is mounted
# copy the private and public keys to the user's directory and set perms

mkdir /root/.ssh
cp /keys/id* /root/.ssh/
chown -R root:root /root/.ssh
chmod 700 /root/.ssh
chmod 400 /root/.ssh/id*

# do the same key setup for the postgres user
# this is required since the database will be started
# by the postgres user and the startup will use
# pg_backrest to pull down the archives, the pull
# is accomplished via ssh
mkdir /var/lib/pgsql/.ssh
cp /keys/id* /var/lib/pgsql/.ssh/
chown -R postgres:postgres /var/lib/pgsql/.ssh
chmod 700 /var/lib/pgsql/.ssh
chmod 400 /var/lib/pgsql/.ssh/id*

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

# do the same for the postgres user, this is required for
# the database startup to work using pg_backrest
su postgres -c '/var/cpm/bin/setupssh.sh $RestoreRemoteHost $RestoreRemoteUser'

pg_backrest restore --type=name --stanza=main --target=$RestoreSet 

chown -R postgres:postgres /tmp/backrest-repo
chown -R postgres:postgres $PGDATA

exit 0
