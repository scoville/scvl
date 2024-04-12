# メタ情報
Name := scvl
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'
ENV ?= development

# Setup
setup:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help
	go get github.com/markbates/refresh

run:
	refresh run

deploy:
	rm -f bin/*
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/scvl

	ssh scvl0001w "supervisorctl stop scvl"
	scp -r bin css js templates scvl0001w:/home/ec2-user/scvl/
	ssh scvl0001w "supervisorctl start scvl"

	ssh scvl0002w "supervisorctl stop scvl"
	scp -r bin css js templates scvl0002w:/home/ec2-user/scvl/
	ssh scvl0002w "supervisorctl start scvl"

deploy-assets:
	scp -r css js templates scvl0001w:/home/ec2-user/scvl/
	scp -r css js templates scvl0002w:/home/ec2-user/scvl/
	ssh scvl0001w "supervisorctl restart scvl"
	ssh scvl0002w "supervisorctl restart scvl"

# Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps help
