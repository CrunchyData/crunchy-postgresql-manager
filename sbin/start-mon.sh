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

export DB_HOST=$DB_HOST
export DB_PORT=$DB_PORT
export DB_USER=$DB_USER

source /var/cpm/bin/setenv.sh

sleep 4

#
# start influx
#
INPID=/tmp/influxdb.pid
/bin/rm -f $INPID

/usr/bin/influxdb -stdout=true -pidfile $INPID -config /var/cpm/conf/config.toml > /cpmlogs/crunchy-mon-influx.log 2> /cpmlogs/crunchy-mon-influx.stderr &

sleep 7

monserver > /cpmlogs/monserver.log

#
# block with the dummy server, allows for hot swapping the backupserver# when needed

#dummyserver > /tmp/dummy.log 
