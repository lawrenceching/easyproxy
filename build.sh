# /bin/bash

PATCH_VERSION=$(cat ./VERSION)
VERSION=0.0.$PATCH_VERSION
OS=`uname`
echo "Building easyproxy $VERSION"
go build
DIST_NAME=easyproxy-$OS.zip
zip $DIST_NAME easyproxy
echo "Built successfully and the package is located in $(pwd)/$DIST_NAME"

let NEXT_VERSION=++PATCH_VERSION
echo $NEXT_VERSION > ./VERSION
echo "Prepare next release $NEXT_VERSION"