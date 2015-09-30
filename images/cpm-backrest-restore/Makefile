
build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/start-restore.sh ./bin
		cp ../../sbin/setupssh.sh ./bin
		cp ../../../../../../bin/backrestrestore ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        sudo docker build -t cpm-backrest-restore .
		sudo docker tag -f cpm-backrest-restore:latest crunchydata/cpm-backrest-restore

