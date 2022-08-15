/*
 * Copyright 2018-2022 the original author or authors.
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

package tomee_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/apache-tomee/tomee"
)

func testHome(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Layers.Path, err = os.MkdirTemp("", "home-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes catalina home", func() {
		dep := libpak.BuildpackDependency{
			ID:     "tomee-microprofile",
			URI:    "https://localhost/stub-tomee.tar.gz",
			SHA256: "26a0f0d1782027e7849389cb975a40cd8e69497d19946442881d61bd5f1756bf",
			PURL:   "pkg:generic/tomee@1.1.1",
			CPEs:   []string{"cpe:2.3:a:apache:tomee:1.1.1:*:*:*:*:*:*:*"},
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		h, _ := tomee.NewHome(dep, dc)

		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = h.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())
		Expect(filepath.Join(layer.Path, "fixture-marker")).To(BeARegularFile())
		Expect(layer.LaunchEnvironment["CATALINA_HOME.default"]).To(Equal(layer.Path))
	})
}
