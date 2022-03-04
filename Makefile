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

pack-tomee-9-microprofile: integration/testdata/test-jaxrs-tomee-jakarta buildpack
	pack build $(IMAGE)-v9-microprofile \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/test-jaxrs-tomee-jakarta/target/test-jaxrs-tomee-jakarta-0.1-SNAPSHOT.war \
		--clear-cache \
		-e BP_JAVA_APP_SERVER=tomee \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=9.0.0-M7

pack-tomee-8-microprofile: integration/testdata/test-jaxrs-tomee buildpack
	pack build $(IMAGE)-v8-microprofile \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war \
		--clear-cache \
		-e BP_JAVA_APP_SERVER=tomee \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=8.*

pack-tomee-8-webprofile: integration/testdata/test-jaxrs-tomee buildpack
	pack build $(IMAGE)-v8-webprofile \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war \
		--clear-cache \
		-e BP_JAVA_APP_SERVER=tomee \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=8.* \
		-e BP_TOMEE_DISTRIBUTION=webprofile

pack-tomee-8-plus: integration/testdata/test-jaxrs-tomee buildpack
	pack build $(IMAGE)-v8-plus \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war \
		--clear-cache \
		-e BP_JAVA_APP_SERVER=tomee \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=8.* \
		-e BP_TOMEE_DISTRIBUTION=webprofile

pack-tomee-8-plume: integration/testdata/test-jaxrs-tomee buildpack
	pack build $(IMAGE)-v8-plume \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war \
		--clear-cache \
		-e BP_JAVA_APP_SERVER=tomee \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=8.* \
		-e BP_TOMEE_DISTRIBUTION=plume

pack-tomee-7-microprofile: integration/testdata/test-jaxrs-tomee buildpack
	pack build $(IMAGE)-v7-microprofile \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/test-jaxrs-tomee/target/test-jaxrs-tomee-0.1-SNAPSHOT.war \
		--clear-cache \
		-e BP_JAVA_APP_SERVER=tomee \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=7.*

release: pack-tomee-7-microprofile pack-tomee-8-microprofile pack-tomee-9-microprofile pack-tomee-8-webprofile pack-tomee-8-plus pack-tomee-8-plume
	docker push garethjevans/test-jaxrs-tomee-v7-microprofile:latest
	docker push garethjevans/test-jaxrs-tomee-v8-webprofile:latest
	docker push garethjevans/test-jaxrs-tomee-v8-microprofile:latest
	docker push garethjevans/test-jaxrs-tomee-v8-plus:latest
	docker push garethjevans/test-jaxrs-tomee-v8-plume:latest
	docker push garethjevans/test-jaxrs-tomee-v9-microprofile:latest

buildpack:
	mkdir -p build
	create-package --destination ./build --version "0.0.1"
	cp package.toml.template ./build/package.toml
	pack buildpack package garethjevans_apache_tomee.cnb --format=file --config ./build/package.toml 

.PHONY: build
build:
	go build ./...

test:
	go test -v -timeout 20m ./...

run: pack
	docker run -p 8080:8080 \
		$(IMAGE)-v8-microprofile

integration: pack-tomee-7-microprofile pack-tomee-8-microprofile pack-tomee-8-webprofile pack-tomee-8-plus pack-tomee-8-plume pack-tomee-9-microprofile
	container-structure-test test --config integrationtests/test-config.yaml --image $(IMAGE)-v7-microprofile
	container-structure-test test --config integrationtests/test-config.yaml --image $(IMAGE)-v8-microprofile
	container-structure-test test --config integrationtests/test-config.yaml --image $(IMAGE)-v8-webprofile
	container-structure-test test --config integrationtests/test-config.yaml --image $(IMAGE)-v8-plus
	container-structure-test test --config integrationtests/test-config.yaml --image $(IMAGE)-v8-plume
	container-structure-test test --config integrationtests/test-config.yaml --image $(IMAGE)-v9-microprofile
	cd integrationtests && docker-compose up -d
	sleep 15
	venom run integrationtests/integration-tests.yaml
	cd integrationtests && docker-compose down

venom:
	curl -o venom -L "https://github.com/ovh/venom/releases/download/v1.0.1/venom.$(OS)-amd64"
	chmod +x venom
	mv venom /usr/local/bin
