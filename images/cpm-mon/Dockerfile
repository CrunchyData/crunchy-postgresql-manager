FROM cpm-base:latest
MAINTAINER crunchy


VOLUME ["/monitordata"]

# open up the monitor server port
EXPOSE 13000

USER root

ADD bin /opt/cpm/bin
ADD conf /opt/cpm/conf

# get influxdb
RUN rpm -Uvh http://s3.amazonaws.com/influxdb/influxdb-latest-1.x86_64.rpm

RUN chown -R root:root /monitordata

CMD ["/opt/cpm/bin/start-mon.sh"]
