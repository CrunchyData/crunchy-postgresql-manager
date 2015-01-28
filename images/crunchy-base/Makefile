VERSION=0.2.0

build:
		rm -rf bin
		mkdir bin
		cp ../../config/.bash_profile conf/
		cp ../../config/cpm.sh bin/
		cp ../../sql/node-setup.sql bin/setup.sql
		cp ../../sql/loadtestfunc.sql bin/
		cp ../../sbin/* bin/
		cp ../../config/postgresql.service bin/
	        docker build -t crunchy-base .

