VERSION=0.2.0

build:
		rm -rf bin
		mkdir bin
		cp ../../bin/backupserver ./bin
		cp ../../bin/dummyserver ./bin
	        docker build -t cpm-backup .
		docker tag cpm-backup:latest crunchydata/cpm-backup

