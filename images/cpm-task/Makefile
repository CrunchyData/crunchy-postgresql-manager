OSFLAVOR=centos7
build:
		cp $(GOBIN)/taskserver ./bin
		cp $(GOBIN)/dummyserver ./bin
		cp ../../sbin/setenv.sh ./bin
	        sudo docker build -t cpm-task -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-task:latest crunchydata/cpm-task

