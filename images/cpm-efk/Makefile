
OSFLAVOR=centos7

build:
	sudo docker build -t cpm-efk -f Dockerfile.$(OSFLAVOR) .
	sudo docker tag -f cpm-efk:latest crunchydata/cpm-efk

download:
	wget https://download.elastic.co/kibana/kibana/kibana-4.3.0-linux-x64.tar.gz
	mv kibana-4.3.0-linux-x64.tar.gz kibana.tar.gz
	wget https://download.elastic.co/elasticsearch/elasticsearch/elasticsearch-2.1.0.tar.gz
	tar xzf elasticsearch-2.1.0.tar.gz
	mv elasticsearch-2.1.0  elasticsearch
	tar czvf elasticsearch.tar.gz ./elasticsearch


