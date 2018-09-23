build: build-shorten build-redirect
.PHONY: build

build-shorten:
	GOARCH=amd64 GOOS=linux go build -o artifact/shorten ./handlers/shorten
.PHONY: build-shorten

build-redirect:
	GOARCH=amd64 GOOS=linux go build -o artifact/redirect ./handlers/redirect
.PHONY: build-redirect

deploy: build
	sam package \
		--template-file template.yml \
		--s3-bucket stack-bucket-for-url-shortener-lambda-go \
		--output-template-file sam.yml
	sam deploy \
		--template-file sam.yml \
		--stack-name url-shortener-lambda-go \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
			LinkTableName=$(LINK_TABLE)
	echo API endpoint URL for Prod environment:
	aws cloudformation describe-stacks \
		--stack-name url-shortener-lambda-go \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text
.PHONY: deploy

delete:
	aws cloudformation delete-stack --stack-name url-shortener-lambda-go
	aws s3 rm s3://stack-bucket-for-url-shortener-lambda-go --recursive
	aws s3 rb s3://stack-bucket-for-url-shortener-lambda-go
.PHONY: delete

test:
	go test ./...
.PHONY: test

DBjar := DynamoDBLocal.jar
DBjar_exists := $(shell find . -name $(DBjar))
DBproc := $(shell lsof -t -i :8000)

db-start:
	java -Djava.library.path=./DynamoDBLocal_lib -jar test/dynamodb_local_latest/DynamoDBLocal.jar -sharedDb
.PHONY: db-start

db-close:
	kill -9 $(DBproc)
.PHONY: db-close

db-create-table:
	aws dynamodb create-table --cli-input-json file://test/link.json --endpoint-url http://localhost:8000
.PHONY: db-create-table
