.PHONY: echo-build
echo-build:
	go build main.go handler.go resp.go

.PHONY: echo-run
echo-run:
	./main

echo-up: echo-build echo-run
