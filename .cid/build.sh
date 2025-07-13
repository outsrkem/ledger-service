#!/bin/bash
# The build script
# 2021-05-09 20:02:15 CST

workspace=$(cd `dirname $0`/..; pwd)
cd $workspace
set -x

app=ledger
version=${1:-build_0}
commit=`git rev-parse HEAD`

repository='swr.cn-north-1.myhuaweicloud.com/onge'

docker build -f $workspace/Dockerfile \
--build-arg APP_NAME=${app} \
--build-arg APP_REVISION=${commit} \
--build-arg APP_VERSION=${version} \
--label org.opencontainers.image.revision=${commit} \
--label org.opencontainers.image.version=${version} \
--tag ${repository}/${app}:${version} .

docker push ${repository}/${app}:${version}
