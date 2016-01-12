docker ps -q > /tmp/runninglist
docker stop `cat /tmp/runninglist`
docker ps -a -q > /tmp/containerlist
docker rm -f `cat /tmp/containerlist`
docker images -q > /tmp/imagelist
docker rmi -f `cat /tmp/imagelist`
