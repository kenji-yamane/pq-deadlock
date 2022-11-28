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

SEVEN_SIMULATOR=10008 10007 10006 10005 10004 10003 10002

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

p1-7:
	$(BIN)/$(PROCESS) 1 $(SEVEN_SIMULATOR)

p2-7:
	$(BIN)/$(PROCESS) 2 $(SEVEN_SIMULATOR)

p3-7:
	$(BIN)/$(PROCESS) 3 $(SEVEN_SIMULATOR)

p4-7:
	$(BIN)/$(PROCESS) 4 $(SEVEN_SIMULATOR)

p5-7:
	$(BIN)/$(PROCESS) 5 $(SEVEN_SIMULATOR)

p6-7:
	$(BIN)/$(PROCESS) 6 $(SEVEN_SIMULATOR)

p7-7:
	$(BIN)/$(PROCESS) 7 $(SEVEN_SIMULATOR)

pup-cycle:
	$(BIN)/$(PUPPETEER) ./fxt/cycle-three.txt $(THREE_SIMULATOR)

pup-hirata-and:
	$(BIN)/$(PUPPETEER) ./fxt/hirata-wfg-4-p-and.txt $(FOUR_SIMULATOR)

pup-hirata-or:
	$(BIN)/$(PUPPETEER) ./fxt/hirata-wfg-4-p-or.txt $(FOUR_SIMULATOR)

pup-book:
	$(BIN)/$(PUPPETEER) ./fxt/book-wfg.txt $(SEVEN_SIMULATOR)
