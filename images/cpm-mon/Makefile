
build:
		rm -rf bin
		mkdir bin
		cp ../../../../../../bin/monserver ./bin
		cp ../../sbin/start-mon.sh ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        docker build -t cpm-mon .
		docker tag cpm-mon:latest crunchydata/cpm-mon

