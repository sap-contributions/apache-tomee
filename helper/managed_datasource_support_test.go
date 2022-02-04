/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helper_test

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/garethjevans/apache-tomee/helper"
)

func testManagedDatasourceSupport(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		a = helper.ManagedDataSourceSupport{}
	)

	it("returns if $BPL_TOMCAT_MANAGED_DATASOURCE_ENABLED is not set", func() {
		Expect(a.Execute()).To(BeNil())
	})

	context("$BPL_TOMCAT_MANAGED_DATASOURCE_ENABLED", func() {
		it.Before(func() {
			Expect(os.Setenv("BPL_TOMCAT_MANAGED_DATASOURCE_ENABLED", "")).To(Succeed())
			Expect(os.Setenv("BPL_TOMCAT_MANAGED_DATASOURCE_DRIVER", "driver")).To(Succeed())
			Expect(os.Setenv("BPL_TOMCAT_MANAGED_DATASOURCE_USERNAME", "username")).To(Succeed())
			Expect(os.Setenv("BPL_TOMCAT_MANAGED_DATASOURCE_PASSWORD", "password")).To(Succeed())
			Expect(os.Setenv("BPL_TOMCAT_MANAGED_DATASOURCE_URL", "url")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BPL_TOMCAT_MANAGED_DATASOURCE_ENABLED")).To(Succeed())
			Expect(os.Unsetenv("BPL_TOMCAT_MANAGED_DATASOURCE_DRIVER")).To(Succeed())
			Expect(os.Unsetenv("BPL_TOMCAT_MANAGED_DATASOURCE_USERNAME")).To(Succeed())
			Expect(os.Unsetenv("BPL_TOMCAT_MANAGED_DATASOURCE_PASSWORD")).To(Succeed())
			Expect(os.Unsetenv("BPL_TOMCAT_MANAGED_DATASOURCE_URL")).To(Succeed())
		})

		it("contributes configuration", func() {
			Expect(a.Execute()).To(Equal(map[string]string{"JAVA_TOOL_OPTIONS": "-Dmanaged.datasource.name=jdbc/datasource " +
				"-Dmanaged.datasource.auth=Container " +
				"-Dmanaged.datasource.type=javax.sql.DataSource " +
				"-Dmanaged.datasource.driver=driver -Dmanaged.datasource.username=username " +
				"-Dmanaged.datasource.password=password " +
				"-Dmanaged.datasource.url=url " +
				"-Dmanaged.datasource.maxtotal=20 " +
				"-Dmanaged.datasource.maxidle=10 " +
				"-Dmanaged.datasource.maxwaitmillis=-1"}))
		})

		context("$JAVA_TOOL_OPTIONS", func() {
			it.Before(func() {
				Expect(os.Setenv("JAVA_TOOL_OPTIONS", "test-java-tool-options")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
			})

			it("contributes configuration appended to existing $JAVA_TOOL_OPTIONS", func() {
				Expect(a.Execute()).To(Equal(map[string]string{
					"JAVA_TOOL_OPTIONS": "test-java-tool-options -Dmanaged.datasource.name=jdbc/datasource " +
						"-Dmanaged.datasource.auth=Container " +
						"-Dmanaged.datasource.type=javax.sql.DataSource " +
						"-Dmanaged.datasource.driver=driver " +
						"-Dmanaged.datasource.username=username " +
						"-Dmanaged.datasource.password=password " +
						"-Dmanaged.datasource.url=url " +
						"-Dmanaged.datasource.maxtotal=20 " +
						"-Dmanaged.datasource.maxidle=10 " +
						"-Dmanaged.datasource.maxwaitmillis=-1",
				}))
			})
		})
	})

}
