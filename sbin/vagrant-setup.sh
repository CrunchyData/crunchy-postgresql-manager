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

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
	echo "This script must be run as root" 1>&2
	exit 1
fi

#
echo " install host deps"
#
rpm -Uvh http://yum.postgresql.org/9.5/redhat/rhel-7-x86_64/pgdg-centos95-9.5-2.noarch.rpm

yum install -y postgresql95 postgresql95-contrib postgresql95-server
yum install -y kernel-devel wget git mercurial net-tools bind-utils golang docker

#
echo " start docker"
#
systemctl start docker.service

#
echo " pull down images from docker hub"
#
docker pull centos:7
docker pull prom/prometheus
docker pull prom/promdash
docker pull crunchydata/skybridge2

#
echo " set up the CPM GOPATH"
#
export GOPATH=/home/vagrant/devproject
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
mkdir -p $GOPATH $GOPATH/bin $GOPATH/src $GOPATH/pkg
export DEVBASE=$GOPATH/src/github.com/crunchydata/crunchy-postgresql-manager

#
echo " set up the host CPM directory, only used by the cpmserver container"
#
export CPMBASE=/var/cpm
mkdir -p $CPMBASE/bin
mkdir -p $CPMBASE/config
mkdir -p $CPMBASE/data/pgsql
mkdir -p $CPMBASE/logs
mkdir -p $CPMBASE/keys
chcon -Rt svirt_sandbox_file_t $CPMBASE

#
echo " build CPM"
#
cd $GOPATH
go get github.com/tools/godep
go get github.com/crunchydata/crunchy-postgresql-manager
cd $DEVBASE
godep restore
make build

cd ./images/cpm-efk
make download
cd ../..

#
echo "build the CPM container images"
#
make buildimages

chown -R vagrant:vagrant $GOPATH

cp $DEVBASE/sbin/* $CPMBASE/bin


#
echo "set env vars used by installation scripts"
#
source $DEVBASE/cpmenv

#
echo "configure docker"
#
$DEVBASE/sbin/configure-docker.sh

#
echo "configure /etc/resolv.conf"
#
$DEVBASE/sbin/configure-resolv.sh

#
echo " restart docker"
#
systemctl stop docker.service
systemctl start docker.service
sleep 5

#
echo "installing swarm binary into /usr/local/bin"
#
$DEVBASE/sbin/install-swarm.sh

#
echo "run swarm"
#
$DEVBASE/sbin/run-swarm.sh
sleep 5

#
echo "starting skybridge container..."
#

$DEVBASE/sbin/run-skybridge.sh

sleep 10

#
echo "run cpm-server container"
#
$DEVBASE/images/cpm-server/run-cpm-server.sh

#
echo "run cpm app containers"
#
$DEVBASE/run-cpm.sh

chown -R vagrant:vagrant $GOPATH
