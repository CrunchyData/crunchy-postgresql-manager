
#!/bin/bash -x

# Copyright 2015 Crunchy Data Solutions, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# script to copy docker server related files to destination
#
adminserver=espresso.crunchy.lab
remoteservers=(jeff1.crunchy.lab espresso.crunchy.lab)

for i in "${remoteservers[@]}"
do
	echo $i
	ssh root@$i "mkdir -p /opt/cpm/bin"
	scp ./bin/*  \
	./sbin/*  \
	./sql/*  \
	root@$i:/opt/cpm/bin/
	scp  ./config/cpm.sh  root@$i:/etc/profile.d/cpm.sh
	scp  ./config/cpmagent.service  \
	 root@$i:/usr/lib/systemd/system
 	ssh root@$i "systemctl enable docker.service"
        ssh root@$i "systemctl enable cpmagent.service"
done

# copy all required admin files to the admin server

ssh root@$adminserver "mkdir -p /opt/cpm/bin"
scp ./bin/* \
./sbin/* \
root@$adminserver:/opt/cpm/bin

scp ./config/cpmagent.service root@$adminserver:/usr/lib/systemd/system

ssh root@$adminserver "systemctl enable cpmagent.service"


