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

# if the syslog config file is found, assume the user wants to start the syslog daemon
if [ -f /syslogconfig/rsyslog.conf ]; then
	cp /syslogconfig/listen.conf /etc/rsyslog.d/
	rsyslogd -f /syslogconfig/rsyslog.conf
fi


echo "about to run startpgpool.sh"
startpgpool.sh

cpmcontainerserver 

#dummyserver
