FROM rhel7
MAINTAINER crunchy

RUN rpm -Uvh http://nginx.org/packages/rhel/7/noarch/RPMS/nginx-release-rhel-7-0.el7.ngx.noarch.rpm
RUN yum install -y procps-ng nginx which hostname wget

# set up cluster directory
#
RUN mkdir -p /cluster/conf
ADD conf /cluster/conf

WORKDIR /cluster

RUN wget http://grafanarel.s3.amazonaws.com/grafana-1.9.1.tar.gz
RUN gunzip grafana-1.9.1.tar.gz
RUN tar -xf grafana-1.9.1.tar.gz
ADD conf/config.js grafana-1.9.1/

RUN chown -R daemon:daemon /cluster

USER daemon

CMD ["/cluster/bin/startnginx.sh"]
