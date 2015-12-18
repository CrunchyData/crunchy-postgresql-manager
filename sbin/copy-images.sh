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

# This install script saves the CPM docker images locally and then
# copies them to a remote server, and finally loads them on the remote
# server
#

# Exit installation on any unexpected error
set -e

# set the istall directory
export REMOTE=bean
export TMPDIR=/tmp

loadImages () {

	echo "loading cpm images on " $REMOTE
	ssh root@$REMOTE 'docker load -i /tmp/cpm.tar; \
	docker load -i /tmp/cpm-pgpool.tar; \
	docker load -i /tmp/cpm-admin.tar; \
	docker load -i /tmp/cpm-restore-job.tar; \
	docker load -i /tmp/cpm-collect.tar; \
	docker load -i /tmp/cpm-task.tar; \
	docker load -i /tmp/cpm-backup-job.tar; \
	docker load -i /tmp/cpm-node.tar; \
	docker load -i /tmp/cpm-node-proxy.tar; \
	docker load -i /tmp/cpm-efk.tar'

}
copyImages () {

	echo "copying cpm images to " $REMOTE
	scp $TMPDIR/cpm.tar \
		$TMPDIR/cpm-pgpool.tar \
		$TMPDIR/cpm-admin.tar \
		$TMPDIR/cpm-restore-job.tar \
		$TMPDIR/cpm-collect.tar \
		$TMPDIR/cpm-task.tar  \
		$TMPDIR/cpm-backup-job.tar  \
		$TMPDIR/cpm-node.tar  \
		$TMPDIR/cpm-node-proxy.tar  \
		$TMPDIR/cpm-efk.tar \
		$REMOTE:$TMPDIR

}

saveImages () {

	echo "saving cpm image"
	sudo docker save crunchydata/cpm:latest > $TMPDIR/cpm.tar

	echo "saving cpm-pgpool image"
	sudo docker save crunchydata/cpm-pgpool:latest > $TMPDIR/cpm-pgpool.tar

	echo "saving cpm-admin image"
	sudo docker save crunchydata/cpm-admin:latest > $TMPDIR/cpm-admin.tar

	echo "saving cpm-base image"
	sudo docker save crunchydata/cpm-restore-job:latest > $TMPDIR/cpm-restore-job.tar

	echo "saving cpm-collect image"
	sudo docker save crunchydata/cpm-collect:latest > $TMPDIR/cpm-collect.tar

	echo "saving cpm-backup image"
	sudo docker save crunchydata/cpm-task:latest > $TMPDIR/cpm-task.tar

	echo "saving cpm-backup-job image"
	sudo docker save crunchydata/cpm-backup-job:latest > $TMPDIR/cpm-backup-job.tar

	echo "saving cpm-node image"
	sudo docker save crunchydata/cpm-node:latest > $TMPDIR/cpm-node.tar

	echo "saving cpm-node-proxy image"
	sudo docker save crunchydata/cpm-node-proxy:latest > $TMPDIR/cpm-node-proxy.tar

	echo "saving cpm-efk image"
	sudo docker save crunchydata/cpm-efk:latest > $TMPDIR/cpm-efk.tar

}

saveImages
copyImages
loadImages
