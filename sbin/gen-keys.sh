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


KEYSDIR=/var/cpm/keys
sudo mkdir -p $KEYSDIR
sudo chown -R postgres:postgres $KEYSDIR
#
# generate self signed cert for cpm-admin's adminapi REST service
#
openssl genrsa 2048 > key.pem
openssl req -new -x509 -key key.pem -out cert.pem -days 1095
chmod +r key.pem
chmod +r cert.pem
sudo cp key.pem $KEYSDIR
sudo cp cert.pem $KEYSDIR
sudo chown -R postgres:postgres $KEYSDIR

#
# generate self signed cert for cpm's nginx web server
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout nginx.key -out nginx.crt
chmod +r nginx.key
chmod +r nginx.crt
sudo cp nginx.key $KEYSDIR
sudo cp nginx.crt $KEYSDIR

