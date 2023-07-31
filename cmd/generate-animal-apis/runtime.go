// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/encoding"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
)

const runtime = "runtime"
const organization = "apigee-apihub-demo"

func runtimeCommand() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "runtime",
		Short: "Generate YAML for mock runtime signals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			animals, err := readAnimals(filename)
			if err != nil {
				return err
			}
			for _, animal := range animals.Animals {
				if err = generateRuntimeMocks(&animal); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&filename, "filename", "animals.yaml", "Animals definition file (YAML)")
	return cmd
}

func generateRuntimeMocks(animal *Animal) error {
	upperPlural := pluralize.NewClient().Plural(animal.Name)
	lowerPlural := strings.ToLower(upperPlural)

	apiID := strings.ToLower(lowerPlural)
	fmt.Printf("generating %+v API\n", apiID)

	// Create output directory.
	err := os.MkdirAll(filepath.Join(runtime, apiID), 0755)
	if err != nil {
		return err
	}

	// Generate info.yaml.
	proxyDependenciesData, err := encoding.NodeForMessage(&apihub.ReferenceList{
		DisplayName: "Apigee Dependencies",
		Description: "Links to dependant Apigee resources.",
		References: []*apihub.ReferenceList_Reference{
			{
				Id:          "petstore",
				DisplayName: "petstore (Apigee)",
				Category:    "",
				Uri:         "https://console.cloud.google.com/apigee/proxies/petstore/overview?project=apigee-apihub-demo",
			},
		},
	})
	if err != nil {
		return err
	}
	proxy := &encoding.Api{
		Header: encoding.Header{
			ApiVersion: "apigeeregistry/v1",
			Kind:       "API",
			Metadata: encoding.Metadata{
				Name: organization + "-" + provider + "-" + apiID + "-proxy",
				Labels: map[string]string{
					"apihub-business-unit": organization,
					"apihub-kind":          "proxy",
					"source":               source,
				},
				Annotations: map[string]string{
					"apigee-proxy": organization + "/apis/" + provider + "-" + apiID,
				},
			},
		},
		Data: encoding.ApiData{
			DisplayName:        organization + " proxy: " + provider + "-" + apiID,
			RecommendedVersion: "v1",
			ApiDeployments: []*encoding.ApiDeployment{
				{
					Header: encoding.Header{
						Metadata: encoding.Metadata{
							Name: "bar-org",
							Labels: map[string]string{
								"apihub-gateway": "apihub-google-cloud-apigee",
							},
							Annotations: map[string]string{
								"apigee-envgroup":       "organizations/apigee-apihub-demo/envgroups/bar",
								"apigee-environment":    "organizations/apigee-apihub-demo/environments/test-env",
								"apigee-proxy-revision": "organizations/apigee-apihub-demo/apis/petstore/revisions/1",
								"message-count-7-days":  "2",
							},
						},
					},
					Data: encoding.ApiDeploymentData{
						DisplayName: "test-env (bar.org)",
						EndpointURI: "bar.org",
					},
				},
			},
			Artifacts: []*encoding.Artifact{
				{
					Header: encoding.Header{
						Kind: "ReferenceList",
						Metadata: encoding.Metadata{
							Name: "apihub-dependencies",
							Labels: map[string]string{
								"source": source,
							},
						},
					},
					Data: *proxyDependenciesData,
				},
			},
		},
	}
	productProxiesData, err := encoding.NodeForMessage(&apihub.ReferenceList{
		DisplayName: "Related Resources",
		Description: "Links to resources in the registry.",
		References: []*apihub.ReferenceList_Reference{
			{
				Id:          "apigee-apihub-demo-petstore-proxy",
				DisplayName: "apigee-apihub-demo proxy: petstore",
				Resource:    "projects/apigee-apihub-demo/locations/global/apis/apigee-apihub-demo-petstore-proxy",
			},
		},
	})
	if err != nil {
		return err
	}
	productDependenciesData, err := encoding.NodeForMessage(&apihub.ReferenceList{
		DisplayName: "Apigee Dependencies",
		Description: "Links to dependant Apigee resources.",
		References: []*apihub.ReferenceList_Reference{
			{
				Id:          "petstore",
				DisplayName: "petstore (Apigee)",
				Category:    "",
				Uri:         "https://console.cloud.google.com/apigee/apiproducts/product/petstore/overview?project=apigee-apihub-demo",
			},
		},
	})
	if err != nil {
		return err
	}
	product := &encoding.Api{
		Header: encoding.Header{
			ApiVersion: "apigeeregistry/v1",
			Kind:       "API",
			Metadata: encoding.Metadata{
				Name: organization + "-" + provider + "-" + apiID + "-product",
				Labels: map[string]string{
					"apihub-business-unit": organization,
					"apihub-kind":          "product",
					"apihub-target-users":  "public",
				},
				Annotations: map[string]string{
					"apigee-product": "organizations/apigee-apihub-demo/apiproducts/petstore",
				},
			},
		},
		Data: encoding.ApiData{
			DisplayName: organization + " product: " + provider + "-" + apiID,
			Artifacts: []*encoding.Artifact{
				{
					Header: encoding.Header{
						Kind: "ReferenceList",
						Metadata: encoding.Metadata{
							Name: "apihub-related",
							Labels: map[string]string{
								"source": source,
							},
						},
					},
					Data: *productProxiesData,
				},
				{
					Header: encoding.Header{
						Kind: "ReferenceList",
						Metadata: encoding.Metadata{
							Name: "apihub-dependencies",
							Labels: map[string]string{
								"source": source,
							},
						},
					},
					Data: *productDependenciesData,
				},
			},
		},
	}
	list := &encoding.List{
		Header: encoding.Header{ApiVersion: encoding.RegistryV1},
		Items:  []interface{}{proxy, product},
	}
	infoBytes, err := encoding.EncodeYAML(list)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(runtime, apiID, "info.yaml"), infoBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
