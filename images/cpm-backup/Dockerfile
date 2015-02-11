FROM cpm-base:latest
MAINTAINER crunchy

# open up the backupserver port
EXPOSE 13010

ADD bin /opt/cpm/bin
ADD conf /opt/cpm/conf

USER root

CMD ["/opt/cpm/bin/start-backupserver.sh"]
