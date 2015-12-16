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
LOCAL_IP=192.168.0.107
SWARM_MANAGER_URL=tcp://$LOCAL_IP:8000 
FLUENT_URL=$LOCAL_IP:24224

echo "setting up log dir..."
LOGDIR=/var/cpm/logs
mkdir $LOGDIR
chcon -Rt svirt_sandbox_file_t $LOGDIR

echo "setting up keys dir..."
KEYSDIR=/var/cpm/keys
mkdir $KEYSDIR
cp $INSTALLDIR/sbin/key.pem $KEYSDIR
cp $INSTALLDIR/sbin/cert.pem $KEYSDIR
chcon -Rt svirt_sandbox_file_t $KEYSDIR

echo "restarting cpm container..."
docker stop cpm-web
docker rm cpm-web
chcon -Rt svirt_sandbox_file_t $INSTALLDIR/images/cpm/www/v3
docker run --name=cpm-web -d \
	-p $LOCAL_IP:13001:13001 \
	-v $LOGDIR:/cpmlogs \
	-v $KEYSDIR:/cpmkeys \
	-v $INSTALLDIR/images/cpm/www/v3:/www \
	crunchydata/cpm:latest

echo "restarting cpm-admin container..."
sleep 2
docker stop cpm-admin
docker rm cpm-admin
DBDIR=/var/cpm/data/pgsql/cpm-admin
mkdir -p $DBDIR
chown postgres:postgres $DBDIR
chcon -Rt svirt_sandbox_file_t $DBDIR
docker run -e DB_HOST=cpm-admin \
	--hostname="cpm-admin" \
	--log-driver=fluentd \
	--log-opt fluentd-address=$FLUENT_URL \
	--log-opt fluentd-tag=docker.cpm-admin \
	-p $LOCAL_IP:14001:13001 \
	-e DOMAIN=crunchy.lab \
	-e SWARM_MANAGER_URL=$SWARM_MANAGER_URL \
	-e CPMBASE=/var/cpm \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-admin -d  \
	-v $KEYSDIR:/cpmkeys \
	-v $DBDIR:/pgdata \
	crunchydata/cpm-admin:latest

echo "restarting cpm-task container..."
sleep 2
docker stop cpm-task
docker rm cpm-task
docker run -e DB_HOST=cpm-admin.crunchy.lab \
	--log-driver=fluentd \
	--log-opt fluentd-address=$FLUENT_URL \
	--log-opt fluentd-tag=docker.cpm-task \
	-e CPMBASE=/var/cpm \
	-e SWARM_MANAGER_URL=$SWARM_MANAGER_URL \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-task \
	-d crunchydata/cpm-task:latest

sleep 2
###############
echo "restarting cpm-promdash container..."
sleep 2
export DATADIR=/var/cpm/data/promdash
mkdir -p  $DATADIR
chmod 777 $DATADIR
cp $INSTALLDIR/config/file.sqlite3 $DATADIR
chcon -Rt svirt_sandbox_file_t $DATADIR

docker stop cpm-promdash
docker rm cpm-promdash
docker run  \
	-v $DATADIR:/tmp/prom \
	-p $LOCAL_IP:3000:3000 \
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
	-p $LOCAL_IP:9090:9090 \
	--name=cpm-prometheus -d prom/prometheus:latest
##############

echo "restarting cpm-collect container..."
sleep 2
docker stop cpm-collect
docker rm cpm-collect
docker run -e DB_HOST=cpm-admin.crunchy.lab \
	--hostname="cpm-collect" \
	--log-driver=fluentd \
	--log-opt fluentd-address=$FLUENT_URL \
	--log-opt fluentd-tag=docker.cpm-collect \
	-e CONT_POLL_INT=4 \
	-e SERVER_POLL_INT=4 \
	-e HC_POLL_INT=4 \
	-e CPMBASE=/var/cpm \
	-e DB_PORT=5432 -e DB_USER=postgres \
	-e SWARM_MANAGER_URL=$SWARM_MANAGER_URL \
	-d --name=cpm-collect \
	crunchydata/cpm-collect:latest 

echo "testing containers for DNS resolution...."

ping -c 2 cpm-web.crunchy.lab
ping -c 2 cpm-admin.crunchy.lab
ping -c 2 cpm-task.crunchy.lab
ping -c 2 cpm-promdash.crunchy.lab
ping -c 2 cpm-prometheus.crunchy.lab
ping -c 2 cpm-collect.crunchy.lab

