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

ROOT=~/swarmproject
mkdir -p  ~/$ROOT/bin ~/$ROOT/pkg ~/$ROOT/src

export GOPATH=$ROOT
GOBIN=$GOPATH/bin
PATH=$PATH:$GOPATH/bin

mkdir -p $GOPATH/src/github.com/docker/
cd $GOPATH/src/github.com/docker/
git clone https://github.com/docker/swarm
cd swarm
git checkout v1.0.0
go get github.com/tools/godep
$GOPATH/bin/godep go install
sudo cp $GOPATH/bin/swarm /usr/local/bin

SWARM_CLUSTER_FILE=/var/cpm/data/swarm_cluster_file
rm $SWARM_CLUSTER_FILE
echo $LOCAL_IP:2375 >> $SWARM_CLUSTER_FILE
chmod +r $SWARM_CLUSTER_FILE

