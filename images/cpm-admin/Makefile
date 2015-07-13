
build:
	        rm -rf bin
	        mkdir bin
	        cp ../../sql/clusteradmin.sql  ./bin/setup.sql
	        cp ../../sql/*  ./bin
	        cp ../../sbin/startadmin.sh ./bin
	        cp ../../../../../../bin/adminapi  ./bin
	        cp ../../../../../../bin/dummyserver  ./bin
	        sudo docker build -t cpm-admin .
		sudo docker tag -f cpm-admin:latest crunchydata/cpm-admin

