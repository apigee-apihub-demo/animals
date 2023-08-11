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

func trafficCommand() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "traffic",
		Short: "Generate YAML for mock traffic signals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			animals, err := readAnimals(filename)
			if err != nil {
				return err
			}
			for i, animal := range animals.Animals {
				if err = generateTraffic(i, &animal); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&filename, "filename", "animals.yaml", "Animals definition file (YAML)")
	return cmd
}

const traffic = "traffic"

func generateTraffic(id int, animal *Animal) error {
	upperPlural := pluralize.NewClient().Plural(animal.Name)
	lowerPlural := strings.ToLower(upperPlural)

	trafficID := fmt.Sprintf("%06d", id)

	enrolledApiID := provider + "-" + strings.ToLower(lowerPlural)

	trafficApiID := source + "-" + encodeName(organization+"-traffic-"+trafficID)
	fmt.Printf("generating %+v API\n", trafficApiID)

	// Create output directory.
	err := os.MkdirAll(filepath.Join(traffic, trafficID), 0755)
	if err != nil {
		return err
	}

	// Generate info.yaml.
	referenceData, err := encoding.NodeForMessage(&apihub.ReferenceList{
		DisplayName: "Related Links",
		Description: "Defines a list of related resources",
		References: []*apihub.ReferenceList_Reference{
			{
				Id:          enrolledApiID,
				DisplayName: "Enrolled API",
				Category:    "enrollment",
				Resource:    "projects/apigee-apihub-demo/locations/global/apis/" + enrolledApiID,
			},
		},
	})
	if err != nil {
		return err
	}
	api := &encoding.Api{
		Header: encoding.Header{
			ApiVersion: "apigeeregistry/v1",
			Kind:       "API",
			Metadata: encoding.Metadata{
				Name: trafficApiID,
				Labels: map[string]string{
					"apihub-kind":   "traffic",
					"apihub-source": source,
				},
				Annotations: map[string]string{},
			},
		},
		Data: encoding.ApiData{
			DisplayName: fmt.Sprintf("Traffic Observation %s", trafficID),
			ApiVersions: []*encoding.ApiVersion{
				{
					Header: encoding.Header{
						Metadata: encoding.Metadata{
							Name: "v1",
							Labels: map[string]string{
								"apihub-source": source,
								"apihub-kind":   "traffic",
							},
						},
					},
					Data: encoding.ApiVersionData{
						DisplayName: "v1",
						PrimarySpec: "openapi",
						ApiSpecs: []*encoding.ApiSpec{
							{
								Header: encoding.Header{
									Metadata: encoding.Metadata{
										Name: "openapi",
										Labels: map[string]string{
											"apihub-source": source,
											"apihub-kind":   "traffic",
										},
									},
								},
								Data: encoding.ApiSpecData{
									MimeType: "application/x.openapi+gzip;version=3.0",
									FileName: "openapi.yaml",
								},
							},
						},
					},
				},
			},
			ApiDeployments: []*encoding.ApiDeployment{
				{
					Header: encoding.Header{
						Metadata: encoding.Metadata{
							Name: "location",
						},
					},
					Data: encoding.ApiDeploymentData{
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
						},
					},
					Data: *referenceData,
				},
			},
		},
	}
	infoBytes, err := encoding.EncodeYAML(api)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(traffic, trafficID, "info.yaml"), infoBytes, 0644)
	if err != nil {
		return err
	}

	// Generate openapi.yaml
	f := newOperationIDFactory()
	openapi := &OpenAPI{
		OpenAPI: "3.1.0",
		Paths: map[string]*OpenAPIPath{
			"/" + lowerPlural: {
				Get: &OpenAPIOperation{
					OperationID: f.next(),
				},
				Post: &OpenAPIOperation{
					OperationID: f.next(),
				},
			},
			"/" + lowerPlural + "/*": {
				Get: &OpenAPIOperation{
					OperationID: f.next(),
				},
			},
		},
	}
	specBytes, err := encoding.EncodeYAML(openapi)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(traffic, trafficID, "openapi.yaml"), specBytes, 0644)
}

type operationIDFactory struct {
	i int
}

func newOperationIDFactory() *operationIDFactory {
	return &operationIDFactory{i: 0}
}

func (f *operationIDFactory) next() string {
	s := fmt.Sprintf("op-%03d", f.i)
	f.i++
	return s
}
