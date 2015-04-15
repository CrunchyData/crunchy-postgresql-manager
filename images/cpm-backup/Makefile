
build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/start-backupserver.sh ./bin
		cp ../../../../../../bin/backupserver ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        sudo docker build -t cpm-backup .
		sudo docker tag -f cpm-backup:latest crunchydata/cpm-backup

