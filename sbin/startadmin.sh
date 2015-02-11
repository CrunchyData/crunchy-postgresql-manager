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

source /opt/cpm/bin/setenv.sh

/opt/cpm/bin/start-pg-wrapper-admin.sh &
export KUBE_URL=$KUBE_URL

# log output will go to /tmp into files created by glog
# named similar to adminapi.7869c0a96e4c.postgres.log.INFO.20150204-192844.2313
/opt/cpm/bin/adminapi -log_dir=/cpmlogs -logtostderr=false &

/opt/cpm/bin/dummyserver > /tmp/dummyserver.log 

