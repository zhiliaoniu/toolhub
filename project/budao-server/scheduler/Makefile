GOBIN=/usr/lib/go-1.9/bin/go
CUR_DIR:=$(shell pwd)
BASE_DIR:=`dirname ${CUR_DIR}`
USER:=$(shell whoami)
APP_NAME:="scheduler"
GOPATH:=${CUR_DIR}:$(BASE_DIR)/toolhub:$(BASE_DIR)/toolhub/_3rd:/home/$(USER)/:/data/gopath/

all: prepare
	@rm -f ${APP_NAME}
	export GOPATH=${GOPATH} && \
	$(GOBIN) build -o ${APP_NAME} ${CUR_DIR}/main.go

prepare :

clean :
	@rm -frv build

.PHONY : prepare clean

