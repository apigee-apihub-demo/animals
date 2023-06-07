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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apigee/registry/pkg/application/apihub"
	"github.com/apigee/registry/pkg/encoding"
	pluralize "github.com/gertd/go-pluralize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const out = "apis"
const provider = "fauna"
const source = "generate-animal-apis"

func main() {
	if err := processAnimals("animals.yaml"); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

type Animals struct {
	Animals []Animal `yaml:"animals"`
}

type Animal struct {
	Name     string `yaml:"name"`
	Class    string `yaml:"class"`
	Legs     string `yaml:"legs"`
	Weight   string `yaml:"weight"`
	Lifespan string `yaml:"lifespan"`
}

func processAnimals(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var animals Animals
	err = yaml.Unmarshal(yamlFile, &animals)
	if err != nil {
		return err
	}
	for _, animal := range animals.Animals {
		if err = animal.generateAPI(); err != nil {
			return err
		}
	}
	return nil
}

func (animal *Animal) generateAPI() error {
	upperSingular := animal.Name
	lowerSingular := strings.ToLower(upperSingular)
	upperPlural := pluralize.NewClient().Plural(animal.Name)
	lowerPlural := strings.ToLower(upperPlural)

	apiID := strings.ToLower(lowerPlural)
	fmt.Printf("generating %+v API\n", apiID)

	// Create output directory.
	err := os.MkdirAll(filepath.Join(out, apiID), 0755)
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
	err = os.WriteFile(filepath.Join(out, apiID, "info.yaml"), infoBytes, 0644)
	if err != nil {
		return err
	}

	// Generate openapi.yaml.
	openapi := &OpenAPI{
		OpenAPI: "3.0.0",
		Info: &OpenAPIInfo{
			Version: "1.0.0",
			Title:   displayName,
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
	return os.WriteFile(filepath.Join(out, apiID, "openapi.yaml"), specBytes, 0644)
}

func errorResponse() *OpenAPIResponse {
	return &OpenAPIResponse{
		Description: "unexpected error",
		Content: map[string]*OpenAPIContent{
			"application/json": {
				Schema: &OpenAPISchema{
					Ref: "#/components/schemas/Error",
				},
			},
		},
	}
}
