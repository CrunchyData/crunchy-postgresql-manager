FROM cpm-base:latest
MAINTAINER crunchy

# open up the cpmagent port
EXPOSE 13000

ADD bin /var/cpm/bin
ADD conf /var/cpm/conf

USER postgres

CMD ["/var/cpm/bin/start-cpmnodeagent.sh"]
