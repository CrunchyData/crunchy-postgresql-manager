FROM cpm-base:latest
MAINTAINER crunchy

ADD bin /var/cpm/bin

USER postgres

CMD ["/var/cpm/bin/start-backupjob.sh"]
