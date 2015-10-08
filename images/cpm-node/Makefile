
build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/start-cpmcontainerapi.sh ./bin
		cp ../../sbin/start-restore.sh ./bin
		cp ../../sbin/removedata.sh ./bin
		cp ../../sbin/fail-over.sh ./bin
		cp ../../sbin/badger-generate.sh ./bin
		cp ../../../../../../bin/cpmcontainerapi ./bin
		cp ../../../../../../bin/dummyserver ./bin
	        sudo docker build -t cpm-node .
		sudo docker tag -f cpm-node:latest crunchydata/cpm-node

