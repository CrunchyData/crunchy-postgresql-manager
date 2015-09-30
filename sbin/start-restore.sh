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
#
# Mounted paths:
# 	/pgdata - the db path to restore files to
# 	/keys - the pg_backrest ssh keys to use for connecting to the backup host
#

source /var/cpm/bin/setenv.sh

#
# function to create the backrest config file
#
createBackrestConfig() {
sed "s+REPO_REMOTE_PATH+$REPO_REMOTE_PATH+g; \
	s/BACKUP_HOST/$BACKUP_HOST/g; \
	s+PGDATA+$PGDATA+g; \
	s/BACKUP_USER/$BACKUP_USER/g" \
	pg_backrest.conf.template  > pg_backrest.conf
}

# PASS IN name of container we will provision
NODE_NAME=restored

# PASS IN pg_backrest env vars
REPO_REMOTE_PATH=/tmp/backrest-repo
BACKUP_HOST=192.168.0.106
BACKUP_USER=backrest
BACKUP_SET=latest

# ssh public key to use for backrest..../tmp/keys is mounted
# copy the private and public keys to the user's directory and set perms
BACKREST_KEYS=/keys

# password for backrest user..you have to use this to get
# no-password ssh login setup
BACKREST_KEY_PASS=backrest

# location to restore database to.../pgdata is mounted
DB_PATH=$PGDATA

# write /etc/pg_backrest.conf
createBackrestConfig

# test ssh
#ssh -i /keys/id_rsa -o StrictHostKeyChecking=no $BACKUP_USER@$BACKUP_HOST
ssh -i /keys/id_rsa -o StrictHostKeyChecking=no $BACKUP_USER@$BACKUP_HOST 'hostname'
#setupssh.sh $BACKUP_HOST $BACKUP_USER $BACKREST_KEY_PASS

# run restore daemon as postgres user
# daemon needs to run 'pg_backrest restore --stanza=main --set=$BACKUP_SET'
# and send back stats to cpm-admin database
# and provision a cpm-node

dummyserver

