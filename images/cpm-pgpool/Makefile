OSFLAVOR=centos7

all:
	cp ../../sbin/setenv.sh  ./bin
	cp $(GOBIN)/cpmcontainerapi  ./bin
	sudo docker build -t cpm-pgpool -f Dockerfile.$(OSFLAVOR) .
	sudo docker tag -f cpm-pgpool:latest crunchydata/cpm-pgpool
