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

#LOCAL_IP=192.168.0.107
#SWARM_MANAGER_URL=tcp://$LOCAL_IP:8000
#CPM_DOMAIN=crunchy.lab

if [ -z "$LOCAL_IP" ]; then
	echo "LOCAL_IP env var is required"
	exit 1
fi
if [ -z "$SWARM_MANAGER_URL" ]; then
	echo "SWARM_MANAGER_URL env var is required"
	exit 1
fi
if [ -z "$CPM_DOMAIN" ]; then
	echo "CPM_DOMAIN env var is required"
	exit 1
fi

DATADIR=/var/cpm/data/etcd
mkdir -p $DATADIR
rm -rf $DATADIR/*
#	-e DOCKER_HOST=http://192.168.0.106:5000 \
# 
#	-v /run/docker.sock:/tmp/docker.sock \

chcon -Rt svirt_sandbox_file_t $DATADIR
echo "restarting skybridge container..."
docker -H $SWARM_MANAGER_URL stop skybridge
docker -H $SWARM_MANAGER_URL rm skybridge
docker \
	-H $SWARM_MANAGER_URL \
	run --name=skybridge -d \
	--hostname="skybridge" \
	--privileged \
	-p $LOCAL_IP:53:53/udp \
	-v /var/cpm/data/etcd:/etcddata \
	-e SWARM_MANAGER_URL=$SWARM_MANAGER_URL \
	-e DNS_DOMAIN=$CPM_DOMAIN \
	-e DNS_NAMESERVER=192.168.0.1 \
	-e constraint:host==$LOCAL_IP \
	crunchydata/skybridge:latest

