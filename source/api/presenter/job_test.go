package presenter_test

import (
	"encoding/json"
	"net/url"

	"code.cloudfoundry.org/korifi/api/presenter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {
	var (
		baseURL *url.URL
		output  []byte
	)

	BeforeEach(func() {
		var err error
		baseURL, err = url.Parse("https://api.example.org")
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("JobFromGUID", func() {
		var (
			job   presenter.Job
			match bool
			guid  string
		)

		BeforeEach(func() {
			guid = "resource.operation~guid"
		})

		JustBeforeEach(func() {
			job, match = presenter.JobFromGUID(guid)
		})

		It("parses a job GUID into a Job struct", func() {
			Expect(match).To(BeTrue())
			Expect(job).To(Equal(presenter.Job{
				GUID:         "resource.operation~guid",
				Type:         "resource.operation",
				ResourceGUID: "guid",
				ResourceType: "Resource",
			}))
		})
	})

	Describe("ForManifestApplyJob", func() {
		JustBeforeEach(func() {
			response := presenter.ForManifestApplyJob(presenter.Job{
				GUID:         "the-job-guid",
				Type:         presenter.SpaceApplyManifestOperation,
				ResourceGUID: "the-space-guid",
			}, *baseURL)
			var err error
			output, err = json.Marshal(response)
			Expect(err).NotTo(HaveOccurred())
		})

		It("renders the job", func() {
			Expect(output).To(MatchJSON(`{
				"created_at": "",
				"errors": [],
				"guid": "the-job-guid",
				"links": {
					"self": {
						"href": "https://api.example.org/v3/jobs/the-job-guid"
					},
					"space": {
						"href": "https://api.example.org/v3/spaces/the-space-guid"
					}
				},
				"operation": "space.apply_manifest",
				"state": "COMPLETE",
				"updated_at": "",
				"warnings": null
			}`))
		})
	})

	Describe("ForJob", func() {
		JustBeforeEach(func() {
			response := presenter.ForJob(presenter.Job{
				GUID: "the-job-guid",
				Type: "the.operation",
			}, []presenter.JobResponseError{{
				Detail: "error detail",
				Title:  "CF-JobErrorTitle",
				Code:   12345,
			}}, "COMPLETE", *baseURL)
			var err error
			output, err = json.Marshal(response)
			Expect(err).NotTo(HaveOccurred())
		})

		It("renders the job", func() {
			Expect(output).To(MatchJSON(`{
				"created_at": "",
				"errors": [
					{
						"code": 12345,
						"detail": "error detail",
						"title": "CF-JobErrorTitle"
					}
				],
				"guid": "the-job-guid",
				"links": {
					"self": {
						"href": "https://api.example.org/v3/jobs/the-job-guid"
					}
				},
				"operation": "the.operation",
				"state": "COMPLETE",
				"updated_at": "",
				"warnings": null
			}`))
		})
	})

	Describe("JobURLForRedirects", func() {
	})
})
