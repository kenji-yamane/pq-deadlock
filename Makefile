.PHONY: build build-cs build-all p1 p2 p3

CMD:=./cmd
BIN:=./bin

PROCESS:=process
SIMULATE_CYCLE:=simulate-cycle

build:
	go build -o $(BIN)/$(PROCESS) $(CMD)/$(PROCESS).go

simulate-cycle:
	go build -o $(BIN)/$(SIMULATE_CYCLE) $(CMD)/puppeteer.go

build-all: build simulate-cycle

THREE_SIMULATOR=10004 10003 10002

p1:
	$(BIN)/$(PROCESS) 1 $(THREE_SIMULATOR)

p2:
	$(BIN)/$(PROCESS) 2 $(THREE_SIMULATOR)

p3:
	$(BIN)/$(PROCESS) 3 $(THREE_SIMULATOR)

pup:
	$(BIN)/$(SIMULATE_CYCLE) $(THREE_SIMULATOR)
