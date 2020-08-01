# /bin/bash

OS=`uname`
echo "Building easyproxy $VERSION"
go build
DIST_NAME=easyproxy-$OS.zip
zip $DIST_NAME easyproxy
echo "Built successfully and the package is located in $(pwd)/$DIST_NAME"