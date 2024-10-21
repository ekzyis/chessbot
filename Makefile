.PHONY: build test

SOURCES := main.go $(shell find chess -name '*.go')

build: chessbot

chessbot: $(SOURCES)
	go build -o chessbot .

test:
	## -count=1 is used to disable cache
	## use -run <regexp> to only run tests that match <regexp>
	go test ./chess -v -count=1
