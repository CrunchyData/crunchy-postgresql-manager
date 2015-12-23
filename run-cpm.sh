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

#
# these env vars are passed down when components are started
# from this master script
#
export CPMROOT=/home/jeffmc/devproject/src/github.com/crunchydata/crunchy-postgresql-manager
export LOCAL_IP=192.168.0.107
export SWARM_MANAGER_URL=tcp://$LOCAL_IP:8000 
export FLUENT_URL=$LOCAL_IP:24224
export CPM_DOMAIN=crunchy.lab
# SERVERNAME is the name we give the CPM server container cpm-$SERVERNAME
export SERVERNAME=server1

# keys dir
echo "setting up keys dir..."
export KEYSDIR=/var/cpm/keys
mkdir $KEYSDIR
cp $CPMROOT/sbin/key.pem $KEYSDIR
cp $CPMROOT/sbin/cert.pem $KEYSDIR
chcon -Rt svirt_sandbox_file_t $KEYSDIR


$CPMROOT/sbin/run-skybridge.sh

echo "sleeping a bit while skybridge starts up...."
sleep 6

$CPMROOT/images/cpm-efk/run-cpm-efk.sh

echo "sleeping a bit while cpm-efk starts up...."
sleep 6

$CPMROOT/images/cpm-server/run-cpm-server.sh

$CPMROOT/images/cpm/run-cpm-web.sh

$CPMROOT/images/cpm-admin/run-cpm-admin.sh

$CPMROOT/images/cpm-task/run-cpm-task.sh

$CPMROOT/images/cpm-prometheus/run-cpm-prometheus.sh

echo "sleeping a bit while cpm-prometheus starts up...."
sleep 6

$CPMROOT/images/cpm-collect/run-cpm-collect.sh

echo "testing containers for DNS resolution...."

ping -c 2 cpm-web.$CPM_DOMAIN
ping -c 2 cpm-admin.$CPM_DOMAIN
ping -c 2 cpm-task.$CPM_DOMAIN
ping -c 2 cpm-promdash.$CPM_DOMAIN
ping -c 2 cpm-prometheus.$CPM_DOMAIN
ping -c 2 cpm-collect.$CPM_DOMAIN
ping -c 2 cpm-$SERVERNAME.$CPM_DOMAIN

