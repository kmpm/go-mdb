
GO=go

utils := mdb-ver mdb-tables


.PHONY: build
build: ${utils}


.PHONY: ${utils}
${utils}:
	go build ./cmd/$@


.PHONY: test
test:
	go test -v ./...


.PHONY: run
run:
	go run ./cmd/mdb-tables dbCustomers.mdb