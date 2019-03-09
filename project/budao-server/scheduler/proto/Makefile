OUTPUT_DIR=$(shell pwd)
PROTOC := /usr/bin/protoc

all : prepare 
	$(PROTOC) --go_out=plugins=grpc:./scheduler --proto_path=./ *.proto
	$(PROTOC) --go_out=plugins=grpc:./task --proto_path=./task ./task/*.proto
	
prepare:
	@rm -rf $(OUTPUT_DIR)/scheduler $(OUTPUT_DIR)/task/*pb.go
	@mkdir -p $(OUTPUT_DIR)/scheduler
