Crunchy Postgresql Manager

Build Requirements
==================

Building CPM requires development tools like the GCC compiler, along with
the Go language.  On Fedora and RedHat Linux distributions, those packages
can be installed as root like this:

    yum install -y gcc
    yum install -y golang

Obtaining some updates to the CPM code may require the Git package manager.
It can be installed on Fedora/RH with this command:

    yum install -y git

CPM also requires the Docker program is installed, running, and will stay
running after a restart:

    yum install -y docker-io
    systemctl start docker
    systemctl enable docker

The user who is building will need to be part of the docker group
to issue docker comments.  Run this command as root, substituting
build userid in for the one at the end of the line:

    usermod -a -G docker userid

You will need to logout and login again as that user for this
change to be useful.

You can confirm that Docker is available to the user you're building as
by running its info command:

    docker info

Installation
============

Install the dnsbridge program before installing this one.

Build and install CPM by running the install.sh script
