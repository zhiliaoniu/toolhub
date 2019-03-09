#!/bin/sh

CURDIR=`dirname $0`
cd $CURDIR
CURDIR=`pwd`
GITVERSION=`/usr/bin/git rev-list --all|wc -l`
SERVER_NAME=budao-server
PKG_PATH=${CURDIR}/${SERVER_NAME}-${GITVERSION}

git pull

rm -rf ${CURDIR}/budao-server
gofmt -w ${CURDIR}/service/* ${CURDIR}/main.go \
        go build -o bin/budao-server ${CURDIR}/main.go

mkdir -p ${PKG_PATH}/bin
mkdir -p ${PKG_PATH}/etc

cp etc/budao-server.json     ${PKG_PATH}/etc
cp etc/location.tally.conf   ${PKG_PATH}/etc
cp -fr bin/*                 ${PKG_PATH}/bin

tar -zcvf ${SERVER_NAME}-${GITVERSION}.tar.gz ${SERVER_NAME}-${GITVERSION}/

