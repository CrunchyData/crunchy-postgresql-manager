OSFLAVOR=centos7

all:
ifeq ($(OSFLAVOR),centos7)
	cp Dockerfile.centos7 Dockerfile
else
	cp Dockerfile.rhel7 Dockerfile
endif
	        sudo docker build -t cpm .
		sudo docker tag -f cpm:latest crunchydata/cpm
