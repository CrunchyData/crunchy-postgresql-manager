Developer Setup
=================

Development Environment
=======================
Here are the steps required to set up a CPM development environment, CPM is 
built using Centos 7.1  and Docker 1.5

### Dependency Setup ###
Set up your developer account with sudo privs, some of the scripts
and docker commands will require you run them as superuser.

~~~~~~~~~~~~~~~~~~~~~~~~~~~~
git clone git@github.com:crunchydata/crunchy-postgresql-manager
cd crunchy-postgresql-manager
cd sbin
./basic-user-install.sh
~~~~~~~~~~~~~~~~~~~~~~~~~~~~


### Setup Go Project Structure ###
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
mkdir -p devproject/src devproject/bin devproject/pkg

export GOPATH=~/devproject
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

### Download and Install godep ###
~~~~~~~~~~~~~~~~~~~~~~~~
cd devproject
go get github.com/tools/godep
~~~~~~~~~~~~~~~~~~~~~~~~

### Download CPM Source ###
~~~~~~~~~~~~~~~~~~~~~~~~
go get github.com/crunchydata/crunchy-postgresql-manager
cd src/github.com/crunchydata/crunchy-postgresql-manager
~~~~~~~~~~~~~~~~~~~~~~~~

### Download and Restore All Dependencies ###
~~~~~~~~~~~~~~~~~~~~~~~~
godep restore
~~~~~~~~~~~~~~~~~~~~~~~~

### Compile CPM ###
~~~~~~~~~~~~~~~~~~~~~~~~
make build
~~~~~~~~~~~~~~~~~~~~~~~~

### Build CPM Docker Images ###
~~~~~~~~~~~~~~~~~~~~~~~~
make buildimages
~~~~~~~~~~~~~~~~~~~~~~~~

Start CPM Server Agent
----------------------
After you have successfully compiled CPM and built the CPM Docker images,
on each server that is to run CPM, you will need to start a CPM Server
Agent, this is started using systemd.  CPM server files are copied to
each server by the sbin/dev-setup.sh script.  Modify this script
if you are going to configure multiple CPM hosts.  The script is currently
setup for a single CPM host installation.
~~~~~~~~~~~~~~~~~~~~~~~~
./sbin/dev-setup.sh
~~~~~~~~~~~~~~~~~~~~~~~~

Network Configuration
------------------------------
CPM requires a reliable IP address of the host on which it is running.
When a VM is created to develop CPM upon, you would create an extra
Ethernet adapter typically so that you can assign it a static IP
address.  In Virtualbox, this adapter would be a Host-Only adapter
for example.

For development, you will likely want to turn off your firewall.  A
script for setting firewall rules is offered if that is not an option
for you.  See firewall.sh for the rules which open the required
CPM ports.

DNS Configuration
------------------------------
CPM uses a dedicated DNS server to orchestrate the Docker container-to-DNS
mapping.  This DNS server is created by installing the dnsbridge or skybridge
projects on your CPM development server.  The simpler installation
of DNS server is the skybridge project.  These projects are found
in github at the following locations:

https://github.com/CrunchyData/skybridge

https://github.com/CrunchyData/dnsbridge

The DNS installation will enable and configure the Docker service
to specify the DNS server as the primary DNS nameserver.  This
DNS server will also be your primary nameserver in your /etc/resolv.conf
configuration.

Running CPM
--------------
After building and deploying CPM, you start CPM up by running the
run-cpm.sh script located in the CPM root directory.  This script
will start several Docker containers that make up CPM.

The CPM web interface is located at:
https://cpm.crunchy.lab:13000

The CPM REST API is located at:
https://cpm-admin.crunchy.lab:13000

