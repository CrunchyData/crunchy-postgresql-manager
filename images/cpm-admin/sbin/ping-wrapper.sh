#!/bin/bash
#
# this is a test of a port, cpm uses this to test whether or not
# pg or another service is up...requires nmap-ncat to be installed
#
# $1 is the host name
# $2 is the port number
# returns 0 if connection is established
#
nc -w 1s $1 $2 < /dev/null
retcode=$?
if [ $retcode == 0 ]; then
	echo -n "connected"
else
	echo -n "not-connected"
fi
