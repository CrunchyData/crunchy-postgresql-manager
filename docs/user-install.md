User Install
=================

Typical Development Environment
-------------------------------
CentOS 7 and RHEL 7 are supported currently, others might work, especially
Fedora or other RHEL variants but you might see differences in
installation.

Obtain the Media
------------------------------
Download the binary install archives for skybridge and CPM
from the following:
http://s3.amazon.com/crunchydata/cpm.beta.tar.gz
http://s3.amazon.com/crunchydata/skybridge.beta.tar.gz

Install
----------
Here is an example of how to perform the installation:

~~~~~~~~~~~~~~~~~
mkdir skybridge-install
cd skybridge-install
tar xvzf ../skybridge.beta.tar.gz
./install.sh

mkdir cpm-install
cd cpm-install
tar xvzf ../cpm.beta.tar.gz
./basic-user-install.sh
~~~~~~~~~~~~~~~~~

You will be prompted for your IP address and the domain name
you want to use for the installation.

The DNS installation will enable and configure the Docker service
to specify the DNS server as the primary DNS nameserver.  This
DNS server will also be your primary nameserver in your /etc/resolv.conf
configuration.

