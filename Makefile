EXE_NAME = nickcam
BIN_OUTPUT_PATH = bin/$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m)

all:
	rm -rf bin
	go build -o $(BIN_OUTPUT_PATH)/$(EXE_NAME)
