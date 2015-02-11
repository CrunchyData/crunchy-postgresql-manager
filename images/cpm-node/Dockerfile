FROM cpm-base:latest
MAINTAINER crunchy

# open up the cpmagent port
EXPOSE 13000

ADD bin /opt/cpm/bin
ADD conf /opt/cpm/conf

USER postgres

CMD ["/opt/cpm/bin/start-cpmagentserver.sh"]
