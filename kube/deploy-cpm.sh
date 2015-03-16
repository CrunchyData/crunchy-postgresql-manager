#!/bin/bash

#
# deploy cpm
#
osc process -n cpm-project -f cpm-template.json   | osc create -n cpm-project -f -


#
# deploy cpm-admin
#
osc process -n cpm-project -f cpm-admin-template.json   | osc create -n cpm-project -f -


#
# deploy cpm-backup
#
osc process -n cpm-project -f cpm-backup-template.json   | osc create -n cpm-project -f -


#
# deploy cpm-mon
#
osc process -n cpm-project -f cpm-mon-template.json   | osc create -n cpm-project -f -
