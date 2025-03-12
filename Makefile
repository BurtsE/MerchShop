include .env
export

.PHONY:  integration_test test lint coverage_report cpu_profile mem_profile migrate_up migrate_down create_migration

help:
	cat Makefile

integration_test:
	echo sdf

lint:
	go fmt ./...
    find . -name '*.go' ! -path "./generated/*" -exec goimports -local go-echo-template/ -w {} +
    find . -name '*.go' ! -path "./generated/*" -exec golines -w {} -m 120 \;
    golangci-lint run ./...
    ./check-go-generate.sh
touch