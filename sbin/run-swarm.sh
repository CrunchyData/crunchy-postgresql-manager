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

LOCALIP=192.168.0.107
SECONDARYIP=192.168.0.105
SWARM_TOKEN=7228941e30510cf52742e1c011d00bca

#swarm manage --host $LOCALIP:8000 nodes://$LOCALIP:2375,$SECONDARYIP:2375 &
swarm manage --host $LOCALIP:8000 nodes://$LOCALIP:2375 &
sleep 4
swarm join --addr=$LOCALIP:2375 token://$SWARM_TOKEN
#swarm join --addr=$SECONDARYIP:2375 token://$SWARM_TOKEN

