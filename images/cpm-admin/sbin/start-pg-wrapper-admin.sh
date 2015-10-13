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

#
# start pg, will initdb if /pgdata is empty as a way to bootstrap
#

#
# clean up any previous PG lock files
#
rm /tmp/.s.PGSQL*

source /var/cpm/bin/setenv.sh

HBA=/var/cpm/conf/admin/pg_hba.conf
#
# the initial start of postgres will create the clusteradmin database
#
if [ ! -f /pgdata/postgresql.conf ]; then
        echo "pgdata is empty"
        initdb -D /pgdata  > /tmp/initdb.log &> /tmp/initdb.err
        echo "setting domain to " $THISDOMAIN >> /tmp/start-db.log

        cp $HBA /tmp
        sed "s/crunchy.lab/$THISDOMAIN/g" /tmp/pg_hba.conf > /pgdata/pg_hba.conf
        cat /pgdata/pg_hba.conf >> /tmp/start-db.log

        sed -i "s/crunchy.lab/$THISDOMAIN/g" /var/cpm/bin/setup.sql

        cp /var/cpm/conf/admin/postgresql.conf /pgdata/
        echo "starting db" >> /tmp/start-db.log

        pg_ctl -D /pgdata start
        sleep 3
        echo "building clusteradmin db" >> /tmp/start-db.log
        psql -U postgres < /var/cpm/bin/setup.sql
        exit
fi

#
# clean up any old pid file that might have remained
# during a bad shutdown of the container/postgres
#
rm /pgdata/postmaster.pid
#
# the normal startup of pg
#
pg_ctl -D /pgdata start 

