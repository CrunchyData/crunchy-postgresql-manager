Developer Setup
=================

Here are the steps required to set up a CPM development environment on a 
clean RHEL or Centos 7.1 minimal installation.

This instruction assumes you are using a static IP address of
192.168.0.107 for your CPM server.

### RHEL Setup  ###
note that for RHEL 7.1, you will need to add the following repos:
~~~~~~~~~~~~~~~
subscription-manager repos --enable=rhel-7-server-extras-rpms
subscription-manager repos --enable=rhel-7-server-optional-rpms
~~~~~~~~~~~~~~~

### Install Dependencies development machine ####
Note that I like to use the PGDG postgres distro instead of the redhat provided postgres!

~~~~~~~~~~~~~~~~~~~~~~~~~~~~
sudo yum -y install net-tools bind-utils install golang git docker mercurial sysstat
rpm -Uvh http://dl.fedoraproject.org/pub/epel/7/x86_64/e/epel-release-7-5.noarch.rpm
rpm -Uvh http://yum.postgresql.org/9.4/redhat/rhel-7-x86_64/pgdg-centos94-9.4-1.noarch.rpm
yum install -y nmap-ncat procps-ng postgresql94 postgresql94-contrib postgresql94-server libxslt unzip openssh-clients hostname bind-utils
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

### Setup Go Project Structure ###
As your development user, create the development directory as follows:
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

make build
~~~~~~~~~~~~~~~~~~~~~~~~
### Compile CPM ###
~~~~~~~~~~~~~~~~~~~~~~~~
make build
~~~~~~~~~~~~~~~~~~~~~~~~

### Configure Docker ###
Edit the docker configuration by editing the OPTIONS parameter as follows:
~~~~~~~~~~~~~~~~~~~~~~~~
vi /etc/sysconfig/docker
OPTIONS='--selinux-enabled --bip=172.17.42.1/16 --dns-search=crunchy.lab --dns=192.168.0.107 --dns=192.168.0.1'
systemctl enable docker.service
systemctl start docker.service
docker info
~~~~~~~~~~~~~~~~~~~~~~~~

### Build CPM Docker Images ###
~~~~~~~~~~~~~~~~~~~~~~~~
make buildimages
docker images
~~~~~~~~~~~~~~~~~~~~~~~~

### Pull down Prometheus images ###
~~~~~~~~~~~~~~~~~~~~~~~~
sudo docker pull prom/promdash
sudo docker pull prom/prometheus
~~~~~~~~~~~~~~~~~~~~~~~~

### Disable Firewalld ###
~~~~~~~~~~~~~~~~~~~~~~~~
systemctl disable firewalld.service
systemctl stop firewalld.service
~~~~~~~~~~~~~~~~~~~~~~~~

There is a document, firewall-setup.md, that shows how the CPM ports
can be opened up.

### Setup skybridge ###

CPM services are found using DNS by the various parts of CPM.  When
a Docker image is started, we need it to be registered with a DNS service
and the local machine configured to resolve using that DNS server.

CPM requires a reliable IP address of the host on which it is running.
When a VM is created to develop CPM upon, you would create an extra
Ethernet adapter typically so that you can assign it a static IP
address.  In Virtualbox, this adapter would be a Host-Only adapter
for example.

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
nameserver 192.168.0.107
nameserver 192.168.0.1
~~~~~~~~~~~~~~~~~

You can make these changes to your /etc/resolv.conf permanent by
adding the following settings to your ethernet adapter configuration
in /etc/syconfig/network-scripts:
~~~~~~~~~~~~~~~~~~~~~~
DNS1=192.168.0.107
DNS2=192.168.0.1
DOMAIN=crunchy.lab
PEERDNS=no
~~~~~~~~~~~~~~~~~~~~~~

This will cause the skybridge DNS nameserver to be queried first.


Pull down skybridge as follows:
~~~~~~~~~~~~~~~~~~~~
sudo docker pull crunchydata/skybridge
~~~~~~~~~~~~~~~~~~~~

Start skybridge by editing the sbin/run-skybridge.sh script
to specify your local IP address, then run the skybridge container:
~~~~~~~~~~~~~~~~~~~~
sudo ./sbin/run-skybridge.sh
~~~~~~~~~~~~~~~~~~~~


Start CPM Server Agent
----------------------
After you have successfully compiled CPM and built the CPM Docker images,
on each server that is to run CPM, you will need to start a CPM Server
Agent.  The server agent is run within the cpm-server container on each
server host that will be configured to be used in CPM.

Each container needs to be started with skybridge running and also
have its container name set to 'cpm-servername' where servername is
the server name you have given the server in the CPM server page once
CPM is running.

For this example, I will name the CPM server, newserver.

So, edit the sbin/run-cpmserver.sh script, and modify the server
name to newserver.

Then run the script which will create a running cpm-server named
cpm-newserver.
~~~~~~~~~~~~~~~
sudo ./run-cpmserver.sh
ping cpm-newserver
~~~~~~~~~~~~~~~

If your networking and skybridge are all correct, then you should be able 
to ping the cpm-server container as follows:
~~~~~~~~~~~~~~~
ping cpm-newserver
~~~~~~~~~~~~~~~

Running CPM
--------------

Modify the run-cpm.sh script by updating the INSTALLDIR 
variable to the path on your host that you are installing CPM 
from.

Also, edit or remove the local host port mapping that is
provided in the example to meet your local requirements
for accessing CPM.

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

CPM Admin API
----------------------
http://cpm-admin.crunchy.lab:13001

Prometheus Dashboard
----------------------
http://cpm-promdash.crunchy.lab:3000

If you are running CPM on a VM (host-only) and
accessing CPM from the VM host (not the guest), then
you will need to edit the dashboard server
configuration via the PromDash user interface
and specify the prometheus server URL
as http://192.168.56.103:16000.

Prometheus DB
----------------------
http://cpm-prometheus.crunchy.lab:9090

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

Initially you will need to first define your CPM server which
is your CPM host (e.g. 192.168.0.107, newserver)

Then you will be ready to start creating PostgreSQL instances.

nginx selinux issues
--------------------

in some cases with selinux enabled, you might see AVC errors, if so, look at this:

http://axilleas.me/en/blog/2013/selinux-policy-for-nginx-and-gitlab-unix-socket-in-fedora-19/


###Godocs

To see the godocs, install godoc, and start up the godoc server, then 
browse to the CPM API documentation:
~~~~~~~~~~~~~~~~~~~~
go get golang.org/x/tools/cmd/godoc
godoc -http=:6060
~~~~~~~~~~~~~~~~~~~~
