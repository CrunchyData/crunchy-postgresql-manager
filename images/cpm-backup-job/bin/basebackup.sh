#!/bin/bash
#
# $1 is the host we are going to do the backup from
#

source /opt/cpm/bin/setenv.sh

# sleep to give DNS time to register the backup job
sleep 7

pg_basebackup -X stream -R --pgdata /pgdata --host=$1 --port=5432 -U postgres > /tmp/basebackup.log 2> /tmp/basebackup.err
