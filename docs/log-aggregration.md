## Log Aggregation

To see how we can aggregate all the CPM and CPM provisioned
database logs into something useful, we will stand
up an EFK docker image.  The EFK includes:
elasticsearch
fluentd
kibana

These tools provide an alternative to splunk for log analysis and
aggregation.

The EFK docker image is run by the sbin/run-efk.sh script.

CPM aggregates logs in a couple of ways....
### CPM Internal Logs
Logs produced by the CPM administration containers
is logged to stdout by each CPM container, we can
specify a Docker log driver to route that stdout log
output to the fluentd system if we desire to aggregate
all of the CPM logs.

In the startup scripts for CPM, run-cpm.sh, you will see
that the Docker log driver is specified to use the fluentd
driver:
~~~~~~
--log-driver=fluentd \
--log-opt fluentd-address=192.168.0.107:24224 \
--log-opt fluentd-tag=docker.cpm-admin \
~~~~~~

This will send the docker log output for the CPM containers
to the fluentd server running at the specified address.

### Postgres Container Logs
On each CPM provisioned postgres container, we have enabled postgres
to send log output to both stdout and syslog.

The stdout log is still required for things like pgbadger to work but
we can use the postgres feature of logging to syslog to aggregrate all
the postgres logs into a central location like fluentd as well.

* Setup

To set up the postgres syslog logging, you first include rsyslogd into
the docker image.  You will see rsyslogd added to the list of packages
that get installed by the container.

We then configure rsyslogd according to the following tutorial by
Dan Walsh:

http://www.projectatomic.io/blog/2014/09/running-syslog-within-a-docker-container/

We also make the containers rsyslogd forward the syslog to a remote
syslog server by adding a section to the /etc/rsyslogd.conf file
as follows:

$WorkDirectory /var/lib/rsyslog
$ActionQueueFileName fwdRule1
$ActionQueueMaxDiskSpace 2g
$ActionQueueSaveOnShutdown on
$ActionQueueType LinkedList
$ActionResumeRetryCount -1
*.* @@192.168.0.105:514

Also, in /etc/rsyslog.d/listen.conf you have to comment out the
line $SystemLogSocketName ....

We then provide this rsyslogd.conf file to the container when it is
provisioned and start up the rsyslogd in the containers startup script.

