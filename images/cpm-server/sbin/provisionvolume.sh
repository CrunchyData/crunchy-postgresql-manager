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
# provision a disk volume for a pg container
# $1 is the full path to the server's PGDataPath/ContainerName
#
FULLPATH=$1

EVENT_LOG=/tmp/server-events.log
echo "volume to provision is ["$FULLPATH"]" >> $EVENT_LOG

echo `date` >> $EVENT_LOG
echo "provisioning volume: removing " $FULLPATH  >> $EVENT_LOG
rm -rf $FULLPATH

echo "mkdir volume"  $FULLPATH >> $EVENT_LOG
#mkdir $VOLUME_BASE/$1
mkdir -p $FULLPATH
echo "chmod volume"  $FULLPATH >> $EVENT_LOG
chmod 0700 $FULLPATH
echo "chown volume"  $FULLPATH >> $EVENT_LOG
chown postgres:postgres $FULLPATH
echo "chcon volume"  $FULLPATH >> $EVENT_LOG
#chcon -Rt svirt_sandbox_file_t  $VOLUME_BASE/$1
chcon -Rt svirt_sandbox_file_t  $FULLPATH


exit 0


