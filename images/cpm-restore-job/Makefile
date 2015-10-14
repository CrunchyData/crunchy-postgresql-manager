OSFLAVOR=centos7

build:
		cp ../../sbin/setenv.sh ./bin
		cp $(GOBIN)/restorecommand ./bin
#		cp $(GOBIN)/dummyserver ./bin
	        sudo docker build -t cpm-restore-job -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-restore-job:latest crunchydata/cpm-restore-job

