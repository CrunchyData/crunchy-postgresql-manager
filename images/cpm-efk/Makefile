
OSFLAVOR=centos7

build:
	        cp ~/Downloads/kibana-4.3.0-linux-x64.tar.gz .
	        cp ~/Downloads/elasticsearch-2.1.0.tar.gz .
	        sudo docker build -t cpm-efk -f Dockerfile.$(OSFLAVOR) .
		sudo docker tag -f cpm-efk:latest crunchydata/cpm-efk

