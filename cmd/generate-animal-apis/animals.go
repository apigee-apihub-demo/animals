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
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

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

func readAnimals(filename string) (*Animals, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var animals Animals
	err = yaml.Unmarshal(yamlFile, &animals)
	if err != nil {
		return nil, err
	}
	return &animals, nil
}
