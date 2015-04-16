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

INSTALLDIR=/var/cpm

echo "enter domain to use..."
read DOMAIN
echo "enter static ip to use..."
read THISIP

echo "setting up log dir..."
LOGDIR=$INSTALLDIR/logs
sudo mkdir -p $LOGDIR
sudo chmod -R 777 $LOGDIR
sudo chcon -Rt svirt_sandbox_file_t $LOGDIR

echo "setting up keys dir..."
KEYSDIR=/var/cpm/keys
sudo chcon -Rt svirt_sandbox_file_t $KEYSDIR

echo "deleting all old log files...."
sudo rm -rf $LOGDIR/*

echo "restarting cpm container..."
sudo chcon -Rt svirt_sandbox_file_t $INSTALLDIR/www/v2
docker rm cpm
docker run --name=cpm -d \
	-p $THISIP:8080:13000 \
	-v $LOGDIR:/cpmlogs \
	-v $KEYSDIR:/cpmkeys \
	-v $INSTALLDIR/www/v2:/www \
	crunchydata/cpm

echo "restarting cpm-admin container..."
sleep 2
DBDIR=$INSTALLDIR/data/pgsql/cpm-admin
sudo mkdir -p $DBDIR
sudo chown postgres:postgres $DBDIR
sudo chcon -Rt svirt_sandbox_file_t $DBDIR
docker rm cpm-admin
docker run -e DB_HOST=127.0.0.1 \
	-p $THISIP:8081:13000 \
	-e DOMAIN=$DOMAIN \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-admin -d -v $LOGDIR:/cpmlogs -v $DBDIR:/pgdata \
	-v $KEYSDIR:/cpmkeys \
	crunchydata/cpm-admin

echo "restarting cpm-backup container..."
sleep 2
docker rm cpm-backup
docker run -e DB_HOST=cpm-admin.$DOMAIN \
	-v $LOGDIR:/cpmlogs \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-backup -d \
	crunchydata/cpm-backup

echo "restarting cpm-mon container..."
sleep 2
INFLUXDIR=$INSTALLDIR/data/influxdb
sudo mkdir -p $INFLUXDIR
sudo chcon -Rt svirt_sandbox_file_t $INFLUXDIR
docker rm cpm-mon
docker run -e DB_HOST=cpm-admin.$DOMAIN \
	-p $THISIP:8083:8083 \
	--hostname="cpm-mon" \
	-e DB_PORT=5432 -e DB_USER=postgres \
	-v $LOGDIR:/cpmlogs \
	-v $INFLUXDIR:/monitordata \
	-d --name=cpm-mon \
	crunchydata/cpm-mon

sleep 2
echo "testing containers for DNS resolution...."
ping -c 2 cpm.$DOMAIN
ping -c 2 cpm-admin.$DOMAIN
ping -c 2 cpm-backup.$DOMAIN
ping -c 2 cpm-mon.$DOMAIN


