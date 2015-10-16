OSFLAVOR=centos7
build:
	        cp $(GOBIN)/cpmserverapi  ./bin
	        cp $(GOBIN)/dummyserver  ./bin
	        sudo docker build -t cpm-server -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-server:latest crunchydata/cpm-server

