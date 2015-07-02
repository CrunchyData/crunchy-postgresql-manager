Developer Setup
=================

Development Environment
=======================
Here are the steps required to set up a CPM development environment, CPM is 
built using Centos 7.1  and Docker 1.5

### Setup Go Project Structure ###
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
yum -y install golang mercurial
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
each server by the sbin/install-cpmserverapi.sh script.  Modify this script
if you are going to configure multiple CPM hosts.  The script is currently
setup for a single CPM host installation.
~~~~~~~~~~~~~~~~~~~~~~~~
./sbin/install-cpmserverapi.sh
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

The DNS installation will enable and configure the Docker service
to specify the DNS server as the primary DNS nameserver.  This
DNS server will also be your primary nameserver in your /etc/resolv.conf
configuration.  You can install the Docker version of skybridge
as follows, you will need to EDIT the run-skybridge.sh file
to specify your IP address it will listen on:

~~~~~~~~~~~~~~~~~
git clone git@github.com/crunchydata/skybridge
cd skybridge/bin
./run-skybridge.sh
~~~~~~~~~~~~~~~~~

For Docker to use the new DNS nameserver, you will need to modify
the docker config file /etc/sysconfig/docker.  Add lines in it
like this:
~~~~~~~~~~~~~~~~~
OPTIONS='--selinux-enabled --bip=172.17.42.1/16 --dns-search=crunchy.lab --dns=192.168.0.106 --dns=192.168.0.1'
~~~~~~~~~~~~~~~~~
This example shows that skybridge is running on 192.168.0.106, I am using
a domain of crunchy.lab, and that my secondary nameserver (from my ISP)
is 192.168.0.1.  This configuration will have all the containers
in CPM trying to use the skybridge DNS nameserver as the primary
nameserver which is required by CPM.

Your /etc/resolv.conf should look similar to this if your network
configuration is set up correctly:
~~~~~~~~~~~~~~~~~
search crunchy.lab
nameserver 192.168.0.106
nameserver 192.168.0.1
~~~~~~~~~~~~~~~~~

This will cause the skybridge DNS nameserver to be queried first.


Running CPM
--------------
After building and deploying CPM, you start CPM up by running the
run-cpm.sh script located in the CPM root directory.  This script
will start several Docker containers that make up CPM.  You will
need to edit the run-cpm.sh script to specify your IP address you
want to bind ports to, or remove the port bindings if you want
to run CPM just locally on your dev instance.

The CPM web interface is located at:
http://cpm.crunchy.lab:13001

Log in with a user id of 'cpm' and a password  of 'cpm', the URL
of the CPM admin API is http://cpm-admin:13001

The CPM REST API is located at:
http://cpm-admin.crunchy.lab:13001

