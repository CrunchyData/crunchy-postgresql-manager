#!/bin/bash

#
# deploy cpm
#
osc process -n default -f cpm-template.json   | osc create -n default -f -

#
# deploy cpm-admin
#
osc process -n default -f cpm-admin-template.json   | osc create -n default -f -

sleep 10

#
# deploy cpm-backup
#
osc process -n default -f cpm-backup-template.json   | osc create -n default -f -

sleep 10

#
# deploy cpm-mon
#
osc process -n default -f cpm-mon-template.json   | osc create -n default -f -
