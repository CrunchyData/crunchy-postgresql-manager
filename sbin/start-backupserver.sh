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
# start up the backup service
#
# the service looks for the following env vars to be set by
# the cpm-admin that provisioned us
#
# $BACKUP_HOST host we are connecting to
# $BACKUP_PORT pg port we are connecting to
# $BACKUP_USER pg user we are connecting with
# $BACKUP_CPMAGENT_URL cpmagent URL we send status to
#

env > /tmp/envvars.out

source /var/cpm/bin/setenv.sh

backupserver  > /cpmlogs/backupserver.log
