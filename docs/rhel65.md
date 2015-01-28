
CPM RHEL 6.5 Setup
=========================

Some customer environments might require CPM to run on RHEL6.5, and
this document is provided to show how you can install CPM into
a RHEL 6.5 only environment.

Infrastructure
---------------
I typically run CPM on a 3 host configuration for simple testing
and to run proof-of-concepts upon.  This is totally flexible and
you can instead run CPM upon a single host, it is up to you and
how realistic you want to be in your deployment.  Obviously for
testing performance or multi-host CPM functions, you will need
at least 2 compute nodes.

Typical configuration is this:
	* 1 node runs CPM web app, the CPM "admin" container,  and the DNS server
	* 2 nodes run all the PostgresSQL containers which allows you to offer HA/scaling across separate servers

##### Example IP Addresses
For this document, I assume the following static IP addresses and host names:

	* rh65-admin.crunchy.lab (192.168.56.110)
	* rh65-server1.crunchy.lab (192.168.56.111)
	* rh65-server2.crunchy.lab (192.168.56.112)
	* public internet gateway (192.168.0.1)
	

##### Step 1

create your hosts using whatever VM or infrastructure you have, I typically
use VirtualBox but CPM doesn't care what the servers are built upon.

For this instruction, we'll assume to be installing RHEL 6.5.  I typically
only install the 'minimal' version of RHEL/CentOS.

Register your new hosts, on rhel6.5, you will enter:
```sh
subscription-manager register
subscription-manager refresh
yum -y update
```

At this point you should verify that you are running at least this level
of kernel:

```sh
uname -a
2.6.32-431.29.2.el6.x86_64
```


At this point in the installation process, each host needs to have a static IP 
address, and each should be able to ping each other. 
	
On each host, install the EPEL rpm and PGDG rpm which contains docker-io and
PostgreSQL 9.3.5 for RHEL 6.5:

```sh
yum -y install http://mirror.metrocast.net/fedora/epel/6/i386/epel-release-6-8.noarch.rpm
yum -y install docker-io golang bind sysstat wget net-tools bind-utils make 
rpm -ivh http://yum.postgresql.org/9.3/redhat/rhel-6.5-x86_64/pgdg-redhat93-9.3-1.noarch.rpm
yum -y install postgresql93-server postgresql93 postgresql93-libs postgresql93-contrib
yum -y install openssh-server openssh-clients
```


Create a user account on each host and add it to the Docker group:
	
```sh
useradd jeffmc
passwd jeffmc
usermod -a -G docker jeffmc
```

Next, on each host, add entries to /etc/hosts for each server:

```sh
192.168.56.110 rh65-admin.crunchy.lab rh65-admin
192.168.56.111 rh65-server1.crunchy.lab rh65-server1
192.168.56.112 rh65-server2.crunchy.lab rh65-server2
```

##### Docker Configuration

Each host needs to have Docker use a different dynamic bridge IP range, this
is done as follows:

edit the /etc/init.d/docker file, and add these options to the docker startup:

	*  (for rh65-admin)

	-H tcp://0.0.0.0:4243 -H unix:///var/run/docker.sock
	--bip=172.17.42.1/16 --dns=192.168.56.110 --dns=192.168.0.1

	*  (for rh65-server1)

	-H tcp://0.0.0.0:4243 -H unix:///var/run/docker.sock
	--bip=172.18.42.1/16 --dns=192.168.56.110 --dns=192.168.0.1

	*  (for rh65-server2)

	-H tcp://0.0.0.0:4243 -H unix:///var/run/docker.sock
	--bip=172.19.42.1/16 --dns=192.168.56.110 --dns=192.168.0.1

On each host, set Docker to start automatically, and reboot the system:

```sh
	chkconfig docker on
```

On each host, you will need to add static routes to allow traffic to
each Docker dynamic ethernet bridge.  On the admin server, you will
have the following static routes configured:
```sh
172.18.0.0/16 via 192.168.56.111 metric 0
172.19.0.0/16 via 192.168.56.112 metric 0
```

On the 'server1' node, you will have the following routes defined:
```sh
172.17.0.0/16 via 192.168.56.110 metric 0
172.19.0.0/16 via 192.168.56.112 metric 0
```

On the 'server2' node you will have the following routes defined:
```sh
172.17.0.0/16 via 192.168.56.110 metric 0
172.18.0.0/16 via 192.168.56.111 metric 0
```

These routes can be automatically added to your network on RHEL systems
by creating a file, /etc/sysconfig/network-scripts/route-eth0, that
contains these route configurations, that is, if your ethernet adapter
is named 'eth0', if not, rename the route file accordingly.


You will likely want to reboot your system around now:
```sh
	reboot
```

