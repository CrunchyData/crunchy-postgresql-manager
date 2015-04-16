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
export VERSION=0.9.3
export WORKDIR=$GOPATH/src/github.com/crunchydata/crunchy-postgresql-manager
export TMPDIR=/tmp/var/cpm
export ARCHIVE=/tmp/cpm.$VERSION-linux-amd64.tar.gz

# verify running as root user

createArchive () {
	mkdir -p $TMPDIR/bin

	cp $WORKDIR/sbin/* $TMPDIR/bin
	cp $GOPATH/bin/* $TMPDIR/bin
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

	echo "saving cpm image"
	sudo docker tag cpm crunchydata/cpm:$VERSION
	sudo docker tag cpm crunchydata/cpm:latest
	sudo docker push crunchydata/cpm:$VERSION
	sudo docker push crunchydata/cpm:latest
#	docker save crunchydata/cpm > /tmp/cpm.tar

	echo "saving cpm-pgpool image"
	sudo docker tag cpm-pgpool crunchydata/cpm-pgpool:$VERSION
	sudo docker tag cpm-pgpool crunchydata/cpm-pgpool:latest
	sudo docker push crunchydata/cpm-pgpool:$VERSION
	sudo docker push crunchydata/cpm-pgpool:latest
#	docker save crunchydata/cpm-pgpool > /tmp/cpm-pgpool.tar

	echo "saving cpm-admin image"
	sudo docker tag cpm-admin crunchydata/cpm-admin:$VERSION
	sudo docker tag cpm-admin crunchydata/cpm-admin:latest
	sudo docker push crunchydata/cpm-admin:$VERSION
	sudo docker push crunchydata/cpm-admin:latest
#	docker save crunchydata/cpm-admin > /tmp/cpm-admin.tar

	echo "saving cpm-base image"
	sudo docker tag cpm-base crunchydata/cpm-base:$VERSION
	sudo docker tag cpm-base crunchydata/cpm-base:latest
	sudo docker push crunchydata/cpm-base:$VERSION
	sudo docker push crunchydata/cpm-base:
#	docker save crunchydata/cpm-base > /tmp/cpm-base.tar

	echo "saving cpm-mon image"
	sudo docker tag cpm-mon crunchydata/cpm-mon:$VERSION
	sudo docker tag cpm-mon crunchydata/cpm-mon:latest
	sudo docker push crunchydata/cpm-mon:$VERSION
	sudo docker push crunchydata/cpm-mon:latest
#	docker save crunchydata/cpm-mon > /tmp/cpm-mon.tar

	echo "saving cpm-backup image"
	sudo docker tag cpm-backup crunchydata/cpm-backup:$VERSION
	sudo docker tag cpm-backup crunchydata/cpm-backup:latest
	sudo docker push crunchydata/cpm-backup:$VERSION
	sudo docker push crunchydata/cpm-backup:latest
#	docker save crunchydata/cpm-backup > /tmp/cpm-backup.tar

	echo "saving cpm-backup-job image"
	sudo docker tag cpm-backup-job crunchydata/cpm-backup-job:$VERSION
	sudo docker tag cpm-backup-job crunchydata/cpm-backup-job:latest
	sudo docker push crunchydata/cpm-backup-job:$VERSION
	sudo docker push crunchydata/cpm-backup-job:latest
#	docker save crunchydata/cpm-backup-job > /tmp/cpm-backup-job.tar

	echo "saving cpm-node image"
	sudo docker tag cpm-node crunchydata/cpm-node:$VERSION
	sudo docker tag cpm-node crunchydata/cpm-node:latest
	sudo docker push crunchydata/cpm-node:$VERSION
	sudo docker push crunchydata/cpm-node:latest
#	docker save crunchydata/cpm-node > /tmp/cpm-node.tar

	echo "saving cpm-dashboard image"
	sudo docker tag cpm-dashboard crunchydata/cpm-dashboard:$VERSION
	sudo docker tag cpm-dashboard crunchydata/cpm-dashboard:latest
	sudo docker push crunchydata/cpm-dashboard:$VERSION
	sudo docker push crunchydata/cpm-dashboard:latest
#	docker save crunchydata/cpm-dashboard > /tmp/cpm-dashboard.tar
}

saveImages () {

	echo "saving cpm image"
	docker save crunchydata/cpm > /tmp/cpm.tar

	echo "saving cpm-pgpool image"
	docker save crunchydata/cpm-pgpool > /tmp/cpm-pgpool.tar

	echo "saving cpm-admin image"
	docker save crunchydata/cpm-admin > /tmp/cpm-admin.tar

	echo "saving cpm-base image"
	docker save crunchydata/cpm-base > /tmp/cpm-base.tar

	echo "saving cpm-mon image"
	docker save crunchydata/cpm-mon > /tmp/cpm-mon.tar

	echo "saving cpm-backup image"
	docker save crunchydata/cpm-backup > /tmp/cpm-backup.tar

	echo "saving cpm-backup-job image"
	docker save crunchydata/cpm-backup-job > /tmp/cpm-backup-job.tar

	echo "saving cpm-node image"
	docker save crunchydata/cpm-node > /tmp/cpm-node.tar

	echo "saving cpm-dashboard image"
	docker save crunchydata/cpm-dashboard > /tmp/cpm-dashboard.tar
}

createArchive
#saveImages
#pushImages
