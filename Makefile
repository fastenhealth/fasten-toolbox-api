.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go mod tidy
	env GOARCH=amd64 GOOS=linux go mod vendor
	env GOARCH=amd64 GOOS=linux go build -mod=mod -ldflags="-s -w" -o bin/health cmd/health/health.go
	env GOARCH=amd64 GOOS=linux go build -mod=mod -ldflags="-s -w" -o bin/source_export cmd/source_export/source_export.go
clean:
	rm -rf ./bin ./vendor go.sum

plugins:
	npx serverless plugin install -n serverless-iam-roles-per-function@^3.2.0
	npx serverless plugin install -n serverless-prune-plugin@2.0.1
	npx serverless plugin install -n serverless-offline@9.2.5
	npx serverless plugin install -n serverless-domain-manager@6.1.0
	npx serverless plugin install -n serverless-plugin-log-retention@2.0.0

deploy-gh: plugins deploy

deploy: clean build
	npx serverless deploy --verbose

remove: clean build
	npx serverless remove --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
