TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

MAKEFLAGS += --silent

build:
	go build -o terraform-provider-keycloak
	mkdir -p providers/registry.terraform.io/joed22636/keycloak/1.0.3/darwin_amd64
	cp terraform-provider-keycloak providers/registry.terraform.io/joed22636/keycloak/1.0.3/darwin_amd64/terraform-provider-keycloak_v1.0.3

build-example: build
	mkdir -p example/.terraform/plugins/terraform.local/joed22636/keycloak/3.0.0/darwin_amd64
	mkdir -p example/terraform.d/plugins/terraform.local/joed22636/keycloak/3.0.0/darwin_amd64
	cp terraform-provider-keycloak example/.terraform/plugins/terraform.local/joed22636/keycloak/3.0.0/darwin_amd64/
	cp terraform-provider-keycloak example/terraform.d/plugins/terraform.local/joed22636/keycloak/3.0.0/darwin_amd64/

build-win: build
	-mkdir providers\registry.terraform.io\joed22636\keycloak\1.0.3\windows_amd64
	copy terraform-provider-keycloak providers\registry.terraform.io\joed22636\keycloak\1.0.3\windows_amd64\terraform-provider-keycloak_v1.0.3.exe
	
local: deps
	docker-compose up --build -d
	./scripts/wait-for-local-keycloak.sh
	./scripts/create-terraform-client.sh

deps:
	./scripts/check-deps.sh

fmt:
	gofmt -w -s $(GOFMT_FILES)

test: fmtcheck vet
	go test -v $(TEST)

testacc: fmtcheck vet
	TF_ACC=1 CHECKPOINT_DISABLE=1 go test -v -timeout 30m -parallel 4 $(TEST) $(TESTARGS)

fmtcheck:
	lineCount=$(shell gofmt -l -s $(GOFMT_FILES) | wc -l | tr -d ' ') && exit $$lineCount

vet:
	go vet ./...

user-federation-example:
	cd custom-user-federation-example && ./gradlew shadowJar
