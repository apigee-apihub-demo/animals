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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const enrollments = "apis"
const provider = "fauna"
const source = "generate-animal-apis"

func enrolledCommand() *cobra.Command {
	var filename string
	var project string
	cmd := &cobra.Command{
		Use:   "enrolled",
		Short: "Generate YAML for enrolled APIs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			animals, err := readAnimals(filename)
			if err != nil {
				return err
			}
			for _, animal := range animals.Animals {
				if err = generateEnrollment(project, &animal); err != nil {
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

func generateEnrollment(project string, animal *Animal) error {
	upperSingular := animal.Name
	lowerSingular := strings.ToLower(upperSingular)
	upperPlural := pluralize.NewClient().Plural(animal.Name)
	lowerPlural := strings.ToLower(upperPlural)

	apiID := strings.ToLower(lowerPlural)
	fmt.Printf("generating %+v API\n", apiID)

	// Create output directory.
	err := os.MkdirAll(filepath.Join(enrollments, apiID), 0755)
	if err != nil {
		return err
	}

	// Generate info.yaml.
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
					"apihub-source":        source,
					"apihub-team":          provider + "-" + strings.ToLower(animal.Class),
					"categories":           strings.ToLower(animal.Class),
				},
				Annotations: map[string]string{
					"apihub-primary-contact":             lowerPlural + "@apigee-apihub-demo.github.io",
					"apihub-primary-contact-description": upperSingular + " Support Team",
				},
			},
		},
		Data: encoding.ApiData{
			DisplayName:           displayName,
			Description:           "The " + cases.Title(language.English).String(provider) + " " + apiTitle + " API allows users to manage a collection of " + lowerPlural + ".",
			RecommendedVersion:    "v1",
			RecommendedDeployment: "prod",
			ApiVersions: []*encoding.ApiVersion{
				{
					Header: encoding.Header{
						Metadata: encoding.Metadata{
							Name: "v1",
							Labels: map[string]string{
								"apihub-source": source,
								"apihub-kind":   "enrolled",
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
											"apihub-kind":   "enrolled",
											"apihub-source": source,
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
							Name: "prod",
							Labels: map[string]string{
								"apihub-source":  source,
								"apihub-kind":    "enrolled",
								"apihub-gateway": "apihub-google-cloud-apigee",
							},
							Annotations: map[string]string{
								"apihub-external-channel-name": fmt.Sprintf("%s Developer Portal", upperPlural),
							},
						},
					},
					Data: encoding.ApiDeploymentData{
						DisplayName:        "Prod",
						Description:        "Production deployment",
						ApiSpecRevision:    "v1/specs/openapi",
						EndpointURI:        fmt.Sprintf("https://%s.example.com", lowerPlural),
						ExternalChannelURI: fmt.Sprintf("https://%s-portal.example.com/%s", provider, lowerPlural),
						IntendedAudience:   "Public",
						AccessGuidance:     fmt.Sprintf("Contact the %s Support Team for an API Key.", upperSingular),
					},
				},
				{
					Header: encoding.Header{
						Metadata: encoding.Metadata{
							Name: "test",
							Labels: map[string]string{
								"apihub-source":  source,
								"apihub-kind":    "enrolled",
								"apihub-gateway": "apihub-google-cloud-apigee",
							},
							Annotations: map[string]string{
								"apihub-external-channel-name": fmt.Sprintf("%s Developer Portal", upperPlural),
							},
						},
					},
					Data: encoding.ApiDeploymentData{
						DisplayName:        "Test",
						Description:        "Test deployment",
						ApiSpecRevision:    "v1/specs/openapi",
						EndpointURI:        fmt.Sprintf("https://%s-test.example.com", lowerPlural),
						ExternalChannelURI: fmt.Sprintf("https://%s-portal.example.com/%s", provider, lowerPlural),
						IntendedAudience:   "Internal",
						AccessGuidance:     fmt.Sprintf("Contact the %s Support Team for an API Key.", upperSingular),
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
					Data: *referenceData,
				},
			},
		},
	}
	infoBytes, err := encoding.EncodeYAML(api)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(enrollments, apiID, "info.yaml"), infoBytes, 0644)
	if err != nil {
		return err
	}

	// Generate openapi.yaml.
	openapi := &OpenAPI{
		OpenAPI: "3.0.0",
		Info: &OpenAPIInfo{
			Version:     "1.0.0",
			Title:       displayName,
			Description: fmt.Sprintf("Manage a collection of %s.", lowerPlural),
			License: &OpenAPIInfoLicense{
				Name: "Apache 2.0",
				Url:  "http://www.apache.org/licenses/LICENSE-2.0.html",
			},
			Contact: &OpenAPIInfoContact{
				Name:  upperSingular + " Support Team",
				Email: lowerPlural + "@apigee-apihub-demo.github.io",
				Url:   "https://github.com/apigee-apihub-demo/animals",
			},
		},
		Tags: []*OpenAPITag{
			{
				Name:        lowerPlural,
				Description: fmt.Sprintf("This API is about %s.", lowerPlural),
			},
		},
		Servers: []*OpenAPIServer{
			{
				Url: "https://apigee-apihub-demo.github.io",
			},
		},
		Paths: map[string]*OpenAPIPath{
			"/" + lowerPlural: {
				Get: &OpenAPIOperation{
					Summary:     "List all " + lowerPlural,
					Description: "List all " + lowerPlural,
					OperationID: "list" + upperPlural,
					Tags:        []string{lowerPlural},
					Parameters: []*OpenAPIParameter{
						{
							Name:        "limit",
							In:          "query",
							Description: "How many " + lowerPlural + " to return at one time (max 100)",
							Required:    false,
							Schema: &OpenAPISchema{
								Type:   "integer",
								Format: "int32",
							},
						},
					},
					Responses: map[string]*OpenAPIResponse{
						"200": {
							Description: "A paged array of " + lowerPlural,
							Headers: map[string]*OpenAPIHeader{
								"x-next": {
									Description: "A link to the next page of responses",
									Schema: &OpenAPISchema{
										Type: "string",
									},
								},
							},
							Content: map[string]*OpenAPIContent{
								"application/json": {
									Schema: &OpenAPISchema{
										Ref: "#/components/schemas/" + upperPlural,
									},
								},
							},
						},
						"default": errorResponse(),
					},
				},
				Post: &OpenAPIOperation{
					Summary:     "Create " + lowerSingular,
					Description: "Create " + lowerSingular,
					OperationID: "create" + upperSingular,
					Tags:        []string{lowerPlural},
					Responses: map[string]*OpenAPIResponse{
						"201": {
							Description: "Null response",
						},
						"default": errorResponse(),
					},
				},
			},
			"/" + lowerPlural + "/{" + lowerSingular + "Id}": {
				Get: &OpenAPIOperation{
					Summary:     "Info for a specific " + lowerSingular,
					Description: "Info for a specific " + lowerSingular,
					OperationID: "get" + upperSingular,
					Tags:        []string{lowerPlural},
					Parameters: []*OpenAPIParameter{
						{
							Name:        lowerSingular + "Id",
							In:          "path",
							Description: "The id of the " + lowerSingular + " to retrieve",
							Required:    true,
							Schema: &OpenAPISchema{
								Type: "string",
							},
						},
					},
					Responses: map[string]*OpenAPIResponse{
						"200": {
							Description: "Expected response to a valid request",
							Content: map[string]*OpenAPIContent{
								"application/json": {
									Schema: &OpenAPISchema{
										Ref: "#/components/schemas/" + upperSingular,
									},
								},
							},
						},
						"default": errorResponse(),
					},
				},
			},
		},
		Components: &OpenAPIComponents{
			Schemas: map[string]*OpenAPISchema{
				upperSingular: {
					Type:     "object",
					Required: []string{"id", "name"},
					Properties: map[string]*OpenAPISchema{
						"id":   {Type: "integer", Format: "int64"},
						"name": {Type: "string"},
						"tag":  {Type: "string"},
					},
				},
				upperPlural: {
					Type:     "array",
					MaxItems: 100,
					Items:    &OpenAPISchema{Ref: "#/components/schemas/" + upperSingular},
				},
				"Error": {
					Type:     "object",
					Required: []string{"code", "message"},
					Properties: map[string]*OpenAPISchema{
						"code":    {Type: "integer", Format: "int32"},
						"message": {Type: "string"},
					},
				},
			},
		},
	}
	specBytes, err := encoding.EncodeYAML(openapi)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(enrollments, apiID, "openapi.yaml"), specBytes, 0644)
}
