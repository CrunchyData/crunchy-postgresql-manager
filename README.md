Crunchy Postgresql Manager (v1.0.2)
==========================
Crunchy Postgresql Manager (CPM) is a Docker-based solution which
provides an on-premise PostgreSQL-as-a-Service platform. CPM utilizes
Docker Swarm for multi-host scaling.

CPM allows for the quick provisioning of PostgreSQL databases
and streaming replication clusters.  

CPM also allows you to monitor and administer PostgreSQL
databases.  Currently CPM only works with databases that have
been provisioned by CPM.

![CPM Web UI](./docs/cpm.png)

A user guide is available at:
[docs/htmldoc/user-guide.html](https://rawgit.com/crunchydata/crunchy-postgresql-manager/master/docs/htmldoc/user-guide.html)

A developer-installation guide is available at:
[docs/htmldoc/doc.html](https://rawgit.com/crunchydata/crunchy-postgresql-manager/master/docs/htmldoc/doc.html)

A description of the REST API is at:
[docs/htmldoc/rest-api.html](https://rawgit.com/crunchydata/crunchy-postgresql-manager/master/docs/htmldoc/rest-api.html)

CPM consists of the following containers:

* cpm - cpm.crunchy.lab - the nginx server that hosts the CPM web app, http://cpm.crunchy.lab:13001

* cpm-admin - cpm-admin.crunchy.lab - the REST API for CPM, http://cpm-admin.crunchy.lab:13001

* cpm-task - cpm-task.crunchy.lab - the task scheduler process used by CPM to schedule and run administrative task jobs

* cpm-collect - cpm-collect.crunchy.lab - the monitoring process used
to collect metrics, these metrics are collected by the Prometheus server
running as cpm-prometheus

* cpm-promdash - cpm-promdash.crunchy.lab - the Prometheus dashboard that can be used to view/query collected CPM metrics , graphs from this dashboard
are displayed within the CPM user interface, the user interface is
found at http://cpm-promdash:3000

* cpm-prometheus - cpm-prometheus.crunchy.lab - the Prometheus database
is found at http://cpm-prometheus:9090

* cpm-server - cpm-yourserver.crunchy.lab - the server agent, one of these is
run on each CPM server, the name of the container is important, it needs to
match the server name you define in CPM
