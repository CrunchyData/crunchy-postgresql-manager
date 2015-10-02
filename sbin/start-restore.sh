#!/bin/bash

#
# restore a database from a pg_backrest backup
#
# Passed in Env Vars:
#	NODE_NAME - the name we give the provisioned db container
#	REPO_REMOTE_PATH - the pg_backrest remote path used for the restore
#	BACKUP_HOST - the pg_backrest host to contact for the restore
#	BACKUP_USER - the pg_backrest user to connect to the backup host
#	BACKUP_SET - the pg_backrest backup set to restore
#	BACKREST_KEY_PASS - the pg_backrest ssh key passphrase
#	BACKUP_SERVERNAME - the CPM server name the restore is executing on
#	BACKUP_SERVERIP - the CPM server ip address the restore is executing on
#	NODE_PROFILENAME - should be a static value of "backrest-restore"
#
# Mounted paths:
# 	/pgdata - the db path to restore files to
# 	/keys - the pg_backrest ssh keys to use for connecting to the backup host
#

source /var/cpm/bin/setenv.sh

echo ***Environment vars***
echo NODE_NAME=$NODE_NAME
echo REPO_REMOTE_PATH=$REPO_REMOTE_PATH
echo BACKUP_HOST=$BACKUP_HOST
echo BACKUP_USER=$BACKUP_USER
echo BACKUP_SET=$BACKUP_SET
echo BACKREST_KEY_PASS=$BACKREST_KEY_PASS
echo BACKUP_SERVERNAME=$BACKUP_SERVERNAME
echo BACKUP_SERVERIP=$BACKUP_SERVERIP
echo NODE_PROFILENAME=$NODE_PROFILENAME
echo BACKUP_SCHEDULEID=$BACKUP_SCHEDULEID
echo PGDATA=$PGDATA
echo ***end of Environment vars***


#
# function to create the backrest config file
#
createBackrestConfig() {
sed "s+REPO_REMOTE_PATH+$REPO_REMOTE_PATH+g; \
	s/BACKUP_HOST/$BACKUP_HOST/g; \
	s+PGDATA+$PGDATA+g; \
	s/BACKUP_USER/$BACKUP_USER/g" \
	/var/cpm/conf/pg_backrest.conf.template  > /etc/pg_backrest.conf
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
#ssh -i /keys/id_rsa -o StrictHostKeyChecking=no $BACKUP_USER@$BACKUP_HOST
#ssh -i /keys/id_rsa -o StrictHostKeyChecking=no $BACKUP_USER@$BACKUP_HOST 'hostname'
setupssh.sh $BACKUP_HOST $BACKUP_USER $BACKREST_KEY_PASS

# run restore daemon as postgres user
# daemon needs to run 'pg_backrest restore --stanza=main --set=$BACKUP_SET'
# and send back stats to cpm-admin database
# and provision a cpm-node

backrestrestore

dummyserver

