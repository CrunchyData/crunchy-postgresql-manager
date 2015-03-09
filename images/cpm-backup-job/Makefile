
build:
		mkdir bin
		cp ../../sbin/cpm.sh ./bin
		cp ../../bin/backupjob ./bin
		cp ../../bin/dummyserver ./bin
	        docker build -t cpm-backup-job .
		docker tag cpm-backup-job:latest crunchydata/cpm-backup-job

