OSFLAVOR=centos7

all:
ifeq ($(OSFLAVOR),centos7)
		cp Dockerfile.centos7 Dockerfile
else
		cp Dockerfile.rhel7 Dockerfile
endif
	rm -rf bin
	mkdir bin
	cp ../../sbin/*  ./bin
	cp ../../../../../../bin/cpmcontainerapi  ./bin
	sudo docker build -t cpm-pgpool .
	sudo docker tag -f cpm-pgpool:latest crunchydata/cpm-pgpool
