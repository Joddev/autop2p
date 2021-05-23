.PHONY: build clean deploy gomodgen

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o ./bin/main ./main
	chmod +x ./bin/main

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose
