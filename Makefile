.PHONY: test
test:
	go test -v internal/test/shopql_test.go

.PHONY: install
install:
	go install github.com/99designs/gqlgen@v0.17.45

.PHONY: init
init:
	go run github.com/99designs/gqlgen init

.PHONY: gen
gen: 
	@echo "-- generatiog graphql files"
	go run github.com/99designs/gqlgen generate --verbose --config configs/gqlgen.yml
