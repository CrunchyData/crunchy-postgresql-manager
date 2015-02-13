Crunchy Postgresql Manager
==========================

Crunchy Postgresql Manager (CPM) is a Docker-based solution that
provides an on-premise PostgreSQL-as-a-Service platform.

CPM allows for the quick provisioning of Postgresql databases
and streaming replication clusters.

CPM also allows you to monitor and administer provisioned Postgres
databases.

Installation
------------

There are 2 installs of CPM available, a user install and a developer
install.

The user install allows you to get CPM up and running quickly by
downloading pre-built binaries and Docker images.

For performing a user install, see the docs/user-install.md 
documentation.

The developer install is more difficult but allows you to build, 
configure, and develop new CPM functionality to suit your needs.

Pre-requisite Installation
============

Install the skybridge or dnsbridge program before installing this one.

skybridge is the preferred DNS bridge solution and is found
at the following location: https://github.com/CrunchyData/skybridge



Developer Build Requirements
==================

Build and install CPM from source by running the install.sh script

Building CPM requires development tools like the GCC compiler, along with
the Go language.  On Fedora and RedHat Linux distributions, those packages
can be installed as root like this:

    yum install -y gcc
    yum install -y golang

CPM also requires the Docker program is installed, running, and will stay
running after a restart:

    yum install -y docker-io
    systemctl start docker
    systemctl enable docker

The user who is building will need to be part of the docker group
to issue docker comments.  Run this command as root, substituting
build userid in for the one at the end of the line:

    usermod -a -G docker userid

You will need to logout and login again as that user for this
change to be useful.

You can confirm that Docker is available to the user you're building as
by running its info command:

    docker info


Running CPM
===========

After a build, run the various CPM containers by running the following
script:

	run-cpm.sh

This should start the the following containers:

	cpm - cpm.crunchy.lab - the nginx server that hosts the CPM
   	      web app, http://cpm.crunchy.lab:10000

	cpm-admin - cluster-admin.crunchy.lab - the REST API
	      for CPM, http://cluster-admin.crunchy.lab:8080

	cpm-backup - cluster-backup.crunchy.lab - the backup process
	      used by CPM to schedule and run backup jobs

	cpm-mon - cluster-mon.crunchy.lab - the monitoring process
	      used by CPM to collect metrics

	dashboard - dashboard.crunchy.lab - the Grafana dashboard that
	      can be used to view/query collected CPM metrics - this is 
	      an optional container

Testing the Install
===========

After starting the CPM containers, you should be able to ping
each one of them and have the DNS name resolve.

You can view the running containers by issuing the following command:

	docker ps

Browse to http://cpm.crunchy.lab:10000 to get started.

You will see several alerts in your browser the first time you use
CPM, it will ask you to enter the Admin Service URL on the Settings
Page, enter http://cluster-admin.crunchy.lab:8080 and press the Save
button.

You will then need to log into CPM, use the default values of cpm
for the user id, and cpm for the password.

If the log in is successful, you are ready to start working with CPM.

Shutting Down CPM
===========

To shut down CPM, run the following commands:

	docker stop cpm
	docker stop cluster-backup
	docker stop cluster-mon
	docker stop cluster-admin
	

To start CPM, run the following commands:

	docker start cpm
	docker start cluster-backup
	docker start cluster-mon
	docker start cluster-admin
	
