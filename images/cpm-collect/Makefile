
build:
		rm -rf bin
		mkdir bin
		cp ../../../../../../bin/collectserver ./bin
		cp ../../sbin/start-collect.sh ./bin
	        sudo docker build -t cpm-collect .
		sudo docker tag -f cpm-collect:latest crunchydata/cpm-collect

