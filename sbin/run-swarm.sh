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

# This install script assumes a Centos7 server is the installation host OS.
# It will download some additional information to build a Docker image for that
# type of host.
#
# Before running, confirm you can run docker commands as this user with
# 'docker info'.  See the README.md file for more information on setup for
# this program and Docker.  You will also need to install the dnsbridge
# program before this one.
#

SWARM_CLUSTER_FILE=/var/cpm/data/swarm_cluster_file

SWARM_URL=$LOCAL_IP:8000
DOCKER_PORT=2375

# use the random strategy just for testing - jeffmc
/usr/local/bin/swarm manage --strategy random --host $SWARM_URL file://$SWARM_CLUSTER_FILE &
sleep 4
/usr/local/bin/swarm join --addr=$LOCAL_IP:$DOCKER_PORT file://$SWARM_CLUSTER_FILE &
#swarm list file:///tmp/my_cluster

