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
chcon -Rt svirt_sandbox_file_t $INSTALLDIR/images/cpm/www/v3
docker run --name=cpm -d \
	-p 192.168.56.103:13001:13001 \
	-v $LOGDIR:/cpmlogs \
	-v $KEYSDIR:/cpmkeys \
	-v $INSTALLDIR/images/cpm/www/v3:/www crunchydata/cpm:latest

echo "restarting cpm-admin container..."
sleep 2
docker stop cpm-admin
docker rm cpm-admin
DBDIR=/var/cpm/data/pgsql/cpm-admin
mkdir -p $DBDIR
chown postgres:postgres $DBDIR
chcon -Rt svirt_sandbox_file_t $DBDIR
docker run -e DB_HOST=127.0.0.1 \
	-p 192.168.56.103:14001:13001 \
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

echo "restarting cpm-collect container..."
sleep 2
docker stop cpm-collect
docker rm cpm-collect
docker run -e DB_HOST=cpm-admin.crunchy.lab \
	--hostname="cpm-collect" \
	-e CONT_POLL_INT=4 \
	-e SERVER_POLL_INT=4 \
	-e HC_POLL_INT=4 \
	-e CPMBASE=/var/cpm \
	-e DB_PORT=5432 -e DB_USER=postgres \
	-v $LOGDIR:/cpmlogs \
	-d --name=cpm-collect crunchydata/cpm-collect:latest 

sleep 2
###############
echo "restarting cpm-promdash container..."
sleep 2
export DATADIR=/var/cpm/data/promdash
mkdir -p  $DATADIR
chmod 777 $DATADIR
chcon -Rt svirt_sandbox_file_t $DATADIR

docker stop cpm-promdash
docker rm cpm-promdash
docker run  \
	-v $DATADIR:/tmp/prom \
	-p 192.168.56.103:15000:3000 \
	-e DATABASE_URL=sqlite3:/tmp/prom/file.sqlite3 \
	--name=cpm-promdash -d prom/promdash
###############
echo "restarting cpm-prometheus container..."
sleep 2
export PROMCONFIG=/var/cpm/config/prometheus.yml
chmod 777 $PROMCONFIG
chcon -Rt svirt_sandbox_file_t $PROMCONFIG

docker stop cpm-prometheus
docker rm cpm-prometheus
docker run  \
	-v $PROMCONFIG:/etc/prometheus/prometheus.yml \
	-p 192.168.56.103:16000:9090 \
	--name=cpm-prometheus -d prom/prometheus:latest
##############

echo "testing containers for DNS resolution...."

ping -c 2 cpm.crunchy.lab
ping -c 2 cpm-admin.crunchy.lab
ping -c 2 cpm-backup.crunchy.lab
ping -c 2 cpm-promdash.crunchy.lab
ping -c 2 cpm-prometheus.crunchy.lab

