
OpenShift 3 Alpha Notes
=========================

here are some notes on setting up the Docker based Openshift 3 alpha
environment I use.  I am working currently on centos7 VMs.

Infrastructure
---------------
I have 3 nodes:

	* registry.crunchy.lab (192.168.56.120)
	* centos7-dev.crunchy.lab (192.168.56.122)
	* openshift.crunchy.lab (192.168.56.121)


##### Step 0 - open ports in firewalld

On the registry host, we open the following ports:
```
firewall-cmd --get-active-zones
firewall-cmd --permanent --zone=public --add-port=5000/tcp
firewall-cmd --permanent --zone=public --add-port=8080/tcp
firewall-cmd --zone=public --list-all
systemctl restart firewall.service
```


##### Step 1 - install docker private registry

On the registry server, I install the docker private registry on it.
I use this to store all the Docker images we care about, namely our
CPM docker images.

Install instructions are found at:
http://the.randomengineer.com/2014/10/12/intro-to-docker-and-private-image-registries/

Using the private registry is documented at:
http://blog.docker.com/2013/07/how-to-use-your-own-registry/

docker login http://registry.crunchy.lab:5000
set up an id and password! then add that to the registry UI so
that it can authenticate to the registry.

I also install a web UI for the registry which listens on port 8080.

```
docker run --privileged  -p 5000:5000 -v /home/jeffmc/registry-storage/:/tmp/registry registry
docker pull  atcol/docker-registry-ui
docker run -d --name=docker-registry-ui -p 8080:8080 -e REG1=http://192.168.56.120:5000/v1/ atcol/docker-registry-ui
```

##### Step 2 - install openshift

I install openshift on it's own server, openshift.crunchy.lab.

Install instructions for openshift are documented here:

https://github.com/openshift/origin

##### Step 3 - set up dev environment

I install all the CPM dev environment on a separate host
called centos7-dev.crunchy.lab

I use this host to build the CPM images.

I am using the centos:centos7 base image found on dockerhub for this
build.


##### Step 4 - push images to registry

After building the CPM images, we push them to the private registry
as follows:

```
docker images
docker tag <imageid> registry.crunchy.lab:5000/crunchy-cpm
docker push registry.crunchy.lab:5000/crunchy-cpm
docker tag <imageid> registry.crunchy.lab:5000/crunchy-admin
docker push registry.crunchy.lab:5000/crunchy-admin
docker tag <imageid> registry.crunchy.lab:5000/crunchy-pgpool
docker push registry.crunchy.lab:5000/crunchy-pgpool
docker tag <imageid> registry.crunchy.lab:5000/crunchy-node
docker push registry.crunchy.lab:5000/crunchy-node
```

You can verify that the image is in the repository by browsing
to http://registry.crunchy.lab:8080/repository/index

##### Step 5 - deploy into openshift

To add a new pod:

```
openshift kube create pods -c ./examples/crunchy-cpm/crunchy-cpm-pod.json
openshift kube list pods
```

To remove a pod:

```
openshift kube list pods
openshift kube delete pods/someid
```

##### steps required to get this to work under kube

I ended up not running systemd inside the containers, this
causes you to HAVE to run the container as 'priviledged' in Docker.

Kube by default doesn't start containers up as 'priviledged'.  There
is a commit that appears to address this as of 2 weeks ago (Sept 2014)
but it also doesn't appear to be in the Openshift fork.  Anyway, I figured
it is not a good idea to have to run as privlidged anyway.

To fix this, I removed systemd, it also is a good practice to not
run sshd in a container so I removed sshd as well.  From now on
you will have to ssh into the Docker host, then use nsenter to
'get into' the running container.

Removing systemd greatly alters the Dockerfile for all images.

Lastly, I found no way to not get errors in the mounted /pgdata
volume from a permission problem, this appeared to be related
to selinux and Docker.  I found that you get permission errors
as you attempt to use the mounted volume no matter what the
permissions of the host volume.  It turns out this is due
to selinux and luckily I found an answer to resolve it
at:

http://stackoverflow.com/questions/24288616/permission-denied-on-accessing-host-directory-in-docker

I followed this path to set the selinux file settings for the
mounted volume and the permission problems went away.

You can temporarily issue

su -c "setenforce 0"
on the host to access or else add an selinux rule by running

chcon -Rt svirt_sandbox_file_t /path/to/volume

I chose the 'chcon' command.  It is my guess at this point that I will
be able to run that command when I provision the volumes under openshift,
if so, we are fine, if not, then I'll probably have to run the
containers as 'priviledged', that is, if Openshift/Kube will allow that!

For Kube deployment:
	cd ~/cpm/images/crunchy-admin/conf
	openshift kube create pods -c ./cpm-admin-pod.json