DNSBridge Build
---------------

CPM depends upon a DNSBRIDGE being run on all the CPM hosts.  The
DNSBRIDGE repo is cloned as follows:

	git clone git@github.com:crunchyds/dnsbridge.git

Set your GOPATH as follows:

	source dnsbridge/bin/setpath.sh

The setpath.sh script assumes you have cloned into your
local /home/$USER directory.

Next, compile the dnsbridge binaries:

	cd dnsbridge/src/crunchy.com
	make

Next, edit the dnsbridge/config/dnsbridgeclient-rhel65-initscript, and modify the ip address to specify your admin/dns server ip address, for example if your
admin host is 192.168.56.113, specify it as follows:

	-d 192.168.56.113:14000


Note:  you will need to edit the deploy-rhel65.sh script to
specify your chosen CPM host names!

Note:  you will need to edit the dnsbridgeclient init script to
specify your chosen CPM Admin host IP address!

Next, deploy the dnsbridge binaries and scripts to your
CPM hosts:

	cd dnsbridge/bin/deploy-rhel65.sh

Reboot the CPM servers.

Verify that dnsbridgeserver is running on the Admin/DNS host.
	ps ax | grep dnsbridgeserver

Verify that dnsbridgeclient is running on the all the CPM hosts.
	ps ax | grep dnsbridgeclient



Docker Base Image Build
-----------------------

This step is required for environments where you can't download a suitable
RHEL 6.5 Docker image.  If you have access to public available RHEL 6.5 Docker
image then you can skip this step.

On one of the RHEL 6.5 hosts, we will build a base Docker RHEL 6.5 image we
will later use as a base for the CPM Docker images.

Start by obtaining the CPM source code, by git this is done as follows:

	yum -y install git

	git clone git@github.com:crunchyds/docker-pg-cluster.git

then, as root, run this script (myimage.sh) as follows to build a RHEL6.5  docker image:

	docker-pg-cluster/pgcluster/bin/make-rhel65-image.sh

of if you are on a CentOS environment:

	docker-pg-cluster/pgcluster/bin/make-centos65-image.sh

You should get a success message if the image builds correctly. as follows towards
the very end of the script output:

	+ tar --numeric-owner -c -C /tmp/myimage .
	+ docker import - crunchy:rhel65
	74195164a0d897fb902a82fb1c8b40595a4f763f1f975d7df2b84bacd9ef3d2e
	+ docker run -i -t crunchy:rhel65 echo success
	success

Confirm the images is in your local Docker repo:
	[root@rh65-admin bin]# docker images
	REPOSITORY          TAG                 IMAGE ID            CREATED             VIRTUAL SIZE
	crunchy             rhel65              74195164a0d8        18 minutes ago      313.8 MB


You then save it off with:

	docker save crunchy:rhel65 > /tmp/rhel65-image.tar
or
	docker save crunchy:centos65 > /tmp/centos65-image.tar

Then to import it on another host:
	
	docker load < /tmp/rhel65-image.tar


Build CPM
------------

Now that we have a base Docker image, we can build our CPM images which
use that base image.

To do this you will need one of your servers to act as a build server.


##### SSH Keys
CPM includes an ssh server in each container to allow you access to
the running container, we place a known key into each container and if you
have the private key, you can ssh into the container.

You generate a unique ssh key using the following script:
```sh
./docker-pg-cluster/keygen.sh
```

This will generate a key pair you can use to ssh into the containers.

##### GOPATH

set your GOPATH env variable:

	export GOPATH=/home/jeffmc/docker-pg-cluster/pgcluster/

##### Compile CPM

	cd docker-pg-cluster/pgcluster/src/crunchy.com

	make


Building CPM Docker Images
-------------------

Now you are ready to build the CPM, CPM Admin, CPM Node, and CPM Pgpool 
Docker images!

in docker-pg-cluster/images/crunchy-cpm, you will build the CPM web
image as follows:

	cp Dockerfile.rhel65 Dockerfile
	make

in docker-pg-cluster/images/crunchy-admin, you will build the Admin
image as follows:

	cp Dockerfile.rhel65 Dockerfile
	make

in docker-pg-cluster/images/crunchy-node, you will build the Node
image as follows:

	cp Dockerfile.rhel65 Dockerfile
	make

in docker-pg-cluster/images/crunchy-pgpool, you will build the pgpool
image as follows:

	cp Dockerfile.rhel65 Dockerfile
	make

