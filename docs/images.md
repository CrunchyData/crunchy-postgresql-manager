Docker Images
=================

The Crunchy Docker images are based on either RHEL 7 or Centos 7, 
you will need to install the base OS images before building
the Crunchy images.

Docker Registry
---------------
Soon, I'll be switching over to using the Redhat Docker Registry
to store the Crunchy Docker images, the steps listed below are
the manual steps taken at the moment.

Installing the CentOS7 Docker Image
-------
download image from:
https://github.com/CentOS/sig-cloud-instance-images/tree/CentOS7

then install the image using this:
cat CentOS-7-x86_64-docker.img.tar.xz  | docker import - crunchy/centos7

Installing the RHEL 7 Docker Image
------------
Download the RHEL 7 image from the following link:

Install the image as follows:



crunchy-node
-------
postgres containers run in the 'crunchy-node' docker image

This image is built by the Dockerfile located here:
	http://github.com/docker-pg-cluster/images/crunchy-node/Dockerfile.rhel7

Build the images as follows:
	cd docker-pg-cluster/images/crunchy-admin
	docker build -t crunchy-node .

A typical development environment will be used to build the image
and then copy the image to a target POC environment as follows:

	docker save crunchy-node > /tmp/crunchy-node.tar
	scp /tmp/crunchy-node.tar server1.crunchy.lab:
	scp /tmp/crunchy-node.tar server2.crunchy.lab:

Then on the server1/server2/admin machines, you load the image
back into the local docker repository as follows:

	docker load -i crunchy-node.tar 

crunchy-admin
-------------

This image is used only for the cluster 'admin' container which runs
the admin REST API.  This image can be built with the following
Dockerfile:
	http://github.com/docker-pg-cluster/images/crunchy-admmin/Dockerfile.rhel7

Build the image as follows:
	cd docker-pg-cluster/images/crunchy-admin
	docker build -t crunchy-admin .

This image is only loaded onto the admin.crunchy.lab POC server.


crunchy-cpm
-----------

This image is used to implement the CPM web application, it includes
the nginx http server and mounts /www which contains the CPM web
app files.

Currently this image is stored on the POC's admin server.

Build the image as follows:
	cd docker-pg-cluster/images/crunchy-cpm
	docker build -t crunchy-cpm .

crunchy-pgpool
------------

This iimage is used to implement the pgpool application.  It includes
the pgpool binary.  Pgpool is used to provide a load-balanced
connection to the PG clusters.  A Pgpool container is created
and then added to a cluster.  When the cluster is configured, pgpool
is configured to serve pg requests to the cluster nodes.

This image is stored on the POC's Docker servers (server1 and server2).

Build the image as follows:
	cd docker-pg-cluster/images/crunchy-pgpool
	docker build -t crunchy-pgpool .
