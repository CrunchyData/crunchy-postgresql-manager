#!/bin/bash -x


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

INSTALLDIR=$HOME/cpm

LOGDIR=/opt/cpm/logs
sudo mkdir -p $LOGDIR
sudo chmod -R 777 $LOGDIR
sudo chcon -Rt svirt_sandbox_file_t $LOGDIR

docker rm cpm
sudo chcon -Rt svirt_sandbox_file_t $INSTALLDIR/images/cpm/www/v2
docker run --name=cpm -d \
	-v $LOGDIR:/cpmlogs \
	-v $INSTALLDIR/images/cpm/www/v2:/www cpm

sleep 2
docker rm cpm-admin
docker run -e DB_HOST=127.0.0.1 \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-admin -d -v $LOGDIR:/cpmlogs -v /var/lib/pgsql/cpm-admin:/pgdata cpm-admin

sleep 2
docker rm cpm-backup
docker run -e DB_HOST=cpm-admin.crunchy.lab \
	-v $LOGDIR:/cpmlogs \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cpm-backup -d cpm-backup

sleep 2
docker rm cpm-mon
INFLUXDIR=/tmp/influxdb
sudo chcon -Rt svirt_sandbox_file_t $INFLUXDIR
docker run -e DB_HOST=cpm-admin.crunchy.lab \
	-e DB_PORT=5432 -e DB_USER=postgres \
	-v $LOGDIR:/cpmlogs \
	-v $INFLUXDIR:/monitordata \
	-d --name=cpm-mon cpm-mon

sleep 2
ping -c 2 cpm.crunchy.lab
ping -c 2 cpm-admin.crunchy.lab
ping -c 2 cpm-backup.crunchy.lab
ping -c 2 cpm-mon.crunchy.lab

exit

docker rm cpm-dashboard
docker run --name=cpm-dashboard -d cpm-dashboard

docker run --name=backup-job-blah \
	-e BACKUP_HOST=blah.crunchy.lab \
	-e BACKUP_PORT=5432 \
	-e BACKUP_USER=postgres \
	-e BACKUP_SERVER_URL=cpm-backup.crunchy.lab:13010 \
	-v $LOGDIR:/opt/cpm/logs \
	-v /var/lib/pgsql/blah-backup-201412181707:/pgdata \
	-d cpm-backup-job

