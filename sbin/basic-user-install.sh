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

# basic user installation script for CPM
#
# This script assumes a registered RHEL 7 server is the installation host OS.
#

# Exit installation on any unexpected error
set -e

# set the istall directory
export INSTALLDIR=/opt/cpm

sudo mkdir -p $INSTALLDIR/bin
sudo mkdir -p $INSTALLDIR/config
sudo mkdir -p $INSTALLDIR/www
export LOGDIR=$INSTALLDIR/logs
sudo mkdir -p $LOGDIR

# prompt for static IP to use
echo "enter static ip to use...(e.g. 192.168.56.103)"
read STATICIP
# prompt for domain name to use
echo "enter domain name to use...(e.g crunchy.lab) "
read DOMAIN

#Check if current user is member to the wheel group
username= whoami
if groups $username | grep &>/dev/null 'wheel'; then
	echo "Group permissions ok"
else
	echo "You must have sudo privledges to run this install"
	exit
fi

# don't error if packages are already installed
set +e

# install deps
sudo yum -y install docker-io
sudo rpm -Uvh http://dl.fedoraproject.org/pub/epel/7/x86_64/e/epel-release-7-5.noarch.rpm
sudo rpm -Uvh http://yum.postgresql.org/9.3/redhat/rhel-7-x86_64/pgdg-redhat93-9.3-1.noarch.rpm
sudo yum install -y postgresql93 postgresql93-contrib postgresql93-server libxslt unzip openssh-clients hostname bind-utils net-tools

set -e
# make sure the user is in the docker group
sudo usermod -a -G docker $USER

# make sure docker is started and enabled
sudo systemctl enable docker.service
sudo systemctl start docker.service

# move the CPM media to the /opt/cpm installation directory
sudo mv `pwd`/bin/* $INSTALLDIR/bin
sudo mv `pwd`/config/* $INSTALLDIR/config
sudo mv `pwd`/www/* $INSTALLDIR/www

echo "TEMPORARY HACK - turning off your firewall!"
sudo systemctl stop firewalld.service

# start up local cpmagent
echo "starting cpmagent...."
sudo cp $INSTALLDIR/config/cpmagent.service /usr/lib/systemd/system
sudo systemctl enable cpmagent.service
sudo systemctl start cpmagent.service

# pull down CPM Docker images from dockerhub
echo "pulling down cpm docker images...."
docker pull crunchydata/cpm
docker pull crunchydata/cpm-pgpool
docker pull crunchydata/cpm-admin
docker pull crunchydata/cpm-base
docker pull crunchydata/cpm-mon
docker pull crunchydata/cpm-backup
docker pull crunchydata/cpm-backup-job
docker pull crunchydata/cpm-node
docker pull crunchydata/cpm-dashboard


echo "starting up CPM containers..."
sudo chmod -R 777 $LOGDIR
sudo chcon -Rt svirt_sandbox_file_t $LOGDIR

echo "starting cpm web..."
sudo chcon -Rt svirt_sandbox_file_t $INSTALLDIR/www/v2
docker run --name=cpm -d \
	-v $LOGDIR:/cpmlogs \
	-v $INSTALLDIR/www/v2:/www crunchydata/cpm

sleep 2
echo "starting cpm-admin..."
export DATADIR=/var/lib/pgsql
export DBDIR=$DATADIR/cpm-admin
sudo mkdir $DBDIR
sudo chown postgres:postgres $DBDIR
sudo chcon -Rt svirt_sandbox_file_t $DBDIR
docker run -e DB_HOST=127.0.0.1 \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-admin -d -v $LOGDIR:/cpmlogs \
	-v $DBDIR:/pgdata crunchydata/cpm-admin

sleep 2
echo "starting cpm-backup..."
docker run -e DB_HOST=cpm-admin.$DOMAIN \
	-v $LOGDIR:/cpmlogs \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-backup -d crunchydata/cpm-backup

sleep 2
echo "starting cpm-mon..."
INFLUXDIR=$INSTALLDIR/data/influxdb
sudo mkdir -p $INFLUXDIR
sudo chcon -Rt svirt_sandbox_file_t $INFLUXDIR
docker run -e DB_HOST=cpm-admin.$DOMAIN \
	-e DB_PORT=5432 -e DB_USER=postgres \
	-v $LOGDIR:/cpmlogs \
	-v $INFLUXDIR:/monitordata \
	-d --name=cpm-mon crunchydata/cpm-mon

sleep 2
echo "starting cpm-dashboard..."
docker run --name=cpm-dashboard -d cpm-dashboard

sleep 2
echo "running ping test...."
ping -c 1 cpm.$DOMAIN
ping -c 1 cpm-admin.$DOMAIN
ping -c 1 cpm-backup.$DOMAIN
ping -c 1 cpm-mon.$DOMAIN
ping -c 1 cpm-dashboard.$DOMAIN

echo "installation complete"

