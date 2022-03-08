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
	"github.com/buildpacks/libcnb"
	"github.com/magiconair/properties"
	"github.com/paketo-buildpacks/libpak/bard"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/garethjevans/apache-tomee/helper"
)

func testDynamicPropertySupport(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		a               = helper.DynamicPropertySupport{Logger: bard.NewLogger(os.Stdout)}
		catalinaBase       string
		existingContent = `existing.property=123
common.loader="${catalina.base}/lib","${catalina.base}/lib/*.jar","${catalina.home}/lib","${catalina.home}/lib/*.jar"
# a comment
multiline.property=a value,\
followed by,\
another value
`
	)

	context("will configure properties", func() {
		it.Before(func() {
			Expect(os.Setenv("BPL_TOMCAT_ENV_POSTGRES_PASSWORD", "password")).To(Succeed())
			Expect(os.Setenv("BPL_TOMCAT_ENV_POSTGRES_USERNAME", "username")).To(Succeed())
			Expect(os.Setenv("BPL_TOMCAT_ENV_POSTGRES_URL", "jdbc:postgresql://localhost:5432/test")).To(Succeed())

			var err error
			catalinaBase, err = ioutil.TempDir("", "catalina-base")
			Expect(err).NotTo(HaveOccurred())

			Expect(os.Mkdir(filepath.Join(catalinaBase, "conf"), 0755)).To(Succeed())
			catalinaProperties := filepath.Join(catalinaBase, "conf", "catalina.properties")
			Expect(ioutil.WriteFile(catalinaProperties, []byte(existingContent),
				0664)).To(Succeed())

			Expect(os.Setenv("CATALINA_BASE", catalinaBase)).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BPL_TOMCAT_ENV_POSTGRES_PASSWORD")).To(Succeed())
			Expect(os.Unsetenv("BPL_TOMCAT_ENV_POSTGRES_USERNAME")).To(Succeed())
			Expect(os.Unsetenv("BPL_TOMCAT_ENV_POSTGRES_URL")).To(Succeed())

			Expect(os.Unsetenv("CATALINA_BASE")).To(Succeed())

			Expect(os.RemoveAll(catalinaBase)).To(Succeed())
		})

		it("contributes configuration", func() {
			_, err := a.Execute()
			Expect(err).To(Not(HaveOccurred()))

			assertFileContents(t, filepath.Join(catalinaBase, "conf", "catalina.properties"), map[string]string{
				"existing.property":  "123",
				"postgres.url":       "jdbc:postgresql://localhost:5432/test",
				"postgres.password":  "password",
				"postgres.username":  "username",
				"common.loader":      `"${catalina.base}/lib","${catalina.base}/lib/*.jar","${catalina.home}/lib","${catalina.home}/lib/*.jar"`,
				"multiline.property": "a value,followed by,another value",
			})
			assertFileContentsDoesNotContain(t, filepath.Join(catalinaBase, "conf", "catalina.properties"), []string{
				"other.property",
			})
		})
	})

	context("will configure properties from bindings", func() {
		it.Before(func() {
			var err error
			catalinaBase, err = ioutil.TempDir("", "catalina-base")
			Expect(err).NotTo(HaveOccurred())

			Expect(os.Mkdir(filepath.Join(catalinaBase, "conf"), 0755)).To(Succeed())
			catalinaProperties := filepath.Join(catalinaBase, "conf", "catalina.properties")
			Expect(ioutil.WriteFile(catalinaProperties, []byte(existingContent),
				0664)).To(Succeed())

			Expect(os.Setenv("CATALINA_BASE", catalinaBase)).To(Succeed())

			a.Bindings = libcnb.Bindings{
				{
					Name:   "postgres",
					Type:   "PostgreSQL",
					Secret: map[string]string{"url": "jdbc:postgresql://localhost:5432/test"},
				},
			}
		})

		it.After(func() {
			Expect(os.Unsetenv("CATALINA_BASE")).To(Succeed())
			Expect(os.RemoveAll(catalinaBase)).To(Succeed())
		})

		it("will not contribute configuration", func() {
			_, err := a.Execute()
			Expect(err).To(Not(HaveOccurred()))

			assertFileContents(t, filepath.Join(catalinaBase, "conf", "catalina.properties"), map[string]string{
				"existing.property":  "123",
				"common.loader":      `"${catalina.base}/lib","${catalina.base}/lib/*.jar","${catalina.home}/lib","${catalina.home}/lib/*.jar"`,
				"multiline.property": "a value,followed by,another value",
			})
			assertFileContentsDoesNotContain(t, filepath.Join(catalinaBase, "conf", "catalina.properties"), []string{
				"other.property",
				"postgres.url",
			})
		})

		context("$BPL_TOMCAT_BINDING_NAME", func() {
			it.Before(func() {
				Expect(os.Setenv("BPL_TOMCAT_BINDING_NAME", "postgresql")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("BPL_TOMCAT_BINDING_NAME")).To(Succeed())
			})

			it("will contribute configuration", func() {
				_, err := a.Execute()
				Expect(err).To(Not(HaveOccurred()))

				assertFileContents(t, filepath.Join(catalinaBase, "conf", "catalina.properties"), map[string]string{
					"existing.property":  "123",
					"postgres.url":       "jdbc:postgresql://localhost:5432/test",
					"common.loader":      `"${catalina.base}/lib","${catalina.base}/lib/*.jar","${catalina.home}/lib","${catalina.home}/lib/*.jar"`,
					"multiline.property": "a value,followed by,another value",
				})
				assertFileContentsDoesNotContain(t, filepath.Join(catalinaBase, "conf", "catalina.properties"), []string{
					"other.property",
				})
			})
		})
	})
}

func assertFileContents(t *testing.T, file string, contents map[string]string) {
	Expect := NewWithT(t).Expect
	Expect(file).Should(BeAnExistingFile())

	loader := properties.Loader{DisableExpansion: true, Encoding: properties.UTF8}
	p, err := loader.LoadFile(file)
	Expect(err).Should(Not(HaveOccurred()))

	for k, v := range contents {
		val, ok := p.Get(k)
		Expect(val).To(Equal(v))
		Expect(ok).To(BeTrue())
	}
}

func assertFileContentsDoesNotContain(t *testing.T, file string, keys []string) {
	Expect := NewWithT(t).Expect
	Expect(file).Should(BeAnExistingFile())

	loader := properties.Loader{DisableExpansion: true, Encoding: properties.UTF8}
	p, err := loader.LoadFile(file)
	Expect(err).Should(Not(HaveOccurred()))

	for _, k := range keys {
		_, ok := p.Get(k)
		Expect(ok).To(BeFalse())
	}
}