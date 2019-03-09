GOBIN=/usr/lib/go-1.9/bin/go
CUR_DIR:=$(shell pwd)
BASE_DIR:=`dirname ${CUR_DIR}`
USER:=$(shell whoami)
APP_NAME:="budao-server"
GOPATH:=${CUR_DIR}:${BASE_DIR}/toolhub:${BASE_DIR}/toolhub/_3rd:/home/${USER}/:/data/gopath/

all:
	@rm -f ${APP_NAME}
	export GOPATH=${GOPATH} && \
	$(GOBIN) build -o ${APP_NAME} ./main.go

transfer:
	export GOPATH=${GOPATH} \
	$(GOBIN) build -o transfer-server ./main.go

test_db:
	export GOPATH=${GOPATH} \
	cd ./db &&  \
	$(GOBIN) test -v
