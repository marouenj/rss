#!/bin/bash

docker run \
       --rm \
       --name go \
       -v $(pwd):/go \
golang:1.6 \
go test ./...
