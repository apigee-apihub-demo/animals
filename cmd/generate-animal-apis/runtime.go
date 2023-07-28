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
	"time"

	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/encoding"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const runtime = "runtime"

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
	upperSingular := animal.Name
	//lowerSingular := strings.ToLower(upperSingular)
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
	updated := time.Now().Format("2006-01-02")
	referenceData, err := encoding.NodeForMessage(&apihub.ReferenceList{
		DisplayName: "Related Links",
		References: []*apihub.ReferenceList_Reference{
			{
				Id:          "wikipedia",
				DisplayName: "Wikipedia",
				Category:    "reference",
				Uri:         "https://en.wikipedia.org/wiki/" + upperSingular,
			},
		},
	})
	if err != nil {
		return err
	}
	apiTitle := cases.Title(language.English).String(apiID)
	displayName := cases.Title(language.English).String(provider) + " " + apiTitle + " API"
	api := &encoding.Api{
		Header: encoding.Header{
			ApiVersion: "apigeeregistry/v1",
			Kind:       "API",
			Metadata: encoding.Metadata{
				Name: provider + "-" + apiID,
				Labels: map[string]string{
					"apihub-business-unit": provider,
					"apihub-lifecycle":     "production",
					"apihub-style":         "apihub-openapi",
					"apihub-target-users":  "public",
					"apihub-kind":          "enrolled",
					"apihub-team":          provider + "-" + strings.ToLower(animal.Class),
					"categories":           strings.ToLower(animal.Class),
					"provider":             provider,
					"updated":              updated,
					"source":               source,
				},
				Annotations: map[string]string{
					"apihub-primary-contact":             lowerPlural + "@apigee-apihub-demo.github.io",
					"apihub-primary-contact-description": upperSingular + " Support Team",
					"legs":                               animal.Legs,
					"weight":                             animal.Weight,
					"lifespan":                           animal.Lifespan,
				},
			},
		},
		Data: encoding.ApiData{
			DisplayName:        displayName,
			Description:        "The " + cases.Title(language.English).String(provider) + " " + apiTitle + " API allows users to manage a collection of " + lowerPlural + ".",
			RecommendedVersion: "v1",
			ApiVersions: []*encoding.ApiVersion{
				{
					Header: encoding.Header{
						Metadata: encoding.Metadata{
							Name: "v1",
							Labels: map[string]string{
								"updated": updated,
								"source":  source,
							},
						},
					},
					Data: encoding.ApiVersionData{
						DisplayName: "v1",
						State:       "production",
						PrimarySpec: "openapi",
						ApiSpecs: []*encoding.ApiSpec{
							{
								Header: encoding.Header{
									Metadata: encoding.Metadata{
										Name: "openapi",
										Labels: map[string]string{
											"updated": updated,
											"source":  source,
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
			Artifacts: []*encoding.Artifact{
				{
					Header: encoding.Header{
						Kind: "ReferenceList",
						Metadata: encoding.Metadata{
							Name: "apihub-related",
							Labels: map[string]string{
								"updated": updated,
								"source":  source,
							},
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
	err = os.WriteFile(filepath.Join(runtime, apiID, "info.yaml"), infoBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
