#!/bin/bash
asciidoc \
-b bootstrap \
-f ./demo.conf \
-o doc.html \
-a toc2 \
-a toc-placement=right \
-a theme=cerulean \
./demo.asciidoc
