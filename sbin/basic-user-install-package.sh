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

# This install script assumes a registered RHEL 7 or CentOS 7 server is the installation host OS.
#

# Exit installation on any unexpected error
set -e

# set the istall directory
export WORKDIR=$HOME/cpm
export TMPDIR=/tmp/opt/cpm
export ARCHIVE=/tmp/cpm.1.0.0-linux-amd64.tar.gz

# verify running as root user

createArchive () {
	mkdir -p $TMPDIR/bin

	cp $WORKDIR/sbin/* $TMPDIR/bin
	cp $WORKDIR/bin/* $TMPDIR/bin
	cp $WORKDIR/sbin/basic-user-install.sh $TMPDIR
	cp $WORKDIR/sbin/bu-*.sh $TMPDIR

	mkdir -p $TMPDIR/config
	cp $WORKDIR/config/* $TMPDIR/config

	mkdir -p $TMPDIR/www
	cp -r $WORKDIR/images/cpm/www/* $TMPDIR/www/

	cd $TMPDIR

	tar cvzf $ARCHIVE .

}

pushImages () {
	# push docker images to dockerhub

	docker tag cpm crunchydata/cpm
	docker push crunchydata/cpm
	docker save crunchydata/cpm > /tmp/cpm.tar

	docker tag cpm-pgpool crunchydata/cpm-pgpool
	docker push crunchydata/cpm-pgpool
	docker save crunchydata/cpm-pgpool > /tmp/cpm-pgpool.tar

	docker tag cpm-admin crunchydata/cpm-admin
	docker push crunchydata/cpm-admin
	docker save crunchydata/cpm-admin > /tmp/cpm-admin.tar

	docker tag cpm-base crunchydata/cpm-base
	docker push crunchydata/cpm-base
	docker save crunchydata/cpm-base > /tmp/cpm-base.tar

	docker tag cpm-mon crunchydata/cpm-mon
	docker push crunchydata/cpm-mon
	docker save crunchydata/cpm-mon > /tmp/cpm-mon.tar

	docker tag cpm-backup crunchydata/cpm-backup
	docker push crunchydata/cpm-backup
	docker save crunchydata/cpm-backup > /tmp/cpm-backup.tar

	docker tag cpm-backup-job crunchydata/cpm-backup-job
	docker push crunchydata/cpm-backup-job
	docker save crunchydata/cpm-backup-job > /tmp/cpm-backup-job.tar

	docker tag cpm-node crunchydata/cpm-node
	docker push crunchydata/cpm-node
	docker save crunchydata/cpm-node > /tmp/cpm-node.tar

	docker tag cpm-dashboard crunchydata/cpm-dashboard
	docker push crunchydata/cpm-dashboard
	docker save crunchydata/cpm-dashboard > /tmp/cpm-dashboard.tar
}

createArchive
#pushImages
