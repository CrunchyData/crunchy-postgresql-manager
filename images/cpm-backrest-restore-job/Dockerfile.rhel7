FROM rhel7
MAINTAINER crunchy

RUN rpm -Uvh http://dl.fedoraproject.org/pub/epel/7/x86_64/e/epel-release-7-5.noarch.rpm
RUN rpm -Uvh http://yum.postgresql.org/9.4/redhat/rhel-7-x86_64/pgdg-redhat94-9.4-1.noarch.rpm
RUN yum install -y procps-ng postgresql94 libxslt unzip openssh-clients hostname bind-utils  && yum clean all -y

RUN mkdir -p /var/cpm/bin
RUN mkdir -p /var/cpm/conf

ENV PGROOT /usr/pgsql-9.4
ENV PGDATA /pgdata

# add path settings for postgres user
ADD conf/.bash_profile /var/lib/pgsql/

VOLUME ["/pgdata"]

VOLUME ["/cpmlogs"]

ADD sbin /var/cpm/bin
ADD bin /var/cpm/bin
ADD conf /var/cpm/conf

USER root

CMD ["/var/cpm/bin/start-backrest-restore-job.sh"]
