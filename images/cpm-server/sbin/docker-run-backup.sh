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
# wraps docker command
#

source /var/cpm/bin/setenv.sh

echo $1 " is docker run pgdatapath " >> $CLUSTER_LOG
echo $2 " is docker run container name " >> $CLUSTER_LOG
echo $3 " is docker run container type " >> $CLUSTER_LOG
echo $4 " is docker run cpu " >> $CLUSTER_LOG
echo $5 " is docker run mem " >> $CLUSTER_LOG
echo $6 " is docker env vars " >> $CLUSTER_LOG
docker rm $2
docker run --name=$2 $6 -c $4 -m $5 -d -v $1:/pgdata $3
#docker run --name=$2 $6 -c $4 -m $5 --rm=true -v $1:/pgdata $3


