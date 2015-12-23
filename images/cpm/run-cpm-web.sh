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

#CPMROOT=/home/jeffmc/devproject/src/github.com/crunchydata/crunchy-postgresql-manager
#LOCAL_IP=192.168.0.107
#SWARM_MANAGER_URL=tcp://$LOCAL_IP:8000 
#FLUENT_URL=$LOCAL_IP:24224

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
	echo "This script must be run as root" 1>&2
	exit 1
fi

if [ -z "$SWARM_MANAGER_URL" ]; then
	echo "This script requires SWARM_MANAGER_URL" 1>&2
	exit 1
fi
if [ -z "$LOCAL_IP" ]; then
	echo "This script requires LOCAL_IP" 1>&2
	exit 1
fi
if [ -z "$CPMROOT" ]; then
	echo "This script requires CPMROOT" 1>&2
	exit 1
fi
if [ -z "$KEYSDIR" ]; then
	echo "This script requires KEYSDIR" 1>&2
	exit 1
fi

#echo "setting up keys dir..."
#export KEYSDIR=/var/cpm/keys
#mkdir $KEYSDIR
#cp $CPMROOT/sbin/key.pem $KEYSDIR
#cp $CPMROOT/sbin/cert.pem $KEYSDIR
#chcon -Rt svirt_sandbox_file_t $KEYSDIR

echo "setting up log dir..."
export LOGDIR=/var/cpm/logs
mkdir $LOGDIR
chcon -Rt svirt_sandbox_file_t $LOGDIR

echo "restarting cpm-web container..."
docker -H $SWARM_MANAGER_URL stop cpm-web
docker -H $SWARM_MANAGER_URL rm cpm-web

chcon -Rt svirt_sandbox_file_t $CPMROOT/images/cpm/www/v3

docker -H $SWARM_MANAGER_URL run --name=cpm-web -d \
	-p $LOCAL_IP:13001:13001 \
	-e constraint:host==$LOCAL_IP \
	-v $LOGDIR:/cpmlogs \
	-v $KEYSDIR:/cpmkeys \
	-v $CPMROOT/images/cpm/www/v3:/www \
	crunchydata/cpm:latest

