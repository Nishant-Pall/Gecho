.PHONY: gecho-build
gecho-build:
	go build main.go gloom_handler.go handler.go resp.go aof.go

.PHONY: gecho-run
gecho-run:
	./main

gecho-up: gecho-build gecho-run
