.PHONY: module
EXE_NAME = nickcam
BIN_OUTPUT_PATH = bin/$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m)

$(BIN_OUTPUT_PATH)/$(EXE_NAME):
	go build -o $(BIN_OUTPUT_PATH)/$(EXE_NAME)

module: 
	go build -o bin/$(EXE_NAME)
	rm -rf module.tar.gz
	tar czf module.tar.gz bin/$(EXE_NAME)
	rm bin/$(EXE_NAME)

clean:
	rm -rf bin module.tar.gz
