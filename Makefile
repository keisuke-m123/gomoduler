.PHONY: init
init:
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: check
check:
	goimports -w .
	go vet ./...
	staticcheck ./...
