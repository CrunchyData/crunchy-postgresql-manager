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
# start up the adminapi agent
#

env > /tmp/envvars.out

source /var/cpm/bin/setenv.sh

HOSTNAME=`hostname`
#
#
#
IP=`grep $HOSTNAME /etc/hosts | cut -f 1`
adminapi  >> /tmp/adminapi.log &
