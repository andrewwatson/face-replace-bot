#!/bin/bash

set -e

DATE=`date +'%Y%m%d%H%M%S'`
IMAGE=kenbot
TAG=$DATE
PUSH=makeandbuildatl

echo "BUILDING DOCKER IMAGE $IMAGE:$TAG"
docker build -t $PUSH/$IMAGE:$TAG  .
docker push $PUSH/$IMAGE:$TAG

docker build -t $PUSH/$IMAGE:latest  .
docker push $PUSH/$IMAGE:latest

echo "PUSHED $IMAGE:$TAG to $PUSH"

if [ -n "$ROLL" ]; then
  kubectl set image deployment/replace-the-face botbrains=$PUSH/$IMAGE:$TAG
fi
