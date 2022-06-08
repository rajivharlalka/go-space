.PHONY : build run fresh test clean

test:
	go test

build:
	go build

run:
	./imagery-v2.bin

clean:
	go clean
	rm -f sample.bin