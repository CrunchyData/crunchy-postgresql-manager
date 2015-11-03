## Swarm Configuration

CPM uses Docker Swarm to virtualize multiple Docker servers into
a single virtual server.  This is a convenient way to implement
multiple host Docker which is necessary to scale out the
CPM containers onto multiple Docker hosts.  In this example
we run the Swarm manager and agent on the same host, this is
they way a developer might run CPM.  In a real setup, you would
have a single manager and multiple swarm agent hosts.

For this example configuration, we start the Swarm Manager
on 192.168.0.103:8000

The Swarm agent is started to listen to 0.0.0.0:2375

###Installation
Swarm is provided by Docker at https://github.com/docker/swarm.  Use the
instructions at the Swarm github page to install a binary version
of Swarm into the /usr/local/bin directory of all the servers you
will be using for CPM.

Swarm needs a single token to define the cluster you are creating.  This
is done one-time as follows, save this token value for future reference:
~~~~~~~~~~~~~~~~~~~~~~~~
swarm create
7b9fb5037919f89bd52c3c4888586be3
~~~~~~~~~~~~~~~~~~~~~~~~

###Docker Configuration
Docker is configured on each server to listen to 0.0.0.0:2375 for API events.  On
Centos/RHEL this is done by adding -H tcp://0.0.0.0:2375 in the /etc/sysconfig/docker
file:
~~~~~~~~~~~~~~~~~~~~~~~~
export SWARM_PORT=2375
/usr/bin/docker -d --selinux-enabled -H tcp://0.0.0.0:$SWARM_PORT
~~~~~~~~~~~~~~~~~~~~~~~~

###Startup
On each server in your cluster, Start the swarm server agent listening to the local Docker API:
~~~~~~~~~~~~~~~~~~~~~~~~
export LOCAL_HOST=192.168.0.103
export SWARM_PORT=2375
swarm join --addr=$LOCAL_HOST:$SWARM_PORT token://7b9fb5037919f89bd52c3c4888586be3
~~~~~~~~~~~~~~~~~~~~~~~~

On one server in your cluster, Start the swarm manager that listens to CPM
requests:
~~~~~~~~~~~~~~~~~~~~~~~~
export MANAGER_HOST=192.168.0.103
export MANAGER_PORT=8000
swarm manage --host $MANAGER_HOST:$MANAGER_PORT token://7b9fb5037919f89bd52c3c4888586be3
~~~~~~~~~~~~~~~~~~~~~~~~

###Test

To see what servers are include in the swarm:
~~~~~~~~~~~~~~~~~~~~~~~~
export MANAGER_HOST=192.168.0.103
export MANAGER_PORT=8000
swarm list token://7b9fb5037919f89bd52c3c4888586be3
docker -H tcp://$MANAGER_HOST:$MANAGER_PORT info
~~~~~~~~~~~~~~~~~~~~~~~~

You now run docker commands via the swarm manager ip:port to interact with swarm:
~~~~~~~~~~~~~~~~~~~~~
export MANAGER_HOST=192.168.0.103
export MANAGER_PORT=8000
docker -H tcp://$MANAGER_HOST:$MANAGER_PORT info
docker -H tcp://$MANAGER_HOST:$MANAGER_PORT run
docker -H tcp://$MANAGER_HOST:$MANAGER_PORT ps
docker -H tcp://$MANAGER_HOST:$MANAGER_PORT logs
~~~~~~~~~~~~~~~~~~~~~

