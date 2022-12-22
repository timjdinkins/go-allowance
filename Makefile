.PHONY: build
build:
	sam build

#build-GetAll:
#	GOOS=linux CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/getAll ./cmd/lambdas/getAll/main.go
#
#build-GetById:
#	GOOS=linux CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/getById ./cmd/lambdas/getById/main.go
#
#build-CreateTransaction:
#	GOOS=linux CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/getAll ./cmd/lambdas/createTransaction/main.go

.PHONY: init
init: build
	sam deploy --guided

.PHONY: deploy
deploy: build
	sam deploy

.PHONY: delete
delete:
	sam delete
