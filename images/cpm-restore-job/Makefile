
build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/start-restore-job.sh ./bin
		cp ../../sbin/setupssh.sh ./bin
		cp ../../../../../../bin/restorecommand ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        sudo docker build -t cpm-restore-job .
		sudo docker tag -f cpm-restore-job:latest crunchydata/cpm-restore-job

