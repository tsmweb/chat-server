#!/bin/bash
image="postgres"
imagelocal="tsmweb/"$image
tag=latest

#build
docker ps | grep $imagelocal && docker stop $(docker ps | grep $imagelocal | awk '{print $1}')
docker image rm -f $imagelocal:$tag
docker build -t $imagelocal:$tag .