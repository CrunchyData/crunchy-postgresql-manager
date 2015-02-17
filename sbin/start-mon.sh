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

#
# start up the monitoring service
#
# the service looks for the following env vars to be set on
# startup
#
# $DB_HOST host we are connecting to
# $DB_PORT pg port we are connecting to
# $DB_USER pg user we are connecting with
#

export LD_LIBRARY_PATH=/usr/pgsql-9.4/lib
export PATH=$PATH:/usr/pgsql-9.4/bin

#
# start influx
#
/usr/bin/influxdb -stdout=true -pidfile /tmp/influxdb.pid -config /opt/cpm/conf/config.toml > /cpmlogs/crunchy-mon-influx.log &

sleep 2

/opt/cpm/bin/monserver -log_dir=/cpmlogs  -logtostderr=false  &

#
# block with the dummy server, allows for hot swapping the backupserver# when needed

/opt/cpm/bin/dummyserver > /tmp/dummy.log 
