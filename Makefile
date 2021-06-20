BIN=$(PWD)/bin

all: pre
	@cd gbk2utf; go build -v -o $(BIN)

pre:
	@if [ ! -x $(BIN) ]; then mkdir $(BIN); fi