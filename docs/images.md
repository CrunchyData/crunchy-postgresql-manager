CPM Docker Images
=========================


The following container images are built and used by CPM:

+ cpm - the web application
+ cpm-admin - the REST API
+ cpm-mon - the monitoring service
+ cpm-backup - the backup/admin task service
+ cpm-backup-job - provisioned on demand to perform a PG backup
+ cpm-pgpool - provisioned on demand to run pgpool
+ cpm-node - provisioned on demand to run Postgres
+ cpm-dashboard - optional Grafana dashboard for advanced metrics graphing

The images are built by running make:
````````````
cd ./images
make
````````````

CPM references these image names within the code so any changes
to these image names would require a change to the CPM code base
currently.
