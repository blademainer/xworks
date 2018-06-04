#!/usr/bin/env bash
arch="${1:-arm}"
SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
${SHELL_FOLDER}/build.sh "$arch"
version="v`date  +"%Y%m%d%H%M%s"`"
REPO=blademainer/xworks:client-${arch}-${version}
echo "repo: $REPO"
docker build -f dockerfile-client -t $REPO .
docker push $REPO
