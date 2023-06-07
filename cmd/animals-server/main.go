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
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// animal represents data about an animal.
type animal struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

// animals slice to seed animal data.
var animals = []animal{
	{ID: "1", Name: "Tardar Sauce", Tag: "cat"},
	{ID: "2", Name: "Bo", Tag: "dog"},
	{ID: "3", Name: "Toto", Tag: "dog"},
}

func main() {
	var collection string
	flag.StringVar(&collection, "animals", "animals", "the animals in this collection")
	flag.Parse()

	router := gin.Default()
	router.GET(fmt.Sprintf("/v1/%s", collection), getAnimals)
	router.GET(fmt.Sprintf("/v1/%s/:id", collection), getAnimalByID)
	router.POST(fmt.Sprintf("/v1/%s", collection), postAnimals)

	router.Run("0.0.0.0:8080")
}

// getAnimals responds with the list of all animals as JSON.
func getAnimals(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, animals)
}

// postAnimals adds an animal from JSON received in the request body.
func postAnimals(c *gin.Context) {
	var newAnimal animal

	// Call BindJSON to bind the received JSON to newAnimal.
	if err := c.BindJSON(&newAnimal); err != nil {
		return
	}

	// Add the new animal to the slice.
	animals = append(animals, newAnimal)
	c.IndentedJSON(http.StatusCreated, newAnimal)
}

// getAnimalByID locates the animal whose ID value matches the id
// parameter sent by the client, then returns that animal as a response.
func getAnimalByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of animals, looking for
	// an animal whose ID value matches the parameter.
	for _, a := range animals {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "animal not found",
	})
}
