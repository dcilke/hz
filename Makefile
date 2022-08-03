.PHONY:all
all: build

.PHONY=build
build:
	@go build -o bin/hz

.PHONY=install
install: build
	@go install