Developer Setup
=================

To develop on CPM, you can do the following on your dev box, typically I set this
up on the 'admin' server's VM but that is your call:

	yum -y install gcc git

clone the project
-----------------

dependencies
---------------------

I use various projects out on github and include their source within
the CPM source structure, these are located on github at the following
locations:

	git@github.com:fsouza/go-dockerclient.git
	git@github.com:ant0ine/go-json-rest.git
	git clone git@github.com:lib/pq.git

I currently copy that code into the CPM src/github.com subdirectory to have our own copy
of the source to build against, I periodically update this when fixes are released.

yes, there are better ways to do this and I'll get to it soon!  example:

	http://git-scm.com/docs/git-submodule



building
---------

I currently use Makefiles typically to build the golang source, each src subdir has a
Makefile.

Here are the steps to building from source:

	cd docker-pg-cluster/pgcluster/src/github.com/lib/pq
	go install github.com/lib/pq

	cd docker-pg-cluster/pgcluster/src/github.com/ant0ine/go-json-rest/rest
	go install github.com/ant0ine/go-json-rest/rest
	
	cd docker-pg-cluster/pgcluster/src/github.com/fsouza/go-dockerclient
	go install github.com/fsouza/go-dockerclient

	cd docker-pg-cluster/pgcluster/src/crunchy.com/adminapi
	make

	cd docker-pg-cluster/pgcluster/src/crunchy.com/cpmagent
	make

	cd docker-pg-cluster/pgcluster/src/crunchy.com/dnsbridge
	make
	
	cd docker-pg-cluster/pgcluster/src/crunchy.com/pgtemplates
	make

After that, your binaries are stored in:

	docker-pg-cluster/pgcluster/bin
	

localhost routing
---------------
On my local dev box, I add network routes as follows to allow me to reach each
Docker container locally:

	ip route add 172.17.0.0/16 via 192.168.56.103 dev vboxnet0

	ip route add 172.18.0.0/16 via 192.168.56.101 dev vboxnet0

	ip route add 172.19.0.0/16 via 192.168.56.102 dev vboxnet0	
