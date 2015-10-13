OSFLAVOR=centos7
build:
		cp $(GOBIN)/backupcommand ./bin
#		cp $(GOBIN)/dummyserver ./bin
		cp ../../sbin/setenv.sh ./bin
	        sudo docker build -t cpm-backup-job -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-backup-job:latest crunchydata/cpm-backup-job

