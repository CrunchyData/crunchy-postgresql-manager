CPM Docker Images
=========================


The following container images are built and used by CPM:

+ cpm - the web application
+ cpm-admin - the REST API
+ cpm-collect - the monitoring service that collects metrics and stores them into prometheus
+ cpm-task - the admin task service that spawns admin task containers, either restore or backup jobs
+ cpm-backup-job - provisioned on demand to perform a PG backup
+ cpm-restore-job - provisioned on demand to perform a backrest restore
+ cpm-pgpool - provisioned on demand to run pgpool
+ cpm-node - provisioned on demand to run Postgres
+ cpm-promdash - prometheus dashboard app, used to host predefined graph templates used by the cpm web app
+ cpm-prometheus - prometheus database used to store metrics
+ cpm-serverapi - server container, run on each CPM server node, configured with a
Docker container name of 'cpm-<<servername>>'

The images are built by running make:
````````````
make buildimages
````````````
