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
# $1 is the master host we are going to do the backup from
# $2 is the username to connect with
# $3 is the password to connect with
#
# destroy whatever data was there....we leave it to the DBA to 
# have done a backup of this database if they wanted it
rm -rf /pgdata/*

source /var/cpm/bin/setenv.sh

export PGPASSFILE=/tmp/pgpass

echo "*:*:*:"$2":"$3  >> $PGPASSFILE

chmod 600 $PGPASSFILE

pg_basebackup -R --pgdata /pgdata --host=$1 --port=5432 -U $2
