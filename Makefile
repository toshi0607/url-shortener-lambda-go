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
		--capabilities CAPABILITY_IAM
	echo API endpoint URL for Prod environment:
	aws cloudformation describe-stacks \
		--stack-name url-shortener-lambda-go \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text
