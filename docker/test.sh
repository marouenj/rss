#!/bin/bash

docker run \
       --rm \
       --name go \
       -v $(pwd)/agent:/go/src/github.com/marouenj/rss/agent \
       -v $(pwd)/util:/go/src/github.com/marouenj/rss/util \
golang:1.6 \
go test ./...
