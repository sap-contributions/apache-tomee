ORG = garethjevans
IMAGE = $(ORG)/test-jaxrs-tomee
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')

clean:
	rm -fr target
	rm -fr build
	rm -fr bin
	cd integrationtests && mvn clean && cd ..

target:
	mkdir -p target

integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war:
	cd integrationtests && mvn clean package --no-transfer-progress && cd ..

integrationtests/test-jaxrs-tomee-jakarta/target/test-jaxrs-tomee-jakarta-0.1-SNAPSHOT.war:
	cd integrationtests && mvn clean package --no-transfer-progress && cd ..

pack-tomee-9: integrationtests/test-jaxrs-tomee-jakarta/target/test-jaxrs-tomee-jakarta-0.1-SNAPSHOT.war buildpack
	pack build $(IMAGE)-v9 \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/test-jaxrs-tomee-jakarta/target/test-jaxrs-tomee-jakarta-0.1-SNAPSHOT.war \
		--clear-cache \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=9.0.0-M7

pack-tomee-8: integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war buildpack
	pack build $(IMAGE)-v8 \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war \
		--clear-cache \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=8.*

pack-tomee-7: integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war buildpack
	pack build $(IMAGE)-v7 \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war \
		--clear-cache \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=7.*

release: pack-tomee-7 pack-tomee-8 pack-tomee-9
	docker push garethjevans/test-jaxrs-tomee-v7:latest
	docker push garethjevans/test-jaxrs-tomee-v8:latest
	docker push garethjevans/test-jaxrs-tomee-v9:latest

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

integration: pack-tomee-7 pack-tomee-8 pack-tomee-9
	container-structure-test test --config integrationtests/test-config.yaml --image $(IMAGE)-v7
	container-structure-test test --config integrationtests/test-config.yaml --image $(IMAGE)-v8
	container-structure-test test --config integrationtests/test-config.yaml --image $(IMAGE)-v9
	cd integrationtests && docker-compose up -d
	sleep 15
	venom run integrationtests/integration-tests.yaml
	cd integrationtests && docker-compose down

venom:
	curl -o venom -L "https://github.com/ovh/venom/releases/download/v1.0.1/venom.$(OS)-amd64"
	chmod +x venom
	mv venom /usr/local/bin
