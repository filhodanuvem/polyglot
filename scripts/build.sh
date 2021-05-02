#!/bin/sh

set -x

GIT_HASH=$(git log --oneline -1 | awk '{print $1}')
IMAGE="cloudson/polyglot"
TAG="${IMAGE}:${GIT_HASH}"
LATEST="${IMAGE}:latest"

docker build . -t $TAG -t $LATEST
docker push $TAG
docker push $LATEST