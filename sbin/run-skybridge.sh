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
MYHOSTIP=192.168.0.107
MYDOMAIN=crunchy.lab
DATADIR=/var/cpm/data/etcd
mkdir -p $DATADIR
rm -rf $DATADIR/*
#	-e DOCKER_HOST=http://192.168.0.106:5000 \

chcon -Rt svirt_sandbox_file_t $DATADIR
echo "restarting skybridge container..."
docker stop skybridge
docker rm skybridge
docker run --name=skybridge -d \
	--hostname="skybridge" \
	--privileged \
	-p $MYHOSTIP:53:53/udp \
	-v /run/docker.sock:/tmp/docker.sock \
	-v /var/cpm/data/etcd:/etcddata \
	-e DNS_DOMAIN=$MYDOMAIN \
	-e DNS_NAMESERVER=192.168.0.1 \
	crunchydata/skybridge:latest

