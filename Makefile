ORG = garethjevans
IMAGE = $(ORG)/test-jaxrs-tomee
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')

clean:
	rm -fr target
	rm -fr build
	rm -fr bin
	cd integration/testdata && mvn clean && cd ..

target:
	mkdir -p target

buildpack:
	mkdir -p build
	create-package --destination ./build --version "0.0.1"
	cp package.toml.template ./build/package.toml
	pack buildpack package garethjevans_apache_tomee.cnb --format=file --config ./build/package.toml 

.PHONY: build
build:
	go build ./...

test:
	richgo test -v ./... -run Unit

.PHONY: integration
integration: integration/testdata/test-jaxrs-tomee
	richgo test -v ./integration/... -run Integration

venom:
	curl -o venom -L "https://github.com/ovh/venom/releases/download/v1.0.1/venom.$(OS)-amd64"
	chmod +x venom
	mv venom /usr/local/bin
