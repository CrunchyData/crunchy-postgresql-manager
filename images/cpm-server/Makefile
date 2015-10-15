OSFLAVOR=centos7
build:
	        cp $(GOBIN)/cpmserverapi  ./bin
	        cp $(GOBIN)/dummyserver  ./bin
	        sudo docker build -t foo -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f foo:latest crunchydata/foo

