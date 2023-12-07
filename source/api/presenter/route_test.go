package presenter_test

import (
	"encoding/json"
	"net/url"
	"time"

	"code.cloudfoundry.org/korifi/api/presenter"
	"code.cloudfoundry.org/korifi/api/repositories"
	. "code.cloudfoundry.org/korifi/tests/matchers"
	"code.cloudfoundry.org/korifi/tools"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Route", func() {
	var (
		baseURL *url.URL
		output  []byte
		record  repositories.RouteRecord
	)

	BeforeEach(func() {
		var err error
		baseURL, err = url.Parse("https://api.example.org")
		Expect(err).NotTo(HaveOccurred())
		record = repositories.RouteRecord{
			GUID:      "test-route-guid",
			SpaceGUID: "test-space-guid",
			Domain: repositories.DomainRecord{
				GUID: "test-domain-guid",
				Name: "example.org",
			},
			Host:     "test-route-host",
			Path:     "/some_path",
			Protocol: "http",
			Destinations: []repositories.DestinationRecord{
				{
					GUID:        "dest-1-guid",
					AppGUID:     "app-1-guid",
					ProcessType: "web",
					Port:        tools.PtrTo(1234),
					Protocol:    tools.PtrTo("http1"),
				},
				{
					GUID:        "dest-2-guid",
					AppGUID:     "app-2-guid",
					ProcessType: "queue",
					Port:        tools.PtrTo(5678),
					Protocol:    tools.PtrTo("http2"),
				},
			},
			Labels:      nil,
			Annotations: nil,
			CreatedAt:   time.UnixMilli(1000),
			UpdatedAt:   tools.PtrTo(time.UnixMilli(2000)),
		}
	})

	Describe("Route Response", func() {
		JustBeforeEach(func() {
			response := presenter.ForRoute(record, *baseURL)
			var err error
			output, err = json.Marshal(response)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the expected JSON", func() {
			Expect(output).To(MatchJSON(`{
				"guid": "test-route-guid",
				"port": null,
				"path": "/some_path",
				"protocol": "http",
				"host": "test-route-host",
				"url": "test-route-host.example.org/some_path",
				"created_at": "1970-01-01T00:00:01Z",
				"updated_at": "1970-01-01T00:00:02Z",
				"destinations": [
					{
						"guid": "dest-1-guid",
						"app": {
							"guid": "app-1-guid",
							"process": {
								"type": "web"
							}
						},
						"weight": null,
						"port": 1234,
						"protocol": "http1"
					},
					{
						"guid": "dest-2-guid",
						"app": {
							"guid": "app-2-guid",
							"process": {
								"type": "queue"
							}
						},
						"weight": null,
						"port": 5678,
						"protocol": "http2"
					}
				],
				"relationships": {
					"space": {
						"data": {
							"guid": "test-space-guid"
						}
					},
					"domain": {
						"data": {
							"guid": "test-domain-guid"
						}
					}
				},
				"metadata": {
					"labels": {},
					"annotations": {}
				},
				"links": {
					"self":{
						"href": "https://api.example.org/v3/routes/test-route-guid"
					},
					"space":{
						"href": "https://api.example.org/v3/spaces/test-space-guid"
					},
					"domain":{
						"href": "https://api.example.org/v3/domains/test-domain-guid"
					},
					"destinations":{
						"href": "https://api.example.org/v3/routes/test-route-guid/destinations"
					}
				}
			}`))
		})

		When("host is empty", func() {
			BeforeEach(func() {
				record.Host = ""
			})

			It("omits it", func() {
				Expect(output).To(MatchJSONPath("$.url", "example.org/some_path"))
			})
		})
	})

	Describe("destinations", func() {
		JustBeforeEach(func() {
			response := presenter.ForRouteDestinations(record, *baseURL)
			var err error
			output, err = json.Marshal(response)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the expected JSON", func() {
			Expect(output).To(MatchJSON(`{
				"destinations": [
					{
						"guid": "dest-1-guid",
						"app": {
							"guid": "app-1-guid",
							"process": {
								"type": "web"
							}
						},
						"weight": null,
						"port": 1234,
						"protocol": "http1"
					},
					{
						"guid": "dest-2-guid",
						"app": {
							"guid": "app-2-guid",
							"process": {
								"type": "queue"
							}
						},
						"weight": null,
						"port": 5678,
						"protocol": "http2"
					}
				],
				"links": {
					"self": {
						"href": "https://api.example.org/v3/routes/test-route-guid/destinations"
					},
					"route": {
						"href": "https://api.example.org/v3/routes/test-route-guid"
					}
				}
			}`))
		})
	})
})
