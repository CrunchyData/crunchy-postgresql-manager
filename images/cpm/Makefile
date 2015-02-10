VERSION=0.2.0
PROJROOT=$(HOME)/cpm
GOPATH=$(HOME)/cpm
OSFLAVOR=centos7

all:
ifeq ($(OSFLAVOR),centos7)
	cp Dockerfile.centos7 Dockerfile
else
	cp Dockerfile.rhel7 Dockerfile
endif
	        docker build -t cpm .