After this, you should see the CPM images in your Docker repository:
```sh
rh65-admin jeffmc~/docker-pg-cluster/images/crunchy-pgpool]: docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             VIRTUAL SIZE
crunchy-pgpool      latest              379dc72216aa        9 seconds ago       332.7 MB
crunchy-node        latest              5ebe1dc75033        5 minutes ago       334.5 MB
crunchy-cpm         latest              fe1acbc01e76        17 minutes ago      313.8 MB
crunchy-admin       latest              6a2c0e07fd7f        19 minutes ago      352.8 MB
crunchy             rhel65              74195164a0d8        About an hour ago   313.8 MB
```

##### Save Images

Next, we need to save off the images so that we can import them
into our CPM servers:
```sh
	docker save crunchy-pgpool:latest > /tmp/crunchy-pgpool.tar
	docker save crunchy-node:latest > /tmp/crunchy-node.tar
```

And optionally, save off the CPM and Admin images if you want to 
run them on a server different than what you built them upon:
```sh
	docker save crunchy-cpm:latest > /tmp/crunchy-cpm.tar
	docker save crunchy-admin:latest > /tmp/crunchy-admin.tar
```


##### Import Images

Remember to add the docker group to whatever user account
you will be using to import the Docker images, example:

```sh
sudo usermod -a -G docker jeffmc
```

On each CPM server that will run pgpool and postgresql containers, you 
need to import the images you have saved off as follows:
```sh
docker load < /tmp/crunchy-pgpool.tar
docker images
docker tag someimageid crunchy-pgpool:latest
docker load < /tmp/crunchy-node.tar
docker images
docker tag someimageid crunchy-node:latest
```

As the infrastructure emerges and is available in our customer's environments
we will likely place these sorts of images into a shared Docker repo,
, something akin to the new Redhat Docker Registry, but for now
we will manually copy them and import them.


Server Setup
------------

Next, you will copy all required server software to each CPM server.  Each
server.

There is a script that helps set up a CPM server, you will need to edit
this script to change the hostnames, but essentially it copies all the
required server files to each CPM host you have configured.  This is done as follows:

```sh
sudo ./docker-pg-cluster/deploy-rhel65.sh
```

Next, configure the CPM server services to start automatically, on each server,
and also open up any CPM used ports.

Ports used by CPM include:
* 8080 used by the CPM REST API (only on the CPM admin host)
* 13000 used by the CPM agent (on each host and in each container)
* 14000 used by the CPM DNS bridge server (only on the DNS host server)
* 5432 used by the PostgreSQL containers (on all containers and hosts)
* 53 used by the CPM DNS server (on the admin or DNS host only)
* 80 used by the CPM web application (on the admin host only)
* 9999 used by the CPM pgpool containers (on the node hosts only)


On the 'admin' server, enter the following commands
```sh
sudo iptables -I INPUT -p tcp --dport 8080 -j ACCEPT
sudo iptables -I INPUT -p tcp --dport 13000 -j ACCEPT
sudo iptables -I INPUT -p tcp --dport 14000 -j ACCEPT
sudo iptables -I INPUT -p tcp --dport 5432 -j ACCEPT
sudo iptables -I INPUT -p udp --dport 53 -j ACCEPT
```

On the 'node' servers, enter the following commands
```sh
sudo iptables -I INPUT -p tcp --dport 8080 -j ACCEPT
sudo iptables -I INPUT -p tcp --dport 5432 -j ACCEPT
sudo iptables -I INPUT -p tcp --dport 13000 -j ACCEPT
```

ICMP
-----
Due to the way the static routes work to the Docker dynamic ethernet
bridge, I found that I needed to allow ICMP traffic in the iptables
on RHEL 6.5.  So, remove the ICMP REJECT rules or work with your
sysadmin to create rules that allow CPM to ping IP addresses
within the Docker bridge IP range.

SAVE
-----
To make the iptables changes permanent, don't forget to add these
rules into the /etc/sysconfig/iptables file!

DNS Setup
---------
I typically run the DNS server used by CPM on the 'admin' server to save 
resources, you don't have to do this but this document assumes you are 
running the DNS (named) on the CPM 'admin' server.  

You will need to edit the dnsbridge/config/zonefiles to match
your chosen IP addresses.

You will need to edit the dns-setup.sh script to use your server host 
IP address, and then execute the following script:

```sh
./docker-pg-cluster/dns-setup.sh
```

The scripts are using an IP address of 192.168.56.103 for the DNS/admin
host, edit the named.conf and zone files if you are using a different
IP address.

The DNS files refer to the DNS host as ns.crunchy.lab, update
your /etc/hosts file to include an entry for ns.crunchy.lab which
is the same IP as the admin/dns host.

Next, add an entry for ns.crunchy.lab to your /etc/hosts file, ns.crunchy.lab
refers to the DNS server host, it is just another name used to 
refer to the DNS host (aka the 'admin' host).

Lastly, restart the DNS server:

```sh
service named restart
```

