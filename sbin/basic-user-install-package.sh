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
#

# Exit installation on any unexpected error
set -e

# set the istall directory
export INSTALLDIR=/opt/cpm

# verify running as root user

# push docker images to dockerhub

docker tag cpm jmccormick2001/cpm
docker push jmccormick2001/cpm

docker tag cpm-pgpool jmccormick2001/cpm-pgpool
docker push jmccormick2001/cpm-pgpool

docker tag cpm-admin jmccormick2001/cpm-admin
docker push jmccormick2001/cpm-admin

docker tag cpm-base jmccormick2001/cpm-base
docker push jmccormick2001/cpm-base

docker tag cpm-mon jmccormick2001/cpm-mon
docker push jmccormick2001/cpm-mon

docker tag cpm-backup jmccormick2001/cpm-backup
docker push jmccormick2001/cpm-backup

docker tag cpm-backup-job jmccormick2001/cpm-backup-job
docker push jmccormick2001/cpm-backup-job

docker tag cpm-node jmccormick2001/cpm-node
docker push jmccormick2001/cpm-node

docker tag cpm-dashboard jmccormick2001/cpm-dashboard
docker push jmccormick2001/cpm-dashboard
