VERSION=0.2.0
PROJROOT=$(HOME)/cpm
GOPATH=$(HOME)/cpm
build:
	cp ../../sbin/*  ./bin
	cp ../../bin/cpmagentclient  ./bin
	cp ../../bin/cpmagentserver  ./bin
	cp ../../config/cpmagent.service  ./bin
	docker build -t crunchy-pgpool .
