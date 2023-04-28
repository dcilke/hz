.PHONY:all
all: build

.PHONY=build
build:
	@go build -o bin/hz

.PHONY=build-test
build-test: ## Build hz with coverage
	@go build -cover -o ./bin/test/hz

.PHONY=install
install: build
	@go install

.PHONY=test
test: build-test
	@rm -rf .covdata
	@mkdir .covdata
	@go test -race ./... -timeout=30s -coverprofile=.covdata/coverage.out -covermode=atomic
	@go tool cover -html=.covdata/coverage.out -o .covdata/coverage.html
	@go tool cover -func=.covdata/coverage.out

.PHONY=format
format:
	gofmt -l -s -w .
	git status
	git diff --exit-code
