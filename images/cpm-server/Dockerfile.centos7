FROM centos:7
MAINTAINER crunchy

RUN rpm -Uvh http://yum.postgresql.org/9.4/redhat/rhel-7-x86_64/pgdg-centos94-9.4-1.noarch.rpm
RUN yum install -y epel-release docker procps-ng postgresql94 postgresql94-server sysstat procps-ng unzip openssh-clients hostname bind-utils && yum clean all -y

RUN mkdir -p /var/cpm/bin

ADD sbin /var/cpm/bin
ADD bin /var/cpm/bin

EXPOSE 10001

USER root

CMD ["/var/cpm/bin/start.sh"]
