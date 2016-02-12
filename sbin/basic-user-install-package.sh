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
export VERSION=1.0.2
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
	sudo docker tag -f cpm crunchydata/cpm:$VERSION
	sudo docker tag -f cpm crunchydata/cpm:latest
	sudo docker push -f crunchydata/cpm:$VERSION
	sudo docker push -f crunchydata/cpm:latest

	echo "saving cpm-pgpool image"
	sudo docker tag -f cpm-pgpool crunchydata/cpm-pgpool:$VERSION
	sudo docker tag -f cpm-pgpool crunchydata/cpm-pgpool:latest
	sudo docker push -f crunchydata/cpm-pgpool:$VERSION
	sudo docker push -f crunchydata/cpm-pgpool:latest

	echo "saving cpm-admin image"
	sudo docker tag -f cpm-admin crunchydata/cpm-admin:$VERSION
	sudo docker tag -f cpm-admin crunchydata/cpm-admin:latest
	sudo docker push -f crunchydata/cpm-admin:$VERSION
	sudo docker push -f crunchydata/cpm-admin:latest

	echo "saving cpm-collect image"
	sudo docker tag -f cpm-collect crunchydata/cpm-collect:$VERSION
	sudo docker tag -f cpm-collect crunchydata/cpm-collect:latest
	sudo docker push -f crunchydata/cpm-collect:$VERSION
	sudo docker push -f crunchydata/cpm-collect:latest

	echo "saving cpm-task image"
	sudo docker tag -f cpm-task crunchydata/cpm-task:$VERSION
	sudo docker tag -f cpm-task crunchydata/cpm-task:latest
	sudo docker push -f crunchydata/cpm-task:$VERSION
	sudo docker push -f crunchydata/cpm-task:latest

	echo "saving cpm-backup-job image"
	sudo docker tag -f cpm-backup-job crunchydata/cpm-backup-job:$VERSION
	sudo docker tag -f cpm-backup-job crunchydata/cpm-backup-job:latest
	sudo docker push -f crunchydata/cpm-backup-job:$VERSION
	sudo docker push -f crunchydata/cpm-backup-job:latest

	echo "saving cpm-restore-job image"
	sudo docker tag -f cpm-restore-job crunchydata/cpm-restore-job:$VERSION
	sudo docker tag -f cpm-restore-job crunchydata/cpm-restore-job:latest
	sudo docker push -f crunchydata/cpm-restore-job:$VERSION
	sudo docker push -f crunchydata/cpm-restore-job:latest


	echo "saving cpm-node image"
	sudo docker tag -f cpm-node crunchydata/cpm-node:$VERSION
	sudo docker tag -f cpm-node crunchydata/cpm-node:latest
	sudo docker push -f crunchydata/cpm-node:$VERSION
	sudo docker push -f crunchydata/cpm-node:latest

	echo "saving cpm-node-proxy image"
	sudo docker tag -f cpm-node-proxy crunchydata/cpm-node-proxy:$VERSION
	sudo docker tag -f cpm-node-proxy crunchydata/cpm-node-proxy:latest
	sudo docker push -f crunchydata/cpm-node-proxy:$VERSION
	sudo docker push -f crunchydata/cpm-node-proxy:latest

	echo "saving cpm-efk image"
	sudo docker tag -f cpm-efk crunchydata/cpm-efk:$VERSION
	sudo docker tag -f cpm-efk crunchydata/cpm-efk:latest
	sudo docker push -f crunchydata/cpm-efk:$VERSION
	sudo docker push -f crunchydata/cpm-efk:latest
}

saveImages () {

	echo "saving cpm image"
	sudo docker save crunchydata/cpm > /tmp/cpm.tar

	echo "saving cpm-pgpool image"
	sudo docker save crunchydata/cpm-pgpool > /tmp/cpm-pgpool.tar

	echo "saving cpm-admin image"
	sudo docker save crunchydata/cpm-admin > /tmp/cpm-admin.tar

	echo "saving cpm-base image"
	sudo docker save crunchydata/cpm-restore-job > /tmp/cpm-restore-job.tar

	echo "saving cpm-collect image"
	sudo docker save crunchydata/cpm-collect > /tmp/cpm-collect.tar

	echo "saving cpm-backup image"
	sudo docker save crunchydata/cpm-task > /tmp/cpm-task.tar

	echo "saving cpm-backup-job image"
	sudo docker save crunchydata/cpm-backup-job > /tmp/cpm-backup-job.tar

	echo "saving cpm-node image"
	sudo docker save crunchydata/cpm-node > /tmp/cpm-node.tar

	echo "saving cpm-node-proxy image"
	sudo docker save crunchydata/cpm-node-proxy > /tmp/cpm-node-proxy.tar

	echo "saving cpm-efk image"
	sudo docker save crunchydata/cpm-efk > /tmp/cpm-efk.tar

}

createArchive
pushImages
