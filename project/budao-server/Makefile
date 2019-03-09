GOBIN=/usr/lib/go-1.9/bin/go

CUR_DIR:=$(shell pwd)
BASE_DIR:=`dirname ${CUR_DIR}`
USER:=$(shell whoami)
GOPATH:=${CUR_DIR}:$(BASE_DIR)/toolhub:$(BASE_DIR)/toolhub/_3rd:/home/$(USER)/:/data/gopath/

all:
	@make -C budao
	@make -C recommend
	@make -C statistic
	@make -C push 
	@make -C share 

test:
