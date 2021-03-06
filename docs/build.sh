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

asciidoc \
-b bootstrap \
-f ./demo.conf \
-o ./htmldoc/doc.html \
-a toc2 \
-a toc-placement=right \
-a theme=cerulean \
./index.asciidoc

asciidoc \
-b bootstrap \
-f ./demo.conf \
-o ./htmldoc/rest-api.html \
-a toc2 \
-a toc-placement=right \
-a theme=cerulean \
./rest-api.asciidoc

asciidoc \
-b bootstrap \
-f ./demo.conf \
-o ./htmldoc/multi-host-setup.html \
-a toc2 \
-a toc-placement=right \
-a theme=cerulean \
./multi-host-setup.asciidoc

asciidoc \
-b bootstrap \
-f ./demo.conf \
-o ./htmldoc/swarm-setup.html \
-a toc2 \
-a toc-placement=right \
-a theme=cerulean \
./swarm-setup.asciidoc

asciidoc \
-b bootstrap \
-f ./demo.conf \
-o ./htmldoc/logging.html \
-a toc2 \
-a toc-placement=right \
-a theme=cerulean \
./logging.asciidoc

asciidoc \
-b bootstrap \
-f ./demo.conf \
-o ./htmldoc/user-guide.html \
-a toc2 \
-a toc-placement=right \
-a theme=cerulean \
./user-guide.asciidoc
