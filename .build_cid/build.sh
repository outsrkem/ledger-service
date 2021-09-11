#!/bin/bash
# The build script
# 2021-05-09 20:02:15 CST

PROJECT_PATH=$(cd `dirname $0`/..; pwd)
cd $PROJECT_PATH
version=$(cat version |tr -d ' \t\r')
images=harbor.hub.com/library/ledger-service:"$version"
docker build --build-arg version="$version" -t "$images" . -f .build_cid/Dockerfile
