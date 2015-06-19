#!/bin/bash
#

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

# install deps

sudo rpm -Uvh http://yum.postgresql.org/9.4/redhat/rhel-7-x86_64/pgdg-centos94-9.4-1.noarch.rpm

sudo yum install -y postgresql94 postgresql94-contrib postgresql94-server

export DEVROOT=/home/jeffmc/devproject
export DEVBASE=$DEVROOT/src/github.com/crunchydata/crunchy-postgresql-manager
export CPMBASE=/var/cpm

sudo mkdir -p $CPMBASE/bin
sudo mkdir -p $CPMBASE/config
sudo mkdir -p $CPMBASE/data/pgsql
sudo mkdir -p $CPMBASE/logs
sudo mkdir -p $CPMBASE/keys

sudo cp $DEVROOT/bin/cpmserveragent $CPMBASE/bin
sudo cp $DEVBASE/sbin/cert.pem $DEVBASE/sbin/key.pem $CPMBASE/keys

sudo cp $DEVBASE/sbin/* $CPMBASE/bin

sudo cp $DEVBASE/config/cpmserveragent.service  /usr/lib/systemd/system

sudo systemctl enable cpmserveragent.service
sudo systemctl start cpmserveragent.service

