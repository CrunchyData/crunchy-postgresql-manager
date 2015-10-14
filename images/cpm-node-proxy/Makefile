OSFLAVOR=centos7

build:
		cp ../../sbin/setenv.sh ./bin
		cp ../../config/.bash_profile ./bin
		cp $(GOBIN)/cpmcontainerapi ./bin
#		cp $(GOBIN)/dummyserver ./bin
	        sudo docker build -t cpm-node-proxy -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-node-proxy:latest crunchydata/cpm-node-proxy

