#!/bin/bash

VERSION="latest";

if [[ $1 != "" ]];
then
  VERSION=$1;
fi

docker build -t marouenj/rss:${VERSION} .
