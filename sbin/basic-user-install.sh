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
sudo mkdir -p $INSTALLDIR/keys
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
sudo rpm -Uvh http://yum.postgresql.org/9.4/redhat/rhel-7-x86_64/pgdg-centos94-9.4-1.noarch.rpm
sudo yum install -y postgresql94 postgresql94-contrib postgresql94-server libxslt unzip openssh-clients hostname bind-utils net-tools sysstat

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

echo "SECURITY WARNING- turning off and disabling your firewall!"
sudo systemctl stop firewalld.service
sudo systemctl disable firewalld.service

# start up local cpmagent
echo "starting cpmagent...."
sudo cp $INSTALLDIR/config/cpmagent.service /usr/lib/systemd/system
sudo systemctl enable cpmagent.service
sudo systemctl start cpmagent.service


sed -i "s/crunchy.lab/$DOMAIN/g" ./bu-init-cpm.sh

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

# generate keys for cpm and cpm-admin
echo "generating keys for cpm and cpm-admin containers..."
sudo $INSTALLDIR/bin/gen-keys.sh
