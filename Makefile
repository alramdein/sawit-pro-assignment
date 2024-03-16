

.PHONY: clean all init generate generate_mocks generate_keys

all: build/main

build/main: cmd/main.go generated generate_keys
	@echo "Building..."
	go build -o $@ $<

clean: clean_keys
	rm -rf generated

init: generate generate_keys
	go mod tidy
	go mod vendor

test:
	go test -short -coverprofile coverage.out -v ./...

generate: generated generate_mocks

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find repository -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

tidy:
	go mod tidy

run: 
	go run cmd/main.go

reinit: clean init

generate_keys:
	@echo "Generating RSA private key..."
	@openssl genrsa -out private.pem 2048
	@echo "Copying private key to PKCS#1 format..."
	@openssl rsa -in private.pem -outform PEM -out private_pkcs1.pem
	@echo "Generating RSA public key..."
	@openssl rsa -in private.pem -pubout -out public.pem

clean_keys:
	@echo "Cleaning up..."
	@rm -f private.pem private_pkcs1.pem public.pem