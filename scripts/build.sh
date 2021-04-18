#!/bin/sh

GIT_HASH=$(git log --oneline -1 | awk '{print $1}')
IMAGE="cloudson/polyglot"
TAG="${IMAGE}:${GIT_HASH}"

docker build . -t $TAG  
docker push $TAG