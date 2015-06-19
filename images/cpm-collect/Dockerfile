FROM cpm-base:latest
MAINTAINER crunchy

# open up the monitor server port
EXPOSE 8080

USER root

ADD bin /var/cpm/bin

CMD ["/var/cpm/bin/start-collect.sh"]
