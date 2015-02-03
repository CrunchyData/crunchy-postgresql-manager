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

docker rm cpm
sudo chcon -Rt svirt_sandbox_file_t $INSTALLDIR/images/crunchy-cpm/www/v2
docker run --name=cpm -d -v $INSTALLDIR/images/crunchy-cpm/www/v2:/www crunchy-cpm

sleep 2
docker rm cluster-admin
docker run -e DB_HOST=127.0.0.1 \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cluster-admin -d -v /var/lib/pgsql/cluster-admin:/pgdata crunchy-admin

sleep 2
docker rm cluster-backup
docker run -e DB_HOST=cluster-admin.crunchy.lab \
	-e DB_PORT=5432 -e DB_USER=postgres \
	--name=cluster-backup -d crunchy-backup

sleep 2
docker rm cluster-mon
INFLUXDIR=/tmp/influxdb
sudo chcon -Rt svirt_sandbox_file_t $INFLUXDIR
docker run -e DB_HOST=cluster-admin.crunchy.lab \
	-e DB_PORT=5432 -e DB_USER=postgres \
	-v $INFLUXDIR:/monitordata \
	-d --name=cluster-mon crunchy-mon

sleep 2
ping -c 2 cpm.crunchy.lab
ping -c 2 cluster-admin.crunchy.lab
ping -c 2 cluster-backup.crunchy.lab
ping -c 2 cluster-mon.crunchy.lab

exit

docker rm dashboard
docker run --name=dashboard -d crunchy-dashboard

docker run --name=backup-job-blah \
	-e BACKUP_HOST=blah.crunchy.lab \
	-e BACKUP_PORT=5432 \
	-e BACKUP_USER=postgres \
	-e BACKUP_SERVER_URL=cluster-backup.crunchy.lab:13010 \
	-v /var/lib/pgsql/blah-backup-201412181707:/pgdata \
	-d crunchy-backup-job

