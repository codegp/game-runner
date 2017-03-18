#!/bin/bash
set -e
REPO_ROOT_FROM_GOPATH_SRC="github.com/codegp/game-runner"
REPO_ROOT="$GOPATH/src/$REPO_ROOT_FROM_GOPATH_SRC"

echo "Generating thrift code ..."
pushd $REPO_ROOT
thrift -r -out $REPO_ROOT --gen go:package_prefix=$REPO_ROOT_FROM_GOPATH_SRC/ thrift/gameObjects.thrift
thrift -r -out $REPO_ROOT --gen go:package_prefix=$REPO_ROOT_FROM_GOPATH_SRC/ thrift/turnInformer.thrift
popd
