#!/bin/bash
#

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

# This install script assumes a registered RHEL 7 server is the installation host OS.
# It will download some additional information to build a Docker image for that
# type of host.
#
# Before running, confirm you can run docker commands as this user with
# 'docker info'.  See the README.md file for more information on setup for
# this program and Docker.  You will also need to install the dnsbridge
# program before this one.
#

# Exit installation on any unexpected error
set -e

# install deps
export INSTALLDIR=$(pwd)

# set the gopath
source $INSTALLDIR/setpath.sh

make

server=$(hostname)

scp bin/* sql/loadtest.sql  \
	root@$server:/cluster/bin/

scp sbin/*  \
	root@$server:/cluster/bin/

sudo scp config/cpmagent.service  \
	root@$server:/usr/lib/systemd/system

ssh root@$server "systemctl enable cpmagent.service"
ssh root@$server "systemctl start cpmagent.service"

sudo mkdir /var/lib/pgsql/cluster-admin
sudo chown postgres:postgres /var/lib/pgsql/cluster-admin
sudo chcon -Rt svirt_sandbox_file_t /var/lib/pgsql/cluster-admin/
sudo chcon -Rt svirt_sandbox_file_t $INSTALLDIR/images/crunchy-cpm/www/v2

