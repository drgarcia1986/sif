PROJECT_FILES:=$(shell go list ./... | grep -v vendor)

test:
	@go test ${PROJECT_FILES}

coverage:
	@go test --cover ${PROJECT_FILES}

coverage-report:
	@./.tests-utils.sh coverage-report

test-e2e:
	@./.tests-utils.sh e2e

build:
	@go build -o sif cmd/sif/main.go
