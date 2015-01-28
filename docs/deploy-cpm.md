
Deploy CPM
==========

A script is provided to help deploy CPM to the servers
as defined for the POC:
	http://github.com/crunchyds/docker-pg-cluster/deploy-server-config.sh

This script copies CPM files to all the servers from your local
github clone.

Registry
========

Right now, images are rebuilt on each server as the image changes
are made.  This is REALLY inefficient!  Therefore we will be moving
to the Redhat Docker Registry for managing our images.  With the registry,
image builds will be done on a single server.
