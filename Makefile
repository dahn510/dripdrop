SHELL := bash
SOUCES := main.go

.PHONY: install
install: build
	mkdir -p $(DESTDIR)$(PREFIX)/share/dripdrop
	install -Dm644 ./resource/rain_loop.flac $(DESTDIR)$(PREFIX)/share/dripdrop/rain_loop.flac
	install -Dm755 ./dripdrop $(DESTDIR)$(PREFIX)/dripdrop

.PHONY: build
build: $(SOURCES)
	go build

.PHONY: clean
clean:
	rm -f ./dripdrop