Next, edit the DNS zone files to use the correct DNS IP address, 
Next, set up the DNS server following the docs/setup-dns.md instructions.
Note that I've assumed an IP address of 192.168.56.110 for the DNS 
server host, you will need to make substitutions in your DNS zone
files if you use a different IP address for your host.  CPM (in it's current
form) requires a DNS host that it can dynamically add DNS records to!

Validate DNS
------------

First, verify that your /etc/resolv.conf on each CPM server
has the following values:
```sh
nameserver 192.168.56.110
nameserver 192.168.0.1
search crunchy.lab
```

This specifies the CPM DNS name server will be searched first
when resolving hostnames and that crunchy.lab will be the default
domain name used.

Validate the DNS configuration by registering a test hostname
as follows:

```sh
# nsupdate
> server ns.crunchy.lab
> zone crunchy.lab.
> update add test1.crunchy.lab 60 A 192.168.56.120
> send
> quit
# nsupdate
> zone 56.168.192.in-addr.arpa
> update add 120.56.168.192.in-addr.arpa 3500 IN PTR test1.crunchy.lab.
> send
> quit
[root@rh65-admin dynamic]# dig test1.crunchy.lab
;; ANSWER SECTION:
test1.crunchy.lab.	60	IN	A	192.168.56.120
[root@rh65-admin dynamic]# dig  -x 192.168.56.120
;; ANSWER SECTION:
120.56.168.192.in-addr.arpa. 3500 IN	PTR	test1.crunchy.lab.
```

This test shows both A and PTR records being created for
test1.crunchy.lab, and the lookups being performed.

Next, try the following command on all other CPM hosts to verify
that the DNS server is accessible and resolves the hostname correctly:

```sh
dig test1.crunchy.lab
```

If not, verify that the iptables on the DNS server is allowing DNS
requests.

Start it all
-----------

Run a simple test before starting up CPM to verify that your
environment is ready, this is done by running the following
script:

```sh
localhost jeffmc~/docker-pg-cluster]: ./net-rhel65.sh 
ping rh65-admin.crunchy.lab
ping rh65-server1.crunchy.lab
ping rh65-server2.crunchy.lab
ping 172.17.42.1
dnsbridgeclient running on rh65-admin.crunchy.lab
docker running on rh65-admin.crunchy.lab
dnsbridgeclient running on rh65-server1.crunchy.lab
docker running on rh65-server1.crunchy.lab
dnsbridgeclient running on rh65-server2.crunchy.lab
docker running on rh65-server2.crunchy.lab
 named running on rh65-admin.crunchy.lab
 dnsbridgeserver running on rh65-admin.crunchy.lab
 cpmagentserver running on rh65-server1.crunchy.lab
 cpmagentserver running on rh65-server2.crunchy.lab
postgres installed on  rh65-admin.crunchy.lab
postgres installed on  rh65-server1.crunchy.lab
postgres installed on  rh65-server2.crunchy.lab
SWEET SUCCESS!!!
```

If you see the 'SWEET SUCCESS' message then you have a valid
CPM configuration, note that this test assumes a certain
CPM host configuration with certain host names, edit if necessary
to match your CPM deployment.


Now you are ready to run CPM and the CPM Admin docker instances. 

On your 'admin' host, as a normal user, create the CPM 'admin'
container as follows:

```sh
sudo rm -rf /var/lib/pgsql/cluster-admin
sudo mkdir /var/lib/pgsql/cluster-admin
sudo chown postgres:postgres /var/lib/pgsql/cluster-admin
docker run --name=cluster-admin -d --privileged \
                -v /var/run/docker.sock:/tmp/docker.sock \
                -v /var/lib/pgsql/cluster-admin:/pgdata crunchy-admin
```

This will create a data volume for the 'cluster-admin' container
to store it's PostgreSQL data within.  
 

Next, create the CPM 'web' service as follows:

```sh
docker run --name=cpmweb -d --privileged \
	-v /home/jeffmc/docker-pg-cluster/images/crunchy-cpm/www:/www crunchy-cpm
```

This command creates the web application and mounts the web content
from the local user's directory where the CPM source code has been
copied to, in this example, /home/jeffmc/docker-pg-cluster/images/crunchy-cpm/www.

At this point, you should see CPM and cluster-admin running:

```sh
docker ps
```

You should also verify that you can ping each newly created container:

```sh
ping cpmweb.crunchy.lab
ping cluster-admin.crunchy.lab
```

If you can't ping these containers, then there is a problem with
the DNS bridging to Docker, contact Crunchy support if this happens.

Typical problems include having the wrong permissions on the
DNS zone files, these need to be owned by 'named' and have 644 permissions.

You can now browse to CPM at http://cpmweb.crunchy.lab

CPM should be ready for use at this point.
