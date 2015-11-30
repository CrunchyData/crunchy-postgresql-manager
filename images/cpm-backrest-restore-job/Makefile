OSFLAVOR=centos7

build:
		cp ../../sbin/setenv.sh ./bin
		cp $(GOBIN)/backrestrestorecommand ./bin
#		cp $(GOBIN)/dummyserver ./bin
	        sudo docker build -t cpm-backrest-restore-job -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-backrest-restore-job:latest crunchydata/cpm-backrest-restore-job

