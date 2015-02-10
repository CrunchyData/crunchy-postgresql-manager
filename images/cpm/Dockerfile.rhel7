FROM rhel7
MAINTAINER crunchy

RUN rpm -Uvh http://nginx.org/packages/rhel/7/noarch/RPMS/nginx-release-rhel-7-0.el7.ngx.noarch.rpm
RUN yum install -y procps-ng nginx which hostname

VOLUME ["/www"]
RUN chown -R daemon:daemon /www

VOLUME ["/cpmlogs"]
RUN chown -R daemon:daemon /cpmlogs

EXPOSE 10000

# set up cluster directory
#
RUN mkdir -p /cluster/bin 
RUN mkdir -p /cluster/conf 
ADD bin /cluster/bin
ADD conf /cluster/conf
RUN chown -R daemon:daemon /cluster

USER daemon

CMD ["/cluster/bin/startnginx.sh"]
