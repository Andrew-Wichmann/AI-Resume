include .env
export

.PHONY: serve
serve: build
	sam local start-api -v /home/twisted/Code/personal/AI-Resume/.aws-sam/build/ --container-host host.docker.internal --debug

.PHONY: build
build:
	sam build -t sam-template.yaml
