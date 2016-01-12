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

#CPMROOT=/home/jeffmc/devproject/src/github.com/crunchydata/crunchy-postgresql-manager
#LOCAL_IP=192.168.0.107

source $CPMROOT/cpmenv

if [ -z "$LOCAL_IP" ]; then
	echo "LOCAL_IP env var required"
	exit 1
fi
if [ -z "$CPMROOT" ]; then
	echo "CPMROOT env var required"
	exit 1
fi


EFKDATA=$CPMROOT/images/cpm-efk/elasticsearch-data
mkdir -p $EFKDATA
chmod 777 $EFKDATA
chcon -Rt svirt_sandbox_file_t $EFKDATA

# port 24224 is the fluentd listen port
FLUENTD_URL=$LOCAL_IP:24224
# port 5140 is the syslog server port that fluentd listens on 
EFK_SYSLOG_URL=$LOCAL_IP:5140
# port 5601  is the kibana http port
KIBANA_URL=$LOCAL_IP:5601

# the presence of these config files locally will cause the
# cpm-node containers to configure rsyslog for remote logging
# comment these lines out if you don't want this
cp $CPMROOT/images/cpm-efk/conf/listen.conf /var/cpm/config
cp $CPMROOT/images/cpm-efk/conf/rsyslog.conf /var/cpm/config

echo "restarting cpm-efk"
docker stop cpm-efk
docker rm cpm-efk
docker run --name=cpm-efk -d \
	-v $EFKDATA:/elasticsearch/data \
	-p $FLUENTD_URL:24224 \
	-p $EFK_SYSLOG_URL:5140 \
	-p $KIBANA_URL:5601 \
	crunchydata/cpm-efk:latest

