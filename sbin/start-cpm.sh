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

if [ -z "$SWARM_MANAGER_URL" ]; then
	echo "SWARM_MANAGER_URL env var is required"
	exit 1
fi

docker -H $SWARM_MANAGER_URL start skybridge
sleep 4
docker -H $SWARM_MANAGER_URL start cpm-efk
sleep 3
docker -H $SWARM_MANAGER_URL start cpm-server1
docker -H $SWARM_MANAGER_URL start cpm-admin
docker -H $SWARM_MANAGER_URL start cpm-web
docker -H $SWARM_MANAGER_URL start cpm-task
docker -H $SWARM_MANAGER_URL start cpm-prometheus
docker -H $SWARM_MANAGER_URL start cpm-promdash
sleep 4
docker -H $SWARM_MANAGER_URL start cpm-collect

