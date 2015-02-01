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

docker tag crunchy-cpm jmccormick2001/crunchy-cpm
docker push jmccormick2001/crunchy-cpm

docker tag crunchy-pgpool jmccormick2001/crunchy-pgpool
docker push jmccormick2001/crunchy-pgpool

docker tag crunchy-admin jmccormick2001/crunchy-admin
docker push jmccormick2001/crunchy-admin

docker tag crunchy-base jmccormick2001/crunchy-base
docker push jmccormick2001/crunchy-base

docker tag crunchy-mon jmccormick2001/crunchy-mon
docker push jmccormick2001/crunchy-mon

docker tag crunchy-backup jmccormick2001/crunchy-backup
docker push jmccormick2001/crunchy-backup

docker tag crunchy-backup-job jmccormick2001/crunchy-backup-job
docker push jmccormick2001/crunchy-backup-job

docker tag crunchy-node jmccormick2001/crunchy-node
docker push jmccormick2001/crunchy-node

docker tag crunchy-dashboard jmccormick2001/crunchy-dashboard
docker push jmccormick2001/crunchy-dashboard
