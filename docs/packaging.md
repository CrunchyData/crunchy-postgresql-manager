Packaging
=================


Compile
-------
compile all the CPM source using
~~~~~~~~~~~~~
make build
~~~~~~~~~~~~~~

Build Images
--------------------
~~~~~~~~~~~~~
make buildimages
~~~~~~~~~~~~~~

Pushing Images to DockerHub
-------------------------------
You have to tag the images you want to push, find the image tag, then
run these commands, currently this is a manual step:
~~~~~~~~~~~~~~~~~~~~~~~~~
sudo docker tag -f c91ba0d8cb98 crunchydata/cpm-dashboard:0.9.3
sudo docker push crunchydata/cpm-dashboard:0.9.3
~~~~~~~~~~~~~~~~~~~~~~~~~

Build Archives
-------------------------------
run the sbin/basic-user-install-package.sh script, it will create an archive
file for CPM, then upload it to the S3 site:
~~~~~~~~~~~~~~~~~
https://s3.amazonaws.com/crunchydata/cpm/cpm.0.9.3-linux-amd64.tar.gz
~~~~~~~~~~~~~~~~~

