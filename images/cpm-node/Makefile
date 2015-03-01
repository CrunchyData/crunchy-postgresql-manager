
build:
		rm -rf bin
		mkdir bin
		cp ../../sbin/reload-cpmagent.sh ./bin
		cp ../../sbin/start-cpmagent.sh ./bin
		cp ../../sbin/stop-cpmagent.sh ./bin
		cp ../../sbin/fail-over.sh ./bin
		cp ../../bin/cpmagentclient ./bin
		cp ../../bin/cpmagentserver ./bin
		cp ../../config/cpmagent.service ./bin
	        docker build -t cpm-node .
		docker tag cpm-node:latest crunchydata/cpm-node

