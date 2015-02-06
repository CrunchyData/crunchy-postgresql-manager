FROM cpm-base:latest
MAINTAINER crunchy

# open up the backupserver port
EXPOSE 13010

ADD bin /cluster/bin
ADD conf /cluster/conf

USER root

CMD ["/cluster/bin/start-backupserver.sh"]
