.PHONY: build launch
build:
	go build -o ./out

launch: build
	./out
