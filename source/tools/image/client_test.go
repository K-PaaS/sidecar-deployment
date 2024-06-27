package image_test

import (
	"os"

	"code.cloudfoundry.org/korifi/tests/helpers/oci"
	"code.cloudfoundry.org/korifi/tools/image"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		zipFile      *os.File
		otherZipFile *os.File
		testErr      error
		pushRef      string
		imgCfg       *v1.ConfigFile
		imgRef       string
		creds        image.Creds
	)

	BeforeEach(func() {
		testErr = nil
		var err error
		zipFile, err = os.Open("fixtures/layer.zip")
		Expect(err).NotTo(HaveOccurred())

		otherZipFile, err = os.Open("fixtures/anotherLayer.zip")
		Expect(err).NotTo(HaveOccurred())

		pushRef = containerRegistry.ImageRef("foo/bar")
		imgCfg = &v1.ConfigFile{
			Config: v1.Config{
				Labels: map[string]string{
					"foo": "bar",
				},
				ExposedPorts: map[string]struct{}{
					"123": {},
					"456": {},
				},
				User: "my-user",
			},
		}

		imgClient = image.NewClient(k8sClientset)
		creds = image.Creds{
			Namespace:   "default",
			SecretNames: []string{secretName},
		}
	})

	Describe("Push", func() {
		JustBeforeEach(func() {
			imgRef, testErr = imgClient.Push(ctx, creds, pushRef, zipFile, "jim")
		})

		It("pushes a zip archive as an image to the registry", func() {
			Expect(testErr).NotTo(HaveOccurred())
			Expect(imgRef).To(HavePrefix(pushRef))

			_, err := imgClient.Config(ctx, creds, imgRef)
			Expect(err).NotTo(HaveOccurred())

			_, err = imgClient.Config(ctx, creds, pushRef+":jim")
			Expect(err).NotTo(HaveOccurred())
		})

		When("pushRef is invalid", func() {
			BeforeEach(func() {
				pushRef += ":bar:baz"
			})

			It("fails", func() {
				Expect(testErr).To(MatchError(ContainSubstring("error parsing repository reference")))
			})
		})

		When("zip input is not valid", func() {
			BeforeEach(func() {
				var err error
				zipFile, err = os.Open("fixtures/not.a.zip")
				Expect(err).NotTo(HaveOccurred())
			})

			It("fails", func() {
				Expect(testErr).To(MatchError(ContainSubstring("not a valid zip file")))
			})
		})

		When("the secret doesn't exist", func() {
			BeforeEach(func() {
				creds.SecretNames = []string{"not-a-secret"}
			})

			It("fails to authenticate", func() {
				Expect(testErr).To(MatchError(ContainSubstring("Unauthorized")))
			})
		})

		When("using a service account for secrets", func() {
			BeforeEach(func() {
				creds.SecretNames = []string{}
				creds.ServiceAccountName = serviceAccountName
			})

			It("pushes a zip archive as an image to the registry", func() {
				Expect(testErr).NotTo(HaveOccurred())
				Expect(imgRef).To(HavePrefix(pushRef))

				_, err := imgClient.Config(ctx, creds, imgRef)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("seret name is empty (simulating ECR)", func() {
			BeforeEach(func() {
				ecrRegistry := oci.NewNoAuthContainerRegistry()
				pushRef = ecrRegistry.ImageRef("foo/bar")
				creds.SecretNames = []string{}
			})

			It("succeeds", func() {
				Expect(testErr).NotTo(HaveOccurred())
				Expect(imgRef).To(HavePrefix(pushRef))
			})
		})
	})

	Describe("Config", func() {
		var config image.Config

		BeforeEach(func() {
			pushRef += "/with/labels"
			containerRegistry.PushImage(pushRef, imgCfg)
		})

		JustBeforeEach(func() {
			config, testErr = imgClient.Config(ctx, creds, pushRef)
		})

		It("fetches the image config", func() {
			Expect(config.Labels).To(Equal(map[string]string{"foo": "bar"}))
			Expect(config.User).To(Equal("my-user"))
			Expect(config.ExposedPorts).To(ConsistOf(int32(123), int32(456)))
		})

		When("the ref is invalid", func() {
			BeforeEach(func() {
				pushRef += "::ads"
			})

			It("fails", func() {
				Expect(testErr).To(MatchError(ContainSubstring("error parsing repository reference")))
			})
		})

		When("the secret doesn't exist", func() {
			BeforeEach(func() {
				creds.SecretNames = []string{"not-a-secret"}
			})

			It("fails to authenticate", func() {
				Expect(testErr).To(MatchError(ContainSubstring("UNAUTHORIZED")))
			})
		})

		When("using a service account for secrets", func() {
			BeforeEach(func() {
				creds.SecretNames = []string{}
				creds.ServiceAccountName = serviceAccountName
			})

			It("fetches the image labels", func() {
				Expect(config.Labels).To(Equal(map[string]string{"foo": "bar"}))
				Expect(config.ExposedPorts).To(ConsistOf(int32(123), int32(456)))
			})
		})

		When("secret name is empty (simulating ECR)", func() {
			BeforeEach(func() {
				ecrRegistry := oci.NewNoAuthContainerRegistry()
				pushRef = ecrRegistry.ImageRef("foo/bar/with/labels")
				ecrRegistry.PushImage(pushRef, imgCfg)
				creds.SecretNames = []string{}
			})

			It("succeeds", func() {
				Expect(testErr).NotTo(HaveOccurred())
				Expect(config.Labels).To(Equal(map[string]string{"foo": "bar"}))
				Expect(config.ExposedPorts).To(ConsistOf(int32(123), int32(456)))
			})
		})

		When("ports are in the format 'port/protocol'", func() {
			BeforeEach(func() {
				imgCfg.Config.ExposedPorts = map[string]struct{}{
					"123/protocol": {},
				}
				containerRegistry.PushImage(pushRef, imgCfg)
			})

			It("succeeds", func() {
				Expect(testErr).NotTo(HaveOccurred())
				Expect(config.ExposedPorts).To(ConsistOf(int32(123)))
			})
		})
	})

	Describe("Delete", func() {
		var tagsToDelete []string

		BeforeEach(func() {
			tagsToDelete = []string{"jim", "bob"}
			var err error
			imgRef, err = imgClient.Push(ctx, creds, pushRef, zipFile, "jim", "bob")
			Expect(err).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			testErr = imgClient.Delete(ctx, creds, imgRef, tagsToDelete...)
		})

		It("deletes an image", func() {
			Expect(testErr).NotTo(HaveOccurred())

			_, err := imgClient.Config(ctx, creds, imgRef)
			Expect(err).To(MatchError(ContainSubstring("MANIFEST_UNKNOWN")))
		})

		When("a tag is omitted", func() {
			BeforeEach(func() {
				tagsToDelete = []string{"bob"}
			})

			It("deletes the tag but not the image", func() {
				Expect(testErr).NotTo(HaveOccurred())

				_, err := imgClient.Config(ctx, creds, pushRef+":bob")
				Expect(err).To(MatchError(ContainSubstring("MANIFEST_UNKNOWN")))

				_, err = imgClient.Config(ctx, creds, imgRef)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("another digest exists in the repo with another tag", func() {
			BeforeEach(func() {
				_, err := imgClient.Push(ctx, creds, pushRef, otherZipFile, "alice")
				Expect(err).NotTo(HaveOccurred())
			})

			It("still deletes the image when all _its_ tags are removed", func() {
				_, err := imgClient.Config(ctx, creds, imgRef)
				Expect(err).To(MatchError(ContainSubstring("MANIFEST_UNKNOWN")))
			})
		})

		When("the secret doesn't exist", func() {
			BeforeEach(func() {
				creds.SecretNames = []string{"not-a-secret"}
			})

			It("fails to authenticate", func() {
				Expect(testErr).To(MatchError(ContainSubstring("UNAUTHORIZED")))
			})
		})

		When("using a service account for secrets", func() {
			BeforeEach(func() {
				creds.SecretNames = []string{}
				creds.ServiceAccountName = serviceAccountName
			})

			It("deletes an image", func() {
				Expect(testErr).NotTo(HaveOccurred())

				_, err := imgClient.Config(ctx, creds, imgRef)
				Expect(err).To(MatchError(ContainSubstring("MANIFEST_UNKNOWN")))
			})
		})

		When("the secret name is empty (simulating ECR)", func() {
			BeforeEach(func() {
				ecrRegistry := oci.NewNoAuthContainerRegistry()
				pushRef = ecrRegistry.ImageRef("foo/bar")
				_, err := zipFile.Seek(0, 0)
				Expect(err).NotTo(HaveOccurred())

				creds.SecretNames = []string{}
				imgRef, err = imgClient.Push(ctx, creds, pushRef, zipFile)
				Expect(err).NotTo(HaveOccurred())
			})

			It("succeeds", func() {
				Expect(testErr).NotTo(HaveOccurred())

				_, err := imgClient.Config(ctx, creds, imgRef)
				Expect(err).To(MatchError(ContainSubstring("MANIFEST_UNKNOWN")))
			})
		})
	})

	for _, reg := range registries {
		// 'Serial` because we use the same image for ECR in this test and above
		// and otherwise they interfere
		Describe("Delete from cloud registries", Serial, func() {
			reg := reg // capture reg so we don't just use the last one!
			var (
				repoName     string
				tagsToDelete []string
			)

			BeforeEach(func() {
				tagsToDelete = []string{"jim", "bob"}
				pushRef = reg.RepoName
				if pushRef == "" {
					repoName = uuid.NewString()
					pushRef = reg.PathPrefix + "/" + repoName
				}
				var err error

				imgRef, err = imgClient.Push(ctx, creds, pushRef, zipFile, "jim", "bob")
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				_ = imgClient.Delete(ctx, creds, imgRef, "jim", "bob")
			})

			JustBeforeEach(func() {
				testErr = imgClient.Delete(ctx, creds, imgRef, tagsToDelete...)
			})

			It("deletes the image", func() {
				Expect(testErr).NotTo(HaveOccurred())

				_, err := imgClient.Config(ctx, creds, imgRef)
				Expect(err).To(MatchError(ContainSubstring("MANIFEST_UNKNOWN")))
			})

			When("a tag is omitted", func() {
				BeforeEach(func() {
					tagsToDelete = []string{"bob"}
				})

				It("does not delete the image", func() {
					Expect(testErr).NotTo(HaveOccurred())

					_, err := imgClient.Config(ctx, creds, imgRef)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			When("another digest exists in the repo with another tag", func() {
				var otherImg string

				BeforeEach(func() {
					var err error
					otherImg, err = imgClient.Push(ctx, creds, pushRef, otherZipFile, "alice")
					Expect(err).NotTo(HaveOccurred())
				})

				AfterEach(func() {
					Expect(imgClient.Delete(ctx, creds, otherImg, "alice")).To(Succeed())
				})

				It("still deletes the image when all _its_ tags are removed", func() {
					_, err := imgClient.Config(ctx, creds, imgRef)
					Expect(err).To(MatchError(ContainSubstring("MANIFEST_UNKNOWN")))
				})
			})
		})
	}
})
