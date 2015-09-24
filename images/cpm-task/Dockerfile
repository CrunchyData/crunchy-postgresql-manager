FROM cpm-base:latest
MAINTAINER crunchy

EXPOSE 13001

ADD bin /var/cpm/bin
ADD conf /var/cpm/conf

USER root

CMD ["/var/cpm/bin/start-taskserver.sh"]
