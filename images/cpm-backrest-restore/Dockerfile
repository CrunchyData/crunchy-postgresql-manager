FROM cpm-base:latest
MAINTAINER crunchy

ADD bin /var/cpm/bin
ADD conf /var/cpm/conf

VOLUME ["/keys"]
VOLUME ["/pgdata"]

USER root

CMD ["/var/cpm/bin/start-restore.sh"]
