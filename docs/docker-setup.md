
Docker Configuration
==================

Each docker container has a dynamic assigned (by docker) IP address
that is within a defined range that is configured when docker starts.

To support multiple Docker servers, we configure each Docker service
to use a different Docker bridge IP address range as follows:

	server1.crunchy.lab - 172.42.18.0/16
	server2.crunchy.lab - 172.42.19.0/16
	admin.crunchy.lab   - 172.42.17.0/16


docker DNS configuration
------------------------

Docker is configured on each server to reference the DNS server
for each container that it starts.  This is done by altering
the Docker service configuration to include the following
Docker command line options:

	--dns=192.168.56.103

If you want the Docker containers to resolve to the public
internet, include the public internet DNS address as well:

	--dns=192.168.0.1

** docker setup
------------

modify /usr/lib/systemd/system/docker.service:
--selinux-enabled -H tcp://0.0.0.0:4243 -H unix:///var/run/docker.sock
--bip=172.1X.42.1/16 --dns=192.168.56.103 --dns=192.168.0.1

NOTE:  adding --dns=192.168.0.1 allows me to touch the public internet
from within the containers, this is something useful for development
but needs to be considered for a real situation, this will be a different 
IP address on your network

** on our POC, the bip values are as follows, each server needs a different value:
server1.crunchy.lab: --bip=172.18.42.1/16
server2.crunchy.lab: --bip=172.19.42.1/16
admin.crunchy.lab:   --bip=172.17.42.1/16

Enable and start docker:

	systemctl enable docker.service
	systemctl start docker.service


for any users that will be creating docker images (for development), run
the following command to add a user to the docker group:

	usermod -a -G docker <your-user>

** docker images
-------------
right now, I am not using a docker image registry, so to get
the images, you will need to build them using the github 
docker files:

to load the base centos7 image:

download from the crunchy repo site the CentOS 7 docker image we will be using
for the demo, then execute the following command to load the image locally:

	cat CentOS-7-x86_64-docker.img.tar.xz  | docker import - crunchy/centos7

OR for rhel7, see https://access.redhat.com/articles/881893#get:

	docker load -i rhel-server-docker-7.0-21.4.x86_64.tar.gz

	git clone git@github.com:crunchyds/docker-pg-cluster.git
	cd docker-pg-cluster/images/

Only on the admin server:
	cd ./crunchy-admin
	docker build -t crunchy-admin .
	cd ../crunchy-cpm
	docker build -t crunchy-cpm .

Only on the node servers:
	cd ./crunchy-node
	docker build -t crunchy-node .

Docker Setup
------------

Docker on each server host is configured as follows:
On the Admin server:
	edit /usr/lib/systemd/system/docker.service

	ExecStart=/usr/bin/docker -d --bip=172.17.42.1/16 --dns=192.168.56.103 --dns=192.168.0.1 --selinux-enabled -H fd://


On the server1 host:
	edit /usr/lib/systemd/system/docker.service

	ExecStart=/usr/bin/docker -d --bip=172.18.42.1/16 --dns=192.168.56.103 --dns=192.168.0.1 --selinux-enabled -H fd://

On the server2 host:
	edit /usr/lib/systemd/system/docker.service

	ExecStart=/usr/bin/docker -d --bip=172.19.42.1/16 --dns=192.168.56.103 --dns=192.168.0.1 --selinux-enabled -H fd://

Docker Configuration Explaination
--------------------

Each VM will run Docker.  Docker creates a dynamic Ethernet bridge it
uses to assign IP addressese to it's containers.
Docker by default assigns IP addresses to containers in a dynamic manner, 
starting with the default address range of 172.17.0.0/16

To avoid IP address conflicts on each host, we override the 
Docker bridge's IP address range for each docker server to be unique.  
The POC assignments are as follows:

server1.crunchy.lab - 172.18.42.1/16
server2.crunchy.lab - 172.19.42.1/16
admin.crunchy.lab   - 172.17.42.1/16

The IP address range is overridden by editing on each server
the docker startup options in /usr/lib/systemd/docker.service

Also, we need to let the admin server have the ability to connect
to the docker HTTP port on both server1 and server2 to provision
containers.  This requires docker to be configured as follows:

	-H tcp://0.0.0.0:4243 -H unix:///var/run/docker.sock

That configuration needs to be made in the /etc/sysconfig/docker file which
gets referenced by the /usr/lib/systemd/docker.service
file used to start docker.

Reference
==========

http://ispyker.blogspot.com/2014/04/accessing-docker-container-private.html
