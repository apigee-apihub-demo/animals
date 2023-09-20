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

func runtimeCommand() *cobra.Command {
	var filename string
	var project string
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
				if err = generateRuntimeMocks(project, &animal); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&filename, "filename", "animals.yaml", "Animals definition file (YAML)")
	cmd.Flags().StringVar(&project, "project", "apigee-apihub-demo", "API Hub GCP project")
	return cmd
}

func generateRuntimeMocks(project string, animal *Animal) error {
	upperPlural := pluralize.NewClient().Plural(animal.Name)
	lowerPlural := strings.ToLower(upperPlural)

	apiID := strings.ToLower(lowerPlural)
	fmt.Printf("generating %+v API\n", apiID)

	enrolledApiID := provider + "-" + strings.ToLower(lowerPlural)
	proxyApiID := source + "-" + encodeName(project+"-proxy-"+enrolledApiID)
	productApiID := source + "-" + encodeName(project+"-product-"+enrolledApiID)

	// Create output directory.
	err := os.MkdirAll(filepath.Join(runtime, apiID), 0755)
	if err != nil {
		return err
	}

	// Generate info.yaml.
	proxyRelatedResourcesData, err := encoding.NodeForMessage(&apihub.ReferenceList{
		DisplayName: "Related Resources",
		Description: "Links to resources in the registry.",
		References: []*apihub.ReferenceList_Reference{
			{
				Id:          enrolledApiID,
				DisplayName: enrolledApiID,
				Category:    "apihub-organization-apis",
				Resource:    "projects/" + project + "/locations/global/apis/" + enrolledApiID,
			},
			{
				Id:          project + "-petstore-product",
				DisplayName: project + " product: petstore",
				Category:    "apihub-organization-apis",
				Resource:    "projects/" + project + "/locations/global/apis/" + productApiID,
			},
		},
	})
	if err != nil {
		return err
	}
	// TODO: Apigee dependencies should link to Apigee management interfaces
	proxyDependenciesData, err := encoding.NodeForMessage(&apihub.ReferenceList{
		DisplayName: "Apigee Dependencies",
		Description: "Links to dependant Apigee resources.",
		References:  []*apihub.ReferenceList_Reference{},
	})
	if err != nil {
		return err
	}
	proxy := &encoding.Api{
		Header: encoding.Header{
			ApiVersion: "apigeeregistry/v1",
			Kind:       "API",
			Metadata: encoding.Metadata{
				Name: proxyApiID,
				Labels: map[string]string{
					"apihub-business-unit": project,
					"apihub-kind":          "proxy",
					"apihub-source":        source,
					"apihub-lifecycle":     "observed-proxy",
				},
				Annotations: map[string]string{
					"proxy": project + "/apis/" + provider + "-" + apiID,
				},
			},
		},
		Data: encoding.ApiData{
			DisplayName:        "Proxy " + project + "-" + provider + "-" + apiID,
			RecommendedVersion: "v1",
			ApiDeployments: []*encoding.ApiDeployment{
				{
					Header: encoding.Header{
						Metadata: encoding.Metadata{
							Name: "bar-org",
							Labels: map[string]string{
								"apihub-gateway": "apihub-google-cloud-apigee",
								"apihub-kind":    "proxy",
							},
							Annotations: map[string]string{
								"organization":         project,
								"envgroup":             "demo",
								"environment":          "test-env",
								"proxy":                "petstore",
								"proxy-revision":       "1",
								"message-count-7-days": "2",
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
							Name: "apihub-related",
							Labels: map[string]string{
								"apihub-source": source,
							},
						},
					},
					Data: *proxyRelatedResourcesData,
				},
				{
					Header: encoding.Header{
						Kind: "ReferenceList",
						Metadata: encoding.Metadata{
							Name: "apihub-dependencies",
							Labels: map[string]string{
								"apihub-source": source,
							},
						},
					},
					Data: *proxyDependenciesData,
				},
			},
		},
	}
	// Related resources link to enrolled APIs and proxies.
	// TODO: Evaluate keeping these in separate lists.
	productRelatedResourcesData, err := encoding.NodeForMessage(&apihub.ReferenceList{
		DisplayName: "Related Resources",
		Description: "Links to resources in the registry.",
		References: []*apihub.ReferenceList_Reference{
			{
				Id:          enrolledApiID,
				DisplayName: enrolledApiID,
				Category:    "apihub-organization-apis",
				Resource:    "projects/" + project + "/locations/global/apis/" + enrolledApiID,
			},
			{
				Id:          project + "-petstore-proxy",
				DisplayName: project + " proxy: petstore",
				Category:    "apihub-organization-apis",
				Resource:    "projects/" + project + "/locations/global/apis/" + proxyApiID,
			},
		},
	})
	if err != nil {
		return err
	}
	// TODO: Apigee dependencies should link to Apigee management interfaces
	productDependenciesData, err := encoding.NodeForMessage(&apihub.ReferenceList{
		DisplayName: "Apigee Dependencies",
		Description: "Links to dependant Apigee resources.",
		References:  []*apihub.ReferenceList_Reference{},
	})
	if err != nil {
		return err
	}
	product := &encoding.Api{
		Header: encoding.Header{
			ApiVersion: "apigeeregistry/v1",
			Kind:       "API",
			Metadata: encoding.Metadata{
				Name: productApiID,
				Labels: map[string]string{
					"apihub-business-unit": project,
					"apihub-kind":          "product",
					"apihub-target-users":  "public",
					"apihub-lifecycle":     "observed-product",
				},
				Annotations: map[string]string{
					"organization": project,
					"product":      "petstore",
				},
			},
		},
		Data: encoding.ApiData{
			DisplayName: "Product " + project + "-" + provider + "-" + apiID,
			Artifacts: []*encoding.Artifact{
				{
					Header: encoding.Header{
						Kind: "ReferenceList",
						Metadata: encoding.Metadata{
							Name: "apihub-related",
							Labels: map[string]string{
								"apihub-source": source,
							},
						},
					},
					Data: *productRelatedResourcesData,
				},
				{
					Header: encoding.Header{
						Kind: "ReferenceList",
						Metadata: encoding.Metadata{
							Name: "apihub-dependencies",
							Labels: map[string]string{
								"apihub-source": source,
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
