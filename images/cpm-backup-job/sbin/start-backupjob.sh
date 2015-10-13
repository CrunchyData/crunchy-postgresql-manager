#!/bin/bash -x

#
# start the backup job
#
# the service looks for the following env vars to be set by
# the cpm-admin that provisioned us
#
# $BACKUP_HOST host we are connecting to
# $BACKUP_PORT pg port we are connecting to
# $BACKUP_USER pg user we are connecting with
# $BACKUP_CPMAGENT_URL cpmagent URL we send status to
#

env > /tmp/envvars.out

source /var/cpm/bin/setenv.sh

backupcommand

#
# next line, is used only for development, block with the dummy server

#dummyserver > /tmp/dummy.log 
