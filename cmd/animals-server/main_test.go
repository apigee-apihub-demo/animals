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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestServer(t *testing.T) {
	port := 9999
	go run(port)

	host := fmt.Sprintf("http://localhost:%d", port)
	arthur := animal{
		ID: "1", Name: "arthur", Tag: "aardvark",
	}
	kermit := animal{
		ID: "2", Name: "kermit", Tag: "frog",
	}
	fozzie := animal{
		ID: "1", Name: "fozzie", Tag: "bear",
	}
	// get animals, should be empty
	{
		resp, err := http.Get(host + "/v1/animals")
		if err != nil {
			t.Fatalf("%s", err)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("%s", err)
		}
		var a []animal
		err = json.Unmarshal(b, &a)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if len(a) != 0 {
			t.Errorf("Error, length should be 0, got %d", len(a))
		}
	}
	// create an animal
	{
		b, err := json.Marshal(arthur)
		if err != nil {
			t.Fatalf("%s", err)
		}
		resp, err := http.Post(host+"/v1/animals", "application/json", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("%s", err)
		}
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("%s", err)
		}
		var a animal
		err = json.Unmarshal(b, &a)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if !cmp.Equal(a, arthur) {
			t.Errorf("expected %+v got %+v", arthur, a)
		}
	}
	// get the animal we created
	{
		resp, err := http.Get(host + "/v1/animals/1")
		if err != nil {
			t.Fatalf("%s", err)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("%s", err)
		}
		var a animal
		err = json.Unmarshal(b, &a)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if !cmp.Equal(a, arthur) {
			t.Errorf("expected %+v got %+v", arthur, a)
		}
	}
	// create another animal
	{
		b, err := json.Marshal(kermit)
		if err != nil {
			t.Fatalf("%s", err)
		}
		resp, err := http.Post(host+"/v1/animals", "application/json", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("%s", err)
		}
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("%s", err)
		}
		var a animal
		err = json.Unmarshal(b, &a)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if !cmp.Equal(a, kermit) {
			t.Errorf("expected %+v got %+v", kermit, a)
		}
	}
	// replace the first animal
	{
		b, err := json.Marshal(fozzie)
		if err != nil {
			t.Fatalf("%s", err)
		}
		resp, err := http.Post(host+"/v1/animals", "application/json", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("%s", err)
		}
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("%s", err)
		}
		var a animal
		err = json.Unmarshal(b, &a)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if !cmp.Equal(a, fozzie) {
			t.Errorf("expected %+v got %+v", fozzie, a)
		}
	}
	// get all animals
	{
		resp, err := http.Get(host + "/v1/animals")
		if err != nil {
			t.Fatalf("%s", err)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("%s", err)
		}
		var a []animal
		err = json.Unmarshal(b, &a)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if len(a) != 2 {
			t.Errorf("Error, length should be 2, got %d", len(a))
		}
	}
}
