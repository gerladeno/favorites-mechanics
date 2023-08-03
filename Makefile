dev-tools:
	go install github.com/daixiang0/gci@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/kazhuravlev/options-gen/cmd/options-gen@latest

gen:
	go generate ./...

lint:
	go mod tidy
	gofumpt -w .
	gci write --custom-order -s standard -s default -s "prefix(github.com/gerladeno/favorites-mechanics)" .
	golangci-lint run ./...

test:
	go test -v -count 10 -race -coverprofile coverage ./...

#build:


do: gen lint test