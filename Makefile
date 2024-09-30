SHELL := bash
SOURCE := main.go

.PHONY: install
install: $(SOURCE)
	go install

.PHONY: build
build: $(SOURCE)
	go build

.PHONY: clean
clean:
	rm -f ./dripdrop
