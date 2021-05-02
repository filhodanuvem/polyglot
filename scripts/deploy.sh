#!/bin/sh

set -x

GIT_HASH=$(git log --oneline -1 | awk '{print $1}')
TEMPLATE_FILE="./resources/kubernetes/deployment-template.yaml"
sed "s/{{RELEASE_TAG}}/${GIT_HASH}/g" $TEMPLATE_FILE >> ./resources/kubernetes/deployment.yaml
rm $TEMPLATE_FILE
kubectl apply -f ./resources/kubernetes/