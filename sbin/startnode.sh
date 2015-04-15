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

source /var/cpm/bin/setenv.sh

sleep 3
initdb.sh > /tmp/initdb.log 2> /tmp/initdb.err
echo "host all all 0.0.0.0/0 md5" >> /pgdata/pg_hba.conf
echo "listen_addresses='*'" >> /pgdata/postgresql.conf
startpg.sh > /tmp/startpg.log
sleep 3
seed.sh > /tmp/seed.log
start-cpmagentserver.sh 
