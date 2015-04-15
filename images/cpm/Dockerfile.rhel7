FROM rhel7
MAINTAINER crunchy

RUN rpm -Uvh http://nginx.org/packages/rhel/7/noarch/RPMS/nginx-release-rhel-7-0.el7.ngx.noarch.rpm
RUN yum install -y openssl procps-ng nginx which hostname && yum clean all -y

VOLUME ["/www"]
RUN chown -R daemon:daemon /www

VOLUME ["/cpmlogs"]
RUN chown -R daemon:daemon /cpmlogs

VOLUME ["/cpmkeys"]
RUN chown -R daemon:daemon /cpmkeys

EXPOSE 13000

# set up cpm directory
#
RUN mkdir -p /var/cpm/bin 
RUN mkdir -p /var/cpm/conf 
ADD bin /var/cpm/bin
ADD conf /var/cpm/conf
RUN chown -R daemon:daemon /var/cpm

USER daemon

CMD ["/var/cpm/bin/startnginx.sh"]
