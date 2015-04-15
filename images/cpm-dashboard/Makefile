OSFLAVOR=centos7

all:
ifeq ($(OSFLAVOR),centos7)
	cp Dockerfile.centos7 Dockerfile
else
	cp Dockerfile.rhel7 Dockerfile
endif
	        sudo docker build -t cpm-dashboard .
		sudo docker tag -f cpm-dashboard:latest crunchydata/cpm-dashboard

