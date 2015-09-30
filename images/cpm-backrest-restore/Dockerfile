FROM cpm-base:latest
MAINTAINER crunchy

ADD bin /var/cpm/bin
ADD conf /var/cpm/conf
ADD backrest/bin  /usr/local/bin
ADD backrest/lib  /usr/local/lib

VOLUME ["/keys"]
VOLUME ["/pgdata"]

USER postgres

CMD ["/var/cpm/bin/start-restore.sh"]
