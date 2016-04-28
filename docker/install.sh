#!/bin/bash

docker run \
       --rm \
       --name go \
       -v $(pwd)/rss.go:/go/src/github.com/marouenj/rss/rss.go \
       -v $(pwd)/agent:/go/src/github.com/marouenj/rss/agent \
       -v $(pwd)/util:/go/src/github.com/marouenj/rss/util \
       -v $(pwd):/go/bin \
golang:1.6 \
go install github.com/marouenj/rss
