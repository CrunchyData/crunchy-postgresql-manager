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

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
	echo "This script must be run as root" 1>&2
	exit 1
fi

rpm -Uvh http://yum.postgresql.org/9.5/redhat/rhel-7-x86_64/pgdg-centos95-9.5-2.noarch.rpm

yum install -y postgresql95 postgresql95-contrib postgresql95-server
yum install -y kernel-devel wget git mercurial net-tools bind-utils golang docker

systemctl start docker.service
docker pull centos:7
docker pull prom/prometheus
docker pull prom/promdash


export GOPATH=/home/vagrant/devproject
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
mkdir -p $GOPATH $GOPATH/bin $GOPATH/src $GOPATH/pkg
export DEVBASE=$GOPATH/src/github.com/crunchydata/crunchy-postgresql-manager
export CPMBASE=/var/cpm

mkdir -p $CPMBASE/bin
mkdir -p $CPMBASE/config
mkdir -p $CPMBASE/data/pgsql
mkdir -p $CPMBASE/logs
mkdir -p $CPMBASE/keys

chcon -Rt svirt_sandbox_file_t $CPMBASE

cd $GOPATH
go get github.com/tools/godep
go get github.com/crunchydata/crunchy-postgresql-manager
cd $DEVBASE
godep restore
make build

cd ./images/cpm-efk
make download
cd ../..

make buildimages

chown -R vagrant:vagrant $GOPATH

cp $DEVBASE/sbin/* $CPMBASE/bin
