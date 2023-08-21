SHELL := bash
SOUCES := main.go

.PHONY: install
install: build
	mkdir -p /usr/local/share/dripdrop
	cp -f ./resource/rain_loop.flac /usr/local/share/dripdrop/rain_loop.flac
	cp -f ./dripdrop /bin/dripdrop

.PHONY: build
build: $(SOURCES)
	go build

.PHONY: clean
clean:
	rm -f ./dripdrop
