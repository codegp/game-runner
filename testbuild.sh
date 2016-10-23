#!/bin/bash
set -e

REPO_ROOT_FROM_GOPATH_SRC="github.com/codegp/game-runner"
REPO_ROOT="$GOPATH/src/$REPO_ROOT_FROM_GOPATH_SRC"

function cleanup {
  rm $REPO_ROOT/testbuild/gameObjects.thrift
  rm $REPO_ROOT/gamerunner/gametypedef.go
  rm -rf $REPO_ROOT/api
  rm  $REPO_ROOT/gamerunner/gamerunner
}
trap cleanup EXIT

cp  $REPO_ROOT/thrift/gameObjects.thrift $REPO_ROOT/testbuild/gameObjects.thrift
cp  $REPO_ROOT/testbuild/testgametype.go.tmpl $REPO_ROOT/gamerunner/gametypedef.go


thrift -r -out $REPO_ROOT --gen go:package_prefix=$REPO_ROOT_FROM_GOPATH_SRC/ $REPO_ROOT/testbuild/api.thrift
pushd $REPO_ROOT/gamerunner
go build
popd
