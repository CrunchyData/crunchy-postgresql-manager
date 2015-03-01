
build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/start-backupserver.sh ./bin
		cp ../../bin/backupserver ./bin
		cp ../../bin/dummyserver ./bin
	        docker build -t cpm-backup .
		docker tag cpm-backup:latest crunchydata/cpm-backup

