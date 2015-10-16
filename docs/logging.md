
Logging
=========================

The CPM services log into the /var/cpm/logs directory on
your docker host.

This directory is initially created by the run-cpm.sh script.

Logs are written to this directory by the following
containers:

+ cpm
+ cpm-admin
+ cpm-collect
+ cpm-task

Each container mounts 
````````````
/cpmlogs
````````````

Each startup script for the services now writes to /cpmlogs which gets
mapped by Docker to /var/cpm/logs on the Docker host.




