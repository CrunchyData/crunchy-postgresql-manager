FROM crunchy-base:latest
MAINTAINER crunchy

# open up the cpmagent port
EXPOSE 13000

ADD bin /cluster/bin
ADD conf /cluster/conf

USER postgres

CMD ["/cluster/bin/start-cpmagentserver.sh"]
