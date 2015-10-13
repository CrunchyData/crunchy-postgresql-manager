OSFLAVOR=centos7
build:
		cp ../../sbin/setenv.sh ./bin
		cp $(GOBIN)/collectserver ./bin
	        sudo docker build -t cpm-collect -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-collect:latest crunchydata/cpm-collect

