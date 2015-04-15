
OSFLAVOR=centos7

all:
ifeq ($(OSFLAVOR),centos7)
	cp Dockerfile.centos7 Dockerfile
else
	cp Dockerfile.rhel7 Dockerfile
endif
		rm -rf bin
		mkdir bin
		cp ../../config/.bash_profile conf/
		cp ../../sql/node-setup.sql bin/setup.sql
		cp ../../sbin/* bin/
		cp ../../config/postgresql.service bin/
	        sudo docker build -t cpm-base .
		sudo docker tag -f cpm-base:latest crunchydata/cpm-base

