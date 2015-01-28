
DNS Configuration
==================

Each docker container has a dynamic assigned (by docker) IP address
that is within a defined range that is configured when docker starts.

To support multiple Docker servers, we configure each Docker service
to use a different Docker bridge IP address range as follows:

	server1.crunchy.lab - 172.42.18.0/16
	server2.crunchy.lab - 172.42.19.0/16
	admin.crunchy.lab   - 172.42.17.0/16

DNS Requirement
---------------

There is no DNS capability with docker out of the box.  Our requirement
for postgres clustering is to have mutliple containers all route to
each other.  If the IP addresses change on container restart, we
would have to reconfigure postgresql.conf and other pg config files
to use the new IP addresses.  To avoid this, we will want to refer
to host names instead of IP addresses and have the host names resolve
to their IP address like DNS offers.

DNS Server
----------

Our solution to Docker DNS is to install a DNS server (named) that runs 
on the admin server.  NOTE: we could have split out the DNS server to another
host but saved running another VM by running DNS on the admin server.

Each Docker host will reference the DNS server with the following
/etc/resolv.conf configuration, in this example. 192.168.56.103 is
the DNS server we will be configuring for Docker:

	domain crunchy.lab
	search crunchy.lab
	nameserver 192.168.56.103
	nameserver 192.168.0.1

NOTE:  if your server has a second Ethernet adapter (DHCP), remember to 
turn off PEERDNS on the dhcp (public internet) ethernet adapter 
configuration file in /etc/sysconfig/network-scripts/ifcfg-XXX, this will
prevent DHCP from overriding your /etc/resolv.conf configuration

Bind Setup
----------------

We are going to run bind on the admin server in our POC.  The other
servers (server1 and server2) and all the containers will resolve
names to addresses using this bind server.

On admin server:

	yum install bind

	vi /etc/named.conf

After much configuration, the /etc/named.conf file looks like this:
	see http://github.com/docker-pg-cluster/...

Edit the /etc/sysconfig/named file and add this line:
	ENABLE_ZONE_WRITE=yes

The following zone files were created in /var/named/dynamic:
	see http://github.com/docker-pg-cluster/pgcluster/config/zonefiles

A script is provided to copy these zone files to a new named install:
	http://github.com/crunchyds/docker-pg-cluster/dns-setup.sh

IMPORTANT!!!!
File Permissions, I found that 770 is the correct file
permissions of the /var/named/dynamic directory and zone files, if not set
correctly, you'll get permission errors in the named log!

Open up the firewall:
	firewall-cmd --permanent --add-port=53/tcp
	firewall-cmd --permanent --add-port=53/udp
	firewall-cmd --reload

Then start up the bind server:
	systemctl start named.service

On RHEL 6.5, the command is:
	
	service named start

Now you should test out the new DNS server:
	dig ns.crunchy.lab

	;; ANSWER SECTION:

	ns.crunchy.lab.		259200	IN	A	192.168.56.103


You should be able to add a record:
	nsupdate 
	> server ns.crunchy.lab
	> zone crunchy.lab.
	> update add test1.crunchy.lab 60 A 172.19.0.3
	> send
	> quit

Test the lookup:
	dig test1.crunchy.lab

You should be able to add a PTR record:
	nsupdate 
	> zone  0.19.172.in-addr.arpa
	> update add 3.0.19.172.in-addr.arpa 3500 IN PTR test1.crunchy.lab.
	> send
	> quit

Test the reverse lookup:
	dig +x 172.19.0.3	

You should be able to delete a PTR record:
	nsupdate
	> zone  0.19.172.in-addr.arpa
	> update delete 3.0.19.172.in-addr.arpa 3500 IN PTR test1.crunchy.lab.
	> send
	> quit

You should be able to delete the A record:
	nsupdate 
	> server ns.crunchy.lab
	> zone crunchy.lab.
	> update delete test1.crunchy.lab 60 A 172.19.0.3
	> send
	> quit
	

Bootstapping the POC Servers DNS
-------------------------------
we need to add entries into our DNS for the 3 POC servers, we do this with
the by running the following commands on the admin server:
	/cluster/bin/add-host.sh 192.168.56.103 admin
	/cluster/bin/add-host.sh 192.168.56.101 server1
	/cluster/bin/add-host.sh 192.168.56.102 server2
