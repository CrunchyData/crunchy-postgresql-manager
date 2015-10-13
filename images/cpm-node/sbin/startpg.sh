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

source /var/cpm/bin/setenv.sh

# clean up any leftover locks that pg might have left
rm /var/run/postgresql/.s.*

su - postgres -c 'source /var/cpm/bin/setenv.sh;pg_ctl -w -D /pgdata start 2> /tmp/startpg.err > /tmp/startpg.log'