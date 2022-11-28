.PHONY: build build-cs build-all p1 p2 p3

CMD:=./cmd
BIN:=./bin

PROCESS:=process
PUPPETEER:=puppeteer

build:
	go build -o $(BIN)/$(PROCESS) $(CMD)/$(PROCESS).go

build-pup:
	go build -o $(BIN)/$(PUPPETEER) $(CMD)/puppeteer.go

build-all: build build-pup

THREE_SIMULATOR=10004 10003 10002

FOUR_SIMULATOR=10005 10004 10003 10002

p1-3:
	$(BIN)/$(PROCESS) 1 $(THREE_SIMULATOR)

p2-3:
	$(BIN)/$(PROCESS) 2 $(THREE_SIMULATOR)

p3-3:
	$(BIN)/$(PROCESS) 3 $(THREE_SIMULATOR)

p1-4:
	$(BIN)/$(PROCESS) 1 $(FOUR_SIMULATOR)

p2-4:
	$(BIN)/$(PROCESS) 2 $(FOUR_SIMULATOR)

p3-4:
	$(BIN)/$(PROCESS) 3 $(FOUR_SIMULATOR)

p4-4:
	$(BIN)/$(PROCESS) 4 $(FOUR_SIMULATOR)

pup-cycle:
	$(BIN)/$(PUPPETEER) ./fxt/cycle-three.txt $(THREE_SIMULATOR)

pup-hirata-and:
	$(BIN)/$(PUPPETEER) ./fxt/hirata-wfg-4-p-and.txt $(FOUR_SIMULATOR)

pup-hirata-or:
	$(BIN)/$(PUPPETEER) ./fxt/hirata-wfg-4-p-or.txt $(FOUR_SIMULATOR)
