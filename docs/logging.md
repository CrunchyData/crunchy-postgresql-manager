
Logging
=========================

glog is used for logging in the CPM code

To trigger the use of glog, your main() functions need
to call flag.Parse if you want to process any glog
command line

The main CPM services log into the /opt/cpm/logs directory.

This directory is created by the run-cpm.sh script.

Logs are written to this directory by the following
containers:

cpm
cpm-admin
cpm-mon
cpm-backup

glog flushes output to these log files perodically or upon demand
by coding glog.Flush()

Each container mounts /cpmlogs

Each script for the services now writes to /cpmlogs which gets
mapped by Docker to /opt/cpm/logs on the Docker host.



