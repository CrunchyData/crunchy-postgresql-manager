#!/bin/bash
/usr/sbin/nginx -c /cluster/conf/nginx.conf > /cpmlogs/nginx.log 2> /cpmlogs/nginx.err
