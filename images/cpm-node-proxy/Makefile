
build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/start-node-proxy.sh ./bin
		cp ../../sbin/badger-generate.sh ./bin
		cp ../../../../../../bin/cpmcontainerapi ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        sudo docker build -t cpm-node-proxy .
		sudo docker tag -f cpm-node-proxy:latest crunchydata/cpm-node-proxy

