IMAGE = example/datasource
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')

clean:
	rm -fr target
	rm -fr build
	rm -fr bin
	rm -fr integrationtests/target

target:
	mkdir -p target

integrationtests/target/datasource-0.1-SNAPSHOT.war:
	cd integrationtests && mvn clean package --no-transfer-progress && cd ..

pack: integrationtests/target/datasource-0.1-SNAPSHOT.war buildpack
	pack build $(IMAGE) \
		--buildpack paketo-buildpacks/syft@1.5.0 \
		--buildpack paketo-buildpacks/ca-certificates@2.4.2 \
		--buildpack paketo-buildpacks/bellsoft-liberica@9.0.2 \
		--buildpack garethjevans_apache_tomee.cnb \
		--path integrationtests/target/datasource-0.1-SNAPSHOT.war \
		-e BP_JVM_VERSION=8 \
		-e BP_TOMEE_VERSION=8.*

release: pack
	docker tag $(IMAGE) garethjevans/datasource:latest
	docker push garethjevans/datasource:latest

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
		-e BPL_TOMCAT_ACCESS_LOGGING_ENABLED=true \
		-e BPL_TOMCAT_MANAGED_DATASOURCE_ENABLED=true \
		-e BPL_TOMCAT_MANAGED_DATASOURCE_DRIVER=oracle.jdbc.Driver \
		-e BPL_TOMCAT_MANAGED_DATASOURCE_URL=jdbc:oracle:thin:@127.0.0.1:1521:mysid \
		-e BPL_TOMCAT_MANAGED_DATASOURCE_USERNAME=scott \
		-e BPL_TOMCAT_MANAGED_DATASOURCE_PASSWORD=tiger \
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
