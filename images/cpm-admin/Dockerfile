FROM cpm-base:latest
MAINTAINER crunchy

ADD bin /var/cpm/bin/
ADD conf /var/cpm/conf/

# the port for the adminapi REST API
EXPOSE 13000

VOLUME ["/kubekeys"]

USER postgres

CMD ["/var/cpm/bin/startadmin.sh"]
