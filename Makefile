echo-build:
	go build main.go handler.go resp.go

echo-run:
	./main

echo-up: echo-build echo-run
