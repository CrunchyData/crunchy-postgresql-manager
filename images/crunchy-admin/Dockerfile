FROM crunchy-base:latest
MAINTAINER crunchy

ADD bin /cluster/bin/
ADD conf /cluster/conf/

# the port for the adminapi REST API
EXPOSE 8080

USER postgres

CMD ["/cluster/bin/startadmin.sh"]
