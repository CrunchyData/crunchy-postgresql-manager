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

#export CPMROOT=/home/jeffmc/devproject/src/github.com/crunchydata/crunchy-postgresql-manager
#export LOCAL_IP=192.168.0.107
#export SWARM_MANAGER_URL=tcp://$LOCAL_IP:8000 

if [ -z "$CPMROOT" ]; then
	echo "CPMROOT is a required env var"
	exit 1
fi
if [ -z "$LOCAL_IP" ]; then
	echo "LOCAL_IP is a required env var"
	exit 1
fi
if [ -z "$SWARM_MANAGER_URL" ]; then
	echo "SWARM_MANAGER_URL is a required env var"
	exit 1
fi

echo "restarting cpm-promdash container..."
sleep 2
export DATADIR=/var/cpm/data/promdash
mkdir -p  $DATADIR
chmod 777 $DATADIR
cp $CPMROOT/config/file.sqlite3 $DATADIR
chcon -Rt svirt_sandbox_file_t $DATADIR

docker -H $SWARM_MANAGER_URL stop cpm-promdash
docker -H $SWARM_MANAGER_URL rm cpm-promdash
docker -H $SWARM_MANAGER_URL run  \
	-v $DATADIR:/tmp/prom \
	-p $LOCAL_IP:3000:3000 \
	-e constraint:host==$LOCAL_IP \
	-e DATABASE_URL=sqlite3:/tmp/prom/file.sqlite3 \
	--name=cpm-promdash -d prom/promdash
###############
echo "restarting cpm-prometheus container..."
sleep 2
export PROMCONFIG=/var/cpm/config/prometheus.yml
chmod 777 $PROMCONFIG
chcon -Rt svirt_sandbox_file_t $PROMCONFIG

docker -H $SWARM_MANAGER_URL stop cpm-prometheus
docker -H $SWARM_MANAGER_URL rm cpm-prometheus
docker -H $SWARM_MANAGER_URL run  \
	-v $PROMCONFIG:/etc/prometheus/prometheus.yml \
	-p $LOCAL_IP:9090:9090 \
	-e constraint:host==$LOCAL_IP \
	--name=cpm-prometheus -d prom/prometheus:latest
