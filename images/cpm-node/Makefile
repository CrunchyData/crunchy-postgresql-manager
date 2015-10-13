
OSFLAVOR=centos7

build:
		cp ../../sbin/setenv.sh ./bin
		cp ../../config/.bash_profile ./bin
		cp $(GOBIN)/cpmcontainerapi ./bin
#		cp $(GOBIN)/dummyserver ./bin
	        sudo docker build -t cpm-node -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-node:latest crunchydata/cpm-node

