FROM cpm-base:latest
MAINTAINER crunchy

ADD bin /cluster/bin

USER postgres

CMD ["/cluster/bin/start-backupjob.sh"]
