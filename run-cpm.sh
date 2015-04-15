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

INSTALLDIR=`pwd`

echo "setting up log dir..."
LOGDIR=/var/cpm/logs
mkdir -p $LOGDIR
chmod -R 777 $LOGDIR
chcon -Rt svirt_sandbox_file_t $LOGDIR

echo "setting up keys dir..."
KEYSDIR=/var/cpm/keys
chcon -Rt svirt_sandbox_file_t $KEYSDIR

echo "deleting all old log files...."
rm -rf $LOGDIR/*

echo "restarting cpm container..."
docker stop cpm
docker rm cpm
chcon -Rt svirt_sandbox_file_t $INSTALLDIR/images/cpm/www/v2
docker run --name=cpm -d \
	-v $LOGDIR:/cpmlogs \
	-v $KEYSDIR:/cpmkeys \
	-v $INSTALLDIR/images/cpm/www/v2:/www crunchydata/cpm:latest

echo "restarting cpm-admin container..."
sleep 2
docker stop cpm-admin
docker rm cpm-admin
DBDIR=/var/cpm/data/pgsql/cpm-admin
mkdir -p $DBDIR
chown postgres:postgres $DBDIR
chcon -Rt svirt_sandbox_file_t $DBDIR
docker run -e DB_HOST=127.0.0.1 \
	-e DOMAIN=crunchy.lab \
	-e CPMBASE=/var/cpm \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-admin -d  \
	-v $KEYSDIR:/cpmkeys \
	-v $LOGDIR:/cpmlogs \
	-v $DBDIR:/pgdata crunchydata/cpm-admin:latest

echo "restarting cpm-backup container..."
sleep 2
docker stop cpm-backup
docker rm cpm-backup
docker run -e DB_HOST=cpm-admin.crunchy.lab \
	-v $LOGDIR:/cpmlogs \
	-e CPMBASE=/var/cpm \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-backup -d crunchydata/cpm-backup:latest

echo "restarting cpm-mon container..."
sleep 2
docker stop cpm-mon
docker rm cpm-mon
INFLUXDIR=/var/cpm/data/influxdb
mkdir -p $INFLUXDIR
chcon -Rt svirt_sandbox_file_t $INFLUXDIR
docker run -e DB_HOST=cpm-admin.crunchy.lab \
	--hostname="cpm-mon" \
	-e CPMBASE=/var/cpm \
	-e DB_PORT=5432 -e DB_USER=postgres \
	-v $LOGDIR:/cpmlogs \
	-v $INFLUXDIR:/monitordata \
	-d --name=cpm-mon crunchydata/cpm-mon:latest

sleep 2
echo "testing containers for DNS resolution...."
ping -c 2 cpm.crunchy.lab
ping -c 2 cpm-admin.crunchy.lab
ping -c 2 cpm-backup.crunchy.lab
ping -c 2 cpm-mon.crunchy.lab

exit

docker rm cpm-dashboard
docker run --name=cpm-dashboard -d crunchydata/cpm-dashboard:latest

