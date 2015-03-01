FROM cpm-base:latest
MAINTAINER crunchy

ADD bin /opt/cpm/bin/
ADD conf /opt/cpm/conf/

# the port for the adminapi REST API
EXPOSE 13000

USER postgres

CMD ["/opt/cpm/bin/startadmin.sh"]
