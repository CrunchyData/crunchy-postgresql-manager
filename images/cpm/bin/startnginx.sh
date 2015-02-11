#!/bin/bash
/usr/sbin/nginx -c /opt/cpm/conf/nginx.conf > /cpmlogs/nginx.log 2> /cpmlogs/nginx.err
