

Setting up Development VM
-------------------------

Use CentOS 7 - minimal as the VM base
specify 2 network adapters
	- bridged (used to get a DHCP address)
	- host-only (used to specify a static IP of 192.168.56.103)
		- specify your host-only addresses to be 192.168.56.XXX

install 
	- docker
	- golang
	- git
	- mercurial
	- net-tools
	- bind-utils
	- sysstat

	- yum -y install docker golang git mercurial net-tools bind-utils sysstat

set up network adapters
	- static address is set using nm-tui on the enp0s8 adapter
		- specify 192.168.56.103 as IPADDR
		- specify 192.168.56.1 as GATEWAY
		- dns 192.168.0.1 (this is my routers dns nameserver addr)
	- set both adapters to connect at startup


disable firewalld (for development)
	- systemctl stop firewalld.service
	- systemctl disable firewalld.service

configure docker
	- edit /etc/sysconfig/docker
		- OPTIONS='--selinux-enabled --bip=172.17.42.1/16 --dns-search=crunchy.lab --dns=192.168.56.103 --dns=192.168.0.1'
	- systemctl enable docker.service
	- systemctl start docker.service

Install Skybridge
-----------------
	- git clone git@github.com:CrunchyData/skybridge.git
	- edit skybridge/bin/run-skybridge.sh
		- change the static IP address to 192.168.56.103
	- sudo ./skybridge/bin/run-skybridge.sh
	- sudo vi /etc/resolv.conf
		- add 192.168.56.103 as a primary DNS nameserver
	

Setup the CPM Source Environment
---------------------
export GOPATH=~/devproject
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH

cd devproject

go get github.com/tools/godep
go get github.com/crunchydata/crunchy-postgresql-manager
cd src/github.com/crunchydata/crunchy-postgresql-manager

godep restore
make build

make buildimages

Run CPM
--------------
edit ./sbin/dev-setup.sh
	- change the home directory to the user you are building under

sudo ./sbin/dev-setup.sh

sudo ./run-cpm.sh

Verify your build
-----------------

You should see the following containers running:
~~~~~~~~~~~~~~~~~~~~~~
CONTAINER ID        IMAGE                            COMMAND                CREATED             STATUS              PORTS                                        NAMES
4b2c682ed159        crunchydata/cpm:latest           "/var/cpm/bin/startn   20 seconds ago      Up 19 seconds       13001/tcp, 192.168.56.103:13050->13000/tcp   cpm                 
58003cd15731        prom/prometheus:latest           "/bin/prometheus -co   2 minutes ago       Up 2 minutes        9090/tcp, 192.168.56.103:16000->3000/tcp     cpm-prometheus      
9a366dab40f3        prom/promdash:latest             "./run ./bin/thin st   2 minutes ago       Up 2 minutes        3000/tcp, 192.168.56.103:15000->8080/tcp     cpm-promdash        
2ef2fa68be12        crunchydata/cpm-collect:latest   "/var/cpm/bin/start-   2 minutes ago       Up 2 minutes        5432/tcp, 8080/tcp                           cpm-collect         
91ac5a3ce888        crunchydata/cpm-backup:latest    "/var/cpm/bin/start-   2 minutes ago       Up 2 minutes        5432/tcp, 13000/tcp                          cpm-backup          
d05f5b7ac614        crunchydata/cpm-admin:latest     "/var/cpm/bin/starta   2 minutes ago       Up 2 minutes        5432/tcp, 192.168.56.103:14000->13000/tcp    cpm-admin           
f08c1afc90b3        crunchydata/skybridge:latest     "/var/cpm/bin/start-   3 hours ago         Up 43 minutes       192.168.56.103:53->53/udp, 53/tcp            skybridge 
~~~~~~~~~~~~~~~~~~~~~~

You should be able to ping the following addresses:
~~~~~~~~~~~~~~~~~~~
ping -c 1 cpm.crunchy.lab
ping -c 1 cpm-admin.crunchy.lab
ping -c 1 cpm-backup.crunchy.lab
ping -c 1 cpm-promdash.crunchy.lab
ping -c 1 cpm-prometheus.crunchy.lab
~~~~~~~~~~~~~~~~~~~

Access CPM
----------------

The CPM web interface is located at:
http://cpm.crunchy.lab:13050
http://192.168.56.103:13050

The CPM REST API is located at:
http://cpm-admin.crunchy.lab:13000
http://192.168.56.103:14000

The PromDash metrics dashboard is at:
http://cpm-promdash:8080
http://192.168.56.103:15000

The Prometheus database console is at:
http://cpm-prometheus:3000
http://192.168.56.103:16000

