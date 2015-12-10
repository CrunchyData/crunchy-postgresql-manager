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

LOCAL_IP=192.168.0.107

EFKDATA=/home/jeffmc/devproject/src/github.com/crunchydata/crunchy-postgresql-manager/images/cpm-efk/elasticsearch-data
chcon -Rt svirt_sandbox_file_t $EFKDATA

# port 24224 is the fluentd listen port
FLUENTD_URL=$LOCAL_IP:24224
# port 5140 is the syslog server port that fluentd listens on 
EFK_SYSLOG_URL=$LOCAL_IP:5140
# port 5601  is the kibana http port
KIBANA_URL=$LOCAL_IP:5601

echo "restarting cpm-efk"
docker stop cpm-efk
docker rm cpm-efk
docker run --name=cpm-efk -d \
	-v $EFKDATA:/elasticsearch/data \
	-p $FLUENTD_URL:24224 \
	-p $EFK_SYSLOG_URL:5140 \
	-p $KIBANA_URL:5601 \
	crunchydata/cpm-efk:latest

