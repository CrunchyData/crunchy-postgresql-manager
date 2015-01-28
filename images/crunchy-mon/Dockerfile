FROM crunchy-base:latest
MAINTAINER crunchy


VOLUME ["/monitordata"]

# open up the monitor server port
EXPOSE 13000

USER root

ADD bin /cluster/bin
ADD conf /cluster/conf

# get influxdb
RUN rpm -Uvh http://s3.amazonaws.com/influxdb/influxdb-latest-1.x86_64.rpm

RUN chown -R root:root /monitordata

CMD ["/cluster/bin/start-mon.sh"]
