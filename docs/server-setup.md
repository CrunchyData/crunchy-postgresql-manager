Server Images
=================

There are 3 servers in the POC (proof-of-concept).

	server1.crunchy.lab - 192.168.56.101 (static IP)
	server2.crunchy.lab - 192.168.56.102 (static IP)
	admin.crunchy.lab   - 192.168.56.103 (static IP)

Each of these servers is running on VirtualBox, the POC has been
tested with both Fedora 20, CentOS7, and RHEL7 as the server hsots.

Here are steps to follow when building a CPM environment from
scratch:

on rhel 7, register your systems:

	subscription-manager repos --enable=rhel-7-server-extras-rpms --enable=rhel-7-server-optional-rpms

install rhel7/centos7 minimal

	yum -y install docker sysstat bind-utils wget git net-tools golang bind gcc docker-registry kernel-devel kernel-headers


network
-------
** configure 2 network adapters: 
Create the VM hosts using host-only vboxnet0 bridge that is defined
to use 192.168.56.1 as the network address space
This allows them to have static IP addresses and route traffic to each other.

Also, add another ethernet adapter which will connect to the br0 bridge
that is defined on the VirtualBox host.  This will allow each VM to
also connect to the public internet.

adjust the host-only adapter with the following settings to set the IP address, domain, and DNS servers:

	BOOTPROTO=none
	ONBOOT=yes
	IPADDR0=192.168.56.103
	DOMAIN=crunchy.lab
	DNS1=192.168.56.103
	DNS2=192.168.0.1


adjust the bridged adapter with the following to avoid overwriting
of /etc/resolv.conf:

	PEERDNS=no

** hostname
--------
as root, run the 'nmtui' utility to change the hostname
** add the hostname to the /etc/hosts file


** minimal desktop
---------
I find that I like to have the Gnome desktop installed on
my admin server when I work on my laptop, to install the Gnome
desktop, perhform the following:

	yum groupinstall 'Server with GUI'
	systemctl enable graphical.target --force

To install chrome:

	cat << EOF > /etc/yum.repos.d/google-chrome.repo
	[google-chrome]
	name=google-chrome - \$basearch
	baseurl=http://dl.google.com/linux/chrome/rpm/stable/\$basearch
	enabled=1
	gpgcheck=1
	gpgkey=https://dl-ssl.google.com/linux/linux_signing_key.pub
	EOF

	yum install google-chrome-stable

PostgreSQL is required on the CPM servers, to install Postgres:

On Centos7:
	rpm -ivh http://yum.postgresql.org/9.3/redhat/rhel-7-x86_64/pgdg-centos93-9.3-1.noarch.rpm

On Rhel7:
	rpm -ivh http://yum.postgresql.org/9.3/redhat/rhel-7-x86_64/pgdg-redhat93-9.3-1.noarch.rpm

Then:
	yum -y install postgresql93-server postgresql93 postgresql93-libs postgresql93-contrib


** firewalld
---------
you might want to develop with the firewall turned off, if so, here
are the commands to shut off the firewall:

	systemctl stop firewalld.service
	systemctl disable firewalld.service


** network route
-------------

For a single CPM server dev environment you can ignore this section, however
for a multi-host CPM deployment, you will define network routes
as follows:

we need to define routes to the other docker servers, we create
a route file in /etc/sysconfig/network-scripts.  The name of the
file is route-xxxx where xxx is the name of the ethernet adapter for our 192.168.56.X network,
for Example:  
	route-enp0s3

For server1, the route values are:
	172.19.0.0/16 via 192.168.56.102 metric 0
	172.17.0.0/16 via 192.168.56.103 metric 0

For admin server, the route values are:
	172.18.0.0/16 via 192.168.56.101 metric 0
	172.19.0.0/16 via 192.168.56.102 metric 0

For server2, the route values are:
	172.18.0.0/16 via 192.168.56.101 metric 0
	172.17.0.0/16 via 192.168.56.103 metric 0

After these routes are in place, each docker container on each
server can route to containers on the other servers in this POC.

Testing
========
A script is provided to help verify that your environment is
configured correctly:
	http://github.com/crunchyds/docker-pg-cluster/network-test.sh


** copy server config files and binaries  to /cluster/bin and /usr/lib/systemd/system:
if you have not already cloned the repo:

	git clone git@github.com:crunchyds/docker-pg-cluster.git

Next, run the script to copy required files to new server:

	sudo ./docker-pg-cluster/deploy-server-config.sh

The DNS zone files are located here:
	http://github.com/docker-pg-cluster/pgcluster/config/zonefiles

** on all servers, start the following services:

	systemctl daemon-reload
	systemctl start dnsbridgeclient.service
	systemctl start cpmagent.service

only on the admin host:

	systemctl start named.service
	systemctl start cpmagent.service
	systemctl start dnsbridgeserver.service


