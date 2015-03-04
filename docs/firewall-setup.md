Firewall Configuration
---------------

This is a set of steps that can be followed to allow you 
to enable firewalld (on Centos7 and RHEL 7) and run CPM.



###disable network manager
~~~~~~~~~~~~~~~~~~~~~~~~
systemctl disable NetworkManager.service
systemctl stop NetworkManager.service
~~~~~~~~~~~~~~~~~~~~~~~~

NOTE: I found the following bugzilla ticket which led me to believe
I should turn off NetworkManager in RHEL/CentoS 7.0, I am assuming
this will be fixed in 7.1 but am not sure:

https://bugzilla.redhat.com/show_bug.cgi?id=1098281

###Enable IP Forwarding
~~~~~~~~~~~~~~~~~~~~~~~~
vi /etc/sysctl.conf
net.ipv4.ip_forward=1
~~~~~~~~~~~~~~~~~~~~~~~~

###Open DNS Port
~~~~~~~~~~~~~~~~~~~~~~~~
firewall-cmd --permanent --zone=public --add-service=dns
firewall-cmd --reload
~~~~~~~~~~~~~~~~~~~~~~~~

###Open CPM Web Port
Create a file in /etc/firewalld/services named cpm.xml
~~~~~~~~~~~~~~~~~~~~~~~~
<?xml version="1.0" encoding="utf-8"?>
<service>
<short>cpm</short>
<description>cpm web interface</description>
<port protocol="tcp" port="13000"/>
</service>
~~~~~~~~~~~~~~~~~~~~~~~~
Set selinux permissions:
~~~~~~~~~~~~~~~~~~~~~~~~
chmod 640 cpm.xml
restorecon cpm.xml
~~~~~~~~~~~~~~~~~~~~~~~~

Add the CPM service:
~~~~~~~~~~~~~~~~~~~~~~~~
firewall-cmd --permanent --add-service=cpm
firewall-cmd --reload
~~~~~~~~~~~~~~~~~~~~~~~~

###Open the Postgresql Port
~~~~~~~~~~~~~~~~~~~~~~~~
firewall-cmd --permanent --zone=public --add-service=postgresql
firewall-cmd --reload
~~~~~~~~~~~~~~~~~~~~~~~~

###Allow Masquerading
You might not want to do this if you want finer grained firewall rules!
This does allow all external hosts to route to the PG containers and 
masquerade as the firewall hosts ip address as the source address.
~~~~~~~~~~~~~~~~~~~~~~~~
firewall-cmd --permanent --zone=public --add-masquerade
firewall-cmd --reload
~~~~~~~~~~~~~~~~~~~~~~~~

###Allow External Host Routing
On external hosts, you need to create a static route that allows
them to reach the Docker containers (e.g. 172.18.0.0/16) via the 
Docker host ip (e.g. 192.168.0.103) using a particular eth interface (e.g. enp2s0):
~~~~~~~~~~~~~~~~~~~~~~~~
ip route add 172.18.0.0/16 via 192.168.0.103 dev enp2s0
~~~~~~~~~~~~~~~~~~~~~~~~

###Allow CPM DNS 
To resolve the Docker container host names, assigned by CPM Skybridge, 
you will need to specify on your remote hosts, the CPM DNS host ip address
as the primary DNS nameserver in your /etc/resolv.conf.

