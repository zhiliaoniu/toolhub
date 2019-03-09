GOBIN=/usr/lib/go-1.9/bin/go
CUR_DIR:=$(shell pwd)
BASE_DIR:=`dirname ${CUR_DIR}`
USER:=$(shell whoami)
APP_NAME:="recommend-server"

all:
	@rm -f ${APP_NAME}
	export GOPATH=${CUR_DIR}:$(BASE_DIR)/toolhub:/home/$(USER)/:/data/gopath/ && \
	$(GOBIN) build -o ${APP_NAME} ./main.go
