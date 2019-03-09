GOBIN=/usr/lib/go-1.9/bin/go
CUR_DIR:=$(shell pwd)
BASE_DIR:=`dirname ${CUR_DIR}`
APP_NAME:="budao-cms-server"
GOPATH:=${CUR_DIR}:${BASE_DIR}/toolhub:${BASE_DIR}/toolhub/_3rd:/data/gopath/

all:
	@rm -f ${APP_NAME}
	export GOPATH=${GOPATH} && \
    $(GOBIN) build -o ${APP_NAME} ./main.go
