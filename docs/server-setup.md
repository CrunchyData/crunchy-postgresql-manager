Server Setup
=================

Single-Server Configuration
---------------------------

The typical CPM developer installation will run CPM on a single
server, both CPM and it's requisite DNS bridge (dnsbridge or skybridge).
This is the single-server configuration that most people will
use to demo, hack, and prototype CPM with.

Multi-Server Configuration
==========================

A more typical example of a real CPM production deployment is
to utilize multiple servers.  In this configuration, CPM
will let you configure Postgres clusters that span more than
one server.

Here is an example of a multi-host scenario of 3 servers:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
server1.crunchy.lab - 192.168.56.101 (static IP)
server2.crunchy.lab - 192.168.56.102 (static IP)
admin.crunchy.lab   - 192.168.56.103 (static IP)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

These can be physical or virtual servers.  All three servers are available to run the CPM postgres and pgpool containers, however only one server will act
as the CPM Admin server (admin.crunchy.lab).

CPM Installation
----------------
You will perform a user-install of CPM, the installer will ask
you which remote servers you want to configure and will copy
the required CPM files to the remote server.  The installer will
also enable and start the required CPM services on each server.

DNS Configuration for Multi-Server CPM
-----------------
On each server in a multi-server configuration, you will need
to specify in your /etc/resolv.conf the CPM DNS server you have
deployed.  In this example, we have chosen to run the CPM skybridge
DNS server on the admin.crunchy.lab server.

So each server would need to specify it's primary DNS nameserver
to be 192.168.56.101.

Also, on each server, the Docker configuration in /etc/sysconfig/docker
would also need to specify the CPM DNS nameserver as follows:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
--dns=192.168.56.101 --dns=192.168.0.1
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Networking for Multi-Server CPM
-------------

For a multi-host CPM deployment, you will define network routes
as follows:

You need to define routes to the other docker servers, we create
a route file in /etc/sysconfig/network-scripts.  The name of the
file is route-xxxx where xxx is the name of the ethernet adapter for our 192.168.56.X network,
for Example:  
	route-enp0s3

For server1, the route values are:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
172.19.0.0/16 via 192.168.56.102 metric 0
172.17.0.0/16 via 192.168.56.103 metric 0
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
For admin server, the route values are:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
172.18.0.0/16 via 192.168.56.101 metric 0
172.19.0.0/16 via 192.168.56.102 metric 0
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
For server2, the route values are:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
172.18.0.0/16 via 192.168.56.101 metric 0
172.17.0.0/16 via 192.168.56.103 metric 0
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

After these routes are in place, each docker container on each
server can route to containers on the other servers.

Route examples on virtualbox:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
ip route add 172.17.0.0/16 via 192.168.56.103 dev vboxnet0
ip route add 172.18.0.0/16 via 192.168.56.101 dev vboxnet0
ip route add 172.19.0.0/16 via 192.168.56.102 dev vboxnet0
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Testing
---------------------
A script is provided to help verify that your environment is
configured correctly:
http://github.com/crunchydata/cpm/network-test.sh
