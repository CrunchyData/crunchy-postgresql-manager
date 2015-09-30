#!/bin/bash 


# Copyright 2015 Crunchy Data Solutions, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
	   echo "This script must be run as root" 1>&2
	      exit 1
fi

INSTALLDIR=/home/jeffmc/devproject/src/github.com/crunchydata/crunchy-postgresql-manager

echo "setting up log dir..."
LOGDIR=/var/cpm/logs

echo "setting up keys dir..."
KEYSDIR=/tmp/keys
chcon -Rt svirt_sandbox_file_t $KEYSDIR

echo "deleting all old log files...."

docker stop cpm-restore-job
docker rm cpm-restore-job
DBDIR=/var/cpm/data/pgsql/restorednode
mkdir -p $DBDIR
chown postgres:postgres $DBDIR
chcon -Rt svirt_sandbox_file_t $DBDIR
docker run -d \
	--hostname="cpm-restore-job" \
	--name=cpm-restore-job   \
	-e NODE_NAME=restorednode \
	-e REPO_REMOTE_PATH=/tmp/backrest-repo \
	-e BACKUP_HOST=192.168.0.106 \
	-e BACKUP_USER=backrest \
	-e BACKUP_SET=latest \
	-e BACKREST_KEY_PASS=backrest \
	-v $KEYSDIR:/keys \
	-v $LOGDIR:/cpmlogs \
	-v $DBDIR:/pgdata \
	crunchydata/cpm-backrest-restore:latest

