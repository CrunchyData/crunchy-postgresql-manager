OSFLAVOR=centos7

all:
	        sudo docker build -t cpm -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm:latest crunchydata/cpm
