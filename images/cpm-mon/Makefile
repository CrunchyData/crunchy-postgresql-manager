
build:
		rm -rf bin
		mkdir bin
		cp ../../../../../../bin/monserver ./bin
		cp ../../sbin/start-mon.sh ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        sudo docker build -t cpm-mon .
		sudo docker tag -f cpm-mon:latest crunchydata/cpm-mon

