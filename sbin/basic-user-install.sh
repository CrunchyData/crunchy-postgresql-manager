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

# prompt for static IP to use
echo "enter static ip to use..."
read STATICIP

# verify running as root user

# download CPM software archive
# extract archive into /opt/cpm

# pull down CPM Docker images from dockerhub

docker pull jmccormick2001/crunchy-cpm
docker pull jmccormick2001/crunchy-pgpool
docker pull jmccormick2001/crunchy-admin
docker pull jmccormick2001/crunchy-base
docker pull jmccormick2001/crunchy-mon
docker pull jmccormick2001/crunchy-backup
docker pull jmccormick2001/crunchy-backup-job
docker pull jmccormick2001/crunchy-node
docker pull jmccormick2001/crunchy-dashboard

# install deps
yum -y install gcc make golang docker-io
rpm -Uvh http://dl.fedoraproject.org/pub/epel/7/x86_64/e/epel-release-7-5.noarch.rpm
rpm -Uvh http://yum.postgresql.org/9.3/redhat/rhel-7-x86_64/pgdg-redhat93-9.3-1.noarch.rpm
yum install -y postgresql93 postgresql93-contrib postgresql93-server libxslt unzip openssh-clients hostname bind-utils net-tools
yum install -y bind

# configure and start docker
echo "setting up docker..."
systemctl enable docker.service
cp $INSTALLDIR/config/docker /etc/sysconfig
systemctl start docker.service

# setup the local pg database used by dnsbridge
echo "setting up local postgres database...."
su - postgres -c '/usr/pgsql-9.3/bin/initdb -D /var/lib/pgsql/9.3/data'
systemctl enable postgresql-9.3.service
systemctl start postgresql-9.3.service
cp $INSTALLDIR/sql/bridge.sql /tmp
su - postgres -c '/usr/pgsql-9.3/bin/psql -U postgres postgres < /tmp/bridge.sql'


# configure DNS
echo "setting up DNS..."
systemctl enable named.service
named-checkzone 0.17.172.in-addr.arpa  /var/named/dynamic/0.17.172.zone.db
named-checkzone crunchy.lab  /var/named/dynamic/crunchy.lab.db
systemctl start named.service

# set up PG data directory
echo "setting up CPM admin directory..."
export DATADIR=/var/lib/pgsql
mkdir $DATADIR/cluster-admin
chown postgres:postgres $DATADIR/cluster-admin
chcon -Rt svirt_sandbox_file_t $DATADIR/cluster-admin/

# set up the CPM www directory
chcon -Rt svirt_sandbox_file_t $INSTALLDIR/www/v2

echo "enabling dnsbridge services.."
systemctl enable dnsbridgeclient.service
systemctl enable dnsbridgeserver.service
echo "starting dnsbridge services.."
systemctl start dnsbridgeserver.service
systemctl start dnsbridgeclient.service

echo "starting up CPM containers..."

docker run --name=cpm -d -v $INSTALLDIR/images/crunchy-cpm/www/v2:/www crunchy-cpm
sleep 2
docker run -e DB_HOST=127.0.0.1 \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cluster-admin -d -v /var/lib/pgsql/cluster-admin:/pgdata crunchy-admin
sleep 2
docker run -e DB_HOST=cluster-admin.crunchy.lab \
	        -e DB_PORT=5432 -e DB_USER=postgres \
		        --name=cluster-backup -d crunchy-backup
sleep 2
INFLUXDIR=$INSTALLDIR/influxdb
chcon -Rt svirt_sandbox_file_t $INFLUXDIR
docker run -e DB_HOST=cluster-admin.crunchy.lab \
	-e DB_PORT=5432 -e DB_USER=postgres \
	-v $INFLUXDIR:/monitordata \
	-d --name=cluster-mon crunchy-mon

sleep 2
docker run --name=dashboard -d crunchy-dashboard

echo "installation complete"

