FROM centos:centos7
MAINTAINER crunchy

#
# this is the stuff that made systemd work on fedora 20
#
RUN yum -y update; yum clean all

# Install postgresql deps
# set up cluster directory
#
RUN mkdir -p /cluster/bin
RUN chown -R daemon:daemon /cluster/bin
RUN rpm -Uvh http://yum.postgresql.org/9.3/redhat/rhel-7-x86_64/pgdg-redhat93-9.3-1.noarch.rpm
RUN yum install -y procps-ng postgresql93 postgresql93-contrib pgpool-II-93 libxslt which unzip openssh-clients hostname


# set environment vars

# open up the pgpool port
EXPOSE 5432
# open up the cpmagentserver port
EXPOSE 13000

ADD bin /cluster/bin
ADD conf/pgpool /cluster/conf/pgpool
ADD conf/pgpool/pool_hba.conf  /etc/pgpool-II-93/pool_hba.conf
ADD conf/pgpool/pool_passwd /etc/pgpool-II-93/pool_passwd

USER daemon

CMD ["/cluster/bin/start-cpmagentserver.sh"]
