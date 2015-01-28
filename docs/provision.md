
Provisioning 
=========================

Bootstrap cluster-amin
===================
We need to create a cluster-admin as part of bootstrapping
CPM.  This is done with the following commands on the Admin
server:
	sudo rm -rf /var/lib/pgsql/cluster-admin
	sudo mkdir /var/lib/pgsql/cluster-admin
	sudo chown postgres:postgres /var/lib/pgsql/cluster-admin
	sudo chcon -Rt svirt_sandbox_file_t /var/lib/pgsql/cluster-admin/
	docker run --name=cluster-admin -d \
		-v /var/lib/pgsql/cluster-admin:/pgdata crunchy-admin

For Kube deployment:
	cd ~/docker-pg-cluster/images/crunchy-admin/conf
	openshift kube create pods -c ./cluster-admin-pod.json


Bootstrap CPM
=============
To start CPM, and have it attach to your local web content for serving:

	sudo chcon -Rt svirt_sandbox_file_t /home/jeffmc/docker-pg-cluster/images/crunchy-cpm/www
	docker run --name=cpm -d \
		-v /home/jeffmc/docker-pg-cluster/images/crunchy-cpm/www:/www crunchy-cpm

In this example, the web content is located at:
	/home/jeffmc/docker-pg-cluster-images/crunchy-cpm/www


Sample API Commands
===================
Typically CPM functionality is accessed via the CPM web application, however
all CPM functions can also be performed via the CPM REST API directly.  
Examples include:

1) create a cluster

	curl -i -d '{"ID":"","Name":"inserted8","ClusterType":"streaming","Status":"uninitialized"}' http://cluster-admin.crunchy.lab:8080/cluster

list the clusters:

	curl http://cluster-admin.crunchy.lab:8080/clusters 

list a single cluster:
	curl http://cluster-admin.crunchy.lab:8080/cluster/1

optional: to delete a cluster
	curl -X DELETE http://cluster-admin.crunchy.lab:8080/cluster/17

1a) create a server to host the nodes
	curl -i -d '{"ID":"","Name":"server1","IPAddress":"192.168.56.101","PGDataPath":"/var/lib/pgsql"}' http://cluster-admin.crunchy.lab:8080/server
	curl -i -d '{"ID":"","Name":"server2","IPAddress":"192.168.56.102","PGDataPath":"/var/lib/pgsql"}' http://cluster-admin.crunchy.lab:8080/server
	curl http://cluster-admin.crunchy.lab:8080/servers 

2) OPTIONAL:  list all nodes for cluster 1

	curl http://cluster-admin.crunchy.lab:8080/nodes/1

2a) OPTIONAL: list all the containers

	curl http://cluster-admin.crunchy.lab:8080/nodes 

2) create a node for pgmaster

	here params are passed in the following order with a '.' between
	each param:
	Role.ServerID.ContainerName.ContainerType

	curl http://cluster-admin.crunchy.lab:8080/provision/master.2.pgmaster.crunchy-node


3) create a node for pgstandby
	curl http://cluster-admin.crunchy.lab:8080/provision/standby.2.pgstandby.crunchy-node

3a) OPTIONAL: to delete a node enter this
	curl -X DELETE http://cluster-admin.crunchy.lab:8080/node/17

3b) OPTIONAL:  see what nodes are not assigned to any cluster
	curl http://cluster-admin.crunchy.lab:8080/nodes/nocluster

4) assign the nodes to the cluster, here ID is the node id
	curl -i -d '{"ID":"1","ClusterID":"1"}' http://cluster-admin.crunchy.lab:8080/event/join-cluster
	curl -i -d '{"ID":"2","ClusterID":"1"}' http://cluster-admin.crunchy.lab:8080/event/join-cluster

6) configure the master node - templating happens here - restarts db, here ID is the node id
	curl -i -d '{"ID":"1"}' http://cluster-admin.crunchy.lab:8080/event/configure-master

7) configure the standby nodes - templating happens here - restarts db
pass in the ID of the node we want to configure

	curl -i -d '{"ID":"2"}' http://cluster-admin.crunchy.lab:8080/event/configure-standby

after this, you should have a simple master-standby set of nodes
replicating....

8) verify that we can stop the standby node, here ID is the node id
of the server we want to stop
	curl -i -d '{"ID":"2"}' http://cluster-admin.crunchy.lab:8080/admin/stop-pg

9) verify that we can start the standby node, here ID is the node id
of the server we want to start
	curl -i -d '{"ID":"2"}' http://cluster-admin.crunchy.lab:8080/admin/start-pg

10) fail over the standby node, turning it into a master node, here
ID is the node id of the fail over node

	curl http://cluster-admin.crunchy.lab:8080/admin/failover/2

11) monitor a server
 	curl http://cluster-admin.crunchy.lab:8080/monitor/server-getinfo/1.cpmiostat
 	curl http://cluster-admin.crunchy.lab:8080/monitor/server-get-info/1.cpmdf

12) ping postgres on a container
 	curl http://cluster-admin.crunchy.lab:8080/monitor/container-getinfo/1.pgstatus

13) get database stats on a container

 	curl http://cluster-admin.crunchy.lab:8080/monitor/container-getinfo/1.bgwriter
 	curl http://cluster-admin.crunchy.lab:8080/monitor/container-getinfo/1.statreplication
 	curl http://cluster-admin.crunchy.lab:8080/monitor/container-getinfo/1.statdatabase


