VERSION=0.2.0
PROJROOT=$(HOME)/cpm
GOPATH=$(HOME)/cpm

build:
	        rm -rf bin
	        mkdir bin
	        cp ../../sql/clusteradmin.sql  ./bin/setup.sql
	        cp ../../sql/*  ./bin
	        cp ../../sbin/startadmin.sh ./bin
	        cp ../../bin/adminapi  ./bin
	        cp ../../bin/cpmagentserver  ./bin
	        cp ../../bin/dummyclient  ./bin
	        cp ../../bin/dummyserver  ./bin
	        docker build -t cpm-admin .

