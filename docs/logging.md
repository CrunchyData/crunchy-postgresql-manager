
Logging
=========================

glog is used for logging in the CPM code

To trigger the use of glog, your main() functions need
to call flag.Parse if you want to process any glog
command line

The CPM services log into the /var/cpm/logs directory on
your docker host.

This directory is initially created by the run-cpm.sh script.

Logs are written to this directory by the following
containers:

+ cpm
+ cpm-admin
+ cpm-mon
+ cpm-backup

glog flushes output to these log files perodically or upon demand
by coding glog.Flush()

Each container mounts 
````````````
/cpmlogs
````````````

Each startup script for the services now writes to /cpmlogs which gets
mapped by Docker to /var/cpm/logs on the Docker host.

Within a container startup script, you will pass the following
glog flags to direct log output to the mounted volume:

````````````
-log_dir=/cpmlogs -logtostderr=false
````````````



