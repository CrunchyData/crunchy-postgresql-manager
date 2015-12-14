OSFLAVOR=centos7
build:
	        cp ../../sbin/setenv.sh  ./bin
	        cp ../../config/.bash_profile  ./bin
	        cp ../../sql/clusteradmin.sql  ./sql/setup.sql
	        cp ../../sql/*  ./sql
	        cp $(GOBIN)/adminserver  ./bin
	        sudo docker build -t cpm-admin -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-admin:latest crunchydata/cpm-admin

