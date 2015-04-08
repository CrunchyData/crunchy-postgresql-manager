
build:
		mkdir -p bin
		cp ../../sbin/start-backupjob.sh ./bin
		cp ../../sbin/basebackup.sh ./bin
		cp ../../sbin/cpm.sh ./bin
		cp ../../../../../../bin/backupcommand ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        docker build -t cpm-backup-job .
		docker tag cpm-backup-job:latest crunchydata/cpm-backup-job

