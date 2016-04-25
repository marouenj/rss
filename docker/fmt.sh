#!/bin/bash

docker run \
       --rm \
       --name go \
       -v $(pwd):/go/src \
golang:1.6 \
go fmt ./...
