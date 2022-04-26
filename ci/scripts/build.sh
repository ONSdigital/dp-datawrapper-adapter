#!/bin/bash -eux

pushd dp-datawrapper-adapter
  make build
  cp build/dp-datawrapper-adapter Dockerfile.concourse ../build
popd
