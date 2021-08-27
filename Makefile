VERSION = $(shell git branch --show-current)

.PHONY: help
help:  ## show this help
	@echo "usage: make [target]"
	@echo ""
	@egrep "^(.+)\:\ .*##\ (.+)" ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

.PHONY: run
run: ## run it will instance server 
	VERSION=$(VERSION) go run main.go

.PHONY: run-watch
run-watch: ## run-watch it will instance server with reload
	VERSION=$(VERSION) nodemon --exec go run main.go --signal SIGTERM

.PHONY: test
test: ## runing unit tests with covarage
	go test -coverprofile=coverage.out ./...

.PHONY: mock
mock: ## mock is command to generate mock using mockgen
	rm -rf ./mocks

	mockgen -source=./store/health/health.go -destination=./mocks/health_mock.go -package=mocks -mock_names=Store=MockHealthStore
	mockgen -source=./util/cache/cache.go -destination=./mocks/cache_mock.go -package=mocks

.PHONY: docs
docs: ## docs is a command to generate doc with swagger
	swag init
