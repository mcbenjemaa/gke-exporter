#!/bin/bash

VERSION=${VERSION:-"latest"}
REPO="medchiheb/gke-info-exporter"

docker build \
--rm \
--file=Dockerfile \
--tag=${REPO}:${VERSION} \
.

docker push ${REPO}:${VERSION}
