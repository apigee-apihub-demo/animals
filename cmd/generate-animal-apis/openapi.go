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

type OpenAPI struct {
	OpenAPI    string                  `yaml:"openapi"`
	Info       *OpenAPIInfo            `yaml:"info,omitempty"`
	Servers    []*OpenAPIServer        `yaml:"servers,omitempty"`
	Paths      map[string]*OpenAPIPath `yaml:"paths"`
	Components *OpenAPIComponents      `yaml:"components,omitempty"`
}

type OpenAPIInfo struct {
	Version     string `yaml:"version,omitempty"`
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type OpenAPIServer struct {
	Url string `yaml:"url,omitempty"`
}

type OpenAPIPath struct {
	Get  *OpenAPIOperation `yaml:"get,omitempty"`
	Post *OpenAPIOperation `yaml:"post,omitempty"`
}

type OpenAPIOperation struct {
	Summary     string                      `yaml:"summary,omitempty"`
	OperationID string                      `yaml:"operationId,omitempty"`
	Tags        []string                    `yaml:"tags,omitempty"`
	Parameters  []*OpenAPIParameter         `yaml:"parameters,omitempty"`
	Responses   map[string]*OpenAPIResponse `yaml:"responses,omitempty"`
}

type OpenAPIParameter struct {
	Name        string         `yaml:"name,omitempty"`
	In          string         `yaml:"in,omitempty"`
	Description string         `yaml:"description,omitempty"`
	Required    bool           `yaml:"required,omitempty"`
	Schema      *OpenAPISchema `yaml:"schema,omitempty"`
}

type OpenAPIResponse struct {
	Description string                     `yaml:"description,omitempty"`
	Headers     map[string]*OpenAPIHeader  `yaml:"headers,omitempty"`
	Content     map[string]*OpenAPIContent `yaml:"content,omitempty"`
}

type OpenAPIHeader struct {
	Description string         `yaml:"description,omitempty"`
	Schema      *OpenAPISchema `yaml:"schema,omitempty"`
}

type OpenAPIContent struct {
	Schema *OpenAPISchema `yaml:"schema,omitempty"`
}

type OpenAPIComponents struct {
	Schemas map[string]*OpenAPISchema `yaml:"schemas,omitempty"`
}

type OpenAPISchema struct {
	Type       string                    `yaml:"type,omitempty"`
	Format     string                    `yaml:"format,omitempty"`
	Ref        string                    `yaml:"$ref,omitempty"`
	Required   []string                  `yaml:"required,omitempty"`
	Properties map[string]*OpenAPISchema `yaml:"properties,omitempty"`
	MaxItems   int                       `yaml:"maxItems,omitempty"`
	Items      *OpenAPISchema            `yaml:"items,omitempty"`
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
