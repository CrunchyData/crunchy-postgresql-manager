#!/bin/bash


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

su - postgres -c '/cluster/bin/start-pg-wrapper.sh'
#
# start up the cpm agent in the foreground for a rhel 65 startup
#

export PATH=$PATH:/cluster/bin
#
/etc/init.d/sshd start
/cluster/bin/start-cpmagentserver.sh &> /tmp/cpmagentserver.log 
