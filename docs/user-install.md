User Install
=================

Typical Development Environment
-------------------------------
CentOS 7.1 and RHEL 7.1 are supported currently, others might work, especially
Fedora or other RHEL variants but you might see differences in
installation.

Obtain the Media
------------------------------
Download the binary install archives for skybridge and CPM
from the following:
~~~~~~~~~~~~~~~~~
https://s3.amazonaws.com/crunchydata/cpm/cpm.0.9.4-linux-amd64.tar.gz
https://s3.amazonaws.com/crunchydata/cpm/skybridge.1.0.2-linux-amd64.tar.gz
~~~~~~~~~~~~~~~~~

Install
----------
Here is an example of how to perform the installation:

~~~~~~~~~~~~~~~~~
mkdir skybridge-install
cd skybridge-install
tar xvzf ../skybridge.1.0.2-linux-amd64.tar.gz
./install.sh

mkdir cpm-install
cd cpm-install
tar xvzf ../cpm.0.9.4-linux-amd64.tar.gz

./basic-user-install.sh
~~~~~~~~~~~~~~~~~

The DNS installation will enable and configure the Docker service
to specify the DNS server as the primary DNS nameserver.  This
DNS server will also be your primary nameserver in your /etc/resolv.conf
configuration.

You will be prompted for your IP address and the domain name
you want to use for the installation.

Self signed certificates will be generated, this means you will
have to accept the certificate warnings in your browser.  With Firefox
you will need to adjust your Cross-Origin security settings to 
allow access:
https://blog.nraboy.com/2014/08/bypass-cors-errors-testing-apis-locally/


You can generate your own self-signed keys by following this:

http://www.hydrogen18.com/blog/your-own-pki-tls-golang.html

You can avoid browser warnings for the self-signed keys by following this:

http://portal.threatpulse.com/docs/sol/Content/03Solutions/ManagePolicy/SSL/ssl_chrome_cert_ta.htm

Next, run this script to create the CPM containers:
~~~~~~~~~~~~~~~~~
./bu-init-cpm.sh
~~~~~~~~~~~~~~~~~

To manually stop the CPM containers, run this script:
~~~~~~~~~~~~~~~~~
./bu-init-stop.sh
~~~~~~~~~~~~~~~~~

To manually start the CPM containers, run this script:
~~~~~~~~~~~~~~~~~
./bu-init-start.sh
~~~~~~~~~~~~~~~~~

systemd unit files for CPM are found in:
~~~~~~~~~~~~~~~~~
/var/cpm/config/cpm.service
~~~~~~~~~~~~~~~~~

Initial Login
-------------
When you first access CPM, you will receive a Login dialog, on that
page enter 'cpm' for the userid, 'cpm' for the password, and
'https://cpm-admin.example.com:13000' for the Admin URL.

Refer to the CPM User Guide for details on how to use the application, but
generally you would first create a Server.  Your current host is
what you would define as your CPM Server.  Use the Docker Bridge value
of 172.17.42.1, and the PG Data Path is /var/cpm/data/pgsql.

After that, you can create a PG container and PG clusters.


