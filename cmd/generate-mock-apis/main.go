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
	"time"

	"github.com/apigee/registry/pkg/encoding"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const numberOfAPIs = 200
const numberOfVersions = 25
const numberOfDeployments = 25

const provider = "sample"
const source = "generate-mock-apis"

var updated = time.Now().Format("2006-01-02")

func main() {
	if err := generate(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func generate() error {
	list := &encoding.List{
		Header: encoding.Header{
			ApiVersion: "apigeeregistry/v1",
		},
	}
	for i := 0; i < numberOfAPIs; i++ {
		list.Items = append(list.Items, generateAPI(i))
	}
	listBytes, err := encoding.EncodeYAML(list)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("list.yaml"), listBytes, 0644)
}

func generateAPI(i int) *encoding.Api {
	apiID := fmt.Sprintf("mock-%03d", i)
	apiTitle := cases.Title(language.English).String(apiID)
	displayName := cases.Title(language.English).String(provider) + " " + apiTitle + " API"
	api := &encoding.Api{
		Header: encoding.Header{
			ApiVersion: "apigeeregistry/v1",
			Kind:       "API",
			Metadata: encoding.Metadata{
				Name: provider + "-" + apiID,
				Labels: map[string]string{
					"provider": provider,
					"updated":  updated,
					"source":   source,
				},
				Annotations: map[string]string{},
			},
		},
		Data: encoding.ApiData{
			DisplayName:    displayName,
			Description:    "A sample API",
			ApiVersions:    generateVersions(),
			ApiDeployments: generateDeployments(),
		},
	}
	return api
}

func generateVersions() []*encoding.ApiVersion {
	versions := []*encoding.ApiVersion{}
	for i := 0; i < numberOfVersions; i++ {
		name := fmt.Sprintf("v%d", i+1)
		versions = append(versions, &encoding.ApiVersion{
			Header: encoding.Header{
				Kind: "Version",
				Metadata: encoding.Metadata{
					Name: name,
					Labels: map[string]string{
						"updated": updated,
						"source":  source,
					},
				},
			},
			Data: encoding.ApiVersionData{
				DisplayName: name,
				State:       "production",
			},
		})
	}
	return versions
}

func generateDeployments() []*encoding.ApiDeployment {
	deployments := []*encoding.ApiDeployment{}
	for i := 0; i < numberOfDeployments; i++ {
		name := fmt.Sprintf("zone-%02d", i)
		deployments = append(deployments, &encoding.ApiDeployment{
			Header: encoding.Header{
				Kind: "Deployment",
				Metadata: encoding.Metadata{
					Name: name,
					Labels: map[string]string{
						"updated": updated,
						"source":  source,
					},
				},
			},
			Data: encoding.ApiDeploymentData{
				DisplayName: name,
				EndpointURI: fmt.Sprintf("https://zone-%02d.example.com", i),
			},
		})
	}
	return deployments
}
