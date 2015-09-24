
build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/start-taskserver.sh ./bin
		cp ../../../../../../bin/taskserver ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        sudo docker build -t cpm-task .
		sudo docker tag -f cpm-task:latest crunchydata/cpm-task

