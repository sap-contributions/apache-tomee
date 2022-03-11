package integration_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testDefault(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
		home = os.Getenv("HOME")
	)

	it.Before(func() {
		pack = occam.NewPack().WithVerbose()
		docker = occam.NewDocker()
	})

	context("when the buildpack is run with pack build", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			source, err = occam.Source(filepath.Join("testdata"))
			Expect(err).NotTo(HaveOccurred())

			m2Repo := fmt.Sprintf("%s/.m2", home)

			if _, err := os.Stat(m2Repo); errors.Is(err, os.ErrNotExist) {
				Expect(os.Mkdir(m2Repo, 0755)).To(Succeed())
			}
		})

		it.After(func() {
			if container.ID != "" {
				Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			}
			if image.ID != "" {
				Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			}
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("builds with the defaults", func() {
			var (
				logs fmt.Stringer
				err  error
			)

			start := time.Now()
			image, logs, err = pack.WithNoColor().WithVerbose().Build.
				WithPullPolicy("if-not-present").
				WithBuilder("paketobuildpacks/builder:base").
				WithBuildpacks("paketo-buildpacks/syft",
					"paketo-buildpacks/ca-certificates@3.0.2",
					"paketo-buildpacks/bellsoft-liberica",
					"paketo-buildpacks/maven",
					buildpack).
				WithEnv(map[string]string{
					"BP_JAVA_APP_SERVER": "tomee",
					"BP_MAVEN_BUILT_ARTIFACT": "test-jaxrs-tomee/target/*.war",
					"BP_MAVEN_BUILT_MODULE": "test-jaxrs-tomee",
					"BP_MAVEN_BUILD_ARGUMENTS": "-Dmaven.test.skip=true package --no-transfer-progress",
				}).
				WithVolumes(fmt.Sprintf("%s/.m2:/home/cnb/.m2:rw", home)).
				Execute(name, source)
			Expect(err).ToNot(HaveOccurred(), logs.String)
			fmt.Println(logs.String())
			elapsed := time.Since(start)
			fmt.Printf("pack build took %s\n", elapsed)

			start = time.Now()
			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())
			elapsed = time.Since(start)
			fmt.Printf("docker run took %s\n", elapsed)

			Eventually(container, time.Second*30).Should(Serve(ContainSubstring("{\"application_status\":\"UP\"}")).OnPort(8080))
		})

		it("builds on tiny", func() {
			var (
				logs fmt.Stringer
				err  error
			)

			start := time.Now()
			image, logs, err = pack.WithNoColor().WithVerbose().Build.
				WithPullPolicy("if-not-present").
				WithBuilder("paketobuildpacks/builder:tiny").
				WithBuildpacks("paketo-buildpacks/syft",
					"paketo-buildpacks/ca-certificates@3.0.2",
					"paketo-buildpacks/bellsoft-liberica",
					"paketo-buildpacks/maven",
					buildpack).
				WithEnv(map[string]string{
					"BP_JAVA_APP_SERVER": "tomee",
					"BP_MAVEN_BUILT_ARTIFACT": "test-jaxrs-tomee/target/*.war",
					"BP_MAVEN_BUILT_MODULE": "test-jaxrs-tomee",
					"BP_MAVEN_BUILD_ARGUMENTS": "-Dmaven.test.skip=true package --no-transfer-progress",
				}).
				WithVolumes(fmt.Sprintf("%s/.m2:/home/cnb/.m2:rw", home)).
				Execute(name, source)
			Expect(err).ToNot(HaveOccurred(), logs.String)
			fmt.Println(logs.String())
			elapsed := time.Since(start)
			fmt.Printf("pack build took %s\n", elapsed)

			start = time.Now()
			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())
			elapsed = time.Since(start)
			fmt.Printf("docker run took %s\n", elapsed)

			Eventually(container, time.Second*30).Should(Serve(ContainSubstring("{\"application_status\":\"UP\"}")).OnPort(8080))
		})
	})
}
