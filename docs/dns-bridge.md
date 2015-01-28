
DNS Bridge Configuration
==================

dnsbridgeclient Configuration
------
dnsbridgeclient listens to the Docker Event API on each docker host
for start and stop events.  It sends add-host and delete-host
messages to the dnsbridgeserver when Docker containers are started
and stopped.

dnsbridgeclient needs to be installed on each Docker host.  There
is a dnsbridgeclient.service systemd script that causes dnsbridgeclient
to be started and stopped upon a reboot of the server.

dnsbridgeclient is necessary because postgres cluster logic depends
upon stable host names or IP addresses to form the cluster.  Since
Docker assigns IP addresses dynamically you can not predict what
IP address your container will be assigned.  Therefore, we
instead register IP addresses to well-known host names that
we can reference in the postgresql config files.

All Docker servers will need to have dnsbridgeclient installed in 

	/cluster/bin/dnsbridgeclient.

Register the dnsbridgeclient service in systemd as follows:

	/usr/lib/systemd/system/dnsbridgeclient.service

	systemctl enable dnsbridgeclient

	systemctl start dnsbridgeclient

dnsbridgeclient can be invoked from the command line as follows:

	/cluster/bin/dnsbridgeclient -d 192.168.56.103:14000

dnsbridgeserver Configuration
------

On the DNS server host, we need to run the dnsbridgeserver program, 
this is installed in /cluster/bin/dnsbridgeserver.  dnsbridgeserver is 
used to receive dnsbridgeclient messages
from each Docker server as containers are started and stopped.

dnsbridgeserver is started using the dnsbridgeserver.service configuration
file located in 

	/usr/lib/systemd/system/dnsbridgeserver.configuration.

dnsbridgeserver is started with the following parameters:

	/cluster/bin/dnsbridgeserver -p 14000

There are scripts that are invoked by the server:

	/cluster/bin/add-host.sh 

	/cluster/bin/delete-host.sh

These scripts are invoked when the dnsbridgeservice (running on the DNS admin server) receives an add or delete event from one of the remote servers (those that are creating docker containers).

These scripts invoke the DNS nsupdate command to dynamically
alter the BIND/named server configuration (zone records).

Register the dnsbridgeserver service in systemd using the following
file:

	/usr/lib/systemd/system/dnsbridgeserver.service

	systemctl enable dnsbridgeserver

	systemctl start dnsbridgeserver
