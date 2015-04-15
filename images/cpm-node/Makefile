
build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/start-cpmnodeagent.sh ./bin
		cp ../../sbin/fail-over.sh ./bin
		cp ../../../../../../bin/cpmnodeagent ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        sudo docker build -t cpm-node .
		sudo docker tag -f cpm-node:latest crunchydata/cpm-node

