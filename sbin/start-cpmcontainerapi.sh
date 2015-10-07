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
# start up the cpm container agent
#

env > /tmp/envvars.out

source /var/cpm/bin/setenv.sh

# when the container starts, see if there is a pg instance we can start up
if [ -f /pgdata/postgresql.conf ]; then
	rm /pgdata/postmaster.pid
	rm /tmp/.s.*
	su - postgres -c 'source /var/cpm/bin/setenv.sh;pg_ctl -D /pgdata start'
fi

# when the container starts, see if there is a pgpool instance we can start up
if [ -f /bin/pgpool ]; then
	echo "about to run startpgpool.sh"
	su - postgres -c 'source /var/cpm/bin/setenv.sh;startpgpool.sh'
fi

cpmcontainerapi 

#dummyserver
