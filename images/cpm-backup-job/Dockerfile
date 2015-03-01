FROM cpm-base:latest
MAINTAINER crunchy

ADD bin /opt/cpm/bin

USER postgres

CMD ["/opt/cpm/bin/start-backupjob.sh"]
