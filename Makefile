IMAGE = garethjevans/rest-service
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')

clean:
	rm -fr target
	rm -fr build
	rm -fr bin
	rm -fr integrationtests/target

target:
	mkdir -p target

integrationtests/target/rest-service-0.1-SNAPSHOT.war:
	cd integrationtests && mvn clean package --no-transfer-progress && cd ..

pack: integrationtests/target/rest-service-0.1-SNAPSHOT.war buildpack
	pack build $(IMAGE) \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/target/rest-service-0.1-SNAPSHOT.war \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=7.* \
		--publish

release: pack
	docker tag $(IMAGE) garethjevans/rest-service:latest
	docker push garethjevans/rest-service:latest

buildpack:
	mkdir -p build
	create-package --destination ./build --version "0.0.1"
	cp package.toml.template ./build/package.toml
	pack buildpack package garethjevans_apache_tomee.cnb --format=file --config ./build/package.toml 

.PHONY: build
build:
	go build ./...

test:
	go test ./...

run: pack
	docker run -p 8080:8080 \
		$(IMAGE)

integration: pack
	container-structure-test test --config test-config.yaml --image $(IMAGE)
	docker-compose up -d
	sleep 30
	venom run integration-tests.yaml
	docker-compose down

venom:
	curl -o venom -L "https://github.com/ovh/venom/releases/download/v1.0.1/venom.$(OS)-amd64"
	chmod +x venom
	mv venom /usr/local/bin
