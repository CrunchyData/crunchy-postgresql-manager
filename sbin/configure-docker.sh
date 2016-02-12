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
if [ -z "$LOCAL_IP" ]; then
	echo "LOCAL_IP is not set and is required for this script"
	exit	
fi
if [[ $EUID -ne 0 ]]; then
           echo "This script must be run as root" 1>&2
              exit 1
fi

NEWOPTIONS="--selinux-enabled --bip=172.17.42.1/16 --dns-search=crunchy.lab --dns=$LOCAL_IP --dns=192.168.0.1 -H unix:///var/run/docker.sock --label host=$LOCAL_IP --label profile=SM -H tcp://$LOCAL_IP:2375"

echo "'"$NEWOPTIONS"'" >> /etc/sysconfig/docker


