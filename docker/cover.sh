#!/bin/bash

PACKAGES=(agent)
PATHS=(./agent)

for IDX in $(seq ${#PATHS[@]});
do
  docker run \
         --rm \
         --name go \
         -v $(pwd):/go \
         -v $(pwd)/resources:/out \
         golang:1.6 \
  go test -coverprofile=/out/${PACKAGES[((IDX - 1))]}.out ${PATHS[((IDX - 1))]}

  sudo sed -i "s/_\/go\///g" ./resources/agent.out

  docker run \
         --rm \
         --name go \
         -v $(pwd):/go/src \
         -v $(pwd)/resources:/out \
         golang:1.6 \
  go tool cover -html=/out/${PACKAGES[((IDX - 1))]}.out -o /out/${PACKAGES[((IDX - 1))]}.html
done
