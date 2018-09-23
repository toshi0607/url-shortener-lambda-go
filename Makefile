build: build-shorten build-redirect

build-shorten:
	GOARCH=amd64 GOOS=linux go build -o artifact/shorten ./handlers/shorten

build-redirect:
	GOARCH=amd64 GOOS=linux go build -o artifact/redirect ./handlers/redirect

deploy: build
	sam package \
		--template-file template.yml \
		--s3-bucket url-shortener-lambda-go \
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

delete:
	aws cloudformation delete-stack --stack-name url-shortener-lambda-go
	aws s3 rm s3://url-shortener-lambda-go --recursive
	aws s3 rb s3://url-shortener-lambda-go

DBjar := DynamoDBLocal.jar
DBjar_exists := $(shell find . -name $(DBjar))
DBproc := $(shell lsof -t -i :8000)

db-start:
	java -Djava.library.path=./DynamoDBLocal_lib -jar test/dynamodb_local_latest/DynamoDBLocal.jar -sharedDb

db-close:
	kill -9 $(DBproc)

db-create-table:
	aws dynamodb create-table --cli-input-json file://test/link.json --endpoint-url http://localhost:8000
