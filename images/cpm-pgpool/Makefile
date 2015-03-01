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
	cp ../../bin/cpmagentclient  ./bin
	cp ../../bin/cpmagentserver  ./bin
	cp ../../config/cpmagent.service  ./bin
	docker build -t cpm-pgpool .
	docker tag cpm-pgpool:latest crunchydata/cpm-pgpool
