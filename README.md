Crunchy Postgresql Manager (Beta v0.9.4)
==========================

Crunchy Postgresql Manager (CPM) is a Docker-based solution which
provides an on-premise PostgreSQL-as-a-Service platform.

CPM allows for the quick provisioning of PostgreSQL databases
and streaming replication clusters.  

CPM also allows you to monitor and administer PostgreSQL
databases.  Currently CPM only works with databases that have
been provisioned by CPM.

![CPM Web UI](./docs/cpm.png)

A user guide is available at:
https://s3.amazonaws.com/crunchydata/cpm/cpm-user-guide.pdf

Installation
------------

There are 2 installs of CPM available, a user install and a developer
install.

The user install allows you to get CPM up and running quickly by
downloading pre-built binaries and Docker images.

The user installation archive can be downloaded from:
[https://s3.amazonaws.com/crunchydata/cpm/cpm.0.9.4-linux-amd64.tar.gz](https://s3.amazonaws.com/crunchydata/cpm/cpm.0.9.4-linux-amd64.tar.gz)

See docs/user-install.md for details on the user installation 
requirements.

For performing a user install, see the [docs/user-install.md](docs/user-install.md)
documentation.

The developer install is more difficult but allows you to build, 
configure, and develop new CPM functionality to suit your needs.

The developer install and setup is documented in [docs/dev-setup.md](docs/dev-setup.md)

Pre-requisite Installation
============

Install the skybridge or dnsbridge program before installing CPM.

skybridge is the preferred DNS bridge solution and is found
at the following location: 
[https://github.com/CrunchyData/skybridge](https://github.com/CrunchyData/skybridge)

dnsbridge is a similar solution to skybridge except that it
supports the BIND DNS server.  dnsbridge information can be found
at: [https://github.com/CrunchyData/dnsbridge](https://github.com/CrunchyData/dnsbridge)



Developer Build Requirements
==================

Build and install CPM from source by running the install.sh script

Building CPM requires development tools like the GCC compiler, along with
the Go language.  On Fedora and RedHat Linux distributions, those packages
can be installed as root like this:
~~~~~~~~~~~~~~~~~~~~~~~~~
yum install -y gcc
yum install -y golang
~~~~~~~~~~~~~~~~~~~~~~~~~

CPM also requires the Docker program is installed, running, and will stay
running after a restart:

~~~~~~~~~~~~~~~~~~~~~~~~~
yum install -y docker-io
systemctl start docker
systemctl enable docker
~~~~~~~~~~~~~~~~~~~~~~~~~

The user who is building will need to be part of the docker group
to issue docker comments.  Run this command as root, substituting
build userid in for the one at the end of the line:

~~~~~~~~~~~~~~~~~~~~~~~~~
usermod -a -G docker userid
~~~~~~~~~~~~~~~~~~~~~~~~~

You will need to logout and login again as that user for this
change to be useful.

You can confirm that Docker is available to the user you're building as
by running its info command:

~~~~~~~~~~~~~~~~~~~~~~~~~
docker info
~~~~~~~~~~~~~~~~~~~~~~~~~


Running CPM
===========

After a build, run the various CPM containers by running the following
script:

~~~~~~~~~~~~~~~~~~~~~~~~~
run-cpm.sh
~~~~~~~~~~~~~~~~~~~~~~~~~

This should start the the following containers:

* cpm - cpm.crunchy.lab - the nginx server that hosts the CPM
   	      web app, https://cpm.crunchy.lab:13000

* cpm-admin - cpm-admin.crunchy.lab - the REST API for CPM, https://cpm-admin.crunchy.lab:13000

* cpm-backup - cpm-backup.crunchy.lab - the backup process used by CPM to schedule and run backup jobs

* cpm-mon - cpm-mon.crunchy.lab - the monitoring process used by CPM to collect metrics, cpm-mon hosts the Influxdb which is used to store metrics collected by CPM, the Influxdb web console is located at http://cpm-mon.crunchy.lab:8083

* cpm-dashboard - dashboard.crunchy.lab - the Grafana dashboard that can be used to view/query collected CPM metrics - this is an optional container

Testing the Install
===========

After starting the CPM containers, you should be able to ping
each one of them and have the DNS name resolve.

You can view the running containers by issuing the following command:

~~~~~~~~~~~~~~~~~~~~~~~~~
docker ps
~~~~~~~~~~~~~~~~~~~~~~~~~

The Beta is built to use https as it's transport protocols, it
includes self-signed certificates.  

To work with the self-signed certificates in your browser, first
access https://cpm-admin.crunchy.lab:13000 and accept the 
certificate.

Then, browse to https://cpm.crunchy.lab:13000 to get started.  You will
need to accept the browser's warnings regarding the untrusted certificates
being used.

Log into the application using and ID of 'cpm' and password of 'cpm'.
Also enter into the Admin URL field the value of:
https://cpm-admin.crunchy.lab:13000

If the log in is successful, you are ready to start working with CPM.

Shutting Down CPM
===========

To shut down CPM, run the following commands:

~~~~~~~~~~~~~~~~~~~~~~~~~
docker stop cpm
docker stop cluster-backup
docker stop cluster-mon
docker stop cluster-admin
~~~~~~~~~~~~~~~~~~~~~~~~~
	

To start CPM, run the following commands:

~~~~~~~~~~~~~~~~~~~~~~~~~
docker start cpm
docker start cluster-backup
docker start cluster-mon
docker start cluster-admin
~~~~~~~~~~~~~~~~~~~~~~~~~
	
