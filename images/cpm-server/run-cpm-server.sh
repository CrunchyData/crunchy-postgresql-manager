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

EFK_IP=192.168.0.107
LOCAL_IP=192.168.0.107
SWARM_MASTER_IP=192.168.0.107
SERVERNAME=espresso
SWARM_MANAGER_URL=tcp://$SWARM_MASTER_IP:8000

echo "restarting cpm-server"
docker stop cpm-$SERVERNAME
docker rm cpm-$SERVERNAME
docker -H $SWARM_MANAGER_URL \
	run --name=cpm-$SERVERNAME -d \
	--privileged \
	--log-driver=fluentd \
	--log-opt fluentd-address=$EFK_IP:24224 \
	--log-opt fluentd-tag=docker.cpm-$SERVERNAME \
	-e constraint:host==$LOCAL_IP \
	-p $LOCAL_IP:10001:10001 \
	-v /:/rootfs \
	-v /var/cpm/data/pgsql:/var/cpm/data/pgsql \
	-v /var/run/docker.sock:/var/run/docker.sock \
	crunchydata/cpm-server:latest

