Developer Setup
=================

Development Environment
=======================
Here are the steps required to set up a CPM development environment, CPM is 
built using Centos 7.1  and Docker 1.5

This instruction assumes you are using a static IP address of
192.168.56.103 for your CPM server.

### Setup Go Project Structure ###
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
yum -y install golang git docker  mercurial sysstat
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

### Disable Firewalld ###
~~~~~~~~~~~~~~~~~~~~~~~~
systemctl disable firewalld.service
systemctl stop firewalld.service
~~~~~~~~~~~~~~~~~~~~~~~~

There is a starter script in CPM /sbin directory
that attempts to open up only the required CPM ports, adjust
this for your local use if you require firewalld.

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
git clone git@github.com:/crunchydata/skybridge
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
You can run CPM by running the following script:
~~~~~~~~~~~~~~~~~~~~~~~~~~
sudo ./run-cpm.sh
~~~~~~~~~~~~~~~~~~~~~~~~~~
This script will start several Docker containers that make up CPM.  You will
need to edit the run-cpm.sh script to specify your IP address of your
server as well as your CPM installation directory path.  You can
also adjust or remove the local port bindings if you want.

On the dev host, the following URLs are useful:

CPM Web User Interface
----------------------
http://cpm.crunchy.lab:13001
http://192.168.56.103:13001

CPM Admin API
----------------------
http://cpm-admin.crunchy.lab:13001
http://192.168.56.103:14001

Prometheus Dashboard
----------------------
http://cpm-promdash.crunchy.lab:3000
http://192.168.56.103:15000

If you are running CPM on a VM (host-only) and
accessing CPM from the VM host (not the guest), then
you will need to edit the dashboard server
configuration via the PromDash user interface
and specify the prometheus server URL
as http://192.168.56.103:16000.

Prometheus DB
----------------------
http://cpm-prometheus.crunchy.lab:9090
http://192.168.56.103:16000

If you are running the CPM user interface from outside the dev host
(e.g.  from your vbox host browser), you will need to update
a couple of javascript files with the promdash URL.  By default
these are specified in the javascript as cpm-promdash:3000, this will
not be accessible from your vbox host unless you specify the 
skybridge DNS server.

The js files to change are:
servers/servers.js
projects/container-logic.js

Look for occurances of cpm-promdash:3000 and change them to
the static IP address and ports listed above.

Login
--------

Browse to the CPM web user interface
user id is cpm
password is cpm
Admin URL is either http://cpm-admin:13001 (on your CPM host)
or http://192.168.56.103:13001

Initially you will need to first define your CPM server which
is your CPM host (e.g. 192.168.56.103)

Then you will be ready to start creating PostgreSQL instances.

