VERSION=0.2.0

build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/reload-cpmagent.sh ./bin
		cp ../../sbin/start-cpmagent.sh ./bin
		cp ../../sbin/stop-cpmagent.sh ./bin
		cp ../../bin/cpmagentclient ./bin
		cp ../../bin/cpmagentserver ./bin
		cp ../../config/cpmagent.service ./bin
	        docker build -t cpm-node .

