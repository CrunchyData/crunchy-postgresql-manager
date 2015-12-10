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
is logged to stdout by each CPM container (except cpm-web), we can
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
to the fluentd server running at the specified address (e.g. 107:24224).

The exception to this is the cpm-web (nginx) container.  Nginx has
a configuration that requires you to send stdout and stderr to
a file.  There is a bug in Docker (1.8.2 and earlier versions), that
prevents a non-root user write access to /dev/stdout and /dev/stderr.
See https://github.com/docker/docker/issues/6880 for details on the bug.

Therefore, until the next release of Docker, we will continue to 
send cpm-web output to log files mounted from the local host (e.g. /var/cpm/logs).

If you want to have the CPM product containers send their logs to the Docker log instead
of using the cpm-efk and fluentd logging, just remove the log driver lines above
from the run-cpm.sh startup script when you start the CPM product containers.  


### Postgres Container Logs
On each CPM provisioned postgres container, we have enabled postgres
to send log output to both stdout and syslog.  We forward the syslog inside
the container to the cpm-efk syslog server if the /syslogconfig mount within
the container has a rsyslog.conf file present, if not, the rsyslogd daemon
will not be started inside each postgres container.

The postgres stdout log is still required for things like pgbadger to work but
we can use the postgres feature of logging to syslog to aggregrate all
the postgres logs into a central location like fluentd as well.

* Setup

To support syslog logging within each Postgres container (cpm-node),
I have included rsyslogd into the installed packages.

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

We then provide this rsyslogd.conf and /etc/rsyslog.d/listen.conf files to the container when it is
provisioned and start up the rsyslogd in the containers startup script.  These config files are
provided by means of the /syslogconf volume mount in the container.  Currently the location
on the local host for these config files is /var/cpm/config, so the mount within
each postgres container looks like this:
~~~~~~~~
-v /var/cpm/config:/syslogconfig
~~~~~~~~

This means on each CPM host, you will need to install these syslog config files at the /var/cpm/config
path if you want to log the postgres logs to the remote (cpm-efk) syslog.

Sample rsyslog.conf and listen.conf files are stored in the github CPM_ROOT/config/syslog/ directory

Postgres Configuration Changes
------------------------------
So, we want to still use pgbadger which requires postgres send output to log files, and
we also want to support sending output to the cpm-efk.  So, we use the postgres logging
capability to also send output to the syslog server running on the cpm-efk container. 

To do this, we modify the postgresql.conf file in the following ways:
~~~~~~~
log_destination = 'stderr,syslog'
syslog_facility = 'LOCAL0'
syslog_ident = 'postgres'
~~~~~~~

If the syslog daemon is not started or configured within the container by choice, your postgres
logs are still stored in the /pgdata/pg_log directory of each container, this is what pgbadger runs
against so it will always be present.

