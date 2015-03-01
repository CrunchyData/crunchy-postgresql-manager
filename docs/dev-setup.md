Developer Setup
=================

Typical Development Environment
-------------------------------
CentOS 7 and RHEL 7 are supported currently, others might work, especially
Fedora or other RHEL variants but you might see differences in
installation.

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

Developer Installation vs User Installation
--------------------------------------------

There is a developer installation which is meant for people
wanting to modify or work with the CPM source code.  There is
also a user installation that has the binaries and Docker images
all ready to pull down and execute.  In this document we will
discuss the developer installation.  See User-Install.md for
details on the user installation.



Run Install.sh
-----------------
The first step of course is to clone the CPM project from
github.  In the root directory is a file named install.sh.
Running install.sh is the first step you would take.  This
script does the following:
* adds your userid to the docker group
* sets your GOPATH to the CPM directory
* creates the /opt/cpm target directory
* runs 'make'
* make compiles the golang source code
* make performs the Docker builds for the CPM containers
* copies the binaries and scripts to the target locations on the local server
* sets up the local postgres data directory
* enables in systemd the cpmagent service


Running CPM
--------------
After building and deploying CPM, you start CPM up by running the
run-cpm.sh script located in the CPM root directory.  This script
will start several Docker containers that make up CPM.

The CPM web interface is located at:
https://cpm.crunchy.lab:13000

The CPM REST API is located at:
https://cpm-admin.crunchy.lab:13000

The CPM influxdb web interface is located at:
http://cpm-mon.crunchy.lab:8083
