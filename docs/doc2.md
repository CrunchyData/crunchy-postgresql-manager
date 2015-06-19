

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

set up network adapters
	- static address is set using nm-tui on the enp0s8 adapter
		- specify 192.168.56.103 as IP address
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
sudo ./sbin/dev-setup.sh

