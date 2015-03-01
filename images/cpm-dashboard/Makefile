OSFLAVOR=centos7

all:
ifeq ($(OSFLAVOR),centos7)
	cp Dockerfile.centos7 Dockerfile
else
	cp Dockerfile.rhel7 Dockerfile
endif
	        docker build -t cpm-dashboard .
		docker tag cpm-dashboard:latest crunchydata/cpm-dashboard

