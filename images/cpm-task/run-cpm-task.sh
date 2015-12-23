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

if [ -z "$LOCAL_IP" ]; then
	echo "LOCAL_IP is a required env var"
	exit 1
fi
if [ -z "$SWARM_MANAGER_URL" ]; then
	echo "SWARM_MANAGER_URL is a required env var"
	exit 1
fi
if [ -z "$FLUENT_URL" ]; then
	echo "FLUENT_URL is a required env var"
	exit 1
fi

#export LOCAL_IP=192.168.0.107
#export SWARM_MANAGER_URL=tcp://$LOCAL_IP:8000 
#export FLUENT_URL=$LOCAL_IP:24224

echo "restarting cpm-task container..."
sleep 2
docker -H $SWARM_MANAGER_URL stop cpm-task
docker -H $SWARM_MANAGER_URL rm cpm-task
docker -H $SWARM_MANAGER_URL run -e DB_HOST=cpm-admin.crunchy.lab \
	--log-driver=fluentd \
	--log-opt fluentd-address=$FLUENT_URL \
	--log-opt fluentd-tag=docker.cpm-task \
	-e constraint:host==$LOCAL_IP \
	-e CPMBASE=/var/cpm \
	-e SWARM_MANAGER_URL=$SWARM_MANAGER_URL \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-task \
	-d crunchydata/cpm-task:latest

